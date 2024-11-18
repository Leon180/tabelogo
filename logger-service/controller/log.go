package controller

import (
	"logger-service/errors"
	"logger-service/model/requestmodel"
	"logger-service/service"
	"logger-service/utility"

	"github.com/gin-gonic/gin"
)

type LogController struct {
	createLogService service.CreateLogHandler
	readLogService   service.ReadLogHandler
	updateLogService service.UpdateLogHandler
	deleteLogService service.DeleteLogHandler
}

func NewLogController(
	createLogService service.CreateLogHandler,
	readLogService service.ReadLogHandler,
	updateLogService service.UpdateLogHandler,
	deleteLogService service.DeleteLogHandler,
) *LogController {
	return &LogController{
		createLogService: createLogService,
		readLogService:   readLogService,
		updateLogService: updateLogService,
		deleteLogService: deleteLogService,
	}
}

func (controller *LogController) CreateLog(c *gin.Context) {
	var (
		requestBody requestmodel.LogEntryRequest
	)
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		utility.CommonErrorResponse(c, errors.HTTPStatusBadRequest, nil)
		return
	}
	if err := controller.createLogService.CreateLog(c.Request.Context(), requestBody.ToEntity()); err != nil {
		utility.CommonErrorResponse(c, err, nil)
		return
	}
	utility.CommonResponse(c, "log created successfully")
}

func (controller *LogController) SearchLogs(c *gin.Context) {
	var (
		requestBody requestmodel.SearchLogRequest
	)
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		utility.CommonErrorResponse(c, errors.HTTPStatusBadRequest, nil)
		return
	}

	logs, err := controller.readLogService.ReadLogsByServiceAndName(c.Request.Context(), requestBody.Service, requestBody.Filter)
	if err != nil {
		utility.CommonErrorResponse(c, err, nil)
		return
	}
	utility.CommonResponse(c, logs)
}

func (controller *LogController) ReadAllLogs(c *gin.Context) {
	logs, err := controller.readLogService.ReadAllLogs(c.Request.Context())
	if err != nil {
		utility.CommonErrorResponse(c, err, nil)
		return
	}
	utility.CommonResponse(c, logs)
}
