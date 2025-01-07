package responsemodel

import (
	"authenticate/grpc/proto"
	"authenticate/model/entitymodel"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserEntityModel entitymodel.User

func (entity UserEntityModel) ToUserProto() *proto.User {
	return &proto.User{
		Nickname: entity.Nickname,
		Email:    entity.Email,
		Active:   entity.Active,
	}
}

type SessionEntityModel entitymodel.Session

func (entity SessionEntityModel) ToLoginUserResponseProto() *proto.LoginUserResponse {
	return &proto.LoginUserResponse{
		UserId:                entity.UserID,
		AccessToken:           entity.AccessToken,
		RefreshToken:          entity.RefreshToken,
		AccessTokenExpiresAt:  timestamppb.New(entity.AccessTokenExpiresAt),
		RefreshTokenExpiresAt: timestamppb.New(entity.RefreshTokenExpiresAt),
	}
}

type UserFavoritePlacesEntityModel entitymodel.UserFavoritePlaces

func (entity UserFavoritePlacesEntityModel) ToGetUserFavoritesResponseProto() *proto.GetUserFavoritesResponse {
	return &proto.GetUserFavoritesResponse{
		User:   UserEntityModel(entity.User).ToUserProto(),
		Places: PlaceSliceEntityModel(entity.Places).ToPlaceProtoSlice(),
	}
}
