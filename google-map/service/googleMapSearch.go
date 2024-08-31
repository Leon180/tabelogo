package service

import (
	"context"
	"encoding/json"
	"google-map/config"
	"google-map/model/entitymodel"
	"google-map/model/enum"
	"google-map/utility"
)

type GooglePlaceSearchHandler interface {
	QuickSearch(ctx context.Context, param entitymodel.QuickSearchRequest, config config.Config) (interface{}, error)
	AdvanceSearch(ctx context.Context, param entitymodel.AdvanceSearchRequest, config config.Config) (interface{}, error)
}

func NewGooglePlaceSearchHandler() GooglePlaceSearchHandler {
	return GooglePlaceSearchHandle{}
}

type GooglePlaceSearchHandle struct{}

func (handle GooglePlaceSearchHandle) QuickSearch(ctx context.Context, param entitymodel.QuickSearchRequest, config config.Config) (interface{}, error) {
	resp, err := utility.RequestToAnotherService(
		enum.RequestMethodGET,
		enum.GoogleMapPlaceV1URL.AddSubDirectory(param.PlaceID).AddParams(map[enum.Param]string{
			enum.APIMask:      param.APIMask,
			enum.Key:          config.GoogleMapAPIKey,
			enum.LanguageCode: param.LanguageCode,
		}),
		nil,
		"",
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var googleRsp interface{}
	json.NewDecoder(resp.Body).Decode(&googleRsp)

	return googleRsp, nil
}

func (handle GooglePlaceSearchHandle) AdvanceSearch(ctx context.Context, param entitymodel.AdvanceSearchRequest, config config.Config) (interface{}, error) {
	resp, err := utility.RequestToAnotherService(
		enum.RequestMethodPOST,
		enum.GoogleMapPlaceSearchTextURL,
		map[enum.RequestHeader]string{
			enum.RequestHeaderContentType:    enum.ContextTypeJSON.ToString(),
			enum.RequestHeaderXGoogAPIKey:    config.GoogleMapAPIKey,
			enum.RequestHeaderXGoogFieldMask: param.APIMask,
		},
		param.ToRequestBody(),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var googleRsp interface{}
	json.NewDecoder(resp.Body).Decode(&googleRsp)

	return googleRsp, nil
}
