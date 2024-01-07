package main

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAllRoutesExist(t *testing.T) {
	// Set the router to the default Gin mode.
	gin.SetMode(gin.TestMode)

	// Add your routes to the router
	key, err := generateRandomString(32)
	require.NoError(t, err, "Error generating random string")
	config := Config{TokenSymmetricKey: key}
	server, err := NewServer(config, nil, nil, CacheInstance{})
	require.NoError(t, err, "Error creating server")

	// Define the expected routes
	expectedRoutes := map[string]struct{}{
		"/regist":                     {},
		"/login":                      {},
		"/renew_access":               {},
		"/find_place":                 {},
		"/set_jp_name":                {},
		"/favorite":                   {},
		"/get_favs":                   {},
		"/get_favs_by_country":        {},
		"/get_favs_by_country_region": {},
		"/get_fav_countries":          {},
		"/get_fav_regions":            {},
		"/check_update_fav":           {},
		"/get_user":                   {},
	}

	// Verify that all expected routes exist
	assertRouteExists(t, server.router, expectedRoutes)

}

func assertRouteExists(t *testing.T, router *gin.Engine, expectedRoute map[string]struct{}) {
	// Iterate over the router's registered routes
	for _, registeredRoute := range router.Routes() {
		_, exist := expectedRoute[registeredRoute.Path]
		if !exist {
			assert.Failf(t, "Route %s does not exist", registeredRoute.Path)
		}
	}
}
