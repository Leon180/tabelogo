package inject

import "github.com/google/wire"

var serviceHandleSet = wire.NewSet(
	provideGetTabelogInfoService,
	provideGetTabelogPhotoService,
)
