package inject

import "github.com/google/wire"

var redisRepositoryHandleSet = wire.NewSet(
	providePlaceRedisRepository,
	providePlaceWithTransactionRedisRepository,
	provideSessionRedisRepository,
	provideSessionWithTransactionRedisRepository,
)
