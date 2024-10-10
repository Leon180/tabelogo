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
	GetPlace(ctx context.Context, key string) (entitymodel.Place, error)
	SetPlace(ctx context.Context, key string, place entitymodel.Place, expiry *time.Duration) error
	DeletePlace(ctx context.Context, key string) error
}

func NewPlaceHandler(redisClient *redis.Client) PlaceHandler {
	return &PlaceHandle{redisClient: redisClient}
}

type PlaceHandle struct {
	redisClient *redis.Client
}

func (p *PlaceHandle) GetPlace(ctx context.Context, key string) (entitymodel.Place, error) {
	var place entitymodel.Place
	st, err := p.redisClient.JSONGet(ctx, key, "$").Result()
	if err != nil {
		if err == redis.Nil {
			return entitymodel.Place{}, nil
		}
		utility.SugarLogger.Error("error during get place, error: %s", err)
		return entitymodel.Place{}, err
	}
	if err = json.Unmarshal([]byte(st), &place); err != nil {
		utility.SugarLogger.Error("error during unmarshal place, error: %s", err)
		return entitymodel.Place{}, err
	}
	return place, nil
}

func (p *PlaceHandle) SetPlace(ctx context.Context, key string, place entitymodel.Place, expiry *time.Duration) error {
	if _, err := p.redisClient.JSONSet(ctx, key, "$", place).Result(); err != nil {
		utility.SugarLogger.Error("error during set place, error: %s", err)
		return err
	}
	if expiry == nil {
		return nil
	}
	if err := p.redisClient.Expire(ctx, key, *expiry).Err(); err != nil {
		utility.SugarLogger.Error("error during set place expiry, error: %s", err)
		return err
	}
	return nil
}

func (p *PlaceHandle) DeletePlace(ctx context.Context, key string) error {
	if err := p.redisClient.Del(ctx, key).Err(); err != nil {
		utility.SugarLogger.Error("error during delete place, error: %s", err)
		return err
	}
	return nil
}
