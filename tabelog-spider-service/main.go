package main

import (
	"bytes"
	"io"
	"log"
	"slices"
	"tabelog-spider/config"
	"tabelog-spider/controller"
	"tabelog-spider/model/enum"
	"tabelog-spider/utility"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/location"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
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
	setRoute(engine)
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

func setRoute(engine *gin.Engine) {
	defaultRouter := engine.Group(engine.BasePath())
	baseRouter := defaultRouter.Group("/tabelogo-spider/api/v1")
	baseRouter.GET("/getTabelogInfo", controller.NewGetTabelogInfoHandle().GetTabelogInfo)
	baseRouter.GET("/getTabelogPhoto", controller.NewGetTabelogInfoHandle().GetTabelogPhoto)
}
