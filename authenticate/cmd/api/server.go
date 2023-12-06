package main

import (
	db "authenticate/cmd/data/sqlc"
	"authenticate/token"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Server struct {
	config     Config
	router     *gin.Engine
	store      db.Store
	rabbitMQ   *amqp.Connection
	tokenMaker token.Maker
}

func NewServer(config Config, store db.Store, rabbitConn *amqp.Connection) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		config:     config,
		tokenMaker: tokenMaker,
		store:      store,
		rabbitMQ:   rabbitConn,
	}
	router := gin.Default()
	// CORS configuration
	router.Use(cors.New(CORSConfig()))
	// Logging and Recovery
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// Routes
	router.POST("/regist", server.Regist)
	router.POST("/login", server.Login)
	router.POST("/renew_access", server.RenewAccessToken)
	// authGroup
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/favorite", server.ToggleFavorite)
	authRoutes.POST("/get_favs", server.GetListFavorites)
	authRoutes.POST("/get_favs_by_country", server.GetListFavoritesByCountry)
	authRoutes.POST("/get_favs_by_country_region", server.GetListFavoritesByCountryAndRegion)
	authRoutes.POST("/get_fav_countries", server.GetFavoritesCountry)
	authRoutes.POST("/get_fav_regions", server.GetFavoritesRegion)
	authRoutes.POST("/check_update_fav", server.CheckAndUpdateFavorite)
	authRoutes.POST("/get_user", server.GetUser)
	// authRoutes.POST("/delete", server.DeletePlace)
	// authRoutes.POST("/get", server.GetPlace)
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
