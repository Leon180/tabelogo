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
	router.POST("/regist", server.TransRequest("POST", authenticateServiceURL+"/regist"))
	router.POST("/login", server.TransRequest("POST", authenticateServiceURL+"/login"))
	router.POST("/renew_access", server.TransRequest("POST", authenticateServiceURL+"/renew_access"))
	router.POST("/favorite", server.TransRequest("POST", authenticateServiceURL+"/favorite")) // toggle favorite
	router.POST("/get_favs", server.TransRequest("POST", authenticateServiceURL+"/get_favs"))
	router.POST("/get_favs_by_country", server.TransRequest("POST", authenticateServiceURL+"/get_favs_by_country"))
	router.POST("/get_favs_by_country_region", server.TransRequest("POST", authenticateServiceURL+"/get_favs_by_country_region"))
	router.POST("get_fav_countries", server.TransRequest("POST", authenticateServiceURL+"/get_fav_countries"))
	router.POST("get_fav_regions", server.TransRequest("POST", authenticateServiceURL+"/get_fav_regions"))
	router.POST("/check_update_fav", server.TransRequest("POST", authenticateServiceURL+"/check_update_fav"))
	router.POST("/get_user", server.TransRequest("POST", authenticateServiceURL+"/get_user"))
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
