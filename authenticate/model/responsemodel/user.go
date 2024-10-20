package responsemodel

import (
	"time"
)

type LoginResponse struct {
	UserID                string    `json:"user_id"`
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time `json:"expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

type User struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Active   bool   `json:"active"`
}

type Place struct {
	GoogleID                 string    `json:"google_id"`
	TWDisplayName            string    `json:"tw_display_name"`
	JPDisplayName            string    `json:"jp_display_name"`
	TWFormattedAddress       string    `json:"tw_formatted_address"`
	TWWeekdayDescriptions    []string  `json:"tw_weekday_descriptions"`
	AdministrativeAreaLevel1 string    `json:"administrative_area_level_1"`
	Country                  string    `json:"country"`
	GoogleMapURI             string    `json:"google_map_uri"`
	InternationalPhoneNumber string    `json:"international_phone_number"`
	Lat                      string    `json:"lat"`
	Lng                      string    `json:"lng"`
	PrimaryType              string    `json:"primary_type"`
	Rating                   string    `json:"rating"`
	Types                    []string  `json:"types"`
	UserRatingCount          int32     `json:"user_rating_count"`
	WebsiteURI               string    `json:"website_uri"`
	PlaceVersion             int32     `json:"place_version"`
	UpdatedAt                time.Time `json:"updated_at"`
}

type PlaceSlice []Place

type GetUserFavoritesResponse struct {
	User   User       `json:"user"`
	Places PlaceSlice `json:"places"`
}
