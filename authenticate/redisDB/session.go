package redisDB

import (
	"github.com/go-redis/redis"
)

type SessionHandler interface {
	// GetSession(ctx context.Context, key string) (string, error)
	// SetSession(ctx context.Context, key string, value string) error
	// DeleteSession(ctx context.Context, key string) error
}

func NewSessionHandler(redisClient *redis.Client) SessionHandler {
	return &SessionHandle{redisClient: redisClient}
}

type SessionHandle struct {
	redisClient *redis.Client
}
