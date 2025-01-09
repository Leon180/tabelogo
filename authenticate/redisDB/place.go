package redisDB

import (
	"authenticate/model/entitymodel"
	"authenticate/utility"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type PlaceHandler interface {
	GetPlace(ctx context.Context, placeGoogleID string) (entitymodel.Place, error)
	SetPlace(ctx context.Context, placeGoogleID string, place entitymodel.Place, expiry *time.Duration) error
	DeletePlace(ctx context.Context, placeGoogleID string) error
}

func NewPlaceHandler(redisClient *redis.Client) PlaceHandler {
	return &PlaceHandle{redisClient: redisClient}
}

type PlaceHandle struct {
	redisClient *redis.Client
}

func (p *PlaceHandle) GetPlace(ctx context.Context, placeGoogleID string) (entitymodel.Place, error) {
	var place entitymodel.Place
	st, err := p.redisClient.JSONGet(ctx, placeGoogleID, "$").Result()
	if err != nil {
		if err == redis.Nil {
			return entitymodel.Place{}, nil
		}
		utility.LogWithTraceID(ctx, "error during get place, error: %s", err)
		return entitymodel.Place{}, err
	}
	if err = json.Unmarshal([]byte(st), &place); err != nil {
		utility.LogWithTraceID(ctx, "error during unmarshal place, error: %s", err)
		return entitymodel.Place{}, err
	}
	return place, nil
}

func (p *PlaceHandle) SetPlace(ctx context.Context, placeGoogleID string, place entitymodel.Place, expiry *time.Duration) error {
	if _, err := p.redisClient.JSONSet(ctx, placeGoogleID, "$", place).Result(); err != nil {
		utility.LogWithTraceID(ctx, "error during set place, error: %s", err)
		return err
	}
	if expiry == nil {
		return nil
	}
	if err := p.redisClient.Expire(ctx, placeGoogleID, *expiry).Err(); err != nil {
		utility.LogWithTraceID(ctx, "error during set place expiry, error: %s", err)
		return err
	}
	return nil
}

func (p *PlaceHandle) DeletePlace(ctx context.Context, placeGoogleID string) error {
	if err := p.redisClient.Del(ctx, placeGoogleID).Err(); err != nil {
		utility.LogWithTraceID(ctx, "error during delete place, error: %s", err)
		return err
	}
	return nil
}

type PlaceWithTransactionHandler interface {
	PlaceHandler
	WithTransaction
}

func NewPlaceWithTransactionHandler(redisClient *redis.Client) PlaceWithTransactionHandler {
	return &PlaceWithTransactionHandle{
		PlaceHandle: PlaceHandle{redisClient: redisClient},
	}
}

type PlaceWithTransactionHandle struct {
	PlaceHandle
}

func (p *PlaceWithTransactionHandle) WithTransaction(ctx context.Context, fn func(tx *redis.Tx) error, keys ...string) error {
	return p.PlaceHandle.redisClient.Watch(ctx, fn, keys...)
}
