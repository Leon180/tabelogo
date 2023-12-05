package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	MailDomain      string `mapstructure:"MAIL_DOMAIN"`
	MailHost        string `mapstructure:"MAIL_HOST"`
	MailPort        int    `mapstructure:"MAIL_PORT"`
	MailUsername    string `mapstructure:"MAIL_USERNAME"`
	MailPassword    string `mapstructure:"MAIL_PASSWORD"`
	MailEncryption  string `mapstructure:"MAIL_ENCRYPTION"`
	MailFromName    string `mapstructure:"MAIL_FROM_NAME"`
	MailFromAddress string `mapstructure:"MAIL_FROM_ADDRESS"`
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
