package repository

import (
	"authenticate/model/entitymodel"
	"context"
)

type UserHandler interface {
	CreateUser(ctx context.Context, user entitymodel.User) error
	GetUserByEmail(ctx context.Context, email string) (entitymodel.User, error)
	UpdateUser(ctx context.Context, userID string, updates map[string]interface{}) error
	DeleteUser(ctx context.Context, userID string) error
}

type UserWithTransactionHandler interface {
	WithTransaction
	UserHandler
}

type UserAndSessionWithTransactionHandler interface {
	WithTransaction
	UserHandler
	SessionHandler
}

type UserAndPlaceAndFavoriteWithTransactionHandler interface {
	WithTransaction
	UserHandler
	PlaceHandler
	FavoriteHandler
}
