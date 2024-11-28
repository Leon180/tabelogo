package inject

import "github.com/google/wire"

var repositoryHandleSet = wire.NewSet(
	provideCreateLogRepository,
	provideReadLogRepository,
	provideUpdateLogRepository,
	provideDeleteLogRepository,
)
