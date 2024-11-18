package controller

import (
	"google-map/config"
	"google-map/errors"
	"google-map/model/requestmodel"
	"google-map/service"
	"google-map/utility"

	"github.com/gin-gonic/gin"
)

type GooglePlaceSearchHandle struct {
	googlePlaceSearchService service.GooglePlaceSearchHandler
	config                   config.Config
}

func NewGooglePlaceSearchHandle(
	googlePlaceSearchService service.GooglePlaceSearchHandler,
	config config.Config,
) *GooglePlaceSearchHandle {
	return &GooglePlaceSearchHandle{
		googlePlaceSearchService: service.NewGooglePlaceSearchHandler(),
		config:                   config,
	}
}

// @Summary クイック検索
// @Description クイック検索
// @Tags google-map-search
// @Accept json
// @Param place_id query string false "キーワード"
// @Param api_mask query string false "APIマスク"
// @Param language_code query string false "言語"
// @Produce  json
// @Success 200 {object} responsemodel.CommonResponse
// @Router /quickSearch [get]
func (handle *GooglePlaceSearchHandle) QuickSearch(c *gin.Context) {
	var (
		req requestmodel.QuickSearchRequest
	)

	if err := c.ShouldBindQuery(&req); err != nil {
		utility.CommonErrorResponse(c, errors.HTTPStatusBadRequest, nil)
		return
	}

	googlePlaceSearch, err := handle.googlePlaceSearchService.QuickSearch(c, req.ToEntity(), handle.config)
	if err != nil {
		utility.CommonErrorResponse(c, err, nil)
		return
	}

	utility.CommonResponse(c, googlePlaceSearch)
}

// @Summary 詳細検索
// @Description 詳細検索
// @Tags google-map-search
// @Accept json
// @Param text_query query string false "テキストクエリ"
// @Param low_latitude query float64 false "低緯度"
// @Param low_longitude query float64 false "低経度"
// @Param high_latitude query float64 false "高緯度"
// @Param high_longitude query float64 false "高経度"
// @Param max_result_count query int false "最大取得件数"
// @Param min_rating query int false "最低評価"
// @Param open_now query bool false "現在開いているか"
// @Param rank_preference query string false "ランク優先度"
// @Param language_code query string false "言語"
// @Param api_mask query string false "APIマスク"
// @Produce  json
// @Success 200 {object} responsemodel.CommonResponse
// @Router /advanceSearch [get]
func (handle *GooglePlaceSearchHandle) AdvanceSearch(c *gin.Context) {
	var (
		req requestmodel.AdvanceSearchRequest
	)

	if err := c.ShouldBindQuery(&req); err != nil {
		utility.CommonErrorResponse(c, errors.HTTPStatusBadRequest, nil)
		return
	}

	googlePlaceSearch, err := handle.googlePlaceSearchService.AdvanceSearch(c, req.ToEntity(), handle.config)
	if err != nil {
		utility.CommonErrorResponse(c, err, nil)
		return
	}

	utility.CommonResponse(c, googlePlaceSearch)
}
