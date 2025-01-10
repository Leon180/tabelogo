//go:build wireinject
// +build wireinject

package inject

import (
	"github.com/google/wire"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

func InitControllerHandle(
	logger *zap.Logger,
	rateLimiter *rate.Limiter,
) *ControllerHandle {
	wire.Build(
		serviceHandleSet,
		controllerHandleSet,
		grpcServiceSet,
		newControllerHandle,
	)
	return &ControllerHandle{}
}
