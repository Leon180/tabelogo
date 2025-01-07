package responsemodel

import (
	"authenticate/grpc/proto"
	"authenticate/model/entitymodel"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PlaceEntityModel entitymodel.Place

func (entity PlaceEntityModel) ToPlaceProto() *proto.Place {
	return &proto.Place{
		GoogleId:                  entity.GoogleID,
		TwDisplayName:             entity.TWDisplayName,
		JpDisplayName:             entity.JPDisplayName,
		TwFormattedAddress:        entity.TWFormattedAddress,
		TwWeekdayDescriptions:     entity.TWWeekdayDescriptions,
		AdministrativeAreaLevel_1: entity.AdministrativeAreaLevel1,
		Country:                   entity.Country,
		GoogleMapUri:              entity.GoogleMapURI,
		InternationalPhoneNumber:  entity.InternationalPhoneNumber,
		Lat:                       entity.Lat,
		Lng:                       entity.Lng,
		PrimaryType:               entity.PrimaryType,
		Rating:                    entity.Rating,
		Types:                     entity.Types,
		UserRatingCount:           entity.UserRatingCount,
		WebsiteUri:                entity.WebsiteURI,
		PlaceVersion:              entity.PlaceVersion,
		UpdatedAt:                 timestamppb.New(entity.UpdatedAt),
	}
}

type PlaceSliceEntityModel entitymodel.PlaceSlice

func (entity PlaceSliceEntityModel) ToPlaceProtoSlice() []*proto.Place {
	return lo.Map(entity, func(place entitymodel.Place, _ int) *proto.Place {
		return PlaceEntityModel(place).ToPlaceProto()
	})
}
