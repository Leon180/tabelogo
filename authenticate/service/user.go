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

type RegistUserServiceHandelr interface {
	RegistUser(ctx context.Context, user entitymodel.User) (entitymodel.User, error)
}

func NewRegistUserServiceHandler(userWithTransactionRepository repository.UserWithTransactionHandler) RegistUserServiceHandelr {
	return RegistUserServiceHandle{userWithTransactionRepository: userWithTransactionRepository}
}

type RegistUserServiceHandle struct {
	userWithTransactionRepository repository.UserWithTransactionHandler
}

func (s RegistUserServiceHandle) RegistUser(ctx context.Context, user entitymodel.User) (entitymodel.User, error) {
	var (
		createUser entitymodel.User
		err        error
	)
	if err = s.userWithTransactionRepository.WithTransaction(ctx, func(tx *gorm.DB) error {
		// check if user already exists
		existedUser, err := s.userWithTransactionRepository.GetUserByEmail(ctx, user.Email)
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
		if err := s.userWithTransactionRepository.CreateUser(ctx, user); err != nil {
			return err
		}
		createUser = user
		return nil
	}); err != nil {
		return entitymodel.User{}, err
	}
	return createUser, nil
}

type LoginUserServiceHandelr interface {
	LoginUser(ctx context.Context, user entitymodel.User, request entitymodel.Request) (entitymodel.Session, error)
}

func NewLoginUserServiceHandler(
	userAndSessionWithTransactionRepository repository.UserAndSessionWithTransactionHandler,
	tokenMaker token.Maker,
	redisSession redisDB.SessionWithTransactionHandler,
	config *config.Config,
) LoginUserServiceHandelr {
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

func (s LoginUserServiceHandle) LoginUser(ctx context.Context, loginUser entitymodel.User, request entitymodel.Request) (entitymodel.Session, error) {
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
	if err = s.SessionWithTransactionRedis.WithTransaction(ctx, func(tx *redis.Tx) error {
		if session, err = s.SessionWithTransactionRedis.GetSession(ctx, loginUser.Email); err != nil {
			return err
		}
		// if session not exist or expired, need update
		if !session.IsExist() || session.IsAccessTokenExpired() {
			needUpdate = true
			if err = s.SessionWithTransactionRedis.DeleteSession(ctx, loginUser.Email); err != nil {
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

	if err = s.userAndSessionWithTransactionRepository.WithTransaction(ctx, func(tx *gorm.DB) error {
		// get exist user
		existedUser, err = s.userAndSessionWithTransactionRepository.GetUserByEmail(ctx, loginUser.Email)
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
		if loginUser.HashedPassword != existedUser.HashedPassword {
			return errors.LoginUserPasswordNotMatchError
		}

		// check session in db
		existedSession, err = s.userAndSessionWithTransactionRepository.GetSessionByUserID(ctx, session.UserID)
		if err != nil {
			return err
		}

		// generate access token
		if accessToken, accessPayload, err = s.tokenMaker.CreateToken(existedUser.Email, s.config.AccessTokenDuration); err != nil {
			return err
		}
		// generate refresh token
		if refreshToken, refreshPayload, err = s.tokenMaker.CreateToken(existedUser.Email, s.config.RefreshTokenDuration); err != nil {
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
			if err = s.userAndSessionWithTransactionRepository.CreateSession(ctx, session); err != nil {
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
			existedSession.AccessTokenExpiresAt = session.AccessTokenExpiresAt.Add(s.config.RefreshTokenDuration)
			existedSession.ClientIP = request.ClientIP
			existedSession.UserAgent = request.UserAgent
			if err = s.userAndSessionWithTransactionRepository.UpdateSession(ctx, existedSession.ID, existedSession.GetUpdates()); err != nil {
				return err
			}
			return nil
		}

		// if session is exist but expired, delete session and create new session
		regenSession = true
		if err = s.userAndSessionWithTransactionRepository.DeleteSession(ctx, existedSession.ID); err != nil {
			return err
		}
		if err = s.userAndSessionWithTransactionRepository.CreateSession(ctx, session); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return entitymodel.Session{}, err
	}

	expiry := time.Until(session.AccessTokenExpiresAt)
	// if exist session, reset exist session in redis
	if !regenSession {
		if err = s.SessionWithTransactionRedis.SetSession(ctx, loginUser.Email, existedSession, &expiry); err != nil {
			return entitymodel.Session{}, err
		}
		return existedSession, nil
	}
	// else set session in redis
	if err = s.SessionWithTransactionRedis.SetSession(ctx, loginUser.Email, session, &expiry); err != nil {
		return entitymodel.Session{}, err
	}
	return session, nil
}
