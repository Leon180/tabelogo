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
	ID           string         `gorm:"primaryKey" json:"id"`
	UserID       string         `gorm:"not null" json:"user_id"`
	RefreshToken string         `json:"refresh_token"`
	UserAgent    string         `json:"user_agent"`
	ClientIp     string         `json:"client_ip"`
	IsBlocked    bool           `json:"is_blocked"`
	ExpiresAt    time.Time      `json:"expires_at"`
	CreatedAt    time.Time      `json:"created_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at"`
}

func (db Session) ToEntityModel() entitymodel.Session {
	return entitymodel.Session{
		ID:           db.ID,
		UserID:       db.UserID,
		RefreshToken: db.RefreshToken,
		UserAgent:    db.UserAgent,
		ClientIp:     db.ClientIp,
		IsBlocked:    db.IsBlocked,
		ExpiresAt:    db.ExpiresAt,
		CreatedAt:    db.CreatedAt,
	}
}
