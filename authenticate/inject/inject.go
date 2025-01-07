package inject

import (
	"authenticate/controller"
	grpcservice "authenticate/grpc/service"
)

// ControllerHandle Controller Handle
type ControllerHandle struct {
	RegistUserController       *controller.RegistUserControllerHandle
	LoginUserController        *controller.LoginUserControllerHandle
	RenewAccessTokenController *controller.RenewAccessTokenControllerHandle
	SaveFavoriteController     *controller.SaveFavoriteControllerHandle
	GetUserFavoritesController *controller.GetUserFavoritesControllerHandle
	SavePlaceController        *controller.SavePlaceControllerHandle
	GetPlaceController         *controller.GetPlaceControllerHandle
	LogoutUserController       *controller.LogoutUserControllerHandle
	UserServiceServer          *grpcservice.UserServiceServer
	PlaceServiceServer         *grpcservice.PlaceServiceServer
}

func newControllerHandle(
	registUserController *controller.RegistUserControllerHandle,
	loginUserController *controller.LoginUserControllerHandle,
	renewAccessTokenController *controller.RenewAccessTokenControllerHandle,
	saveFavoriteController *controller.SaveFavoriteControllerHandle,
	getUserFavoritesController *controller.GetUserFavoritesControllerHandle,
	savePlaceController *controller.SavePlaceControllerHandle,
	getPlaceController *controller.GetPlaceControllerHandle,
	logoutUserController *controller.LogoutUserControllerHandle,
	userServiceServer *grpcservice.UserServiceServer,
	placeServiceServer *grpcservice.PlaceServiceServer,
) *ControllerHandle {
	return &ControllerHandle{
		RegistUserController:       registUserController,
		LoginUserController:        loginUserController,
		RenewAccessTokenController: renewAccessTokenController,
		SaveFavoriteController:     saveFavoriteController,
		GetUserFavoritesController: getUserFavoritesController,
		SavePlaceController:        savePlaceController,
		GetPlaceController:         getPlaceController,
		LogoutUserController:       logoutUserController,
		UserServiceServer:          userServiceServer,
		PlaceServiceServer:         placeServiceServer,
	}
}
