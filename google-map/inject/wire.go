//go:build wireinject
// +build wireinject

package inject

import (
	"google-map/config"

	"github.com/google/wire"
	"go.uber.org/zap"
)

func InitControllerHandle(
	config config.Config,
	logger *zap.Logger,
) *ControllerHandle {
	wire.Build(
		serviceHandleSet,
		controllerHandleSet,
		newControllerHandle,
	)
	return &ControllerHandle{}
}
