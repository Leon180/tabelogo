package inject

import (
	"authenticate/config"
	"authenticate/controller"
	"authenticate/redisDB"
	"authenticate/repository"
	"authenticate/repository/postgresqlRepository"
	"authenticate/service"
	"authenticate/token"

	"gorm.io/gorm"
)

func provideUserRepository(db *gorm.DB) repository.UserHandler {
	return postgresqlRepository.NewUserHandler(db)
}

func provideUserWithTransactionRepository(db *gorm.DB) repository.UserWithTransactionHandler {
	return postgresqlRepository.NewUserWithTransactionHandler(db)
}

func provideUserAndSessionWithTransactionRepository(db *gorm.DB) repository.UserAndSessionWithTransactionHandler {
	return postgresqlRepository.NewUserAndSessionWithTransactionHandler(db)
}

func provideUserAndPlaceAndFavoriteWithTransactionRepository(db *gorm.DB) repository.UserAndPlaceAndFavoriteWithTransactionHandler {
	return postgresqlRepository.NewUserAndPlaceAndFavoriteWithTransactionHandler(db)
}

func provideRegistUserService(userWithTransactionRepository repository.UserWithTransactionHandler) service.RegistUserServiceHandler {
	return service.NewRegistUserServiceHandler(userWithTransactionRepository)
}

func provideLoginUserService(
	userAndSessionWithTransactionRepository repository.UserAndSessionWithTransactionHandler,
	tokenMaker token.Maker,
	redisSession redisDB.SessionWithTransactionHandler,
	config *config.Config,
) service.LoginUserServiceHandler {
	return service.NewLoginUserServiceHandler(
		userAndSessionWithTransactionRepository,
		tokenMaker,
		redisSession,
		config,
	)
}

func provideRenewAccessTokenService(
	userAndSessionWithTransactionRepository repository.UserAndSessionWithTransactionHandler,
	tokenMaker token.Maker,
	SessionWithTransactionRedis redisDB.SessionWithTransactionHandler,
	config *config.Config,
) service.RenewAccessTokenServiceHandler {
	return service.NewRenewAccessTokenServiceHandler(
		userAndSessionWithTransactionRepository,
		tokenMaker,
		SessionWithTransactionRedis,
		config,
	)
}

func provideSaveFavoriteService(
	userAndPlaceAndFavoriteWithTransactionRepository repository.UserAndPlaceAndFavoriteWithTransactionHandler,
	redisPlace redisDB.PlaceWithTransactionHandler,
	config *config.Config,
) service.SaveFavoriteServiceHandler {
	return service.NewSaveFavoriteServiceHandler(
		userAndPlaceAndFavoriteWithTransactionRepository,
		redisPlace,
		config,
	)
}

func provideGetUserFavoritesService(
	favoriteRepository repository.FavoriteHandler,
) service.GetUserFavoritesServiceHandler {
	return service.NewGetUserFavoritesServiceHandler(favoriteRepository)
}

func provideRegistUserController(registUserServiceHandler service.RegistUserServiceHandler) *controller.RegistUserControllerHandle {
	return controller.NewRegistUserControllerHandle(registUserServiceHandler)
}

func provideLoginUserController(loginUserServiceHandler service.LoginUserServiceHandler) *controller.LoginUserControllerHandle {
	return controller.NewLoginUserControllerHandle(loginUserServiceHandler)
}

func provideRenewAccessTokenController(renewAccessTokenServiceHandler service.RenewAccessTokenServiceHandler) *controller.RenewAccessTokenControllerHandle {
	return controller.NewRenewAccessTokenControllerHandle(renewAccessTokenServiceHandler)
}

func provideSaveFavoriteController(saveFavoriteServiceHandler service.SaveFavoriteServiceHandler) *controller.SaveFavoriteControllerHandle {
	return controller.NewSaveFavoriteControllerHandle(saveFavoriteServiceHandler)
}

func provideGetUserFavoritesController(getUserFavoritesServiceHandler service.GetUserFavoritesServiceHandler) *controller.GetUserFavoritesControllerHandle {
	return controller.NewGetUserFavoritesControllerHandle(getUserFavoritesServiceHandler)
}
