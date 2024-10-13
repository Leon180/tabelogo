package dbmodel

import (
	"authenticate/model/entitymodel"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID             string         `gorm:"primaryKey" json:"user_id"`
	Nickname       string         `gorm:"not null" json:"nickname"`
	Email          string         `gorm:"not null" json:"email"`
	HashedPassword string         `gorm:"not null" json:"hashed_password"`
	Active         bool           `gorm:"timestamp with time zone" json:"active"`
	CreatedAt      time.Time      `gorm:"timestamp with time zone" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"timestamp with time zone" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"timestamp with time zone" json:"deleted_at"`
}

func (db User) ToEntityModel() entitymodel.User {
	return entitymodel.User{
		ID:             db.ID,
		Nickname:       db.Nickname,
		Email:          db.Email,
		HashedPassword: db.HashedPassword,
		Active:         db.Active,
		CreatedAt:      db.CreatedAt,
		UpdatedAt:      db.UpdatedAt,
	}
}

type UserSlice []User

func (db UserSlice) ToEntityModel() entitymodel.UserSlice {
	entityModel := make(entitymodel.UserSlice, len(db))
	for i, user := range db {
		entityModel[i] = user.ToEntityModel()
	}
	return entityModel
}

type Session struct {
	ID                    string         `gorm:"primaryKey" json:"id"`
	UserID                string         `gorm:"not null" json:"user_id"`
	AccessToken           string         `json:"access_token"`
	RefreshToken          string         `json:"refresh_token"`
	UserAgent             string         `json:"user_agent"`
	ClientIP              string         `json:"client_ip"`
	IsBlocked             bool           `json:"is_blocked"`
	AccessTokenExpiresAt  time.Time      `json:"expires_at"`
	RefreshTokenExpiresAt time.Time      `json:"refresh_token_expires_at"`
	CreatedAt             time.Time      `json:"created_at"`
	DeletedAt             gorm.DeletedAt `json:"deleted_at"`
}

func (db Session) ToEntityModel() entitymodel.Session {
	return entitymodel.Session{
		ID:                    db.ID,
		UserID:                db.UserID,
		AccessToken:           db.AccessToken,
		RefreshToken:          db.RefreshToken,
		UserAgent:             db.UserAgent,
		ClientIP:              db.ClientIP,
		IsBlocked:             db.IsBlocked,
		AccessTokenExpiresAt:  db.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: db.RefreshTokenExpiresAt,
		CreatedAt:             db.CreatedAt,
	}
}

type SessionPreloadUser struct {
	Session
	User User `gorm:"foreignKey:UserID"`
}

func (db SessionPreloadUser) TableName() string {
	return "session"
}

func (db SessionPreloadUser) ToEntityModel() entitymodel.SessionPreloadUser {
	return entitymodel.SessionPreloadUser{
		Session: db.Session.ToEntityModel(),
		User:    db.User.ToEntityModel(),
	}
}
