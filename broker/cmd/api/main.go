package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

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
	rabbitConn, err := connectToRabbitMQ()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	server := NewServer(rabbitConn)
	err = server.Run(":" + webPort)
	if err != nil {
		panic(err)
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
