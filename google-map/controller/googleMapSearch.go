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
