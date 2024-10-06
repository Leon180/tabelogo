package main

import (
	"authenticate/config"
	"authenticate/inject"
	"authenticate/model/enum"
	"authenticate/postgresqldb"
	"authenticate/postgresqldb/postgresqldbMigrate"
	"authenticate/redisDB"
	"authenticate/utility"
	"bytes"
	"context"
	"io"
	"log"
	"os/signal"
	"slices"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/location"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

func main() {
	var (
		cfg         config.Config
		engine      *gin.Engine
		db          *gorm.DB
		redisClient *redis.Client
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
	redisClient = redisDB.ConnectToRedis(cfg.RedisConnectHost)

	// inject
	controllerHandle := inject.InitControllerHandle(db, cfg, utility.Logger, redisClient)

	// graceful shutdown setup
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	engine = gin.Default()
	engine.SetTrustedProxies(nil)
	engine.Use(
		ginzap.Ginzap(utility.Logger, time.RFC3339, true),
		ginzap.RecoveryWithZap(utility.Logger, true),
		gzip.Gzip(gzip.DefaultCompression),
		location.Default(),
		cors.New(cfg.GenCORSConfig()),
		LogRequest(),
	)
	setRoute(engine, controllerHandle)
	if err := engine.Run(":" + cfg.ConnWebPort); err != nil {
		utility.SugarLogger.Fatal(err)
	}

	// graceful shutdown
	<-ctx.Done()
	stop()
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

func setRoute(engine *gin.Engine, controllerHandle *inject.ControllerHandle) {
	// defaultRouter := engine.Group(engine.BasePath())
	// baseRouter := defaultRouter.Group("/tabelogo-spider/api/v1")
	// getTabelogInfoHandle := controller.NewGetTabelogInfoHandle(service.NewGetTabelogInfoHandler(), service.NewGetTabelogPhotoHandler())
	// baseRouter.GET("/getTabelogInfo", getTabelogInfoHandle.GetTabelogInfo)
	// baseRouter.GET("/getTabelogPhoto", getTabelogInfoHandle.GetTabelogPhoto)
}
