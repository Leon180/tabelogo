package inject

import "github.com/google/wire"

var serviceHandleSet = wire.NewSet(
// provideHealthCheckService,
// provideReportService,
)

// func provideReportService(
// 	dbRepo repository.DBRepositoryHandler,
// 	reportRepo repository.ReportRepoHandler,
// 	cacheRepo repository.CacheRepoHandler,
// ) service.ReportServiceHandler {
// 	return service.NewReportServiceHandler(dbRepo, reportRepo, cacheRepo)
// }
