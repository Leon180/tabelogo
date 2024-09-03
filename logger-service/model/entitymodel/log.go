package entitymodel

import (
	"logger-service/model/enum"
	"time"
)

type LogEntry struct {
	ID        string
	Name      string
	Service   enum.Service
	Data      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

type LogEntrySlice []LogEntry
