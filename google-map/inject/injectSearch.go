package inject

import (
	"google-map/config"
	"google-map/controller"
	grpcservice "google-map/grpc/service"
	"google-map/service"
)

func provideSearchService(config *config.Config) service.GooglePlaceSearchHandler {
	return service.NewGooglePlaceSearchHandler(config)
}

func provideSearchController(
	searchService service.GooglePlaceSearchHandler,
	config *config.Config,
) *controller.GooglePlaceSearchHandle {
	return controller.NewGooglePlaceSearchHandle(
		searchService,
		config,
	)
}

func provideGoogleMapServiceServer(
	googlePlaceSearchService service.GooglePlaceSearchHandler,
	config *config.Config,
) *grpcservice.GoogleMapServiceServer {
	return grpcservice.NewGoogleMapServiceServer(googlePlaceSearchService, config)
}
