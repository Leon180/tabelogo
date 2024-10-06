package config

import (
	"authenticate/model/enum"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	GoogleMapAPIKey                            string        `mapstructure:"GOOGLE_MAP_API_KEY"`
	MaxSize                                    int           `mapstructure:"MAX_SIZE"`
	MaxAge                                     int           `mapstructure:"MAX_AGE"`
	MaxBackups                                 int           `mapstructure:"MAX_BACKUPS"`
	Compress                                   bool          `mapstructure:"COMPRESS"`
	UseRedis                                   bool          `mapstructure:"USE_REDIS"`
	Environment                                string        `mapstructure:"-"`
	CORSMaxAge                                 int           `mapstructure:"CORS_MAX_AGE"`
	CORSAllowAllOrigins                        bool          `mapstructure:"CORS_ALLOW_ALL_ORIGINS"`
	CORSAllowCredentials                       bool          `mapstructure:"CORS_ALLOW_CREDENTIALS"`
	ConnWebPort                                string        `mapstructure:"CONN_WEB_PORT"`
	ConnMaxCollectLinks                        int           `mapstructure:"CONN_MAX_COLLECT_LINKS"`
	CORSAllowMethods                           string        `mapstructure:"CORS_ALLOW_METHODS"`
	CORSAllowHeaders                           string        `mapstructure:"CORS_ALLOW_HEADERS"`
	CORSExposeHeaders                          string        `mapstructure:"CORS_EXPOSE_HEADERS"`
	RABBITMQConnect                            string        `mapstructure:"RABBITMQ_CONNECT"`
	RedisConnectHost                           string        `mapstructure:"REDIS_CONNECT_HOST"`
	TokenSymmetricKey                          string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration                        time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration                       time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	RefreshDuration                            time.Duration `mapstructure:"REFRESH_DURATION"`
	DSNTest                                    string        `mapstructure:"DSN_TEST"`
	DSNDeployment                              string        `mapstructure:"DSN_DEPLOYMENT"`
	DBMaxIdle                                  int           `mapstructure:"DB_MAX_IDLE"`
	DBMaxOpen                                  int           `mapstructure:"DB_MAX_OPEN"`
	DBMaxLifetimeMinute                        int           `mapstructure:"DB_MAX_LIFETIME_MINUTE"`
	DBDisableForeignKeyConstraintWhenMigrating bool          `mapstructure:"DB_DISABLE_FOREIGN_KEY_CONSTRAINT_WHEN_MIGRATING"`
}

func LoadConfig(config *Config, path string) error {
	viper.SetConfigType("env")
	viper.SetConfigName("config")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&config); err != nil {
		return err
	}
	if env := os.Getenv("CONFIG_ENV"); env != "" {
		config.Environment = env
	} else {
		config.Environment = "development"
	}
	return nil
}

func (c Config) GenCORSConfig() cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = c.CORSAllowAllOrigins
	corsConfig.AllowCredentials = c.CORSAllowCredentials
	corsConfig.MaxAge = time.Duration(c.CORSMaxAge)
	corsConfig.AllowMethods = enum.ContextTypeGroup(c.CORSAllowMethods).GetSlice().ToStringSlice()
	corsConfig.AllowHeaders = enum.RequestHeaderGroup(c.CORSAllowHeaders).GetSlice().ToStringSlice()
	corsConfig.ExposeHeaders = enum.ExposeHeaderGroup(c.CORSExposeHeaders).GetSlice().ToStringSlice()
	return corsConfig
}

type LogConfig struct {
	Level      zapcore.Level `json:"level"`
	FileName   string        `json:"file_name"`
	MaxSize    int           `json:"max_size"`
	MaxAge     int           `json:"max_age"`
	MaxBackups int           `json:"max_backups"`
	Compress   bool          `json:"compress"`
}

func (c Config) GenLogConfig() LogConfig {
	return LogConfig{
		FileName:   fmt.Sprintf("./log/tabelogo-google-search-service-%s.log", strings.ToLower(c.Environment)),
		Level:      zapcore.DebugLevel,
		MaxSize:    c.MaxSize,
		MaxAge:     c.MaxAge,
		MaxBackups: c.MaxBackups,
		Compress:   c.Compress,
	}
}
