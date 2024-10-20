package inject

import "authenticate/controller"

// ControllerHandle Controller Handle
type ControllerHandle struct {
	registUserController       *controller.RegistUserControllerHandle
	loginUserController        *controller.LoginUserControllerHandle
	renewAccessTokenController *controller.RenewAccessTokenControllerHandle
	saveFavoriteController     *controller.SaveFavoriteControllerHandle
	getUserFavoritesController *controller.GetUserFavoritesControllerHandle
	savePlaceController        *controller.SavePlaceControllerHandle
	getPlaceController         *controller.GetPlaceControllerHandle
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
		registUserController:       registUserController,
		loginUserController:        loginUserController,
		renewAccessTokenController: renewAccessTokenController,
		saveFavoriteController:     saveFavoriteController,
		getUserFavoritesController: getUserFavoritesController,
		savePlaceController:        savePlaceController,
		getPlaceController:         getPlaceController,
	}
}
