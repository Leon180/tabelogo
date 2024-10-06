package inject

import (
	"github.com/google/wire"
)

var repositoryHandleSet = wire.NewSet(
// providePostgresqlRepository,
)

// func providePostgresqlRepository(settings config.Config, logger *zap.Logger) repository.DBRepositoryHandle {
// 	return postgresqlrep.NewPostgresqlRepoHandler(settings, logger)
// }
