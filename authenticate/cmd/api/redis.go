package main

import (
	db "authenticate/cmd/data/sqlc"
	"authenticate/token"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RedisInstance struct {
	Session  *redis.Client
	Place    *redis.Client
	Tabelogo *redis.Client
}

// 3 kind of redis cache:
// 1. cache for user

func (s *Server) GetSessionInRedisOrDatabase(c *gin.Context, authPayload *token.Payload) (db.Session, error) {
	// get user sessio  in redis
	session, err := s.getSessionInRedis(c, authPayload.Email)
	if err != nil {
		return db.Session{}, err
	}
	// if empty session, get session in database
	if session == (db.Session{}) {
		session, err = s.store.GetSession(c, authPayload.ID)
		if err != nil {

			return session, err
		}
	}

	return session, nil
}

func (s *Server) setSessionInRedisWithExpiry(ctx *gin.Context, data db.Session, expiry time.Duration) error {
	_, err := s.redisInstance.Session.JSONSet(ctx, data.Email, "$", data).Result()
	if err != nil {
		return err
	}
	err = s.redisInstance.Session.Expire(ctx, data.Email, expiry).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) getSessionInRedis(ctx *gin.Context, email string) (db.Session, error) {
	var session []db.Session
	st, err := s.redisInstance.Session.JSONGet(ctx, email, "$").Result()
	if err != nil {
		return db.Session{}, err
	}
	if st == "" || st == "null" {
		return db.Session{}, nil
	}
	err = json.Unmarshal([]byte(st), &session)
	if err != nil {
		return db.Session{}, err
	}
	return session[0], nil
}

func (s *Server) deleteSessionInRedis(ctx *gin.Context, email string) error {
	err := s.redisInstance.Session.Del(ctx, email).Err()
	if err != nil {
		return err
	}
	return nil
}

// 2. cache for place
// func (s *Server) checkPlaceInRedis(place string) db.Place {}
// func (s *Server) savePlaceInRedis(place string) db.Place    {}
// 3. cache for tabelogo
// func (s *Server) checkTabelogoInRedis(tabelogo string) db.Tabelogo {}
// func (s *Server) saveTabelogoInRedis(tabelogo string) db.Tabelogo    {}
