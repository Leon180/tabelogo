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
	if err := server.logEventViaRabbit("authenticated", fmt.Sprintf("user %s regist", user.Email), "log.INFO"); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
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
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUser(c, request.Email)
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
		if err := server.logEventViaRabbit("unauthenticated", fmt.Sprintf("user %s log in failed: mistype password, from IP: %s, UserAgent: %s", user.Email, c.ClientIP(), c.Request.UserAgent()), "log.ERROR"); err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		// if err := server.LogRequest("unauthenticated", fmt.Sprintf("user %s log in failed: mistype password, from IP: %s, UserAgent: %s", user.Email, c.ClientIP(), c.Request.UserAgent())); err != nil {
		// 	c.JSON(http.StatusInternalServerError, errorResponse(err))
		// 	return
		// }
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
	if err := server.logEventViaRabbit("authenticated", fmt.Sprintf("user %s log in", user.Email), "log.INFO"); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// if err := server.LogRequest("authenticated", fmt.Sprintf("user %s log in", user.Email)); err != nil {
	// 	c.JSON(http.StatusInternalServerError, errorResponse(err))
	// 	return
	// }

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

	session, err := server.store.GetSession(c, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows { // user not found
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.IsBlocked {
		c.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("session is blocked")))
		return
	}
	if session.Email != refreshPayload.Email {
		c.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("incorrect user")))
		return
	}
	if session.RefreshToken != request.RefreshToken {
		c.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("incorrect refresh token")))
		return
	}
	if time.Now().After(session.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("session has expired")))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(session.Email, server.config.AccessTokenDuration)
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

type VarifyAccessTokenRequest struct {
	AccessToken string `json:"access_token" binding:"required"`
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
	PlaceID       int64  `json:"place_id"`
	GoogleID      string `json:"google_id"`
	TwDisplayName string `json:"tw_display_name"`
}

func (server *Server) SaveFavorite(c *gin.Context) {
	var request SaveFavoriteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
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
		arg.AdministrativeAreaLevel1 = sql.NullString{
			String: request.AdministrativeAreaLevel1,
			Valid:  true,
		}
	} else {
		arg.AdministrativeAreaLevel1 = sql.NullString{
			String: "",
			Valid:  false,
		}
	}
	if len(request.Country) != 0 {
		arg.Country = sql.NullString{
			String: request.Country,
			Valid:  true,
		}
	} else {
		arg.Country = sql.NullString{
			String: "",
			Valid:  false,
		}
	}
	if len(request.InternationalPhoneNumber) != 0 {
		arg.InternationalPhoneNumber = sql.NullString{
			String: request.InternationalPhoneNumber,
			Valid:  true,
		}
	} else {
		arg.InternationalPhoneNumber = sql.NullString{
			String: "",
			Valid:  false,
		}
	}
	if len(request.PrimaryType) != 0 {
		arg.PrimaryType = sql.NullString{
			String: request.PrimaryType,
			Valid:  true,
		}
	} else {
		arg.PrimaryType = sql.NullString{
			String: "",
			Valid:  false,
		}
	}
	if len(request.Rating) != 0 {
		arg.Rating = sql.NullString{
			String: request.Rating,
			Valid:  true,
		}
	} else {
		arg.Rating = sql.NullString{
			String: "",
			Valid:  false,
		}
	}
	if request.UserRatingCount != 0 {
		arg.UserRatingCount = sql.NullInt32{
			Int32: request.UserRatingCount,
			Valid: true,
		}
	} else {
		arg.UserRatingCount = sql.NullInt32{
			Int32: 0,
			Valid: false,
		}
	}
	if len(request.WebsiteUri) != 0 {
		arg.WebsiteUri = sql.NullString{
			String: request.WebsiteUri,
			Valid:  true,
		}
	} else {
		arg.WebsiteUri = sql.NullString{
			String: "",
			Valid:  false,
		}
	}

	// check if user exists, if not, return
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	user, err := server.store.GetUser(c, authPayload.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// check if place exists, if not, create one
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
		PlaceID: place.PlaceID,
		GoogleID: sql.NullString{
			String: place.GoogleID,
			Valid:  true,
		},
		TwDisplayName: sql.NullString{
			String: request.TwDisplayName,
			Valid:  true,
		},
		TwFormattedAddress: sql.NullString{
			String: request.TwFormattedAddress,
			Valid:  true,
		},
		TwWeekdayDescriptions: pq.StringArray(request.TwWeekdayDescriptions),
		GoogleMapUri: sql.NullString{
			String: request.GoogleMapUri,
			Valid:  true,
		},
		Lat: sql.NullString{
			String: request.Lat,
			Valid:  true,
		},
		Lng: sql.NullString{
			String: request.Lng,
			Valid:  true,
		},
		Types:                    pq.StringArray(request.Types),
		AdministrativeAreaLevel1: arg.AdministrativeAreaLevel1,
		Country:                  arg.Country,
		InternationalPhoneNumber: arg.InternationalPhoneNumber,
		PrimaryType:              arg.PrimaryType,
		Rating:                   arg.Rating,
		UserRatingCount:          arg.UserRatingCount,
		WebsiteUri:               arg.WebsiteUri,
	}

	place, err = server.store.UpdatePlace(c, updateArg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// check if favorite exists, if not, create one
	fav, err := server.store.GetFavorite(c, db.GetFavoriteParams{
		UserID:  user.UserID,
		PlaceID: place.PlaceID,
	})
	if fav.UserID != 0 {
		c.JSON(http.StatusConflict, errorResponse(fmt.Errorf("favorite already exists")))
		return
	}
	if err != nil {
		if err == sql.ErrNoRows { // favorite not found
			fav, err = server.store.CreateFavorite(c, db.CreateFavoriteParams{
				UserID:  user.UserID,
				PlaceID: place.PlaceID,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"User": NewUserResponse(user), "Favorite": fav, "Place": place})
}

type RemoveFavoriteRequest struct {
	GoogleID string `json:"google_id" binding:"required"`
}

func (server *Server) RemoveFavorite(c *gin.Context) {
	var request RemoveFavoriteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// check if user exists, if not, return
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	user, err := server.store.GetUser(c, authPayload.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// check if place exists, if not, create one
	place, err := server.store.GetPlaceByGoogleId(c, request.GoogleID)
	if err != nil {
		if err == sql.ErrNoRows { // place not found
			c.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("place not found")))
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}
	// check if favorite exists, if not, create one
	fav, err := server.store.GetFavorite(c, db.GetFavoriteParams{
		UserID:  user.UserID,
		PlaceID: place.PlaceID,
	})
	if err != nil {
		if err == sql.ErrNoRows { // favorite not found
			c.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("favorite not found")))
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}
	if err := server.store.RemoveFavorite(c, db.RemoveFavoriteParams{
		UserID:  user.UserID,
		PlaceID: place.PlaceID,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"User": NewUserResponse(user), "Favorite": fav, "Place": place})
}
