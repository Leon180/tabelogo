package service

import (
	"authenticate/config"
	"authenticate/errors"
	"authenticate/model/entitymodel"
	"authenticate/redisDB"
	"authenticate/repository"
	"authenticate/token"
	"authenticate/utility"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type RegistUserServiceHandler interface {
	RegistUser(ctx context.Context, user entitymodel.User) (entitymodel.User, error)
}

func NewRegistUserServiceHandler(userWithTransactionRepository repository.UserWithTransactionHandler) RegistUserServiceHandler {
	return &RegistUserServiceHandle{userWithTransactionRepository: userWithTransactionRepository}
}

type RegistUserServiceHandle struct {
	userWithTransactionRepository repository.UserWithTransactionHandler
}

func (handle *RegistUserServiceHandle) RegistUser(ctx context.Context, user entitymodel.User) (entitymodel.User, error) {
	var (
		createUser entitymodel.User
		err        error
	)
	if err = handle.userWithTransactionRepository.WithTransaction(ctx, func(tx *gorm.DB) error {
		// check if user already exists
		existedUser, err := handle.userWithTransactionRepository.GetUserByEmail(ctx, user.Email)
		if err != nil {
			return err
		}
		if existedUser.IsExist() {
			if existedUser.IsActive() {
				return errors.UserNotExistsError
			} else {
				return errors.UserExistsButNotActive
			}
		}
		if err := handle.userWithTransactionRepository.CreateUser(ctx, user); err != nil {
			return err
		}
		createUser = user
		return nil
	}); err != nil {
		return entitymodel.User{}, err
	}
	return createUser, nil
}

type LoginUserServiceHandler interface {
	LoginUser(ctx context.Context, user entitymodel.User, request entitymodel.Request) (entitymodel.Session, error)
}

func NewLoginUserServiceHandler(
	userAndSessionWithTransactionRepository repository.UserAndSessionWithTransactionHandler,
	tokenMaker token.Maker,
	redisSession redisDB.SessionWithTransactionHandler,
	config *config.Config,
) LoginUserServiceHandler {
	return &LoginUserServiceHandle{
		userAndSessionWithTransactionRepository: userAndSessionWithTransactionRepository,
		tokenMaker:                              tokenMaker,
		SessionWithTransactionRedis:             redisSession,
		config:                                  config,
	}
}

type LoginUserServiceHandle struct {
	userAndSessionWithTransactionRepository repository.UserAndSessionWithTransactionHandler
	tokenMaker                              token.Maker
	SessionWithTransactionRedis             redisDB.SessionWithTransactionHandler
	config                                  *config.Config
}

func (handle *LoginUserServiceHandle) LoginUser(ctx context.Context, loginUser entitymodel.User, request entitymodel.Request) (entitymodel.Session, error) {
	var (
		existedUser    entitymodel.User
		err            error
		accessToken    string
		accessPayload  *token.Payload
		refreshToken   string
		refreshPayload *token.Payload
		session        entitymodel.Session
		existedSession entitymodel.Session
		needUpdate     bool
		regenSession   bool
	)

	// get session in redis
	if err = handle.SessionWithTransactionRedis.WithTransaction(ctx, func(tx *redis.Tx) error {
		if session, err = handle.SessionWithTransactionRedis.GetSession(ctx, loginUser.Email); err != nil {
			return err
		}
		// if session not exist or expired, need update
		if !session.IsExist() || session.IsAccessTokenExpired() {
			needUpdate = true
			if err = handle.SessionWithTransactionRedis.DeleteSession(ctx, loginUser.Email); err != nil {
				return err
			}
			return nil
		}
		return nil
	}, loginUser.Email); err != nil {
		return entitymodel.Session{}, err
	}

	// if session not need update, return session directly
	if !needUpdate {
		return session, nil
	}

	if err = handle.userAndSessionWithTransactionRepository.WithTransaction(ctx, func(tx *gorm.DB) error {
		// get exist user
		existedUser, err = handle.userAndSessionWithTransactionRepository.GetUserByEmail(ctx, loginUser.Email)
		if err != nil {
			return err
		}
		// check if user exists
		if !existedUser.IsExist() {
			return errors.UserNotExistsError
		}
		// check if user is active
		if !existedUser.IsActive() {
			return errors.UserExistsButNotActive
		}
		// check if password match
		if err = utility.CompareHashAndPassword(existedUser.HashedPassword, loginUser.Password); err != nil {
			return errors.LoginUserPasswordNotMatchError
		}

		// check session in db
		existedSession, err = handle.userAndSessionWithTransactionRepository.GetSessionByUserID(ctx, session.UserID)
		if err != nil {
			return err
		}

		// generate access token
		if accessToken, accessPayload, err = handle.tokenMaker.CreateToken(existedUser.Email, handle.config.AccessTokenDuration); err != nil {
			return err
		}
		// generate refresh token
		if refreshToken, refreshPayload, err = handle.tokenMaker.CreateToken(existedUser.Email, handle.config.RefreshTokenDuration); err != nil {
			return err
		}
		session = entitymodel.Session{
			ID:                    utility.GenDefaultUUID(),
			UserID:                existedUser.ID,
			AccessToken:           accessToken,
			RefreshToken:          refreshToken,
			UserAgent:             request.UserAgent,
			ClientIP:              request.ClientIP,
			IsBlocked:             false,
			AccessTokenExpiresAt:  accessPayload.ExpiresAt,
			RefreshTokenExpiresAt: refreshPayload.ExpiresAt,
			CreatedAt:             time.Now(),
		}

		// if session not exist, create session
		if !existedSession.IsExist() {
			regenSession = true
			if err = handle.userAndSessionWithTransactionRepository.CreateSession(ctx, session); err != nil {
				return err
			}
			return nil
		}

		// if session exist:

		// if session is blocked, return error
		if existedSession.Blocked() {
			return errors.SessionBlockedError
		}

		// if session can be updated, update session
		if (existedSession.IsAccessTokenExpired() && !existedSession.IsRefreshTokenExpired()) ||
			existedSession.IsAccessTokenExpired() {
			existedSession.AccessTokenExpiresAt = session.AccessTokenExpiresAt.Add(handle.config.RefreshTokenDuration)
			existedSession.ClientIP = request.ClientIP
			existedSession.UserAgent = request.UserAgent
			if err = handle.userAndSessionWithTransactionRepository.UpdateSession(ctx, existedSession.ID, existedSession.GetUpdates()); err != nil {
				return err
			}
			return nil
		}

		// if session is exist but expired, delete session and create new session
		regenSession = true
		if err = handle.userAndSessionWithTransactionRepository.DeleteSession(ctx, existedSession.ID); err != nil {
			return err
		}
		if err = handle.userAndSessionWithTransactionRepository.CreateSession(ctx, session); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return entitymodel.Session{}, err
	}

	expiry := time.Until(session.AccessTokenExpiresAt)
	// if exist session, reset exist session in redis
	if !regenSession {
		if err = handle.SessionWithTransactionRedis.SetSession(ctx, loginUser.Email, existedSession, &expiry); err != nil {
			return entitymodel.Session{}, err
		}
		return existedSession, nil
	}
	// else set session in redis
	if err = handle.SessionWithTransactionRedis.SetSession(ctx, loginUser.Email, session, &expiry); err != nil {
		return entitymodel.Session{}, err
	}

	return session, nil
}

type RenewAccessTokenServiceHandler interface {
	RenewAccessToken(ctx context.Context, refreshToken string) (entitymodel.Session, error)
}

func NewRenewAccessTokenServiceHandler(
	userAndSessionWithTransactionRepository repository.UserAndSessionWithTransactionHandler,
	tokenMaker token.Maker,
	SessionWithTransactionRedis redisDB.SessionWithTransactionHandler,
	config *config.Config,
) RenewAccessTokenServiceHandler {
	return &RenewAccessTokenServiceHandle{
		userAndSessionWithTransactionRepository: userAndSessionWithTransactionRepository,
		tokenMaker:                              tokenMaker,
		SessionWithTransactionRedis:             SessionWithTransactionRedis,
		config:                                  config,
	}
}

func (handle *RenewAccessTokenServiceHandle) RenewAccessToken(ctx context.Context, refreshToken string) (entitymodel.Session, error) {
	var (
		sessionWithUser entitymodel.SessionPreloadUser
		err             error
		regenSession    bool
	)

	if err = handle.userAndSessionWithTransactionRepository.WithTransaction(ctx, func(tx *gorm.DB) error {
		// get exist session with user
		sessionWithUser, err = handle.userAndSessionWithTransactionRepository.GetSessionWithUserByRefreshToken(ctx, refreshToken)
		if err != nil {
			return err
		}
		// check if user exists
		if !sessionWithUser.IsExist() {
			return errors.UserNotExistsError
		}
		// check if user is active
		if !sessionWithUser.User.IsActive() {
			return errors.UserExistsButNotActive
		}

		// if session expired, return error
		if sessionWithUser.IsRefreshTokenExpired() {
			return errors.SessionExpiredError
		}

		// if session is blocked, return error
		if sessionWithUser.Blocked() {
			return errors.SessionBlockedError
		}

		// if session can be updated, update session
		if (sessionWithUser.IsAccessTokenExpired() && !sessionWithUser.IsRefreshTokenExpired()) ||
			!sessionWithUser.IsAccessTokenExpired() {
			sessionWithUser.AccessTokenExpiresAt = sessionWithUser.AccessTokenExpiresAt.Add(handle.config.RefreshTokenDuration)
			updates := sessionWithUser.GetUpdates()
			if err = handle.userAndSessionWithTransactionRepository.UpdateSession(ctx, sessionWithUser.Session.ID, updates); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return entitymodel.Session{}, err
	}

	expiry := time.Until(sessionWithUser.AccessTokenExpiresAt)
	// if exist session, reset exist session in redis
	if !regenSession {
		if err = handle.SessionWithTransactionRedis.SetSession(ctx, sessionWithUser.User.Email, sessionWithUser.Session, &expiry); err != nil {
			return entitymodel.Session{}, err
		}
		return sessionWithUser.Session, nil
	}
	// else set session in redis
	if err = handle.SessionWithTransactionRedis.SetSession(ctx, sessionWithUser.User.Email, sessionWithUser.Session, &expiry); err != nil {
		return entitymodel.Session{}, err
	}

	return sessionWithUser.Session, nil
}

type RenewAccessTokenServiceHandle struct {
	userAndSessionWithTransactionRepository repository.UserAndSessionWithTransactionHandler
	tokenMaker                              token.Maker
	SessionWithTransactionRedis             redisDB.SessionWithTransactionHandler
	config                                  *config.Config
}

type SaveFavoriteServiceHandler interface {
	SaveFavorite(ctx context.Context, favorite entitymodel.Favorite) error
}

func NewSaveFavoriteServiceHandler(
	userAndPlaceAndFavoriteWithTransactionRepository repository.UserAndPlaceAndFavoriteWithTransactionHandler,
	redisPlace redisDB.PlaceWithTransactionHandler,
	config *config.Config,
) SaveFavoriteServiceHandler {
	return &SaveFavoriteServiceHandle{
		userAndPlaceAndFavoriteWithTransactionRepository: userAndPlaceAndFavoriteWithTransactionRepository,
		redisPlace: redisPlace,
		config:     config,
	}
}

type SaveFavoriteServiceHandle struct {
	userAndPlaceAndFavoriteWithTransactionRepository repository.UserAndPlaceAndFavoriteWithTransactionHandler
	redisPlace                                       redisDB.PlaceWithTransactionHandler
	config                                           *config.Config
}

func (handle *SaveFavoriteServiceHandle) SaveFavorite(ctx context.Context, favorite entitymodel.Favorite) error {
	var (
		err   error
		place entitymodel.Place
	)
	// check if place is exist (redis first, then db if not exist in redis)
	if err = handle.redisPlace.WithTransaction(ctx, func(tx *redis.Tx) error {
		if place, err = handle.redisPlace.GetPlace(ctx, favorite.PlaceGoogleID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	// if place exist in redis, directly create/update favorite and return
	if place.IsExist() {
		if err = handle.userAndPlaceAndFavoriteWithTransactionRepository.WithTransaction(ctx, func(tx *gorm.DB) error {
			return handle.checkExistedFavoriteAndCreateOrUpdate(ctx, favorite)
		}); err != nil {
			return err
		}
		return nil
	}

	// if place not exist in redis, get place from db and check if place is exist,
	// if exist, create favorite and set place in redis
	// if not exist, return error
	if err = handle.userAndPlaceAndFavoriteWithTransactionRepository.WithTransaction(ctx, func(tx *gorm.DB) error {
		if place, err = handle.userAndPlaceAndFavoriteWithTransactionRepository.GetPlaceByGoogleID(ctx, favorite.PlaceGoogleID); err != nil {
			return err
		}
		if !place.IsExist() {
			return errors.PlaceNotExistsError
		}
		if err = handle.redisPlace.SetPlace(ctx, place.GoogleID, place, &handle.config.PlaceRedisExpiry); err != nil {
			return err
		}
		return handle.checkExistedFavoriteAndCreateOrUpdate(ctx, favorite)
	}); err != nil {
		return err
	}
	return nil
}

func (handle *SaveFavoriteServiceHandle) checkExistedFavoriteAndCreateOrUpdate(ctx context.Context, favorite entitymodel.Favorite) error {
	var (
		existedFavorite entitymodel.Favorite
		err             error
	)
	if existedFavorite, err = handle.userAndPlaceAndFavoriteWithTransactionRepository.GetFavoriteByUserIDAndPlaceGoogleID(ctx, favorite.UserID, favorite.PlaceGoogleID); err != nil {
		return err
	}
	if existedFavorite.IsExist() {
		favorite.ID = existedFavorite.ID
		if err = handle.userAndPlaceAndFavoriteWithTransactionRepository.UpdateFavorite(ctx, favorite); err != nil {
			return err
		}
		return nil
	}
	if err = handle.userAndPlaceAndFavoriteWithTransactionRepository.CreateFavorite(ctx, favorite); err != nil {
		return err
	}
	return nil
}

type GetUserFavoritesServiceHandler interface {
	GetUserFavorites(ctx context.Context, request entitymodel.GetUserFavoritesRequest) (entitymodel.UserFavoritePlaces, error)
}

func NewGetUserFavoritesServiceHandler(
	favoriteRepository repository.FavoriteHandler,
) GetUserFavoritesServiceHandler {
	return &GetUserFavoritesServiceHandle{favoriteRepository: favoriteRepository}
}

type GetUserFavoritesServiceHandle struct {
	favoriteRepository repository.FavoriteHandler
}

func (handle *GetUserFavoritesServiceHandle) GetUserFavorites(ctx context.Context, request entitymodel.GetUserFavoritesRequest) (entitymodel.UserFavoritePlaces, error) {
	return handle.favoriteRepository.GetUserFavoritePlaces(ctx, request.UserID, request.Country, request.AdministrativeAreaLevel1, request.OrderBy)
}
