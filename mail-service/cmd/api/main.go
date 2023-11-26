package main

import (
	"fmt"
	"log"
)

const (
	webPort = "80"
)

func main() {
	// load config
	config, err := LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}
	mail := createMail(config)
	// create a new server
	server, err := NewServer(config, mail)
	if err != nil {
		log.Fatal(err)
	}

	// run server
	err = server.Run(":" + webPort)
	if err != nil {
		panic(err)
	}

}

func createMail(config Config) Mail {
	m := Mail{
		Domain:      config.MailDomain,
		Host:        config.MailHost,
		Port:        config.MailPort,
		Username:    config.MailUsername,
		Password:    config.MailPassword,
		Encryption:  config.MailEncryption,
		FromName:    config.MailFromName,
		FromAddress: config.MailFromAddress,
	}
	fmt.Println(m)
	return m
}
