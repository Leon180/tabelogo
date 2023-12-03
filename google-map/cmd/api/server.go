package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config Config
	router *gin.Engine
}

func NewServer(config Config) (*Server, error) {

	server := &Server{
		config: config,
	}
	router := gin.Default()
	// CORS configuration
	router.Use(cors.New(CORSConfig()))
	// Logging and Recovery
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// Routes
	router.POST("/", server.Hello)
	router.POST("/quick_search", server.QuickSearch)
	router.POST("/advance_search", server.AdvanceSearch)
	server.router = router
	return server, nil
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

func CORSConfig() cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true // for testing
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"X-PINGOTHER", "Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Upgrade", "Origin",
		"Connection", "Accept-Encoding", "Accept-Language", "Host", "Access-Control-Request-Method", "Access-Control-Request-Headers"}
	corsConfig.AllowCredentials = true
	corsConfig.ExposeHeaders = []string{"Content-Length", "Link"}
	corsConfig.MaxAge = 500
	return corsConfig
}
