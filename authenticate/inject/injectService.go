package inject

import "github.com/google/wire"

var serviceHandleSet = wire.NewSet(
	provideSavePlaceService,
	provideGetPlaceService,
	provideRegistUserService,
	provideLoginUserService,
	provideRenewAccessTokenService,
	provideSaveFavoriteService,
	provideGetUserFavoritesService,
	provideLogoutUserService,
)
