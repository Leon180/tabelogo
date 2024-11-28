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
	QuickSearch(ctx context.Context, param entitymodel.QuickSearchRequest) (interface{}, error)
	AdvanceSearch(ctx context.Context, param entitymodel.AdvanceSearchRequest) (interface{}, error)
}

func NewGooglePlaceSearchHandler(config *config.Config) GooglePlaceSearchHandler {
	return &GooglePlaceSearchHandle{
		config: config,
	}
}

type GooglePlaceSearchHandle struct {
	config *config.Config
}

func (handle *GooglePlaceSearchHandle) QuickSearch(ctx context.Context, param entitymodel.QuickSearchRequest) (interface{}, error) {
	resp, err := utility.RequestToAnotherService(
		ctx,
		enum.RequestMethodGET,
		enum.GoogleMapPlaceV1URL.AddSubDirectory(param.PlaceID).AddParams(map[enum.Param]string{
			enum.APIMask:      param.APIMask,
			enum.Key:          handle.config.GoogleMapAPIKey,
			enum.LanguageCode: param.LanguageCode,
		}),
		nil,
		"",
	)
	if err != nil {
		utility.SugarLogger.Error("error during request to google map quick search, error: %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	var googleRsp interface{}
	if err := json.NewDecoder(resp.Body).Decode(&googleRsp); err != nil {
		utility.SugarLogger.Error("error during decode google map quick search response, error: %s", err)
		return nil, err
	}

	return googleRsp, nil
}

func (handle *GooglePlaceSearchHandle) AdvanceSearch(ctx context.Context, param entitymodel.AdvanceSearchRequest) (interface{}, error) {
	resp, err := utility.RequestToAnotherService(
		ctx,
		enum.RequestMethodPOST,
		enum.GoogleMapPlaceSearchTextURL,
		map[enum.RequestHeader]string{
			enum.RequestHeaderContentType:    enum.ContextTypeJSON.ToString(),
			enum.RequestHeaderXGoogAPIKey:    handle.config.GoogleMapAPIKey,
			enum.RequestHeaderXGoogFieldMask: param.APIMask,
		},
		param.ToRequestBody(),
	)
	if err != nil {
		utility.SugarLogger.Error("error during request to google map advance search, error: %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	var googleRsp interface{}
	if err := json.NewDecoder(resp.Body).Decode(&googleRsp); err != nil {
		utility.SugarLogger.Error("error during decode google map advance search response, error: %s", err)
		return nil, err
	}

	return googleRsp, nil
}
