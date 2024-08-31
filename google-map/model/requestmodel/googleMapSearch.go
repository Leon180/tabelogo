package requestmodel

import "google-map/model/entitymodel"

type QuickSearchRequest struct {
	PlaceID      string `form:"place_id" binding:"required"`
	APIMask      string `form:"api_mask" binding:"required"`
	LanguageCode string `form:"language_code" binding:"required"`
}

func (r QuickSearchRequest) ToEntity() entitymodel.QuickSearchRequest {
	return entitymodel.QuickSearchRequest{
		PlaceID:      r.PlaceID,
		APIMask:      r.APIMask,
		LanguageCode: r.LanguageCode,
	}
}

type AdvanceSearchRequest struct {
	TextQuery string `form:"text_query" binding:"required"`
	// For location bias
	LowLatitude   float64 `form:"low_latitude" binding:"required"`
	LowLongitude  float64 `form:"low_longitude" binding:"required"`
	HighLatitude  float64 `form:"high_latitude" binding:"required"`
	HighLongitude float64 `form:"high_longitude" binding:"required"`
	//
	MaxResultCount int    `form:"max_result_count" binding:"required"`
	MinRating      int    `form:"min_rating" binding:"required"`
	OpenNow        bool   `form:"open_now"`
	RankPreference string `form:"rank_preference" binding:"required"`
	LanguageCode   string `form:"language_code" binding:"required"`
	APIMask        string `form:"api_mask"`
}

func (r AdvanceSearchRequest) ToEntity() entitymodel.AdvanceSearchRequest {
	return entitymodel.AdvanceSearchRequest{
		TextQuery:      r.TextQuery,
		LowLatitude:    r.LowLatitude,
		LowLongitude:   r.LowLongitude,
		HighLatitude:   r.HighLatitude,
		HighLongitude:  r.HighLongitude,
		MaxResultCount: r.MaxResultCount,
		MinRating:      r.MinRating,
		OpenNow:        r.OpenNow,
		RankPreference: r.RankPreference,
		LanguageCode:   r.LanguageCode,
		APIMask:        r.APIMask,
	}
}
