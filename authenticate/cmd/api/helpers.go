package main

import (
	db "authenticate/cmd/data/sqlc"
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func HashedPassword(password string) (string, error) {
	hpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hpw), nil
}

func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func NewUserResponse(user db.User) gin.H {
	return gin.H{
		"user_id":    user.UserID,
		"email":      user.Email,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}
}
