package inject

import "github.com/google/wire"

var grpcServerSet = wire.NewSet(
	provideLogServiceServer,
)
