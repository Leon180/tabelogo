package repository

import (
	"authenticate/model/entitymodel"
	"authenticate/model/enum"
	"context"
)

type FavoriteHandler interface {
	CreateFavorite(ctx context.Context, favorite entitymodel.Favorite) error
	GetFavoriteByUserIDAndPlaceID(ctx context.Context, userID, placeID string) (entitymodel.Favorite, error)
	GetUserFavoritePlaces(ctx context.Context, userID string, country *string, administrativeAreaLevel1 *string, orderBy *enum.FavoriteOrderBy) (entitymodel.UserFavoritePlaces, error)
	// GetPlaceFavoriteUsers(ctx context.Context, placeID string, country *string, administrativeAreaLevel1 *string, orderBy *enum.FavoriteOrderBy) (entitymodel.PlaceFavoriteUsers, error)
	UpdateFavorite(ctx context.Context, favorite entitymodel.Favorite) error
}
