package inject

import "tabelog-spider/controller"

// ControllerHandle Controller Handle
type ControllerHandle struct {
	GetTabelogInfoController *controller.GetTabelogInfoHandle
}

func newControllerHandle(
	getTabelogInfoController *controller.GetTabelogInfoHandle,
) *ControllerHandle {
	return &ControllerHandle{
		GetTabelogInfoController: getTabelogInfoController,
	}
}
