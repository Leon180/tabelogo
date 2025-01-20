package grpcservice

import (
	"context"
	"google-map/config"
	"google-map/service"

	convertrequest "google-map/grpc/convert/requestmodel"
	convertresponse "google-map/grpc/convert/responsemodel"
	"google-map/grpc/proto"
)

func NewGoogleMapServiceServer(
	googlePlaceSearchService service.GooglePlaceSearchHandler,
	config *config.Config,
) *GoogleMapServiceServer {
	return &GoogleMapServiceServer{
		googlePlaceSearchService: googlePlaceSearchService,
		config:                   config,
	}
}

type GoogleMapServiceServer struct {
	proto.UnimplementedGoogleMapServiceServer
	googlePlaceSearchService service.GooglePlaceSearchHandler
	config                   *config.Config
}

func (s *GoogleMapServiceServer) QuickSearch(ctx context.Context, req *proto.QuickSearchRequest) (*proto.InterfaceResponse, error) {
	requestModel := convertrequest.ConvertQuickSearchRequest(req)
	res, err := s.googlePlaceSearchService.QuickSearch(ctx, requestModel.ToEntity())
	if err != nil {
		return nil, err
	}
	return convertresponse.ConvertInterfaceResponse(res), nil
}

func (s *GoogleMapServiceServer) AdvanceSearch(ctx context.Context, req *proto.AdvanceSearchRequest) (*proto.InterfaceResponse, error) {
	requestModel := convertrequest.ConvertAdvanceSearchRequest(req)
	res, err := s.googlePlaceSearchService.AdvanceSearch(ctx, requestModel.ToEntity())
	if err != nil {
		return nil, err
	}
	return convertresponse.ConvertInterfaceResponse(res), nil
}
