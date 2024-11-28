package inject

import "logger-service/controller"

// ControllerHandle Controller Handle
type ControllerHandle struct {
	LogController *controller.LogController
}

func newControllerHandle(
	logController *controller.LogController,
) *ControllerHandle {
	return &ControllerHandle{
		LogController: logController,
	}
}
