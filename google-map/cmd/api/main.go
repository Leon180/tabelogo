package main

import (
	"log"
)

const webPort = "80"

func main() {
	config, err := LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config")
	}

	server, err := NewServer(config)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	err = server.Run(":" + webPort)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
}
