package redisDB

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type WithTransaction interface {
	WithTransaction(ctx context.Context, fn func(tx *redis.Tx) error, keys ...string) error
}
