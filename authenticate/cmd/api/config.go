package main

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DSN                  string        `mapstructure:"DSN"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	RefreshDuration      time.Duration `mapstructure:"REFRESH_DURATION"`
	RabbitMQConnect      string        `mapstructure:"RABBITMQ_CONNECT"`
	RedisConnectSession  string        `mapstructure:"REDIS_CONNECT_SESSION"`
	RedisConnectPlace    string        `mapstructure:"REDIS_CONNECT_PLACE"`
	RedisConnectTabelogo string        `mapstructure:"REDIS_CONNECT_TABELOGO"`
}

func LoadConfig(path string) (config Config, err error) {

	// set the config file type
	viper.SetConfigType("env")
	viper.SetConfigName("app")
	viper.AddConfigPath(path)

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
