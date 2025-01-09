package requestmodel

import (
	"logger-service/grpc/proto"
	"logger-service/model/enum"
	"logger-service/model/requestmodel"
)

func ConvertCreateLogRequest(req *proto.CreateLogRequest) *requestmodel.LogEntryRequest {
	return &requestmodel.LogEntryRequest{
		Name:    req.Name,
		Service: enum.Service(req.Service),
		Data:    req.Data,
	}
}

func ConvertSearchLogRequest(req *proto.SearchLogRequest) *requestmodel.SearchLogRequest {
	return &requestmodel.SearchLogRequest{
		Service: enum.Service(req.Service),
		Filter:  req.Filter,
	}
}
