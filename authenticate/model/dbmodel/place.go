package dbmodel

import (
	"authenticate/model/entitymodel"
	"time"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

type Place struct {
	ID                       string         `gorm:"primaryKey" json:"id"`
	GoogleID                 string         `gorm:"not null" json:"google_id"`
	TWDisplayName            string         `gorm:"not null" json:"tw_display_name"`
	JPDisplayName            string         `json:"jp_display_name"`
	TWFormattedAddress       string         `json:"tw_formatted_address"`
	TWWeekdayDescriptions    []string       `json:"tw_weekday_descriptions"`
	AdministrativeAreaLevel1 string         `json:"administrative_area_level_1"`
	Country                  string         `json:"country"`
	GoogleMapURI             string         `json:"google_map_uri"`
	InternationalPhoneNumber string         `json:"international_phone_number"`
	Lat                      string         `json:"lat"`
	Lng                      string         `json:"lng"`
	PrimaryType              string         `json:"primary_type"`
	Rating                   string         `json:"rating"`
	Types                    []string       `json:"types"`
	UserRatingCount          int32          `json:"user_rating_count"`
	WebsiteURI               string         `json:"website_uri"`
	PlaceVersion             int32          `json:"place_version"`
	CreatedAt                time.Time      `json:"created_at"`
	UpdatedAt                time.Time      `json:"updated_at"`
	DeletedAt                gorm.DeletedAt `json:"deleted_at"`
}

func (db Place) ToEntityModel() entitymodel.Place {
	return entitymodel.Place{
		ID:                       db.ID,
		GoogleID:                 db.GoogleID,
		TWDisplayName:            db.TWDisplayName,
		JPDisplayName:            db.JPDisplayName,
		TWFormattedAddress:       db.TWFormattedAddress,
		TWWeekdayDescriptions:    db.TWWeekdayDescriptions,
		AdministrativeAreaLevel1: db.AdministrativeAreaLevel1,
		Country:                  db.Country,
		GoogleMapURI:             db.GoogleMapURI,
		InternationalPhoneNumber: db.InternationalPhoneNumber,
		Lat:                      db.Lat,
		Lng:                      db.Lng,
		PrimaryType:              db.PrimaryType,
		Rating:                   db.Rating,
		Types:                    db.Types,
		UserRatingCount:          db.UserRatingCount,
		WebsiteURI:               db.WebsiteURI,
		PlaceVersion:             db.PlaceVersion,
		CreatedAt:                db.CreatedAt,
		UpdatedAt:                db.UpdatedAt,
	}
}

type PlacePreloadFavoritesUsers struct {
	Place
	FavoritePreloadUsers FavoritePreloadUserSlice `gorm:"foreignKey:PlaceID;references:ID"`
}

func (db PlacePreloadFavoritesUsers) TableName() string {
	return "place"
}

func (db PlacePreloadFavoritesUsers) ToEntityModel() entitymodel.PlaceFavoritesUsersInfo {
	return entitymodel.PlaceFavoritesUsersInfo{
		Place:         db.Place.ToEntityModel(),
		FavoriteUsers: db.FavoritePreloadUsers.ToEntityModel(),
	}
}

type FavoritePreloadUser struct {
	Favorite
	User User `gorm:"foreignKey:ID;references:UserID"`
}

func (db FavoritePreloadUser) TableName() string {
	return "favorite"
}

func (db FavoritePreloadUser) ToEntityModel() entitymodel.FavoriteUser {
	return entitymodel.FavoriteUser{
		Favorite: db.Favorite.ToEntityModel(),
		User:     db.User.ToEntityModel(),
	}
}

type FavoritePreloadUserSlice []FavoritePreloadUser

func (db FavoritePreloadUserSlice) ToEntityModel() entitymodel.FavoriteUserSlice {
	return lo.Map(db, func(item FavoritePreloadUser, _ int) entitymodel.FavoriteUser {
		return item.ToEntityModel()
	})
}
