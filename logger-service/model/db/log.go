package db

import (
	"time"

	"logger-service/model/entitymodel"
	"logger-service/model/enum"

	"github.com/samber/lo"
)

type LogEntry struct {
	ID        string       `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string       `bson:"name" json:"name"`
	Service   enum.Service `bson:"service" json:"service"`
	Data      string       `bson:"data" json:"data"`
	CreatedAt time.Time    `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time    `bson:"updated_at" json:"updated_at"`
}

func (entry LogEntry) ToEntity() entitymodel.LogEntry {
	return entitymodel.LogEntry{
		ID:        entry.ID,
		Name:      entry.Name,
		Service:   entry.Service,
		Data:      []byte(entry.Data),
		CreatedAt: entry.CreatedAt,
		UpdatedAt: entry.UpdatedAt,
	}
}

type LogEntrySlice []LogEntry

func (slice LogEntrySlice) ToEntitySlice() entitymodel.LogEntrySlice {
	return lo.Map(slice, func(logEntry LogEntry, _ int) entitymodel.LogEntry {
		return logEntry.ToEntity()
	})
}

type LogEntryEntity entitymodel.LogEntry

func (entity LogEntryEntity) ToDBModel() LogEntry {
	return LogEntry{
		ID:        entity.ID,
		Name:      entity.Name,
		Service:   entity.Service,
		Data:      string(entity.Data),
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
