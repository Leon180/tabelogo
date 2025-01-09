package repository

import (
	"context"
	"logger-service/model/db"
	"logger-service/model/entitymodel"
	"logger-service/model/enum"
	"logger-service/utility"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CreateLogHandler interface {
	CreateLog(ctx context.Context, logEntry entitymodel.LogEntry) error
}

func NewCreateLogRepository(db *mongo.Client) CreateLogHandler {
	return &CreateLogHandle{db: db}
}

type CreateLogHandle struct {
	db *mongo.Client
}

func (handle *CreateLogHandle) CreateLog(ctx context.Context, logEntry entitymodel.LogEntry) error {
	_, err := handle.db.Database("logs").
		Collection("logs").
		InsertOne(ctx, db.LogEntryEntity(logEntry).ToDBModel())
	if err != nil {
		utility.LogWithTraceID(ctx, "error inserting log entry", err)
		return err
	}
	return nil
}

type ReadLogHandler interface {
	ReadLogByID(ctx context.Context, id string) (entitymodel.LogEntry, error)
	ReadLogByIDs(ctx context.Context, ids []string) (entitymodel.LogEntrySlice, error)
	ReadAllLogs(ctx context.Context) (entitymodel.LogEntrySlice, error)
	ReadLogsByService(ctx context.Context, service enum.Service) (entitymodel.LogEntrySlice, error)
	ReadLogsByServiceAndName(ctx context.Context, service enum.Service, name string) (entitymodel.LogEntrySlice, error)
	ReadLogsSearch(ctx context.Context, service enum.Service, filter string) (entitymodel.LogEntrySlice, error)
}

func NewReadLogRepository(db *mongo.Client) ReadLogHandler {
	return &ReadLogHandle{db: db}
}

type ReadLogHandle struct {
	db *mongo.Client
}

func (handle *ReadLogHandle) ReadLogByID(ctx context.Context, id string) (entitymodel.LogEntry, error) {
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utility.LogWithTraceID(ctx, "error converting string to objectID", err)
		return entitymodel.LogEntry{}, err
	}

	var logEntry db.LogEntry
	if err := handle.db.Database("logs").
		Collection("logs").
		FindOne(ctx, bson.M{"_id": docID}).
		Decode(&logEntry); err != nil {
		utility.LogWithTraceID(ctx, "error getting log entry", err)
		return entitymodel.LogEntry{}, err
	}

	return logEntry.ToEntity(), nil
}

func (handle *ReadLogHandle) ReadLogByIDs(ctx context.Context, ids []string) (entitymodel.LogEntrySlice, error) {
	docIDs := make([]primitive.ObjectID, len(ids))
	for i, id := range ids {
		docID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			utility.LogWithTraceID(ctx, "error converting string to objectID", err)
			return []entitymodel.LogEntry{}, err
		}
		docIDs[i] = docID
	}

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := handle.db.Database("logs").
		Collection("logs").
		Find(ctx, bson.M{"_id": bson.M{"$in": docIDs}}, opts)
	if err != nil {
		utility.LogWithTraceID(ctx, "error getting log entries", err)
		return []entitymodel.LogEntry{}, err
	}

	var logEntries db.LogEntrySlice
	if err := cursor.All(ctx, &logEntries); err != nil {
		utility.LogWithTraceID(ctx, "error decoding log entries", err)
		return []entitymodel.LogEntry{}, err
	}

	return logEntries.ToEntitySlice(), nil
}

func (handle *ReadLogHandle) ReadAllLogs(ctx context.Context) (entitymodel.LogEntrySlice, error) {
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := handle.db.Database("logs").
		Collection("logs").
		Find(ctx, bson.D{}, opts)
	if err != nil {
		utility.LogWithTraceID(ctx, "error getting all log entries", err)
		return entitymodel.LogEntrySlice{}, err
	}

	var logEntries db.LogEntrySlice
	if err := cursor.All(ctx, &logEntries); err != nil {
		utility.LogWithTraceID(ctx, "error decoding log entries", err)
		return entitymodel.LogEntrySlice{}, err
	}

	return logEntries.ToEntitySlice(), nil
}

