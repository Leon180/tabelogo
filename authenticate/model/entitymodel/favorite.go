package entitymodel

import (
	"time"
)

type Favorite struct {
	ID            string    `json:"id"`
	IsFavorite    bool      `json:"is_favorite"`
	UserID        string    `json:"user_id"`
	PlaceGoogleID string    `json:"place_google_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type FavoriteSlice []Favorite

type PlaceFavoriteUsers struct {
	Place Place     `json:"place"`
	Users UserSlice `json:"users"`
}

type PlaceFavoriteUsersSlice []PlaceFavoriteUsers

type UserFavoritePlaces struct {
	User   User       `json:"user"`
	Places PlaceSlice `json:"places"`
}

func (entity UserFavoritePlaces) GetUserFavoritePlaceAdministrativeAreaLevel1s() []string {
	administrativeAreaLevel1 := make([]string, 0, len(entity.Places))
	m := make(map[string]interface{})

	for _, v := range entity.Places {
		if _, ok := m[v.AdministrativeAreaLevel1]; !ok {
			m[v.AdministrativeAreaLevel1] = nil
			administrativeAreaLevel1 = append(administrativeAreaLevel1, v.AdministrativeAreaLevel1)
		}
	}

	return administrativeAreaLevel1
}

func (entity UserFavoritePlaces) GetUserFavoritePlaceCountries() []string {
	countries := make([]string, 0, len(entity.Places))
	m := make(map[string]interface{})

	for _, v := range entity.Places {
		if _, ok := m[v.Country]; !ok {
			m[v.Country] = nil
			countries = append(countries, v.Country)
		}
	}

	return countries
}

type UserFavoritePlacesSlice []UserFavoritePlaces
