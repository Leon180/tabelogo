package responsemodel

import (
	"logger-service/model/enum"
	"time"
)

type LogEntryResponse struct {
	ID        string
	Name      string
	Service   enum.Service
	Data      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type LogEntrySliceResponse []LogEntryResponse
