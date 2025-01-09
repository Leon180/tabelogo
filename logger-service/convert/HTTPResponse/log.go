package HTTPResponse

import (
	"logger-service/model/entitymodel"
	"logger-service/model/responsemodel"

	"github.com/samber/lo"
)

type LogResponseRef entitymodel.LogEntry

func (ref LogResponseRef) ToResponse() responsemodel.LogEntryResponse {
	return responsemodel.LogEntryResponse{
		ID:        ref.ID,
		Name:      ref.Name,
		Service:   ref.Service,
		Data:      string(ref.Data),
		CreatedAt: ref.CreatedAt,
		UpdatedAt: ref.UpdatedAt,
	}
}

type LogEntrySliceResponseRef entitymodel.LogEntrySlice

func (ref LogEntrySliceResponseRef) ToResponse() responsemodel.LogEntrySliceResponse {
	return lo.Map(ref, func(log entitymodel.LogEntry, _ int) responsemodel.LogEntryResponse {
		return LogResponseRef(log).ToResponse()
	})
}
