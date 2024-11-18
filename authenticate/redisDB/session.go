package redisDB

import (
	"authenticate/model/entitymodel"
	"authenticate/utility"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionHandler interface {
	GetSession(ctx context.Context, key string) (entitymodel.Session, error)
	SetSession(ctx context.Context, key string, session entitymodel.Session, expiry *time.Duration) error
	DeleteSession(ctx context.Context, key string) error
}

func NewSessionHandler(redisClient *redis.Client) SessionHandler {
	return &SessionHandle{redisClient: redisClient}
}

type SessionHandle struct {
	redisClient *redis.Client
}

func (s *SessionHandle) GetSession(ctx context.Context, key string) (entitymodel.Session, error) {
	var session []entitymodel.Session
	st, err := s.redisClient.JSONGet(ctx, key, "$").Result()
	if err != nil {
		if err == redis.Nil {
			return entitymodel.Session{}, nil
		}
		utility.SugarLogger.Error("error during get session, error: %s", err)
		return entitymodel.Session{}, err
	}
	if st == "" {
		return entitymodel.Session{}, nil
	}
	if err = json.Unmarshal([]byte(st), &session); err != nil {
		utility.SugarLogger.Error("error during unmarshal session, error: %s", err)
		return entitymodel.Session{}, err
	}
	return session[0], nil
}

func (s *SessionHandle) SetSession(ctx context.Context, key string, session entitymodel.Session, expiry *time.Duration) error {
	if _, err := s.redisClient.JSONSet(ctx, key, "$", session).Result(); err != nil {
		utility.SugarLogger.Error("error during set session, error: %s", err)
		return err
	}
	if expiry == nil {
		return nil
	}
	if err := s.redisClient.Expire(ctx, key, *expiry).Err(); err != nil {
		utility.SugarLogger.Error("error during set session expiry, error: %s", err)
		return err
	}
	return nil
}

func (s *SessionHandle) DeleteSession(ctx context.Context, key string) error {
	if err := s.redisClient.Del(ctx, key).Err(); err != nil {
		utility.SugarLogger.Error("error during delete session, error: %s", err)
		return err
	}
	return nil
}

type SessionWithTransactionHandler interface {
	SessionHandler
	WithTransaction
}

func NewSessionWithTransactionHandler(redisClient *redis.Client) SessionWithTransactionHandler {
	return &SessionWithTransactionHandle{
		SessionHandle: SessionHandle{redisClient: redisClient},
	}
}

type SessionWithTransactionHandle struct {
	SessionHandle
}

func (s *SessionWithTransactionHandle) WithTransaction(ctx context.Context, fn func(tx *redis.Tx) error, keys ...string) error {
	return s.SessionHandle.redisClient.Watch(ctx, fn, keys...)
}
