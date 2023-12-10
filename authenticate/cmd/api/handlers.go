package main

import (
	db "authenticate/cmd/data/sqlc"
	"authenticate/token"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
)

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	UserID    int64  `json:"user_id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// if error, return errorResponse:
//
//	{
//		"error": "error message"
//	}
//
// if no error, return UserResponse:
//
//	{
//		"user": {
//			"user_id": 1,
//
//			...
//
//	}
func (server *Server) Regist(c *gin.Context) {
	var request CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	hashedPassword, err := HashedPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Email:          request.Email,
		HashedPassword: hashedPassword,
	}
	user, err := server.store.CreateUser(c, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation": // duplicate email
				c.JSON(http.StatusConflict, errorResponse(err))
				return
			}
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	go server.logEventViaRabbit("authenticated", fmt.Sprintf("user %s regist", user.Email), "log.INFO")

	c.JSON(http.StatusOK, gin.H{"user": NewUserResponse(user)})
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  UserResponse `json:"user"`
}

// if error, return errorResponse:
//
//	{
//		"error": "error message"
//	}
//
// if no error, return LoginResponse:
//
//	{
//		"session": {
//			"session_id": "uuid",
//			"access_token": "string",
//			"access_token_expires_at": "time",
//			"refresh_token": "string",
//			"refresh_token_expires_at": "time",
//			"user": {
//				"user_id": 1,
//
//				...
//
//			}
//		}
func (server *Server) Login(c *gin.Context) {
	var request LoginRequest
	var user db.User
	var err error
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err = server.store.GetUser(c, request.Email)
	if err != nil {
		if err == sql.ErrNoRows { // user not found
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if err := ComparePassword(user.HashedPassword, request.Password); err != nil {
		// auth log
		go server.logEventViaRabbit("unauthenticated", fmt.Sprintf("user %s log in failed: mistype password, from IP: %s, UserAgent: %s", user.Email, c.ClientIP(), c.Request.UserAgent()), "log.ERROR")
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// generate paseto token
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Email, server.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// generate refresh token
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Email,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.CreateSession(c, db.CreateSessionParams{
		SessionID:    refreshPayload.ID,
		Email:        user.Email,
		RefreshToken: refreshToken,
		UserAgent:    c.Request.UserAgent(),
		ClientIp:     c.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiresAt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// log auth event
	go server.logEventViaRabbit("authenticated", fmt.Sprintf("user %s log in", user.Email), "log.INFO")

	// set user in redis
	err = server.setSessionInRedisWithExpiry(c, session, server.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := LoginResponse{
		SessionID:             uuid.UUID(session.SessionID),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiresAt,
		User:                  NewUserResponse(user),
	}

	c.JSON(http.StatusOK, gin.H{"session": resp})
}

type RenewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RenewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

// if error, return errorResponse:
//
//	{
//		"error": "error message"
//	}
//
// if no error, return RenewAccessTokenResponse:
//
//	{
//		"renew": {
//			"access_token": "string",
//			"access_token_expires_at": "time"
//		}
//	}
func (server *Server) RenewAccessToken(c *gin.Context) {
	var request RenewAccessTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshPayload, err := server.tokenMaker.VerifyToken(request.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// get user in redis
	session, err := server.getSessionInRedis(c, refreshPayload.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// if empty session, get user in database
	if session == (db.Session{}) {
		session, err = server.store.GetSession(c, refreshPayload.ID)
		if err != nil {
			if err == sql.ErrNoRows { // user not found
				c.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	if session.IsBlocked {
		err = server.deleteSessionInRedis(c, session.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		c.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("session is blocked")))
		return
	}
	if session.Email != refreshPayload.Email {
		err = server.deleteSessionInRedis(c, session.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		c.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("incorrect user")))
		return
	}
	if session.RefreshToken != request.RefreshToken {
		err = server.deleteSessionInRedis(c, session.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		c.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("incorrect refresh token")))
		return
	}
	if time.Now().After(session.ExpiresAt) {
		err = server.deleteSessionInRedis(c, session.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		c.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("session has expired")))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(session.Email, server.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// set user in redis
	err = server.setSessionInRedisWithExpiry(c, session, server.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := RenewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiresAt,
	}

	c.JSON(http.StatusOK, gin.H{"renew": resp})
}

type SaveFavoriteRequest struct {
	GoogleID                 string   `json:"google_id" binding:"required"`
	TwDisplayName            string   `json:"tw_display_name" binding:"required"`
	TwFormattedAddress       string   `json:"tw_formatted_address" binding:"required"`
	TwWeekdayDescriptions    []string `json:"tw_weekday_descriptions" binding:"required"`
	AdministrativeAreaLevel1 string   `json:"administrative_area_level_1"`
	Country                  string   `json:"country"`
	GoogleMapUri             string   `json:"google_map_uri" binding:"required"`
	InternationalPhoneNumber string   `json:"international_phone_number"`
	Lat                      string   `json:"lat" binding:"required"`
	Lng                      string   `json:"lng" binding:"required"`
	PrimaryType              string   `json:"primary_type"`
	Rating                   string   `json:"rating"`
	Types                    []string `json:"types" binding:"required"`
	UserRatingCount          int32    `json:"user_rating_count"`
	WebsiteUri               string   `json:"website_uri"`
}

type SaveFavoriteResponse struct {
	UserID        int64  `json:"user_id"`
	GoogleID      string `json:"google_id"`
	TwDisplayName string `json:"tw_display_name"`
}

func (server *Server) ToggleFavorite(c *gin.Context) {
	var request SaveFavoriteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// check if user exists, if not, return
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	session, err := server.GetSessionInRedisOrDatabase(c, authPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// todo: check place in redis
	// place, err := s.checkPlaceInRedis(c, request.GoogleID)
	// if err != nil { // redis server error
	// 	c.JSON(http.StatusInternalServerError, errorResponse(err))
	// 	return
	// }
	// if place!=nil { // place exists in redis
	// 	c.JSON(http.StatusOK, gin.H{"Place": place})
	// 	return
	// }
	// else: place not exists in redis, then check place in database
	// check if place exists, if yes, update it; if not, create one
	// check if place exists, if not, create one
	arg := db.CreatePlaceParams{
		GoogleID:              request.GoogleID,
		TwDisplayName:         request.TwDisplayName,
		TwFormattedAddress:    request.TwFormattedAddress,
		TwWeekdayDescriptions: pq.StringArray(request.TwWeekdayDescriptions),
		GoogleMapUri:          request.GoogleMapUri,
		Lat:                   request.Lat,
		Lng:                   request.Lng,
		Types:                 pq.StringArray(request.Types),
	}
	if len(request.AdministrativeAreaLevel1) != 0 {
		arg.AdministrativeAreaLevel1 = request.AdministrativeAreaLevel1
	}
	if len(request.Country) != 0 {
		arg.Country = request.Country
	}
	if len(request.InternationalPhoneNumber) != 0 {
		arg.InternationalPhoneNumber = request.InternationalPhoneNumber
	}
	if len(request.PrimaryType) != 0 {
		arg.PrimaryType = request.PrimaryType
	}
	if len(request.Rating) != 0 {
		arg.Rating = request.Rating
	}
	if request.UserRatingCount != 0 {
		arg.UserRatingCount = request.UserRatingCount
	}
	if len(request.WebsiteUri) != 0 {
		arg.WebsiteUri = request.WebsiteUri
	}
	place, err := server.store.GetPlaceByGoogleId(c, request.GoogleID)
	if err != nil {
		if err == sql.ErrNoRows { // place not found
			place, err = server.store.CreatePlace(c, arg)
			if err != nil {
				c.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	// toggle favorite
	fav, err := server.store.ToggleFavorite(c, db.ToggleFavoriteParams{
		UserEmail: session.Email,
		GoogleID:  place.GoogleID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// todo: save place in redis
	// err = s.savePlaceInRedis(c, place)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "Place": place})
	// 	return
	// }
	if fav.IsFavorite {
		c.JSON(http.StatusOK, gin.H{"User": session.Email, "action": "add", "Place": place})
	}
	c.JSON(http.StatusOK, gin.H{"User": session.Email, "action": "remove", "Place": place})
}

type GetListFavoritesRequest struct {
	Limit  int32 `json:"limit" binding:"required"`
	Offset int32 `json:"offset"`
}

func (server *Server) GetListFavorites(c *gin.Context) {
	var request GetListFavoritesRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// check if user exists, if not, return
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	session, err := server.GetSessionInRedisOrDatabase(c, authPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	favs, err := server.store.ListFavoritesByCreateTime(c, db.ListFavoritesByCreateTimeParams{
		UserEmail: session.Email,
		Limit:     request.Limit,
		Offset:    request.Offset,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if len(favs) == 0 {
		c.JSON(http.StatusOK, gin.H{"User": session.Email, "Favorites": []string{}})
	} else {
		c.JSON(http.StatusOK, gin.H{"User": session.Email, "Favorites": favs})
	}
}

type GetListFavoritesByCountryRequest struct {
	Country string `json:"country" binding:"required"`
	Limit   int32  `json:"limit" binding:"required"`
	Offset  int32  `json:"offset"`
}

func (server *Server) GetListFavoritesByCountry(c *gin.Context) {
	var request GetListFavoritesByCountryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// check if user exists, if not, return
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	session, err := server.GetSessionInRedisOrDatabase(c, authPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	favs, err := server.store.ListFavoritesByCountry(c, db.ListFavoritesByCountryParams{
		UserEmail: session.Email,
		Country:   request.Country,
		Limit:     request.Limit,
		Offset:    request.Offset,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if len(favs) == 0 {
		c.JSON(http.StatusOK, gin.H{"User": session.Email, "Favorites": []string{}})
	} else {
		c.JSON(http.StatusOK, gin.H{"User": session.Email, "Favorites": favs})
	}
}

type GetListFavoritesByCountryAndRegionRequest struct {
	Country                  string `json:"country" binding:"required"`
	AdministrativeAreaLevel1 string `json:"administrative_area_level_1"`
	Limit                    int32  `json:"limit" binding:"required"`
	Offset                   int32  `json:"offset"`
}

func (server *Server) GetListFavoritesByCountryAndRegion(c *gin.Context) {
	var request GetListFavoritesByCountryAndRegionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// check if user exists, if not, return
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	session, err := server.GetSessionInRedisOrDatabase(c, authPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	favs, err := server.store.ListFavoritesByCountrAndRegion(c, db.ListFavoritesByCountrAndRegionParams{
		UserEmail:                session.Email,
		Country:                  request.Country,
		AdministrativeAreaLevel1: request.AdministrativeAreaLevel1,
		Limit:                    request.Limit,
		Offset:                   request.Offset,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if len(favs) == 0 {
		c.JSON(http.StatusOK, gin.H{"User": session.Email, "Favorites": []string{}})
	} else {
		c.JSON(http.StatusOK, gin.H{"User": session.Email, "Favorites": favs})
	}
}

func (server *Server) GetFavoritesCountry(c *gin.Context) {
	// check if user exists, if not, return
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	session, err := server.GetSessionInRedisOrDatabase(c, authPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	countries, err := server.store.GetCountryList(c, session.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if len(countries) == 0 {
		c.JSON(http.StatusOK, gin.H{"User": session.Email, "Countries": []string{}})
	} else {
		c.JSON(http.StatusOK, gin.H{"User": session.Email, "Countries": countries})
	}
}

type GetFavoritesRegionRequest struct {
	Country string `json:"country" binding:"required"`
}

func (server *Server) GetFavoritesRegion(c *gin.Context) {
	var request GetFavoritesRegionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// check if user exists, if not, return
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	session, err := server.GetSessionInRedisOrDatabase(c, authPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	regions, err := server.store.GetRegionList(c, db.GetRegionListParams{
		UserEmail: session.Email,
		Country:   request.Country,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"User": session.Email, "Country": request.Country, "Regions": regions})
}

func (server *Server) GetSession(c *gin.Context) {
	// check if user exists, if not, return
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	session, err := server.GetSessionInRedisOrDatabase(c, authPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"User": session.Email})
}

// while load in serch result, call this api to check if the place is in favorite list and update the info if exist.
func (server *Server) CheckAndUpdateFavorite(c *gin.Context) {
	var request SaveFavoriteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// check if user exists, if not, return
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	session, err := server.GetSessionInRedisOrDatabase(c, authPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// todo: check place in redis
	// place = s.checkPlaceInRedis(c, request.GoogleID)
	// if place!=nil { // place exists in redis
	// 	c.JSON(http.StatusOK, gin.H{"isFavorite": true, "User": NewUserResponse(user), "Place": place})
	// 	return
	// }
	// else: place not exists in redis, then check place in database
	arg := db.CreatePlaceParams{
		GoogleID:              request.GoogleID,
		TwDisplayName:         request.TwDisplayName,
		TwFormattedAddress:    request.TwFormattedAddress,
		TwWeekdayDescriptions: pq.StringArray(request.TwWeekdayDescriptions),
		GoogleMapUri:          request.GoogleMapUri,
		Lat:                   request.Lat,
		Lng:                   request.Lng,
		Types:                 pq.StringArray(request.Types),
	}
	if len(request.AdministrativeAreaLevel1) != 0 {
		arg.AdministrativeAreaLevel1 = request.AdministrativeAreaLevel1
	}
	if len(request.Country) != 0 {
		arg.Country = request.Country
	}
	if len(request.InternationalPhoneNumber) != 0 {
		arg.InternationalPhoneNumber = request.InternationalPhoneNumber
	}
	if len(request.PrimaryType) != 0 {
		arg.PrimaryType = request.PrimaryType
	}
	if len(request.Rating) != 0 {
		arg.Rating = request.Rating
	}
	if request.UserRatingCount != 0 {
		arg.UserRatingCount = request.UserRatingCount
	}
	if len(request.WebsiteUri) != 0 {
		arg.WebsiteUri = request.WebsiteUri
	}
	place, err := server.store.GetPlaceByGoogleId(c, request.GoogleID)
	if err != nil {
		if err == sql.ErrNoRows { // place not found
			place, err = server.store.CreatePlace(c, arg)
			if err != nil {
				c.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}
	// update place
	updateArg := db.UpdatePlaceParams{
		PlaceVersion: place.PlaceVersion,
	}
	if len(request.GoogleID) != 0 {
		updateArg.GoogleID = sql.NullString{String: request.GoogleID, Valid: true}
	} else {
		updateArg.GoogleID = sql.NullString{String: "", Valid: false}
	}
	if len(request.TwDisplayName) != 0 {
		updateArg.TwDisplayName = sql.NullString{String: request.TwDisplayName, Valid: true}
	} else {
		updateArg.TwDisplayName = sql.NullString{String: "", Valid: false}
	}
	if len(request.TwFormattedAddress) != 0 {
		updateArg.TwFormattedAddress = sql.NullString{String: request.TwFormattedAddress, Valid: true}
	} else {
		updateArg.TwFormattedAddress = sql.NullString{String: "", Valid: false}
	}
	if len(request.TwWeekdayDescriptions) != 0 {
		updateArg.TwWeekdayDescriptions = pq.StringArray(request.TwWeekdayDescriptions)
	} else {
		updateArg.TwWeekdayDescriptions = pq.StringArray{}
	}
	if len(request.GoogleMapUri) != 0 {
		updateArg.GoogleMapUri = sql.NullString{String: request.GoogleMapUri, Valid: true}
	} else {
		updateArg.GoogleMapUri = sql.NullString{String: "", Valid: false}
	}
	if len(request.Lat) != 0 {
		updateArg.Lat = sql.NullString{String: request.Lat, Valid: true}
	} else {
		updateArg.Lat = sql.NullString{String: "", Valid: false}
	}
	if len(request.Lng) != 0 {
		updateArg.Lng = sql.NullString{String: request.Lng, Valid: true}
	} else {
		updateArg.Lng = sql.NullString{String: "", Valid: false}
	}
	if len(request.Types) != 0 {
		updateArg.Types = pq.StringArray(request.Types)
	} else {
		updateArg.Types = pq.StringArray{}
	}
	if len(request.AdministrativeAreaLevel1) != 0 {
		updateArg.AdministrativeAreaLevel1 = sql.NullString{String: request.AdministrativeAreaLevel1, Valid: true}
	} else {
		updateArg.AdministrativeAreaLevel1 = sql.NullString{String: "", Valid: false}
	}
	if len(request.Country) != 0 {
		updateArg.Country = sql.NullString{String: request.Country, Valid: true}
	} else {
		updateArg.Country = sql.NullString{String: "", Valid: false}
	}
	if len(request.InternationalPhoneNumber) != 0 {
		updateArg.InternationalPhoneNumber = sql.NullString{String: request.InternationalPhoneNumber, Valid: true}
	} else {
		updateArg.InternationalPhoneNumber = sql.NullString{String: "", Valid: false}
	}
	if len(request.PrimaryType) != 0 {
		updateArg.PrimaryType = sql.NullString{String: request.PrimaryType, Valid: true}
	} else {
		updateArg.PrimaryType = sql.NullString{String: "", Valid: false}
	}
	if len(request.Rating) != 0 {
		updateArg.Rating = sql.NullString{String: request.Rating, Valid: true}
	} else {
		updateArg.Rating = sql.NullString{String: "", Valid: false}
	}
	if request.UserRatingCount != 0 {
		updateArg.UserRatingCount = sql.NullInt32{Int32: request.UserRatingCount, Valid: true}
	} else {
		updateArg.UserRatingCount = sql.NullInt32{Int32: 0, Valid: false}
	}
	if len(request.WebsiteUri) != 0 {
		updateArg.WebsiteUri = sql.NullString{String: request.WebsiteUri, Valid: true}
	} else {
		updateArg.WebsiteUri = sql.NullString{String: "", Valid: false}
	}
	fmt.Println(updateArg)

	updatePlace, err := server.store.UpdatePlace(c, updateArg)
	fmt.Println(place)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	} else {
		// if updatePlace is empty, then place is not updated since place version is not matched
		// redefine updatePlace to place
		updatePlace = place
	}
	// check if favorite exists,
	fav, err := server.store.GetFavorite(c, db.GetFavoriteParams{
		UserEmail: session.Email,
		GoogleID:  updatePlace.GoogleID,
	})
	if err != nil {
		if err == sql.ErrNoRows { // favorite not found
			c.JSON(http.StatusOK, gin.H{"isFavorite": false, "User": session.Email, "Place": updatePlace})
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"isFavorite": true, "User": session.Email, "Favorite": fav, "Place": updatePlace})
}
