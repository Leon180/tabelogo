package repository

import (
	"authenticate/model/entitymodel"
	"context"
)

type SessionHandler interface {
	CreateSession(ctx context.Context, session entitymodel.Session) error
	GetSessionByID(ctx context.Context, sessionID string) (entitymodel.Session, error)
	GetSessionByUserID(ctx context.Context, userID string) (entitymodel.Session, error)
	GetSessionWithUserByRefreshToken(ctx context.Context, refreshToken string) (entitymodel.SessionPreloadUser, error)
	UpdateSession(ctx context.Context, sessionID string, updates map[string]interface{}) error
	DeleteSession(ctx context.Context, sessionID string) error
}
