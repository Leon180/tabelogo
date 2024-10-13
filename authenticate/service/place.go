package service

import (
	"authenticate/config"
	"authenticate/model/entitymodel"
	"authenticate/redisDB"
	"authenticate/repository"
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SavePlaceServiceHandler interface {
	SavePlace(ctx context.Context, place entitymodel.Place) error
}

func NewSavePlaceServiceHandler(
	placeWithTransactionRepository repository.PlaceWithTransactionHandler,
) SavePlaceServiceHandler {
	return &SavePlaceServiceHandle{
		placeWithTransactionRepository: placeWithTransactionRepository,
	}
}

type SavePlaceServiceHandle struct {
	placeWithTransactionRepository repository.PlaceWithTransactionHandler
	redisPlace                     redisDB.PlaceWithTransactionHandler
	config                         *config.Config
}

func (handle *SavePlaceServiceHandle) SavePlace(ctx context.Context, place entitymodel.Place) error {
	// check if place already in redis, if yes, update place
	// if not exist, check if place exist in db,
	//	 if yes, update place in db
	//	 if no, create place in db
	// set created/updated place in redis
	var (
		err          error
		existedPlace entitymodel.Place
	)
	if err = handle.redisPlace.WithTransaction(ctx, func(tx *redis.Tx) error {
		if existedPlace, err = handle.redisPlace.GetPlace(ctx, place.GoogleID); err != nil {
			return err
		}
		if !existedPlace.IsExist() {
			return nil
		}
		if err = handle.redisPlace.DeletePlace(ctx, place.GoogleID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	if !existedPlace.IsExist() {
		if err = handle.placeWithTransactionRepository.WithTransaction(ctx, func(tx *gorm.DB) error {
			if existedPlace, err = handle.placeWithTransactionRepository.GetPlaceByGoogleID(ctx, place.GoogleID); err != nil {
				return err
			}
			if existedPlace.IsExist() {
				if err = handle.placeWithTransactionRepository.UpdatePlace(ctx, existedPlace.ID, place.GetUpdates()); err != nil {
					return err
				}
			}
			return handle.placeWithTransactionRepository.CreatePlace(ctx, place)
		}); err != nil {
			return err
		}
	} else {
		if err = handle.placeWithTransactionRepository.UpdatePlace(ctx, existedPlace.ID, place.GetUpdates()); err != nil {
			return err
		}
	}

	return handle.redisPlace.SetPlace(ctx, place.GoogleID, place, &handle.config.PlaceRedisExpiry)
}

type GetPlaceServiceHandler interface {
	GetPlace(ctx context.Context, googleID string) (entitymodel.Place, error)
}

func NewGetPlaceServiceHandler(
	placeRepository repository.PlaceHandler,
	redisPlace redisDB.PlaceHandler,
) GetPlaceServiceHandler {
	return &GetPlaceServiceHandle{
		placeRepository: placeRepository,
		redisPlace:      redisPlace,
	}
}

type GetPlaceServiceHandle struct {
	placeRepository repository.PlaceHandler
	redisPlace      redisDB.PlaceHandler
}

func (handle *GetPlaceServiceHandle) GetPlace(ctx context.Context, googleID string) (entitymodel.Place, error) {
	var (
		place entitymodel.Place
		err   error
	)
	if place, err = handle.redisPlace.GetPlace(ctx, googleID); err != nil {
		return entitymodel.Place{}, err
	}
	if !place.IsExist() {
		if place, err = handle.placeRepository.GetPlaceByGoogleID(ctx, googleID); err != nil {
			return entitymodel.Place{}, err
		}
	}
	return place, nil
}
