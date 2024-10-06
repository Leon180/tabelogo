package dbmodel

import (
	"authenticate/model/entitymodel"
	"time"
)

type Favorite struct {
	ID         string    `gorm:"primaryKey" json:"id"`
	IsFavorite bool      `json:"is_favorite"`
	UserID     string    `json:"user_id"`
	PlaceID    string    `json:"place_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (db Favorite) ToEntityModel() entitymodel.Favorite {
	return entitymodel.Favorite{
		ID:         db.ID,
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

type PreloadFavoritePlaceUser struct {
	Favorite
	Place Place `gorm:"foreignKey:ID;references:PlaceID"`
	User  User  `gorm:"foreignKey:ID;references:UserID"`
}

func (db PreloadFavoritePlaceUser) TableName() string {
	return "favorite"
}

type PreloadFavoritePlaceUserSlice []PreloadFavoritePlaceUser

func (db PreloadFavoritePlaceUserSlice) ToPlaceFavoriteUsersSliceEntityModel() entitymodel.PlaceFavoriteUsersSlice {
	m := map[string]int{} // record the index of the place in the slice
	entityModels := make(entitymodel.PlaceFavoriteUsersSlice, 0, len(db))
	for i, v := range db {
		if !v.IsFavorite {
			continue
		}
		if _, ok := m[v.PlaceID]; !ok {
			m[v.PlaceID] = i
			entityModels[i] = entitymodel.PlaceFavoriteUsers{
				Place: v.Place.ToEntityModel(),
				Users: []entitymodel.User{},
			}
		}
		entityModels[m[v.PlaceID]].Users = append(entityModels[m[v.PlaceID]].Users, v.User.ToEntityModel())
	}
	return entityModels
}

func (db PreloadFavoritePlaceUserSlice) ToUserFavoritePlacesSliceEntityModel() entitymodel.UserFavoritePlacesSlice {
	m := map[string]int{} // record the index of the user in the slice
	entityModels := make(entitymodel.UserFavoritePlacesSlice, 0, len(db))
	for i, v := range db {
		if !v.IsFavorite {
			continue
		}
		if _, ok := m[v.UserID]; !ok {
			m[v.UserID] = i
			entityModels[i] = entitymodel.UserFavoritePlaces{
				User:   v.User.ToEntityModel(),
				Places: []entitymodel.Place{},
			}
		}
		entityModels[m[v.UserID]].Places = append(entityModels[m[v.UserID]].Places, v.Place.ToEntityModel())
	}
	return entityModels
}
