package main

import (
	db "authenticate/cmd/data/sqlc"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"

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

func NewUserResponse(user db.User) UserResponse {
	return UserResponse{
		UserID:    user.UserID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

func generateRandomNumber(max int64) (int64, error) {
	randomNumber, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		return 0, err
	}
	return (randomNumber.Int64()) + 1, nil
}
