package utility

import (
	"net/http"
	"tabelog-spider/errors"
	"tabelog-spider/model/enum"
	"tabelog-spider/model/responsemodel"

	"github.com/gin-gonic/gin"
)

func CommonErrorResponse(c *gin.Context, err error, result interface{}) {
	if eventIDKey, ok := c.Get(enum.MiddleWareEventIDKey); ok {
		SugarLogger.Warnf("Http Response EventID: %s, Error: %v", eventIDKey.(string), err)
	}
	if apiErr, ok := err.(errors.APIErr); ok {
		c.AbortWithStatusJSON(apiErr.GetHttpStatus(), responsemodel.CommonErrorResponse{
			ErrorCode:    apiErr.GetErrorCode(),
			ErrorMessage: apiErr.Error(),
			Result:       result,
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, responsemodel.CommonErrorResponse{
		ErrorCode:    http.StatusInternalServerError,
		ErrorMessage: err.Error(),
		Result:       result,
	})
}

func CommonResponse(c *gin.Context, result interface{}) {
	if eventIDKey, ok := c.Get(enum.MiddleWareEventIDKey); ok {
		SugarLogger.Infof("Http Response EventID: %s, Payload: %+v", eventIDKey.(string), result)
	}
	c.AbortWithStatusJSON(http.StatusOK, responsemodel.CommonResponse{Result: result})
}
