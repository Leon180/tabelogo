package inject

import "github.com/google/wire"

var controllerHandleSet = wire.NewSet(
	provideGetTabelogInfoController,
)
