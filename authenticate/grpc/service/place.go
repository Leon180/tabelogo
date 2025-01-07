package grpcservice

import (
	"authenticate/grpc/proto"
	"authenticate/service"
	"context"

	convertrequest "authenticate/grpc/convert/requestmodel"
	convertresponse "authenticate/grpc/convert/responsemodel"
)

func NewPlaceServiceServer(
	savePlaceServiceHandler service.SavePlaceServiceHandler,
	getPlaceServiceHandler service.GetPlaceServiceHandler,
) *PlaceServiceServer {
	return &PlaceServiceServer{
		savePlaceServiceHandler: savePlaceServiceHandler,
		getPlaceServiceHandler:  getPlaceServiceHandler,
	}
}

type PlaceServiceServer struct {
	proto.UnimplementedPlaceServiceServer
	savePlaceServiceHandler service.SavePlaceServiceHandler
	getPlaceServiceHandler  service.GetPlaceServiceHandler
}

func (s *PlaceServiceServer) GetPlace(ctx context.Context, req *proto.GetPlaceRequest) (*proto.Place, error) {
	request := convertrequest.ConvertGetPlaceRequest(req)
	place, err := s.getPlaceServiceHandler.GetPlace(ctx, request.GoogleID)
	if err != nil {
		return nil, err
	}
	return convertresponse.PlaceEntityModel(place).ToPlaceProto(), nil
}

func (s *PlaceServiceServer) SavePlace(ctx context.Context, req *proto.SavePlaceRequest) (*proto.CommonResponse, error) {
	request := convertrequest.ConvertSavePlaceRequest(req)
	if err := s.savePlaceServiceHandler.SavePlace(ctx, request.ToEntity()); err != nil {
		return nil, err
	}
	return &proto.CommonResponse{Message: "success"}, nil
}
