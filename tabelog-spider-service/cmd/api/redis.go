package main

import (
	"github.com/redis/go-redis/v9"
)

type RedisInstance struct {
	Tabelogo *redis.Client
}

// 1 kind of redis cache:

// 1. cache for tabelogo
// pipeline
// func (s *Server) setTabelogoInRedisWithExpiry(ctx *gin.Context, collection []TabelogInfo, expiry time.Duration) error {

// }

// func (s *Server) getTabelogoInRedis(ctx *gin.Context, googleId string) (db.Place, error) {

// }
