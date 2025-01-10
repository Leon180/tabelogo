package grpcservice

import (
	"context"
	convertrequest "tabelog-spider/grpc/convert/requestmodel"
	convertresponse "tabelog-spider/grpc/convert/responsemodel"
	"tabelog-spider/grpc/proto"
	"tabelog-spider/service"
)

func NewTabelogoSpiderServiceServer(
	getTabelogInfoService service.GetTabelogInfoHandler,
	getTabelogPhotoService service.GetTabelogPhotoHandler,
) *TabelogoSpiderServiceServer {
	return &TabelogoSpiderServiceServer{
		getTabelogInfoService:  getTabelogInfoService,
		getTabelogPhotoService: getTabelogPhotoService,
	}
}

type TabelogoSpiderServiceServer struct {
	proto.UnimplementedTabelogoSpiderServiceServer
	getTabelogInfoService  service.GetTabelogInfoHandler
	getTabelogPhotoService service.GetTabelogPhotoHandler
}

func (s *TabelogoSpiderServiceServer) GetTabelogInfo(ctx context.Context, req *proto.GetTabelogInfoRequest) (*proto.TabelogInfoResponseList, error) {
	requestModel := convertrequest.ConvertGetTabelogInfoRequest(req)
	tablogoInfo, err := s.getTabelogInfoService.GetTabelogInfo(ctx, requestModel.ToEntity())
	if err != nil {
		return nil, err
	}
	return convertresponse.TabelogElementInfoSlice(tablogoInfo).TabelogInfoResponseListProto(), nil
}

func (s *TabelogoSpiderServiceServer) GetTabelogPhoto(ctx context.Context, req *proto.GetTabelogPhotoRequest) (*proto.TabelogPhotoResponse, error) {
	requestModel := convertrequest.ConvertGetTabelogPhotoRequest(req)
	tablogoPhoto, err := s.getTabelogPhotoService.GetTabelogPhoto(ctx, requestModel.Link)
	if err != nil {
		return nil, err
	}
	return convertresponse.GetTabelogPhoto(tablogoPhoto).ToTabelogPhotoResponseProto(), nil
}
