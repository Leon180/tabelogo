package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"tabelog-spider/config"
	"tabelog-spider/grpc/proto"
	"tabelog-spider/inject"
	"tabelog-spider/middleware"
	"tabelog-spider/utility"
	"time"

	_ "tabelog-spider/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/location"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// @title           Tabelog Spider API
// @version         1.0
// @description     Tabelog Spider service API documentation
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  your-email@domain.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:80
// @BasePath /tabelogo-spider
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
	rateLimiter := rate.NewLimiter(rate.Every(100*time.Millisecond), 1)
	// inject
	controllerHandle := inject.InitControllerHandle(utility.Logger, rateLimiter)

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
	utility.SugarLogger.Info("[TABELOG-SERVICE] Shutting down gracefully")
	utility.SugarLogger.Info("[TABELOG-SERVICE] Server shutdown")
}

func setRoute(engine *gin.Engine, controllerHandle *inject.ControllerHandle) {
	// Swagger docs
	// Use this URL config for swagger
	url := ginSwagger.URL("/swagger/doc.json") // The url pointing to API definition
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	defaultRouter := engine.Group(engine.BasePath())
	baseRouter := defaultRouter.Group("/tabelogo-spider")
	baseRouter.GET("/getTabelogInfo", controllerHandle.GetTabelogInfoController.GetTabelogInfo)
	baseRouter.GET("/getTabelogPhoto", controllerHandle.GetTabelogInfoController.GetTabelogPhoto)
}

func setGRPCService(s *grpc.Server, controllerHandle *inject.ControllerHandle) {
	proto.RegisterTabelogoSpiderServiceServer(s, controllerHandle.TabelogoSpiderServiceServer)
}
