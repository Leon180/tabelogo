package inject

import "github.com/google/wire"

var serviceHandleSet = wire.NewSet(
	provideCreateLogService,
	provideReadLogService,
	provideUpdateLogService,
	provideDeleteLogService,
)
