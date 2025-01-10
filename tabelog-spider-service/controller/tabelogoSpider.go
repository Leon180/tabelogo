package controller

import (
	"tabelog-spider/convert"
	"tabelog-spider/errors"
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

// @Summary タブログ情報取得
// @Description タブログ情報取得
// @Tags tabelogo-spider
// @Accept json
// @Param area query string false "エリア"
// @Param place_name query string false "店名"
// @Param max_result_amount query int false "最大取得件数"
// @Produce  json
// @Success 200 {object} responsemodel.TabelogInfoResponseList
// @Router /getTabelogInfo [get]
func (handle *GetTabelogInfoHandle) GetTabelogInfo(c *gin.Context) {
	var (
		req requestmodel.GetTabelogInfoRequest
		err error
	)
	if err = c.ShouldBindQuery(&req); err != nil {
		utility.CommonErrorResponse(c, errors.HTTPStatusBadRequest, nil)
		return
	}
	tablogoInfo, err := handle.getTabelogInfoService.GetTabelogInfo(c, req.ToEntity())
	if err != nil {
		utility.CommonErrorResponse(c, err, nil)
		return
	}
	utility.CommonResponse(c, convert.GetTabelogInfo(tablogoInfo).ToResponse())
}

// @Summary タブログ写真取得
// @Description タブログ写真取得
// @Tags tabelogo-spider
// @Accept json
// @Param link query string true "リンク"
// @Produce  json
// @Success 200 {object} responsemodel.TabelogPhotoResponse
// @Router /getTabelogPhoto [get]
func (handle *GetTabelogInfoHandle) GetTabelogPhoto(c *gin.Context) {
	var (
		req requestmodel.GetTabelogPhotoRequest
		err error
	)
	if err = c.ShouldBindQuery(&req); err != nil {
		utility.CommonErrorResponse(c, errors.HTTPStatusBadRequest, nil)
		return
	}
	tablogoPhoto, err := handle.getTabelogPhotoService.GetTabelogPhoto(c, req.Link)
	if err != nil {
		utility.CommonErrorResponse(c, err, nil)
		return
	}
	utility.CommonResponse(c, convert.GetTabelogPhoto(tablogoPhoto).ToResponse())
}
