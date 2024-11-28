package inject

import (
	"logger-service/controller"
	"logger-service/repository"
	"logger-service/service"

	"go.mongodb.org/mongo-driver/mongo"
)

func provideCreateLogRepository(db *mongo.Client) repository.CreateLogHandler {
	return repository.NewCreateLogRepository(db)
}

func provideReadLogRepository(db *mongo.Client) repository.ReadLogHandler {
	return repository.NewReadLogRepository(db)
}

func provideUpdateLogRepository(db *mongo.Client) repository.UpdateLogHandler {
	return repository.NewUpdateLogRepository(db)
}

func provideDeleteLogRepository(db *mongo.Client) repository.DeleteLogHandler {
	return repository.NewDeleteLogRepository(db)
}

func provideCreateLogService(createLogRepository repository.CreateLogHandler) service.CreateLogHandler {
	return service.NewCreateLogService(createLogRepository)
}

func provideReadLogService(readLogRepository repository.ReadLogHandler) service.ReadLogHandler {
	return service.NewReadLogService(readLogRepository)
}

func provideUpdateLogService(updateLogRepository repository.UpdateLogHandler) service.UpdateLogHandler {
	return service.NewUpdateLogService(updateLogRepository)
}

func provideDeleteLogService(deleteLogRepository repository.DeleteLogHandler) service.DeleteLogHandler {
	return service.NewDeleteLogService(deleteLogRepository)
}

func provideLogController(
	createLogService service.CreateLogHandler,
	readLogService service.ReadLogHandler,
	updateLogService service.UpdateLogHandler,
	deleteLogService service.DeleteLogHandler,
) *controller.LogController {
	return controller.NewLogController(createLogService, readLogService, updateLogService, deleteLogService)
}