func (handle *ReadLogHandle) ReadLogsByService(ctx context.Context, service enum.Service) (entitymodel.LogEntrySlice, error) {
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := handle.db.Database("logs").
		Collection("logs").
		Find(ctx, bson.M{"service": service}, opts)
	if err != nil {
		utility.LogWithTraceID(ctx, "error getting log entries by service", err)
		return entitymodel.LogEntrySlice{}, err
	}

	var logEntries db.LogEntrySlice
	if err := cursor.All(ctx, &logEntries); err != nil {
		utility.LogWithTraceID(ctx, "error decoding log entries", err)
		return entitymodel.LogEntrySlice{}, err
	}

	return logEntries.ToEntitySlice(), nil
}

func (handle *ReadLogHandle) ReadLogsByServiceAndName(ctx context.Context, service enum.Service, name string) (entitymodel.LogEntrySlice, error) {
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := handle.db.Database("logs").
		Collection("logs").
		Find(ctx, bson.M{"service": service, "name": name}, opts)
	if err != nil {
		utility.LogWithTraceID(ctx, "error getting log entries by service and name", err)
		return entitymodel.LogEntrySlice{}, err
	}

	var logEntries db.LogEntrySlice
	if err := cursor.All(ctx, &logEntries); err != nil {
		utility.LogWithTraceID(ctx, "error decoding log entries", err)
		return entitymodel.LogEntrySlice{}, err
	}

	return logEntries.ToEntitySlice(), nil
}

func (handle *ReadLogHandle) ReadLogsSearch(ctx context.Context, service enum.Service, filter string) (entitymodel.LogEntrySlice, error) {
	entity, err := handle.ReadLogsByService(ctx, service)
	if err != nil {
		return entitymodel.LogEntrySlice{}, err
	}
	res := entitymodel.LogEntrySlice{}
	for _, log := range entity {
		if strings.Contains(string(log.Data), filter) {
			res = append(res, log)
			continue
		}
		if strings.Contains(log.Name, filter) {
			res = append(res, log)
			continue
		}
		if strings.Contains(log.Service.ToString(), filter) {
			res = append(res, log)
			continue
		}
	}
	return res, nil
}

type UpdateLogHandler interface {
	UpdateLog(ctx context.Context, logEntry entitymodel.LogEntry) error
}

func NewUpdateLogRepository(db *mongo.Client) UpdateLogHandler {
	return &UpdateLogHandle{db: db}
}

type UpdateLogHandle struct {
	db *mongo.Client
}

func (handle *UpdateLogHandle) UpdateLog(ctx context.Context, logEntry entitymodel.LogEntry) error {
	docID, err := primitive.ObjectIDFromHex(logEntry.ID)
	if err != nil {
		utility.LogWithTraceID(ctx, "error converting string to objectID", err)
		return err
	}

	_, err = handle.db.Database("logs").
		Collection("logs").
		UpdateOne(ctx, bson.M{"_id": docID}, bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "name", Value: logEntry.Name},
				{Key: "data", Value: logEntry.Data},
				{Key: "updated_at", Value: time.Now()},
			}}})
	if err != nil {
		utility.LogWithTraceID(ctx, "error updating log entry", err)
		return err
	}

	return nil
}

type DeleteLogHandler interface {
	DeleteLogByID(ctx context.Context, id string) error
}

func NewDeleteLogRepository(db *mongo.Client) DeleteLogHandler {
	return &DeleteLogHandle{db: db}
}

type DeleteLogHandle struct {
	db *mongo.Client
}

func (handle *DeleteLogHandle) DeleteLogByID(ctx context.Context, id string) error {
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utility.LogWithTraceID(ctx, "error converting string to objectID", err)
		return err
	}

	_, err = handle.db.Database("logs").
		Collection("logs").
		DeleteOne(ctx, bson.M{"_id": docID})
	if err != nil {
		utility.LogWithTraceID(ctx, "error deleting log entry", err)
		return err
	}

	return nil
}
