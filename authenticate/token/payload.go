package token

import (
	"authenticate/model/entitymodel"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	Issuer = "authenticate-service"
)

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

type Payload struct {
	ID        uuid.UUID        `json:"uuid"`
	UserInfo  entitymodel.User `json:"user_info"`
	IssuedAt  time.Time        `json:"issued_at"`
	ExpiresAt time.Time        `json:"expires_at"`
	Issuer    string           `json:"issuer"`
	Subject   string           `json:"subject"`
	NotBefore time.Time        `json:"not_before"`
	Audience  string           `json:"audience"`
}

func NewPayload(userInfo entitymodel.User, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:        tokenID,
		UserInfo:  userInfo,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
		Issuer:    Issuer,
		Subject:   "user token",
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return ErrExpiredToken
	}
	if !payload.UserInfo.IsExist() || (payload.Issuer != Issuer) {
		return ErrInvalidToken
	}
	return nil
}

func (payload *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(payload.ExpiresAt), nil
}

func (payload *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(payload.IssuedAt), nil
}

func (payload *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(payload.NotBefore), nil
}

func (payload *Payload) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{payload.Audience}, nil
}

func (payload *Payload) GetIssuer() (string, error) {
	return payload.Issuer, nil
}

func (payload *Payload) GetSubject() (string, error) {
	return payload.Subject, nil
}
