package inject

import (
	"tabelog-spider/controller"
	grpcservice "tabelog-spider/grpc/service"
)

// ControllerHandle Controller Handle
type ControllerHandle struct {
	GetTabelogInfoController    *controller.GetTabelogInfoHandle
	TabelogoSpiderServiceServer *grpcservice.TabelogoSpiderServiceServer
}

func newControllerHandle(
	getTabelogInfoController *controller.GetTabelogInfoHandle,
	tabelogoSpiderServiceServer *grpcservice.TabelogoSpiderServiceServer,

) *ControllerHandle {
	return &ControllerHandle{
		GetTabelogInfoController:    getTabelogInfoController,
		TabelogoSpiderServiceServer: tabelogoSpiderServiceServer,
	}
}
