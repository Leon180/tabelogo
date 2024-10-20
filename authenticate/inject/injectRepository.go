package inject

import (
	"github.com/google/wire"
)

var repositoryHandleSet = wire.NewSet(
	providePlaceRepository,
	providePlaceWithTransactionRepository,
	provideUserRepository,
	provideUserWithTransactionRepository,
	provideUserAndSessionWithTransactionRepository,
	provideUserAndPlaceAndFavoriteWithTransactionRepository,
	provideFavoriteRepository,
)
