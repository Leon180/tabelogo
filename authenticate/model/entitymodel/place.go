package entitymodel

import (
	"time"
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
}

type PlaceSlice []Place
