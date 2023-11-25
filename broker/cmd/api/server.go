package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	server := &Server{}
	router := gin.Default()
	// CORS configuration
	router.Use(cors.New(CORSConfig()))
	// Logging and Recovery
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// Routes
	router.POST("/", server.Broker)
	router.POST("/tabelogo", server.TransRequest("POST", tabelogSpiderServiceURL))
	// Authenticate service:
	router.POST("/regist", server.TransRequest("POST", authenticateServiceURL+"/regist"))
	router.POST("/login", server.TransRequest("POST", authenticateServiceURL+"/login"))
	router.POST("/renew_access", server.TransRequest("POST", authenticateServiceURL+"/renew_access"))
	// Google API service:
	router.POST("/quick_search", server.TransRequest("POST", googleMapServiceURL+"/quick_search"))
	router.POST("/advance_search", server.TransRequest("POST", googleMapServiceURL+"/advance_search"))
	server.router = router
	return server
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
