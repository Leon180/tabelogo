package inject

import (
	"authenticate/controller"
	"authenticate/redisDB"
	"authenticate/repository"
	"authenticate/repository/postgresqlRepository"
	"authenticate/service"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func providePlaceRepository(db *gorm.DB) repository.PlaceHandler {
	return postgresqlRepository.NewPlaceHandle(db)
}

func providePlaceWithTransactionRepository(db *gorm.DB) repository.PlaceWithTransactionHandler {
	return postgresqlRepository.NewPlaceWithTransactionHandler(db)
}

func providePlaceRedisRepository(redisClient *redis.Client) redisDB.PlaceHandler {
	return redisDB.NewPlaceHandler(redisClient)
}

func providePlaceWithTransactionRedisRepository(redisClient *redis.Client) redisDB.PlaceWithTransactionHandler {
	return redisDB.NewPlaceWithTransactionHandler(redisClient)
}

func provideSavePlaceService(placeWithTransactionRepository repository.PlaceWithTransactionHandler) service.SavePlaceServiceHandler {
	return service.NewSavePlaceServiceHandler(placeWithTransactionRepository)
}

func provideGetPlaceService(placeRepository repository.PlaceHandler, redisPlace redisDB.PlaceHandler) service.GetPlaceServiceHandler {
	return service.NewGetPlaceServiceHandler(placeRepository, redisPlace)
}

func provideSavePlaceController(savePlaceServiceHandler service.SavePlaceServiceHandler) *controller.SavePlaceControllerHandle {
	return controller.NewSavePlaceControllerHandle(savePlaceServiceHandler)
}

func provideGetPlaceController(getPlaceServiceHandler service.GetPlaceServiceHandler) *controller.GetPlaceControllerHandle {
	return controller.NewGetPlaceControllerHandle(getPlaceServiceHandler)
}
