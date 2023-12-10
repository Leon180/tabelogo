package main

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

const (
	webPort         = "80"
	maxCollectLinks = 4
)

func main() {
	config, err := LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	// connect to redis
	redisTabelogo := connectToRedis(config.RedisConnectTabelogo)

	server := NewServer(RedisInstance{
		Tabelogo: redisTabelogo,
	})
	err = server.Run(":" + webPort)
	if err != nil {
		panic(err)
	}
}

func connectToRedis(redis_conn string) *redis.Client {
	opts, err := redis.ParseURL(redis_conn)
	if err != nil {
		log.Panic(err)
	}
	rdb := redis.NewClient(opts)
	err = rdb.Ping(context.Background()).Err()
	if err != nil {
		log.Panic(err)
	}
	return rdb
}
