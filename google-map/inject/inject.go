package inject

import "google-map/controller"

// ControllerHandle Controller Handle
type ControllerHandle struct {
	GooglePlaceSearchController *controller.GooglePlaceSearchHandle
}

func newControllerHandle(
	googlePlaceSearchController *controller.GooglePlaceSearchHandle,
) *ControllerHandle {
	return &ControllerHandle{
		GooglePlaceSearchController: googlePlaceSearchController,
	}
}
