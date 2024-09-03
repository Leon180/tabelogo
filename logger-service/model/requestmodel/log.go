package requestmodel

import (
	"logger-service/model/entitymodel"
	"logger-service/model/enum"
	"logger-service/utility"
	"time"
)

type LogEntryRequest struct {
	Name    string `json:"name"`
	Service string `json:"service"`
	Data    string `json:"data"`
}

func (req LogEntryRequest) ToEntity() entitymodel.LogEntry {
	return entitymodel.LogEntry{
		ID:        utility.GenDefaultUUID(),
		Name:      req.Name,
		Service:   enum.Service(req.Service),
		Data:      []byte(req.Data),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

type SearchLogRequest struct {
	Service enum.Service `json:"service"`
	Filter  string       `json:"filter"`
}
