//go:build wireinject
// +build wireinject

package inject

import (
	"logger-service/config"

	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func InitControllerHandle(
	config config.Config,
	logger *zap.Logger,
	db *mongo.Client,
) *ControllerHandle {
	wire.Build(
		repositoryHandleSet,
		serviceHandleSet,
		controllerHandleSet,
		grpcServerSet,
		newControllerHandle,
	)
	return &ControllerHandle{}
}
