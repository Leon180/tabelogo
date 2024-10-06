package entitymodel

import (
	"time"
)

type User struct {
	ID             string
	Nickname       string
	Email          string
	HashedPassword string
	Active         bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type UserSlice []User

type Session struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	UserID       string    `gorm:"not null" json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}
