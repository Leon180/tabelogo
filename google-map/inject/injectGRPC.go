package inject

import "github.com/google/wire"

var grpcServiceSet = wire.NewSet(
	provideGoogleMapServiceServer,
)
