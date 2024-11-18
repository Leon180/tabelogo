package inject

import (
	"tabelog-spider/controller"
	"tabelog-spider/service"
)

func provideGetTabelogInfoService() service.GetTabelogInfoHandler {
	return service.NewGetTabelogInfoHandler()
}

func provideGetTabelogPhotoService() service.GetTabelogPhotoHandler {
	return service.NewGetTabelogPhotoHandler()
}

func provideGetTabelogInfoController(
	getTabelogInfoService service.GetTabelogInfoHandler,
	getTabelogPhotoService service.GetTabelogPhotoHandler,
) *controller.GetTabelogInfoHandle {
	return controller.NewGetTabelogInfoHandle(
		getTabelogInfoService,
		getTabelogPhotoService,
	)
}
