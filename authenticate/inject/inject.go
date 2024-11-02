package inject

import "authenticate/controller"

// ControllerHandle Controller Handle
type ControllerHandle struct {
	RegistUserController       *controller.RegistUserControllerHandle
	LoginUserController        *controller.LoginUserControllerHandle
	RenewAccessTokenController *controller.RenewAccessTokenControllerHandle
	SaveFavoriteController     *controller.SaveFavoriteControllerHandle
	GetUserFavoritesController *controller.GetUserFavoritesControllerHandle
	SavePlaceController        *controller.SavePlaceControllerHandle
	GetPlaceController         *controller.GetPlaceControllerHandle
}

func newControllerHandle(
	registUserController *controller.RegistUserControllerHandle,
	loginUserController *controller.LoginUserControllerHandle,
	renewAccessTokenController *controller.RenewAccessTokenControllerHandle,
	saveFavoriteController *controller.SaveFavoriteControllerHandle,
	getUserFavoritesController *controller.GetUserFavoritesControllerHandle,
	savePlaceController *controller.SavePlaceControllerHandle,
	getPlaceController *controller.GetPlaceControllerHandle,
) *ControllerHandle {
	return &ControllerHandle{
		RegistUserController:       registUserController,
		LoginUserController:        loginUserController,
		RenewAccessTokenController: renewAccessTokenController,
		SaveFavoriteController:     saveFavoriteController,
		GetUserFavoritesController: getUserFavoritesController,
		SavePlaceController:        savePlaceController,
		GetPlaceController:         getPlaceController,
	}
}
