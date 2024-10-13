package repository

import (
	"authenticate/model/entitymodel"
	"context"
)

type PlaceHandler interface {
	CreatePlace(ctx context.Context, place entitymodel.Place) error
	GetPlaceByGoogleID(ctx context.Context, googleID string) (entitymodel.Place, error)
	GetPlaceByID(ctx context.Context, placeID string) (entitymodel.Place, error)
	UpdatePlace(ctx context.Context, placeID string, updates map[string]interface{}) error
	DeletePlace(ctx context.Context, placeID string) error
}
