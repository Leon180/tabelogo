package main

import (
	db "authenticate/cmd/data/sqlc"
	"authenticate/token"
	"database/sql"
	"encoding/json"
	"math/rand"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type CacheInstance struct {
	Session *redis.Client
	Place   *redis.Client
}

// 2 kind of redis cache:
// 1. cache for session
func (s *Server) GetSessionInCacheOrDatabase(c *gin.Context, authPayload *token.Payload) (db.Session, string, error) {
	// get user sessio  in redis
	session, err := s.getSessionInCache(c, authPayload.Email)
	// if error, return error
	if err != nil {
		return db.Session{}, "error", err
	}
	// if empty session, get session in database
	if session == (db.Session{}) {
		session, err = s.store.GetSession(c, authPayload.ID)
		if err != nil {
			return db.Session{}, "error", err
		}
		return session, "database", nil
	}
	return session, "cache", nil
}

func (s *Server) setSessionInCacheWithExpiry(ctx *gin.Context, data db.Session, expiry time.Duration) error {
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

func (s *Server) getSessionInCache(ctx *gin.Context, email string) (db.Session, error) {
	var session []db.Session
	st, err := s.redisInstance.Session.JSONGet(ctx, email, "$").Result()
	if err != nil {
		if err == redis.Nil {
			return db.Session{}, nil
		}
		return db.Session{}, err
	}
	if st == "" || st == "null" {
		s.redisInstance.Session.Del(ctx, email)
		return db.Session{}, nil
	}
	err = json.Unmarshal([]byte(st), &session)
	if err != nil {
		return db.Session{}, err
	}
	return session[0], nil
}

func (s *Server) deleteSessionInCache(ctx *gin.Context, email string) error {
	err := s.redisInstance.Session.Del(ctx, email).Err()
	if err != nil {
		return err
	}
	return nil
}

// 2. cache for place
// if place not exists in redis, then check place in database
// if both not exists, then create place in database and set place in redis, and return with a bool value: true
func (s *Server) GetPlaceInCacheOrDatabaseAndCreateIfNotExist(c *gin.Context, request SaveFavoriteRequest) (db.Place, string, error) {
	// check place in redis
	place, err := s.getPlaceInCache(c, request.GoogleID)
	if err != nil {
		return db.Place{}, "error", err
	}
	// place not exists in redis, then check place in database
	if reflect.DeepEqual(place, db.Place{}) {
		place, err = s.store.GetPlaceByGoogleId(c, request.GoogleID)
		if err != nil {
			// if place not found
			if err == sql.ErrNoRows {
				arg := db.CreatePlaceParams{
					GoogleID:              request.GoogleID,
					TwDisplayName:         request.TwDisplayName,
					TwFormattedAddress:    request.TwFormattedAddress,
					TwWeekdayDescriptions: pq.StringArray(request.TwWeekdayDescriptions),
					GoogleMapUri:          request.GoogleMapUri,
					Lat:                   request.Lat,
					Lng:                   request.Lng,
					Types:                 pq.StringArray(request.Types),
				}
				if len(request.AdministrativeAreaLevel1) != 0 {
					arg.AdministrativeAreaLevel1 = request.AdministrativeAreaLevel1
				}
				if len(request.Country) != 0 {
					arg.Country = request.Country
				}
				if len(request.InternationalPhoneNumber) != 0 {
					arg.InternationalPhoneNumber = request.InternationalPhoneNumber
				}
				if len(request.PrimaryType) != 0 {
					arg.PrimaryType = request.PrimaryType
				}
				if len(request.Rating) != 0 {
					arg.Rating = request.Rating
				}
				if request.UserRatingCount != 0 {
					arg.UserRatingCount = request.UserRatingCount
				}
				if len(request.WebsiteUri) != 0 {
					arg.WebsiteUri = request.WebsiteUri
				}
				// create place
				place, err = s.store.CreatePlace(c, arg)
				if err != nil {
					return db.Place{}, "error", err // create place failed
				}
				// set place in redis with expiry
				err = s.setPlaceInCacheWithExpiry(c, place, 10*time.Minute+time.Duration(rand.Intn(5))*time.Minute)
				if err != nil {
					return place, "database_create", err
				}
				return place, "database_create", nil
			} else {
				return db.Place{}, "error", err
			}
		}
		// place exists in database and not in redis, handler will update place in db and redis
		return place, "database_exist_cache_not_exist", nil
	}
	return place, "cache_exist", nil
}

func (s *Server) setPlaceInCacheWithExpiry(ctx *gin.Context, data db.Place, expiry time.Duration) error {
	_, err := s.redisInstance.Place.JSONSet(ctx, data.GoogleID, "$", data).Result()
	if err != nil {
		return err
	}
	err = s.redisInstance.Place.Expire(ctx, data.GoogleID, expiry).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) getPlaceInCache(ctx *gin.Context, googleId string) (db.Place, error) {
	var place []db.Place
	st, err := s.redisInstance.Place.JSONGet(ctx, googleId, "$").Result()
	if err != nil {
		if err == redis.Nil {
			return db.Place{}, nil
		}
		return db.Place{}, err
	}
	if st == "" || st == "null" {
		s.redisInstance.Place.Del(ctx, googleId)
		return db.Place{}, nil
	}
	err = json.Unmarshal([]byte(st), &place)
	if err != nil {
		return db.Place{}, err
	}
	return place[0], nil
}
