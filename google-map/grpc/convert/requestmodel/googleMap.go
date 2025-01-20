package requestmodel

import (
	"google-map/grpc/proto"
	"google-map/model/requestmodel"
)

func ConvertQuickSearchRequest(req *proto.QuickSearchRequest) requestmodel.QuickSearchRequest {
	return requestmodel.QuickSearchRequest{
		PlaceID:      req.PlaceId,
		APIMask:      req.ApiMask,
		LanguageCode: req.LanguageCode,
	}
}

func ConvertAdvanceSearchRequest(req *proto.AdvanceSearchRequest) requestmodel.AdvanceSearchRequest {
	return requestmodel.AdvanceSearchRequest{
		TextQuery:      req.TextQuery,
		LowLatitude:    float64(req.LowLatitude),
		LowLongitude:   float64(req.LowLongitude),
		HighLatitude:   float64(req.HighLatitude),
		HighLongitude:  float64(req.HighLongitude),
		MaxResultCount: int(req.MaxResultCount),
		MinRating:      int(req.MinRating),
		OpenNow:        req.OpenNow,
		RankPreference: req.RankPreference,
		LanguageCode:   req.LanguageCode,
		APIMask:        req.ApiMask,
	}
}
