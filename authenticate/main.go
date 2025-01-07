package main

import (
	"authenticate/config"
	"authenticate/grpc/proto"
	"authenticate/inject"
	"authenticate/middleware"
	"authenticate/model/enum"
	"authenticate/postgresqldb"
	"authenticate/postgresqldb/postgresqldbMigrate"
	"authenticate/redisDB"
	"authenticate/utility"
	"bytes"
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"os/signal"
	"slices"
	"syscall"
	"time"

	_ "authenticate/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/location"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

// @title           Authentication API
// @version         1.0
// @description     Authentication service API documentation
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  your-email@domain.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:80
// @BasePath /authenticate

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT authorization header
func main() {
	var (
		cfg              config.Config
		engine           *gin.Engine
		db               *gorm.DB
		redisClient      *redis.Client
		controllerHandle *inject.ControllerHandle
	)
	if err := config.LoadConfig(&cfg, "./config"); err != nil {
		log.Fatalf("cannot load config: %+v", err)
	}
	utility.InitLogger(cfg.GenLogConfig())
	defer utility.Logger.Sync()
	defer utility.SugarLogger.Sync()
	utility.SugarLogger.Debugf("Env========= %s", cfg.Environment)

	// init postgresql db
	db = postgresqldb.InitGormDB(cfg, utility.Logger)
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			utility.SugarLogger.Error(err)
		}
		if err := sqlDB.Close(); err != nil {
			utility.SugarLogger.Error(err)
		}
	}()
	if err := postgresqldbMigrate.MigrateDB(db); err != nil {
		utility.SugarLogger.Fatal(err)
	}

	// connect to redis
	redisClient = redisDB.ConnectToRedis(context.Background(), cfg.RedisConnectHost)

	// inject
	controllerHandle = inject.InitControllerHandle(db, &cfg, utility.Logger, redisClient, cfg.TokenSymmetricKey)

	// init middleware
	authMiddleware := middleware.NewTransactionAuthMiddleware(redisClient, db)
	traceIDMiddleware := middleware.NewTraceIDMiddleware()

	engine = gin.Default()
	engine.SetTrustedProxies(nil)
	engine.Use(
		ginzap.Ginzap(utility.Logger, time.RFC3339, true),
		ginzap.RecoveryWithZap(utility.Logger, true),
		gzip.Gzip(gzip.DefaultCompression),
		location.Default(),
		cors.New(cfg.GenCORSConfig()),
		LogRequest(),
		traceIDMiddleware.Handler(),
	)
	setRoute(engine, controllerHandle, authMiddleware)
	server := &http.Server{
		Addr:    ":" + cfg.ConnWebPort,
		Handler: engine,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			utility.SugarLogger.Fatal(err)
		}
	}()

	lis, err := net.Listen("tcp", ":"+cfg.ConnGRPCPort) // Add cfg.GRPCPort to your config
	if err != nil {
		utility.SugarLogger.Fatalf("failed to listen: %v", err)
	}
	authInterceptor := middleware.NewTransactionAuthInterceptor(redisClient, db)
	traceIDInterceptor := middleware.NewTraceIDInterceptor()
	s := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
		grpc.UnaryInterceptor(traceIDInterceptor.Unary()),
	)
	setGRPCService(s, controllerHandle)
	reflection.Register(s)
	go func() {
		if err := s.Serve(lis); err != nil {
			utility.SugarLogger.Fatal(err)
		}
	}()

	// graceful shutdown setup
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// graceful shutdown
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Graceful shutdown
	if err := server.Shutdown(shutdownCtx); err != nil {
		utility.SugarLogger.Errorf("Error during server shutdown: %v", err)
	}
	s.GracefulStop()
	utility.SugarLogger.Info("[AUTHENTICATE-SERVICE] Shutting down gracefully")
	utility.SugarLogger.Info("[AUTHENTICATE-SERVICE] Server shutdown")
}

func LogRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		eventID := utility.GenDefaultUUID()
		c.Set(enum.MiddleWareEventIDKey, eventID)
		if contextType := c.Request.Header.Get("Content-type"); slices.Contains(enum.ContextTypeGroupDefault.GetSlice().ToStringSlice(), contextType) {
			buf, err := io.ReadAll(c.Request.Body)
			if err != nil {
				utility.SugarLogger.Error(err)
				c.Next()
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(buf))
			utility.SugarLogger.Infof("Http Request: %+v, EventID: %s, Body: %s", c.Request, eventID, string(buf))
		} else {
			utility.SugarLogger.Infof("Http Request: %+v, EventID: %s", c.Request, eventID)
		}
		c.Next()
	}
}

func setRoute(engine *gin.Engine, controllerHandle *inject.ControllerHandle, authMiddleware *middleware.AuthMiddleware) {
	// Swagger docs
	// Use this URL config for swagger
	url := ginSwagger.URL("/swagger/doc.json") // The url pointing to API definition
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	defaultRouter := engine.Group(engine.BasePath())
	baseRouter := defaultRouter.Group("/authenticate")
	placeRouter := baseRouter.Group("/place").Use(authMiddleware.Handler())
	{
		placeRouter.POST("/savePlace", controllerHandle.SavePlaceController.SavePlace)
		placeRouter.GET("/getPlace", controllerHandle.GetPlaceController.GetPlace)
	}
	userRouter := baseRouter.Group("/user")
	{
		userRouter.POST("/registUser", controllerHandle.RegistUserController.RegistUser)
		userRouter.POST("/loginUser", controllerHandle.LoginUserController.LoginUser)
		userRouter.POST("/renewAccessToken", controllerHandle.RenewAccessTokenController.RenewAccessToken)
		userRouter.POST("/logoutUser", authMiddleware.Handler(), controllerHandle.LogoutUserController.LogoutUser)
		userRouter.POST("/saveFavorite", authMiddleware.Handler(), controllerHandle.SaveFavoriteController.SaveFavorite)
		userRouter.GET("/getUserFavorites", authMiddleware.Handler(), controllerHandle.GetUserFavoritesController.GetUserFavorites)
	}
}

func setGRPCService(s *grpc.Server, controllerHandle *inject.ControllerHandle) {
	proto.RegisterUserServiceServer(s, controllerHandle.UserServiceServer)
	proto.RegisterPlaceServiceServer(s, controllerHandle.PlaceServiceServer)
}
