package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"logger-service/config"
	"logger-service/controller"
	"logger-service/model/enum"
	"logger-service/repository"
	"logger-service/service"
	"logger-service/utility"
	"slices"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/location"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var (
		cfg         config.Config
		engine      *gin.Engine
		mongoClient *mongo.Client
	)
	if err := config.LoadConfig(&cfg, "./config"); err != nil {
		log.Fatalf("cannot load config: %+v", err)
	}
	utility.InitLogger(cfg.GenLogConfig())
	defer utility.Logger.Sync()
	defer utility.SugarLogger.Sync()
	utility.SugarLogger.Debugf("Env========= %s", cfg.Environment)
	mongoClient, err := connectToMongoDB(cfg)
	if err != nil {
		utility.SugarLogger.Fatal(err)
	}
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
	setRoute(engine, mongoClient)
	if err := engine.Run(":" + cfg.ConnWebPort); err != nil {
		utility.SugarLogger.Fatal(err)
	}
}

func connectToMongoDB(config config.Config) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(config.ConnMongoDB)
	clientOptions.SetAuth(options.Credential{
		Username: config.MongoUser,
		Password: config.MongoPassword,
	})

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		utility.SugarLogger.Error("Error connecting to MongoDB", err)
		return nil, err
	}

	return client, nil
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

func setRoute(engine *gin.Engine, mongoClient *mongo.Client) {
	defaultRouter := engine.Group(engine.BasePath())
	baseRouter := defaultRouter.Group("/logger/api/v1")
	logController := controller.NewLogController(
		service.NewCreateLogService(
			repository.NewCreateLogRepository(mongoClient),
		),
		service.NewReadLogService(
			repository.NewReadLogRepository(mongoClient),
		),
		service.NewUpdateLogService(
			repository.NewUpdateLogRepository(mongoClient),
		),
		service.NewDeleteLogService(
			repository.NewDeleteLogRepository(mongoClient),
		),
	)
	baseRouter.POST("/createLog", logController.CreateLog)
	baseRouter.GET("/readAllLogs", logController.ReadAllLogs)
	baseRouter.POST("/searchLogs", logController.SearchLogs)
}
