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

	Password string
}

func (u User) IsActive() bool {
	return u.Active
}

func (u User) IsExist() bool {
	return u.ID != ""
}

type UserSlice []User

type Session struct {
	ID                    string    `gorm:"primaryKey" json:"id"`
	UserID                string    `gorm:"not null" json:"user_id"`
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	UserAgent             string    `json:"user_agent"`
	ClientIP              string    `json:"client_ip"`
	IsBlocked             bool      `json:"is_blocked"`
	AccessTokenExpiresAt  time.Time `json:"expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	CreatedAt             time.Time `json:"created_at"`
}

func (s Session) IsExist() bool {
	return s.ID != ""
}

func (s Session) IsAccessTokenExpired() bool {
	return s.AccessTokenExpiresAt.UTC().Before(time.Now().UTC())
}

func (s Session) IsRefreshTokenExpired() bool {
	return s.RefreshTokenExpiresAt.UTC().Before(time.Now().UTC())
}

func (s Session) Blocked() bool {
	return s.IsBlocked
}

func (s Session) GetUpdates() map[string]interface{} {
	return map[string]interface{}{
		"user_agent":              s.UserAgent,
		"client_ip":               s.ClientIP,
		"is_blocked":              s.IsBlocked,
		"access_token_expires_at": s.AccessTokenExpiresAt,
	}
}

type Request struct {
	UserAgent string
	ClientIP  string
}

type SessionPreloadUser struct {
	Session
	User User
}

func (entity SessionPreloadUser) IsUserExist() bool {
	return entity.User.IsExist()
}

func (entity SessionPreloadUser) Blocked() bool {
	return entity.Session.Blocked()
}

func (entity SessionPreloadUser) GetUpdates() map[string]interface{} {
	return entity.Session.GetUpdates()
}

func (entity SessionPreloadUser) IsAccessTokenExpired() bool {
	return entity.Session.IsAccessTokenExpired()
}

func (entity SessionPreloadUser) IsExist() bool {
	return entity.Session.IsExist()
}

func (entity SessionPreloadUser) IsRefreshTokenExpired() bool {
	return entity.Session.IsRefreshTokenExpired()
}
