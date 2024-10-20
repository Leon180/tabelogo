//go:build wireinject
// +build wireinject

package inject

import (
	"authenticate/config"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func InitControllerHandle(
	db *gorm.DB,
	config *config.Config,
	logger *zap.Logger,
	redisClient *redis.Client,
	symmetricKey string,
) *ControllerHandle {
	wire.Build(
		repositoryHandleSet,
		redisRepositoryHandleSet,
		tokenMakerHandleSet,
		serviceHandleSet,
		controllerHandleSet,
		newControllerHandle,
	)
	return &ControllerHandle{}
}
