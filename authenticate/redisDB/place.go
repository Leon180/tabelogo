package redisDB

import "github.com/go-redis/redis"

type PlaceHandler interface {
	// GetSession(ctx context.Context, key string) (string, error)
	// SetSession(ctx context.Context, key string, value string) error
	// DeleteSession(ctx context.Context, key string) error
}

func NewPlaceHandler(redisClient *redis.Client) PlaceHandler {
	return &PlaceHandle{redisClient: redisClient}
}

type PlaceHandle struct {
	redisClient *redis.Client
}
