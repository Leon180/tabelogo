//go:build wireinject
// +build wireinject

package inject

import (
	"github.com/google/wire"
	"go.uber.org/zap"
)

func InitControllerHandle(
	logger *zap.Logger,
) *ControllerHandle {
	wire.Build(
		serviceHandleSet,
		controllerHandleSet,
		newControllerHandle,
	)
	return &ControllerHandle{}
}
