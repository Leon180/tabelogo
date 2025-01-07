package requestmodel

import (
	"authenticate/grpc/proto"
	"authenticate/model/requestmodel"
)

func ConvertRegistUserRequest(req *proto.RegistUserRequest) requestmodel.RegistUserRequest {
	return requestmodel.RegistUserRequest{
		Nickname: req.Nickname,
		Email:    req.Email,
		Password: req.Password,
	}
}

func ConvertLoginUserRequest(req *proto.LoginUserRequest) requestmodel.LoginUserRequest {
	return requestmodel.LoginUserRequest{
		Email:    req.Email,
		Password: req.Password,
	}
}

func ConvertRenewAccessTokenRequest(req *proto.RenewAccessTokenRequest) requestmodel.RenewAccessTokenRequest {
	return requestmodel.RenewAccessTokenRequest{
		RefreshToken: req.RefreshToken,
	}
}

func ConvertSaveFavoriteRequest(req *proto.SaveFavoriteRequest) requestmodel.SaveFavoriteRequest {
	return requestmodel.SaveFavoriteRequest{
		IsFavorite:    req.IsFavorite,
		PlaceGoogleID: req.PlaceGoogleId,
	}
}

func ConvertGetUserFavoritesRequest(req *proto.GetUserFavoritesRequest) requestmodel.GetUserFavoritesRequest {
	return requestmodel.GetUserFavoritesRequest{
		Country: func() *string {
			if req.Country == "" {
				return nil
			}
			return &req.Country
		}(),
		AdministrativeAreaLevel1: func() *string {
			if req.AdministrativeAreaLevel_1 == "" {
				return nil
			}
			return &req.AdministrativeAreaLevel_1
		}(),
		OrderBy: func() *int {
			if req.OrderBy == 0 {
				return nil
			}
			order := int(req.OrderBy)
			return &order
		}(),
	}
}
