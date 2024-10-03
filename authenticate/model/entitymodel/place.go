package entitymodel

import (
	"time"

	"gorm.io/gorm"
)

type Place struct {
	ID                       string
	GoogleID                 string
	TWDisplayName            string
	JPDisplayName            string
	TWFormattedAddress       string
	TWWeekdayDescriptions    []string
	AdministrativeAreaLevel1 string
	Country                  string
	GoogleMapURI             string
	InternationalPhoneNumber string
	Lat                      string
	Lng                      string
	PrimaryType              string
	Rating                   string
	Types                    []string
	UserRatingCount          int32
	WebsiteURI               string
	PlaceVersion             int32
	CreatedAt                time.Time
	UpdatedAt                time.Time
	DeletedAt                gorm.DeletedAt
}

type PlaceFavoritesUsersInfo struct {
	Place
	FavoriteUsers FavoriteUserSlice
}
