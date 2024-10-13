package entitymodel

import (
	"strings"
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

func (entity Place) IsExist() bool {
	return strings.TrimSpace(entity.ID) != "" && strings.TrimSpace(entity.GoogleID) != ""
}

func (entity Place) GetUpdates() map[string]interface{} {
	return map[string]interface{}{
		"google_id":                   entity.GoogleID,
		"tw_display_name":             entity.TWDisplayName,
		"jp_display_name":             entity.JPDisplayName,
		"tw_formatted_address":        entity.TWFormattedAddress,
		"tw_weekday_descriptions":     entity.TWWeekdayDescriptions,
		"administrative_area_level_1": entity.AdministrativeAreaLevel1,
		"country":                     entity.Country,
		"google_map_uri":              entity.GoogleMapURI,
		"international_phone_number":  entity.InternationalPhoneNumber,
		"lat":                         entity.Lat,
		"lng":                         entity.Lng,
		"primary_type":                entity.PrimaryType,
		"rating":                      entity.Rating,
		"types":                       entity.Types,
		"user_rating_count":           entity.UserRatingCount,
		"website_uri":                 entity.WebsiteURI,
		"place_version":               entity.PlaceVersion,
		"updated_at":                  entity.UpdatedAt,
	}
}

type PlaceSlice []Place
