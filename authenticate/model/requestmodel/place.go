package requestmodel

import (
	"authenticate/model/entitymodel"
	"authenticate/utility"
	"time"
)

type SavePlaceRequest struct {
	GoogleID                 string   `json:"google_id"`
	TWDisplayName            string   `json:"tw_display_name"`
	JPDisplayName            string   `json:"jp_display_name"`
	TWFormattedAddress       string   `json:"tw_formatted_address"`
	TWWeekdayDescriptions    []string `json:"tw_weekday_descriptions"`
	AdministrativeAreaLevel1 string   `json:"administrative_area_level_1"`
	Country                  string   `json:"country"`
	GoogleMapURI             string   `json:"google_map_uri"`
	InternationalPhoneNumber string   `json:"international_phone_number"`
	Lat                      string   `json:"lat"`
	Lng                      string   `json:"lng"`
	PrimaryType              string   `json:"primary_type"`
	Rating                   string   `json:"rating"`
	Types                    []string `json:"types"`
	UserRatingCount          int32    `json:"user_rating_count"`
	WebsiteURI               string   `json:"website_uri"`
	PlaceVersion             int32    `json:"place_version"`
}

func (req SavePlaceRequest) ToEntity() entitymodel.Place {
	return entitymodel.Place{
		ID:                       utility.GenDefaultUUID(),
		GoogleID:                 req.GoogleID,
		TWDisplayName:            req.TWDisplayName,
		JPDisplayName:            req.JPDisplayName,
		TWFormattedAddress:       req.TWFormattedAddress,
		TWWeekdayDescriptions:    req.TWWeekdayDescriptions,
		AdministrativeAreaLevel1: req.AdministrativeAreaLevel1,
		Country:                  req.Country,
		GoogleMapURI:             req.GoogleMapURI,
		InternationalPhoneNumber: req.InternationalPhoneNumber,
		Lat:                      req.Lat,
		Lng:                      req.Lng,
		PrimaryType:              req.PrimaryType,
		Rating:                   req.Rating,
		Types:                    req.Types,
		UserRatingCount:          req.UserRatingCount,
		WebsiteURI:               req.WebsiteURI,
		PlaceVersion:             req.PlaceVersion,
		CreatedAt:                time.Time{},
		UpdatedAt:                time.Time{},
	}
}

type GetPlaceRequest struct {
	GoogleID string `json:"google_id"`
}
