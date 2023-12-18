package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	BrokerURL            string `mapstructure:"BROKER_URL"`
	WebsiteURL           string `mapstructure:"WEBSITE_URL"`
	BrokerURLDeployment  string `mapstructure:"BROKER_URL_DEPLOYMENT"`
	WebsiteURLDeployment string `mapstructure:"WEBSITE_URL_DEPLOYMENT"`
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
