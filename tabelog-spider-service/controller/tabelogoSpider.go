package controller

import (
	"tabelog-spider/convert"
	"tabelog-spider/errors"
	"tabelog-spider/model/entitymodel"
	"tabelog-spider/model/requestmodel"
	"tabelog-spider/service"
	"tabelog-spider/utility"

	"github.com/gin-gonic/gin"
)

type GetTabelogInfoHandle struct {
	getTabelogInfoService  service.GetTabelogInfoHandler
	getTabelogPhotoService service.GetTabelogPhotoHandler
}

func NewGetTabelogInfoHandle(
	getTabelogInfoService service.GetTabelogInfoHandler,
	getTabelogPhotoService service.GetTabelogPhotoHandler,
) *GetTabelogInfoHandle {
	return &GetTabelogInfoHandle{
		getTabelogInfoService:  getTabelogInfoService,
		getTabelogPhotoService: getTabelogPhotoService,
	}
}

func (handle *GetTabelogInfoHandle) GetTabelogInfo(c *gin.Context) {
	var (
		req requestmodel.GetTabelogInfoRequest
		err error
	)
	if err = c.ShouldBindQuery(&req); err != nil {
		utility.CommonErrorResponse(c, errors.HTTPStatusBadRequest, nil)
	}
	tablogoInfo, err := handle.getTabelogInfoService.GetTabelogInfo(c, entitymodel.GetTabelogInfoParam{
		Area:          req.Area,
		PlaceName:     req.PlaceName,
		MaxLinkAmount: req.MaxResultAmount,
	})
	if err != nil {
		utility.CommonErrorResponse(c, err, nil)
	}
	utility.CommonResponse(c, convert.GetTabelogInfo(tablogoInfo).ToResponse())
}

func (handle *GetTabelogInfoHandle) GetTabelogPhoto(c *gin.Context) {
	var (
		req requestmodel.GetTabelogPhotoRequest
		err error
	)
	if err = c.ShouldBindQuery(&req); err != nil {
		utility.CommonErrorResponse(c, errors.HTTPStatusBadRequest, nil)
	}
	tablogoPhoto, err := handle.getTabelogPhotoService.GetTabelogPhoto(c, req.Link)
	if err != nil {
		utility.CommonErrorResponse(c, err, nil)
	}
	utility.CommonResponse(c, convert.GetTabelogPhoto(tablogoPhoto).ToResponse())
}
