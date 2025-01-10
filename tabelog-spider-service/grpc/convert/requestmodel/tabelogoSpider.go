package requestmodel

import (
	"tabelog-spider/grpc/proto"
	"tabelog-spider/model/enum"
	"tabelog-spider/model/requestmodel"
)

func ConvertGetTabelogInfoRequest(req *proto.GetTabelogInfoRequest) requestmodel.GetTabelogInfoRequest {
	return requestmodel.GetTabelogInfoRequest{
		Area:      req.Area,
		PlaceName: req.PlaceName,
		MaxResultAmount: func() *int {
			if req.MaxResultAmount == 0 {
				return nil
			}
			amount := int(req.MaxResultAmount)
			return &amount
		}(),
	}
}

func ConvertGetTabelogPhotoRequest(req *proto.GetTabelogPhotoRequest) requestmodel.GetTabelogPhotoRequest {
	return requestmodel.GetTabelogPhotoRequest{
		Link: enum.URL(req.Link),
	}
}
