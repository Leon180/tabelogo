package main

import (
	"context"
	"google-map/config"
	"google-map/grpc/proto"
	"google-map/inject"
	"google-map/middleware"
	"google-map/utility"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	_ "google-map/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/location"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// @title           Google Map API
// @version         1.0
// @description     Google Map service API documentation
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  your-email@domain.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:80
// @BasePath /google-map-search
func main() {
	var (
		cfg    config.Config
		engine *gin.Engine
	)
	if err := config.LoadConfig(&cfg, "./config"); err != nil {
		log.Fatalf("cannot load config: %+v", err)
	}
	utility.InitLogger(cfg.GenLogConfig())
	defer utility.Logger.Sync()
	defer utility.SugarLogger.Sync()
	utility.SugarLogger.Debugf("Env========= %s", cfg.Environment)
	// inject
	controllerHandle := inject.InitControllerHandle(&cfg, utility.Logger)

	// init middleware
	traceIDMiddleware := middleware.NewTraceIDMiddleware()

	engine = gin.Default()
	engine.SetTrustedProxies(nil)
	engine.Use(
		ginzap.Ginzap(utility.Logger, time.RFC3339, true),
		ginzap.RecoveryWithZap(utility.Logger, true),
		gzip.Gzip(gzip.DefaultCompression),
		location.Default(),
		cors.New(cfg.GenCORSConfig()),
		traceIDMiddleware.Handler(),
	)
	setRoute(engine, controllerHandle)
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
	traceIDInterceptor := middleware.NewTraceIDInterceptor()
	s := grpc.NewServer(
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
	utility.SugarLogger.Info("[GOOGLE-MAP-SERVICE] Shutting down gracefully")
	utility.SugarLogger.Info("[GOOGLE-MAP-SERVICE] Server shutdown")
}

func setRoute(engine *gin.Engine, controllerHandle *inject.ControllerHandle) {
	// Swagger docs
	// Use this URL config for swagger
	url := ginSwagger.URL("/swagger/doc.json") // The url pointing to API definition
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	defaultRouter := engine.Group(engine.BasePath())
	baseRouter := defaultRouter.Group("/google-map-search")
	baseRouter.GET("/quickSearch", controllerHandle.GooglePlaceSearchController.QuickSearch)
	baseRouter.GET("/advanceSearch", controllerHandle.GooglePlaceSearchController.AdvanceSearch)
}

func setGRPCService(s *grpc.Server, controllerHandle *inject.ControllerHandle) {
	proto.RegisterGoogleMapServiceServer(s, controllerHandle.GoogleMapServiceServer)
}
