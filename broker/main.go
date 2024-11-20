package main

import (
	"broker/config"
	"broker/model/enum"
	"broker/server"
	"broker/utility"
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"slices"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/location"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	webPort                 = "8080"
	tabelogSpiderServiceURL = "http://tabelog-spider-service" // service's name
	authenticateServiceURL  = "http://authenticate-service"   // service's name
	googleMapServiceURL     = "http://google-map-service"     // service's name
	loggerServiceURL        = "http://logger-service"         // service's name
	mailServiceURL          = "http://mail-service"           // service's name
)

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
	rabbitConn, err := connectToRabbitMQ()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	setRoute(engine, rabbitConn)
	if err := engine.Run(":" + cfg.ConnWebPort); err != nil {
		utility.SugarLogger.Fatal(err)
	}
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

func connectToRabbitMQ() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}

func setRoute(engine *gin.Engine, rabbitConn *amqp.Connection) {

	server := server.NewServer(engine, rabbitConn)

	defaultRouter := engine.Group(engine.BasePath())
	baseRouter := defaultRouter.Group("/tabelogo")

	// Routes
	baseRouter.POST("/", server.Broker)

	// Tabelog spider service:
	tabelogSpiderRouter := baseRouter.Group("/tabelogo-spider")
	tabelogSpiderRouter.POST("/getTabelogInfo", server.HTTPTransRequest("POST", tabelogSpiderServiceURL+"/getTabelogInfo"))
	tabelogSpiderRouter.POST("/getTabelogPhoto", server.HTTPTransRequest("POST", tabelogSpiderServiceURL+"/getTabelogPhoto"))

	// Authenticate service:
	authenticateRouter := baseRouter.Group("/authenticate")
	placeRouter := authenticateRouter.Group("/place")
	placeRouter.POST("/savePlace", server.HTTPTransRequest("POST", authenticateServiceURL+"/savePlace"))
	placeRouter.GET("/getPlace", server.HTTPTransRequest("GET", authenticateServiceURL+"/getPlace"))
	userRouter := authenticateRouter.Group("/user")
	userRouter.POST("/registUser", server.HTTPTransRequest("POST", authenticateServiceURL+"/registUser"))
	userRouter.POST("/loginUser", server.HTTPTransRequest("POST", authenticateServiceURL+"/loginUser"))
	userRouter.POST("/renewAccessToken", server.HTTPTransRequest("POST", authenticateServiceURL+"/renewAccessToken"))
	userRouter.POST("/saveFavorite", server.HTTPTransRequest("POST", authenticateServiceURL+"/saveFavorite"))
	userRouter.GET("/getUserFavorites", server.HTTPTransRequest("GET", authenticateServiceURL+"/getUserFavorites"))

	// Google API service:
	googleMapRouter := baseRouter.Group("/google-map-search")
	googleMapRouter.POST("/quickSearch", server.HTTPTransRequest("POST", googleMapServiceURL+"/quickSearch"))
	googleMapRouter.POST("/advanceSearch", server.HTTPTransRequest("POST", googleMapServiceURL+"/advanceSearch"))

	// logger service:(for testing)
	loggerRouter := baseRouter.Group("/logger")
	loggerRouter.POST("/createLog", server.HTTPTransRequest("POST", loggerServiceURL+"/createLog"))
	loggerRouter.GET("/readAllLogs", server.HTTPTransRequest("GET", loggerServiceURL+"/readAllLogs"))
	loggerRouter.POST("/searchLogs", server.HTTPTransRequest("POST", loggerServiceURL+"/searchLogs"))
}
