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

// @Summary ログ作成
// @Description ログ作成
// @Tags logger
// @Accept json
// @Param logEntry body requestmodel.LogEntryRequest true "ログ作成リクエスト"
// @Produce json
// @Success 200 {object} responsemodel.CommonResponse
// @Router /createLog [post]
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

// @Summary ログ検索
// @Description ログ検索
// @Tags logger
// @Accept json
// @Param searchLogRequest body requestmodel.SearchLogRequest true "ログ検索リクエスト"
// @Produce json
// @Success 200 {object} responsemodel.CommonResponse
// @Router /searchLogs [post]
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

// @Summary 全ログ取得
// @Description 全ログ取得
// @Tags logger
// @Accept json
// @Produce json
// @Success 200 {object} responsemodel.CommonResponse
// @Router /readAllLogs [get]
func (controller *LogController) ReadAllLogs(c *gin.Context) {
	logs, err := controller.readLogService.ReadAllLogs(c.Request.Context())
	if err != nil {
		utility.CommonErrorResponse(c, err, nil)
		return
	}
	utility.CommonResponse(c, logs)
}
