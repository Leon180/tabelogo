package inject

import (
	"logger-service/controller"
	grpcservice "logger-service/grpc/service"
)

// ControllerHandle Controller Handle
type ControllerHandle struct {
	LogController    *controller.LogController
	LogServiceServer *grpcservice.LogServiceServer
}

func newControllerHandle(
	logController *controller.LogController,
	logServiceServer *grpcservice.LogServiceServer,
) *ControllerHandle {
	return &ControllerHandle{
		LogController:    logController,
		LogServiceServer: logServiceServer,
	}
}
