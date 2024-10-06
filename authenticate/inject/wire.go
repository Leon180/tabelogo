//go:build wireinject
// +build wireinject

package inject

import (
	"authenticate/config"

	"github.com/go-redis/redis"
	"github.com/google/wire"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func InitControllerHandle(
	db *gorm.DB,
	config config.Config,
	logger *zap.Logger,
	redisClient *redis.Client,
) *ControllerHandle {
	wire.Build(
		controllerHandleSet,
		serviceHandleSet,
		repositoryHandleSet,
		newControllerHandle,
	)
	return &ControllerHandle{}
}
