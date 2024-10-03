package dbmodel

import (
	"authenticate/model/entitymodel"
	"time"
)

type Favorite struct {
	IsFavorite bool      `json:"is_favorite"`
	UserID     string    `json:"user_id"`
	PlaceID    string    `json:"place_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (db Favorite) ToEntityModel() entitymodel.Favorite {
	return entitymodel.Favorite{
		IsFavorite: db.IsFavorite,
		UserID:     db.UserID,
		PlaceID:    db.PlaceID,
		CreatedAt:  db.CreatedAt,
		UpdatedAt:  db.UpdatedAt,
	}
}

type FavoriteSlice []Favorite

func (db FavoriteSlice) ToEntityModel() entitymodel.FavoriteSlice {
	entityModels := make(entitymodel.FavoriteSlice, len(db))
	for i, v := range db {
		entityModels[i] = v.ToEntityModel()
	}
	return entityModels
}
