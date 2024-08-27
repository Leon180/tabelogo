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

type GetTabelogInfoHandle struct{}

func NewGetTabelogInfoHandle() GetTabelogInfoHandle {
	return GetTabelogInfoHandle{}
}

func (handle GetTabelogInfoHandle) GetTabelogInfo(c *gin.Context) {
	var (
		req requestmodel.GetTabelogInfoRequest
		err error
	)
	if err = c.ShouldBindQuery(&req); err != nil {
		utility.CommonErrorResponse(c, errors.HTTPStatusBadRequest, nil)
	}
	tablogoInfo, err := service.NewGetTabelogInfoHandler().GetTabelogInfo(c, entitymodel.GetTabelogInfoParam{
		Area:          req.Area,
		PlaceName:     req.PlaceName,
		MaxLinkAmount: req.MaxResultAmount,
	})
	if err != nil {
		utility.CommonErrorResponse(c, err, nil)
	}
	utility.CommonResponse(c, convert.GetTabelogInfo(tablogoInfo).ToResponse())
}

func (handle GetTabelogInfoHandle) GetTabelogPhoto(c *gin.Context) {
	var (
		req requestmodel.GetTabelogPhotoRequest
		err error
	)
	if err = c.ShouldBindQuery(&req); err != nil {
		utility.CommonErrorResponse(c, errors.HTTPStatusBadRequest, nil)
	}
	tablogoPhoto, err := service.NewGetTabelogPhotoHandler().GetTabelogPhoto(c, req.Link)
	if err != nil {
		utility.CommonErrorResponse(c, err, nil)
	}
	utility.CommonResponse(c, convert.GetTabelogPhoto(tablogoPhoto).ToResponse())
}
