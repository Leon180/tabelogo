package inject

import (
	"google-map/controller"
	grpcservice "google-map/grpc/service"
)

// ControllerHandle Controller Handle
type ControllerHandle struct {
	GooglePlaceSearchController *controller.GooglePlaceSearchHandle
	GoogleMapServiceServer      *grpcservice.GoogleMapServiceServer
}

func newControllerHandle(
	googlePlaceSearchController *controller.GooglePlaceSearchHandle,
	googleMapServiceServer *grpcservice.GoogleMapServiceServer,
) *ControllerHandle {
	return &ControllerHandle{
		GooglePlaceSearchController: googlePlaceSearchController,
		GoogleMapServiceServer:      googleMapServiceServer,
	}
}
