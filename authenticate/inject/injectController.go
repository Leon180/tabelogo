package inject

import (
	"github.com/google/wire"
)

var controllerHandleSet = wire.NewSet(
	provideRegistUserController,
	provideLoginUserController,
	provideRenewAccessTokenController,
	provideSaveFavoriteController,
	provideGetUserFavoritesController,
	provideSavePlaceController,
	provideGetPlaceController,
	provideLogoutUserController,
)
