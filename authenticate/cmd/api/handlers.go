package main

import (
	db "authenticate/cmd/data/sqlc"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// if error, return errorResponse:
//
//	{
//		"error": "error message"
//	}
//
// if no error, return user
//
//	{
//		"user": {
//			"user_id": 1,
//			"email": "test@gmail",
//			"created_at": "2021-08-01T00:00:00Z",
//			"updated_at": "2021-08-01T00:00:00Z"
//		}
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
	c.JSON(http.StatusOK, gin.H{"user": NewUserResponse(user)})
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (server *Server) Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUser(c, request.Email)
	if err != nil {
		if err == sql.ErrNoRows {
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

	// TODO: generate paseto token

	c.JSON(http.StatusOK, gin.H{"user": NewUserResponse(user)})
}
