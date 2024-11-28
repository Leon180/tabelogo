package inject

import (
	"tabelog-spider/controller"
	"tabelog-spider/service"

	"golang.org/x/time/rate"
)

func provideGetTabelogInfoService(rateLimiter *rate.Limiter) service.GetTabelogInfoHandler {
	return service.NewGetTabelogInfoHandler(rateLimiter)
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
