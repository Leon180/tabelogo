package token

import (
	"authenticate/model/entitymodel"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTMaker is a JSON Web Token maker
type JWTMaker struct {
	symmetricKey []byte
}

// NewJWTMaker creates a new JWTMaker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < 32 {
		return nil, fmt.Errorf("invalid key size: must be at least 32 characters")
	}

	return &JWTMaker{symmetricKey: []byte(secretKey)}, nil
}

func (maker *JWTMaker) CreateToken(userInfo entitymodel.User, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userInfo, duration)
	if err != nil {
		return "", payload, fmt.Errorf("failed to create payload: %w", err)
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString(maker.symmetricKey)
	if err != nil {
		return "", payload, fmt.Errorf("failed to create token: %w", err)
	}

	return token, payload, nil
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return maker.symmetricKey, nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
