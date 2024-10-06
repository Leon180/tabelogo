package inject

import (
	"github.com/google/wire"
)

var controllerHandleSet = wire.NewSet(
// provideHealthCheckControllerHandleController,
// provideTestControllerHandleController,
)

// func provideTestControllerHandleController(
// 	reportService service.ReportServiceHandler,
// ) controller.TestControllerHandle {
// 	return controller.NewTestControllerHandle(reportService)
// }
