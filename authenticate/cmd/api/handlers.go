package main

import (
	db "authenticate/cmd/data/sqlc"
	"authenticate/token"
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
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

func (server *Server) Regist(ctx *gin.Context) {
	var request CreateUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	hashedPassword, err := HashedPassword(request.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Email:          request.Email,
		HashedPassword: hashedPassword,
	}
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation": // duplicate email
				ctx.JSON(http.StatusConflict, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if server.rabbitMQ != nil {
		go server.logEventViaRabbit("authenticated", fmt.Sprintf("user %s regist", user.Email), "log.INFO")
	}

	ctx.JSON(http.StatusOK, gin.H{"user": NewUserResponse(user)})
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

func (server *Server) Login(c *gin.Context) {
	var request LoginRequest
	var user db.User
	var err error
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get user in database
	user, err = server.store.GetUser(c, request.Email)
	if err != nil {
		if err == sql.ErrNoRows { // user not found
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// check password
	if err := ComparePassword(user.HashedPassword, request.Password); err != nil {
		// if password not match, log event and return
		go server.logEventViaRabbit("unauthenticated", fmt.Sprintf("user %s log in failed: mistype password, from IP: %s, UserAgent: %s", user.Email, c.ClientIP(), c.Request.UserAgent()), "log.ERROR")
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// generate access token
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Email, server.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// generate refresh token
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Email, server.config.RefreshTokenDuration)
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
	_ = server.setSessionInCacheWithExpiry(c, session, server.config.AccessTokenDuration)

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

func (server *Server) RenewAccessToken(c *gin.Context) {
	var request RenewAccessTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// verify refresh token
	refreshPayload, err := server.tokenMaker.VerifyToken(request.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// get user in redis or database
	session, st, err := server.GetSessionInCacheOrDatabase(c, refreshPayload)
	if st == "error" {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// check session status
	if session.IsBlocked {
		err = server.deleteSessionInCache(c, session.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		c.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("session is blocked")))
		return
	}
	if session.Email != refreshPayload.Email {
		err = server.deleteSessionInCache(c, session.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		c.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("incorrect user")))
		return
	}
	if session.RefreshToken != request.RefreshToken {
		err = server.deleteSessionInCache(c, session.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		c.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("incorrect refresh token")))
		return
	}
	if time.Now().After(session.ExpiresAt) {
		err = server.deleteSessionInCache(c, session.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		c.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("session has expired")))
		return
	}

	// generate new access token
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(session.Email, server.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// set session in redis
	err = server.setSessionInCacheWithExpiry(c, session, server.config.AccessTokenDuration)
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

// 1. check if user exists, if not, return
// 2-1. check place in redis or db, if not, create place in db and set place in redis
// 2-2. if place exists in redis or db, use the place to toggle favorite
// 3. toggle favorite
// 4. return response
func (server *Server) ToggleFavorite(c *gin.Context) {
	var request SaveFavoriteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// check if user exists, if not, return
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	session, st, err := server.GetSessionInCacheOrDatabase(c, authPayload)
	if st == "error" {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	place, _, err := server.GetPlaceInCacheOrDatabaseAndCreateIfNotExist(c, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
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

	// return response
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
	session, st, err := server.GetSessionInCacheOrDatabase(c, authPayload)
	if st == "error" {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// get favorites
	favs, err := server.store.ListFavoritesByCreateTime(c, db.ListFavoritesByCreateTimeParams{
		UserEmail: session.Email,
		Limit:     request.Limit,
		Offset:    request.Offset,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// return response
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
	session, st, err := server.GetSessionInCacheOrDatabase(c, authPayload)
	if st == "error" {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
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
	session, st, err := server.GetSessionInCacheOrDatabase(c, authPayload)
	if st == "error" {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
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
	session, st, err := server.GetSessionInCacheOrDatabase(c, authPayload)
	if st == "error" {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
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
	session, st, err := server.GetSessionInCacheOrDatabase(c, authPayload)
	if st == "error" {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
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
	session, st, err := server.GetSessionInCacheOrDatabase(c, authPayload)
	if st == "error" {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
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
	session, st, err := server.GetSessionInCacheOrDatabase(c, authPayload)
	if st == "error" {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	place, st, err := server.GetPlaceInCacheOrDatabaseAndCreateIfNotExist(c, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// if not created, and redis expiry, update place in database and set place in redis
	if st == "database_exist_cache_not_exist" {
		// update place
		updateArg := db.UpdatePlaceParams{
			PlaceVersion: place.PlaceVersion,
		}
		if len(request.GoogleID) != 0 {
			updateArg.GoogleID = sql.NullString{String: request.GoogleID, Valid: true}
		}
		if len(request.TwDisplayName) != 0 {
			updateArg.TwDisplayName = sql.NullString{String: request.TwDisplayName, Valid: true}
		}
		if len(request.TwFormattedAddress) != 0 {
			updateArg.TwFormattedAddress = sql.NullString{String: request.TwFormattedAddress, Valid: true}
		}
		if len(request.TwWeekdayDescriptions) != 0 {
			updateArg.TwWeekdayDescriptions = pq.StringArray(request.TwWeekdayDescriptions)
		}
		if len(request.GoogleMapUri) != 0 {
			updateArg.GoogleMapUri = sql.NullString{String: request.GoogleMapUri, Valid: true}
		}
		if len(request.Lat) != 0 {
			updateArg.Lat = sql.NullString{String: request.Lat, Valid: true}
		}
		if len(request.Lng) != 0 {
			updateArg.Lng = sql.NullString{String: request.Lng, Valid: true}
		}
		if len(request.Types) != 0 {
			updateArg.Types = pq.StringArray(request.Types)
		}
		if len(request.AdministrativeAreaLevel1) != 0 {
			updateArg.AdministrativeAreaLevel1 = sql.NullString{String: request.AdministrativeAreaLevel1, Valid: true}
		}
		if len(request.Country) != 0 {
			updateArg.Country = sql.NullString{String: request.Country, Valid: true}
		}
		if len(request.InternationalPhoneNumber) != 0 {
			updateArg.InternationalPhoneNumber = sql.NullString{String: request.InternationalPhoneNumber, Valid: true}
		}
		if len(request.PrimaryType) != 0 {
			updateArg.PrimaryType = sql.NullString{String: request.PrimaryType, Valid: true}
		}
		if len(request.Rating) != 0 {
			updateArg.Rating = sql.NullString{String: request.Rating, Valid: true}
		}
		if request.UserRatingCount != 0 {
			updateArg.UserRatingCount = sql.NullInt32{Int32: request.UserRatingCount, Valid: true}
		}
		if len(request.WebsiteUri) != 0 {
			updateArg.WebsiteUri = sql.NullString{String: request.WebsiteUri, Valid: true}
		}

		place, err = server.store.UpdatePlace(c, updateArg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		// set place in redis with expiry
		err = server.setPlaceInCacheWithExpiry(c, place, 10*time.Minute+time.Duration(rand.Intn(5))*time.Minute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	// check if favorite exists,
	fav, err := server.store.GetFavorite(c, db.GetFavoriteParams{
		UserEmail: session.Email,
		GoogleID:  place.GoogleID,
	})
	if err != nil {
		if err == sql.ErrNoRows { // favorite not found
			c.JSON(http.StatusOK, gin.H{"isFavorite": false, "User": session.Email, "Place": place})
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"isFavorite": true, "User": session.Email, "Favorite": fav, "Place": place})
}

type FindPlaceInCacheRequest struct {
	PlaceID string `json:"place_id" binding:"required"`
}

type FindPlaceInCacheResponse struct {
	Place db.Place `json:"place"`
	Found bool     `json:"found"`
}

func (s *Server) FindPlaceInCache(c *gin.Context) {
	var request FindPlaceInCacheRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	place, err := s.getPlaceInCache(c, request.PlaceID)
	if err != nil && err != redis.Nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if reflect.DeepEqual(place, db.Place{}) {
		c.JSON(http.StatusOK, FindPlaceInCacheResponse{
			Place: place,
			Found: false,
		})
		return
	}
	c.JSON(http.StatusOK, FindPlaceInCacheResponse{
		Place: place,
		Found: true,
	})
}

type SetJPDisplayNameInCacheAndDataBaseRequest struct {
	PlaceID       string `json:"place_id" binding:"required"`
	JPDisplayName string `json:"jp_display_name" binding:"required"`
}

func (s *Server) SetJPDisplayNameInCacheAndDataBase(c *gin.Context) {
	var request SetJPDisplayNameInCacheAndDataBaseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	place, err := s.store.GetPlaceByGoogleId(c, request.PlaceID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	place, err = s.store.UpdatePlace(c, db.UpdatePlaceParams{
		GoogleID:      sql.NullString{String: request.PlaceID, Valid: true},
		JpDisplayName: sql.NullString{String: request.JPDisplayName, Valid: true},
		PlaceVersion:  place.PlaceVersion,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = s.setPlaceInCacheWithExpiry(c, place, 10*time.Minute+time.Duration(rand.Intn(5))*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"Place": place})
}
