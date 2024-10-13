package postgresqlRepository

import (
	"authenticate/model/dbmodel"
	"authenticate/model/entitymodel"
	"authenticate/repository"
	"authenticate/utility"
	"context"

	"gorm.io/gorm"
)

func NewSessionHandle(db *gorm.DB) repository.SessionHandler {
	return &SessionHandle{db: db}
}

type SessionHandle struct {
	db *gorm.DB
}

func (handle *SessionHandle) CreateSession(ctx context.Context, session entitymodel.Session) error {
	dbModel := dbmodel.Session{
		ID:                    session.ID,
		UserID:                session.UserID,
		AccessToken:           session.AccessToken,
		RefreshToken:          session.RefreshToken,
		UserAgent:             session.UserAgent,
		ClientIP:              session.ClientIP,
		IsBlocked:             session.IsBlocked,
		AccessTokenExpiresAt:  session.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: session.RefreshTokenExpiresAt,
		CreatedAt:             session.CreatedAt,
	}

	if err := handle.db.WithContext(ctx).
		Create(&dbModel).Error; err != nil {
		utility.SugarLogger.Errorln("error creating session:", err)
		return err
	}

	return nil
}

func (handle *SessionHandle) GetSessionByID(ctx context.Context, sessionID string) (entitymodel.Session, error) {
	var session dbmodel.Session

	if err := handle.db.WithContext(ctx).
		Where("id = ?", sessionID).
		Find(&session).Error; err != nil {
		utility.SugarLogger.Errorln("error getting session by ID:", err)
		return entitymodel.Session{}, err
	}

	return session.ToEntityModel(), nil
}

func (handle *SessionHandle) GetSessionByUserID(ctx context.Context, userID string) (entitymodel.Session, error) {
	var session dbmodel.Session

	if err := handle.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&session).Error; err != nil {
		utility.SugarLogger.Errorln("error getting session by user ID:", err)
		return entitymodel.Session{}, err
	}

	return session.ToEntityModel(), nil
}

func (handle *SessionHandle) GetSessionWithUserByRefreshToken(ctx context.Context, refreshToken string) (entitymodel.SessionPreloadUser, error) {
	var dbModel dbmodel.SessionPreloadUser

	if err := handle.db.WithContext(ctx).
		Preload("User").
		Where("refresh_token = ?", refreshToken).
		Find(&dbModel).Error; err != nil {
		utility.SugarLogger.Errorln("error getting session with user by refresh token:", err)
		return entitymodel.SessionPreloadUser{}, err
	}

	return dbModel.ToEntityModel(), nil
}

func (handle *SessionHandle) UpdateSession(ctx context.Context, sessionID string, updates map[string]interface{}) error {
	if err := handle.db.WithContext(ctx).
		Model(&dbmodel.Session{}).
		Where("id = ?", sessionID).
		Updates(updates).Error; err != nil {
		utility.SugarLogger.Errorln("error updating session:", err)
		return err
	}

	return nil
}

func (handle *SessionHandle) DeleteSession(ctx context.Context, sessionID string) error {
	if err := handle.db.WithContext(ctx).
		Where("id = ?", sessionID).
		Delete(&dbmodel.Session{}).Error; err != nil {
		utility.SugarLogger.Errorln("error deleting session:", err)
		return err
	}

	return nil
}
