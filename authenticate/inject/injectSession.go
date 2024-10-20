package inject

import (
	"authenticate/redisDB"

	"github.com/redis/go-redis/v9"
)

func provideSessionRedisRepository(redisClient *redis.Client) redisDB.SessionHandler {
	return redisDB.NewSessionHandler(redisClient)
}

func provideSessionWithTransactionRedisRepository(redisClient *redis.Client) redisDB.SessionWithTransactionHandler {
	return redisDB.NewSessionWithTransactionHandler(redisClient)
}
