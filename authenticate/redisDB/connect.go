package redisDB

import (
	"authenticate/utility"

	"github.com/go-redis/redis"
)

type RedisSession *redis.Client
type RedisPlace *redis.Client

func ConnectToRedis(redisConn string) *redis.Client {
	var (
		opts *redis.Options
		rdb  *redis.Client
		err  error
	)

	if opts, err = redis.ParseURL(redisConn); err != nil {
		utility.SugarLogger.Fatal("Error during redis connection, error: %s", err)
		return nil
	}
	rdb = redis.NewClient(opts)
	if err = rdb.Ping().Err(); err != nil {
		utility.SugarLogger.Fatal("Error during redis connection, error: %s", err)
	}

	return rdb
}
