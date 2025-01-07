package requestmodel

import (
	"authenticate/grpc/proto"
	"authenticate/model/requestmodel"
)

func ConvertGetPlaceRequest(req *proto.GetPlaceRequest) requestmodel.GetPlaceRequest {
	return requestmodel.GetPlaceRequest{
		GoogleID: req.GoogleId,
	}
}

func ConvertSavePlaceRequest(req *proto.SavePlaceRequest) requestmodel.SavePlaceRequest {
	return requestmodel.SavePlaceRequest{
		GoogleID:                 req.GoogleId,
		TWDisplayName:            req.TwDisplayName,
		JPDisplayName:            req.JpDisplayName,
		TWFormattedAddress:       req.TwFormattedAddress,
		TWWeekdayDescriptions:    req.TwWeekdayDescriptions,
		AdministrativeAreaLevel1: req.AdministrativeAreaLevel_1,
		Country:                  req.Country,
		GoogleMapURI:             req.GoogleMapUri,
		InternationalPhoneNumber: req.InternationalPhoneNumber,
		Lat:                      req.Lat,
		Lng:                      req.Lng,
		PrimaryType:              req.PrimaryType,
		Rating:                   req.Rating,
		Types:                    req.Types,
		UserRatingCount:          req.UserRatingCount,
		WebsiteURI:               req.WebsiteUri,
		PlaceVersion:             req.PlaceVersion,
	}
}
