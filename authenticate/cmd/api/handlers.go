package main

import (
	db "authenticate/cmd/data/sqlc"
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

	resp := LoginResponse{
		SessionID:             uuid.UUID(session.SessionID),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiresAt,
		User:                  NewUserResponse(user),
	}

	c.JSON(http.StatusOK, resp)
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

	c.JSON(http.StatusOK, resp)
}

func (server *Server) Message(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello World!"})
}
