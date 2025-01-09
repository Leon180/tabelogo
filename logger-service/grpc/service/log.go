package grpcservice

import (
	"context"
	convertrequest "logger-service/grpc/convert/requestmodel"
	convertresponse "logger-service/grpc/convert/responsemodel"
	"logger-service/grpc/proto"
	"logger-service/service"
)

func NewLogServiceServer(
	createLogService service.CreateLogHandler,
	readLogService service.ReadLogHandler,
	updateLogService service.UpdateLogHandler,
	deleteLogService service.DeleteLogHandler,
) *LogServiceServer {
	return &LogServiceServer{
		createLogServiceHandler: createLogService,
		readLogServiceHandler:   readLogService,
		updateLogServiceHandler: updateLogService,
		deleteLogServiceHandler: deleteLogService,
	}
}

type LogServiceServer struct {
	proto.UnimplementedLogServiceServer
	createLogServiceHandler service.CreateLogHandler
	readLogServiceHandler   service.ReadLogHandler
	updateLogServiceHandler service.UpdateLogHandler
	deleteLogServiceHandler service.DeleteLogHandler
}

func (s *LogServiceServer) CreateLog(ctx context.Context, req *proto.CreateLogRequest) (*proto.CommonResponse, error) {
	request := convertrequest.ConvertCreateLogRequest(req)
	if err := s.createLogServiceHandler.CreateLog(ctx, request.ToEntity()); err != nil {
		return nil, err
	}
	return &proto.CommonResponse{Message: "success"}, nil
}

func (s *LogServiceServer) SearchLogs(ctx context.Context, req *proto.SearchLogRequest) (*proto.LogEntrySliceResponse, error) {
	request := convertrequest.ConvertSearchLogRequest(req)
	logs, err := s.readLogServiceHandler.ReadLogsByServiceAndName(ctx, request.Service, request.Filter)
	if err != nil {
		return nil, err
	}
	return convertresponse.LogEntitySliceModel(logs).ToLogSliceProto(), nil
}

func (s *LogServiceServer) ReadAllLogs(ctx context.Context, req *proto.CommonRequest) (*proto.LogEntrySliceResponse, error) {
	logs, err := s.readLogServiceHandler.ReadAllLogs(ctx)
	if err != nil {
		return nil, err
	}
	return convertresponse.LogEntitySliceModel(logs).ToLogSliceProto(), nil
}
