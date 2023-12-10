package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	RedisConnectTabelogo string `mapstructure:"REDIS_CONNECT_TABELOGO"`
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
