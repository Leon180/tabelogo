package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"logger-service/config"
	"logger-service/inject"
	"logger-service/model/enum"
	"logger-service/utility"
	"slices"
	"time"

	_ "logger-service/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/location"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @title           Logger Service API
// @version         1.0
// @description     Logger service API documentation
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  your-email@domain.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:80
// @BasePath /logger
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
	controllerHandle := inject.InitControllerHandle(cfg, utility.Logger, mongoClient)
	setRoute(engine, controllerHandle)
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

func setRoute(engine *gin.Engine, controllerHandle *inject.ControllerHandle) {
	// Swagger docs
	// Use this URL config for swagger
	url := ginSwagger.URL("/swagger/doc.json") // The url pointing to API definition
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	defaultRouter := engine.Group(engine.BasePath())
	baseRouter := defaultRouter.Group("/logger")
	baseRouter.POST("/createLog", controllerHandle.LogController.CreateLog)
	baseRouter.GET("/readAllLogs", controllerHandle.LogController.ReadAllLogs)
	baseRouter.POST("/searchLogs", controllerHandle.LogController.SearchLogs)
}
