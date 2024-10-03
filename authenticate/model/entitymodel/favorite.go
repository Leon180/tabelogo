package entitymodel

import "time"

type Favorite struct {
	IsFavorite bool      `json:"is_favorite"`
	UserID     string    `json:"user_id"`
	PlaceID    string    `json:"place_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type FavoriteSlice []Favorite

type FavoriteUser struct {
	Favorite
	User User `json:"user"`
}

type FavoriteUserSlice []FavoriteUser
