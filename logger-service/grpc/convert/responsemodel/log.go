package responsemodel

import (
	"logger-service/grpc/proto"
	"logger-service/model/entitymodel"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type LogEntryModel entitymodel.LogEntry

func (entity LogEntryModel) ToLogProto() *proto.LogEntryResponse {
	return &proto.LogEntryResponse{
		Id:        entity.ID,
		Name:      entity.Name,
		Service:   string(entity.Service),
		Data:      string(entity.Data),
		CreatedAt: timestamppb.New(entity.CreatedAt),
		UpdatedAt: timestamppb.New(entity.UpdatedAt),
	}
}

type LogEntitySliceModel []entitymodel.LogEntry

func (entity LogEntitySliceModel) ToLogSliceProto() *proto.LogEntrySliceResponse {
	return &proto.LogEntrySliceResponse{
		LogEntries: lo.Map(entity, func(log entitymodel.LogEntry, _ int) *proto.LogEntryResponse {
			return LogEntryModel(log).ToLogProto()
		}),
	}
}
