package controller

import (
	"authenticate/convert/HTTPResponse"
	"authenticate/errors"
	"authenticate/middleware"
	"authenticate/model/requestmodel"
	"authenticate/service"
	"authenticate/utility"

	"github.com/gin-gonic/gin"
)

type SavePlaceControllerHandle struct {
	savePlaceServiceHandler service.SavePlaceServiceHandler
}

func NewSavePlaceControllerHandle(savePlaceServiceHandler service.SavePlaceServiceHandler) *SavePlaceControllerHandle {
	return &SavePlaceControllerHandle{savePlaceServiceHandler: savePlaceServiceHandler}
}

// @Summary 場所保存
// @Description 場所保存
// @Tags place
// @Accept json
// @Security BearerAuth
// @Param param body requestmodel.SavePlaceRequest true "json"
// @Produce  json
// @Success 200 "success"
// @Router /place/savePlace [post]
func (handle *SavePlaceControllerHandle) SavePlace(c *gin.Context) {
	var req requestmodel.SavePlaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utility.CommonErrorResponse(c, err, nil)
		return
	}
	if _, ok := middleware.GetSession(c); !ok {
		utility.CommonErrorResponse(c, errors.HTTPStatusUnauthorized, nil)
		return
	}
	if err := handle.savePlaceServiceHandler.SavePlace(c.Request.Context(), req.ToEntity()); err != nil {
		utility.CommonErrorResponse(c, err, nil)
		return
	}
	utility.CommonResponse(c, "success")
}

type GetPlaceControllerHandle struct {
	getPlaceServiceHandler service.GetPlaceServiceHandler
}

func NewGetPlaceControllerHandle(getPlaceServiceHandler service.GetPlaceServiceHandler) *GetPlaceControllerHandle {
	return &GetPlaceControllerHandle{getPlaceServiceHandler: getPlaceServiceHandler}
}

// @Summary 場所取得
// @Description 場所取得
// @Tags place
// @Accept json
// @Security BearerAuth
// @Param param query requestmodel.GetPlaceRequest true "json"
// @Produce  json
// @Success 200 "success"
// @Router /place/getPlace [get]
func (handle *GetPlaceControllerHandle) GetPlace(c *gin.Context) {
	var req requestmodel.GetPlaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utility.CommonErrorResponse(c, err, nil)
		return
	}
	_, ok := middleware.GetSession(c)
	if !ok {
		utility.CommonErrorResponse(c, errors.HTTPStatusUnauthorized, nil)
		return
	}
	place, err := handle.getPlaceServiceHandler.GetPlace(c.Request.Context(), req.GoogleID)
	if err != nil {
		utility.CommonErrorResponse(c, err, nil)
		return
	}
	utility.CommonResponse(c, HTTPResponse.PlaceRef(place).ToResponse())
}
