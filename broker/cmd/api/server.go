package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Server struct {
	router   *gin.Engine
	rabbitMQ *amqp.Connection
}

func NewServer(rabbitConn *amqp.Connection) *Server {

	server := &Server{
		rabbitMQ: rabbitConn,
	}
	router := gin.Default()
	// CORS configuration
	router.Use(cors.New(CORSConfig()))
	// Logging and Recovery
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// Routes
	router.POST("/", server.Broker)
	// Tabelog spider service:
	router.POST("/tabelogo", server.TransRequest("POST", tabelogSpiderServiceURL))
	router.POST("/tabephoto", server.TransRequest("POST", tabelogSpiderServiceURL+"/photo"))
	// Authenticate service:
	router.POST("/user/registUser", server.TransRequest("POST", authenticateServiceURL+"/user/registUser"))
	router.POST("/user/loginUser", server.TransRequest("POST", authenticateServiceURL+"/user/loginUser"))
	router.POST("/user/renewAccessToken", server.TransRequest("POST", authenticateServiceURL+"/user/renewAccessToken"))
	router.POST("/user/saveFavorite", server.TransRequest("POST", authenticateServiceURL+"/user/saveFavorite")) // toggle favorite
	router.POST("/user/getUserFavorites", server.TransRequest("POST", authenticateServiceURL+"/user/getUserFavorites"))
	// Google API service:
	router.POST("/quick_search", server.TransRequest("POST", googleMapServiceURL+"/quick_search"))
	router.POST("/advance_search", server.TransRequest("POST", googleMapServiceURL+"/advance_search"))
	// logger service:(for testing)
	// router.POST("/write_log", server.TransRequest("POST", loggerServiceURL+"/write_log"))
	router.POST("/write_log", server.logEventViaRabbit)
	// mail
	router.POST("/send_mail", server.TransRequest("POST", mailServiceURL+"/send"))

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
