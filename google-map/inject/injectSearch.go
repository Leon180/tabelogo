package inject

import (
	"google-map/config"
	"google-map/controller"
	"google-map/service"
)

func provideSearchService() service.GooglePlaceSearchHandler {
	return service.NewGooglePlaceSearchHandler()
}

func provideSearchController(
	searchService service.GooglePlaceSearchHandler,
	config config.Config,
) *controller.GooglePlaceSearchHandle {
	return controller.NewGooglePlaceSearchHandle(
		searchService,
		config,
	)
}
