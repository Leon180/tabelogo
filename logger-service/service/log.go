package service

import (
	"context"
	"logger-service/model/entitymodel"
	"logger-service/model/enum"
	"logger-service/repository"
)

type CreateLogHandler interface {
	CreateLog(ctx context.Context, logEntry entitymodel.LogEntry) error
}

func NewCreateLogService(createLogRepository repository.CreateLogHandler) CreateLogHandler {
	return CreateLogHandle{createLogRepository: createLogRepository}
}

type CreateLogHandle struct {
	createLogRepository repository.CreateLogHandler
}

func (handle CreateLogHandle) CreateLog(ctx context.Context, logEntry entitymodel.LogEntry) error {
	return handle.createLogRepository.CreateLog(ctx, logEntry)
}

type ReadLogHandler interface {
	ReadLogByID(ctx context.Context, id string) (entitymodel.LogEntry, error)
	ReadLogByIDs(ctx context.Context, ids []string) (entitymodel.LogEntrySlice, error)
	ReadAllLogs(ctx context.Context) (entitymodel.LogEntrySlice, error)
	ReadLogsByService(ctx context.Context, service enum.Service) (entitymodel.LogEntrySlice, error)
	ReadLogsByServiceAndName(ctx context.Context, service enum.Service, name string) (entitymodel.LogEntrySlice, error)
	ReadLogsSearch(ctx context.Context, service enum.Service, filter string) (entitymodel.LogEntrySlice, error)
}

func NewReadLogService(readLogRepository repository.ReadLogHandler) ReadLogHandler {
	return ReadLogHandle{readLogRepository: readLogRepository}
}

type ReadLogHandle struct {
	readLogRepository repository.ReadLogHandler
}

func (handle ReadLogHandle) ReadLogByID(ctx context.Context, id string) (entitymodel.LogEntry, error) {
	return handle.readLogRepository.ReadLogByID(ctx, id)
}

func (handle ReadLogHandle) ReadLogByIDs(ctx context.Context, ids []string) (entitymodel.LogEntrySlice, error) {
	return handle.readLogRepository.ReadLogByIDs(ctx, ids)
}

func (handle ReadLogHandle) ReadAllLogs(ctx context.Context) (entitymodel.LogEntrySlice, error) {
	return handle.readLogRepository.ReadAllLogs(ctx)
}

func (handle ReadLogHandle) ReadLogsByService(ctx context.Context, service enum.Service) (entitymodel.LogEntrySlice, error) {
	return handle.readLogRepository.ReadLogsByService(ctx, service)
}

func (handle ReadLogHandle) ReadLogsByServiceAndName(ctx context.Context, service enum.Service, name string) (entitymodel.LogEntrySlice, error) {
	return handle.readLogRepository.ReadLogsByServiceAndName(ctx, service, name)
}

func (handle ReadLogHandle) ReadLogsSearch(ctx context.Context, service enum.Service, filter string) (entitymodel.LogEntrySlice, error) {
	return handle.readLogRepository.ReadLogsSearch(ctx, service, filter)
}

type UpdateLogHandler interface {
	UpdateLog(ctx context.Context, logEntry entitymodel.LogEntry) error
}

func NewUpdateLogService(updateLogRepository repository.UpdateLogHandler) UpdateLogHandler {
	return UpdateLogHandle{updateLogRepository: updateLogRepository}
}

type UpdateLogHandle struct {
	updateLogRepository repository.UpdateLogHandler
}

func (handle UpdateLogHandle) UpdateLog(ctx context.Context, logEntry entitymodel.LogEntry) error {
	return handle.updateLogRepository.UpdateLog(ctx, logEntry)
}

type DeleteLogHandler interface {
	DeleteLogByID(ctx context.Context, id string) error
}

func NewDeleteLogService(deleteLogRepository repository.DeleteLogHandler) DeleteLogHandler {
	return DeleteLogHandle{deleteLogRepository: deleteLogRepository}
}

type DeleteLogHandle struct {
	deleteLogRepository repository.DeleteLogHandler
}

func (handle DeleteLogHandle) DeleteLogByID(ctx context.Context, id string) error {
	return handle.deleteLogRepository.DeleteLogByID(ctx, id)
}
