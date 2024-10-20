package HTTPResponse

import (
	"authenticate/model/entitymodel"
	"authenticate/model/responsemodel"
)

type LoginResponseRef entitymodel.Session

func (ref LoginResponseRef) ToResponse() responsemodel.LoginResponse {
	return responsemodel.LoginResponse{
		UserID:                ref.UserID,
		AccessToken:           ref.AccessToken,
		RefreshToken:          ref.RefreshToken,
		AccessTokenExpiresAt:  ref.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: ref.RefreshTokenExpiresAt,
	}
}

type UserRef entitymodel.User

func (ref UserRef) ToResponse() responsemodel.User {
	return responsemodel.User{
		Nickname: ref.Nickname,
		Email:    ref.Email,
		Active:   ref.Active,
	}
}

type PlaceRef entitymodel.Place

func (ref PlaceRef) ToResponse() responsemodel.Place {
	return responsemodel.Place{
		GoogleID:                 ref.GoogleID,
		TWDisplayName:            ref.TWDisplayName,
		JPDisplayName:            ref.JPDisplayName,
		TWFormattedAddress:       ref.TWFormattedAddress,
		TWWeekdayDescriptions:    ref.TWWeekdayDescriptions,
		AdministrativeAreaLevel1: ref.AdministrativeAreaLevel1,
		Country:                  ref.Country,
		GoogleMapURI:             ref.GoogleMapURI,
		InternationalPhoneNumber: ref.InternationalPhoneNumber,
		Lat:                      ref.Lat,
		Lng:                      ref.Lng,
		PrimaryType:              ref.PrimaryType,
		Rating:                   ref.Rating,
		Types:                    ref.Types,
		UserRatingCount:          ref.UserRatingCount,
		WebsiteURI:               ref.WebsiteURI,
		PlaceVersion:             ref.PlaceVersion,
		UpdatedAt:                ref.UpdatedAt,
	}
}

type PlaceSliceRef entitymodel.PlaceSlice

func (ref PlaceSliceRef) ToResponse() responsemodel.PlaceSlice {
	response := make(responsemodel.PlaceSlice, len(ref))
	for i, v := range ref {
		response[i] = PlaceRef(v).ToResponse()
	}
	return response
}

type UserFavoritePlacesRef entitymodel.UserFavoritePlaces

func (ref UserFavoritePlacesRef) ToResponse() responsemodel.GetUserFavoritesResponse {
	return responsemodel.GetUserFavoritesResponse{
		User:   UserRef(ref.User).ToResponse(),
		Places: PlaceSliceRef(ref.Places).ToResponse(),
	}
}
