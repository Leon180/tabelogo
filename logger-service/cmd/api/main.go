package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
)

func main() {
	// load config
	config, err := LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}
	// connect to mongoDB
	mongoClient, err := connectToMongoDB(config)
	if err != nil {
		log.Fatal(err)
	}

	// create a context to disconnect from mongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// create a new server
	server, err := NewServer(config, mongoClient)
	if err != nil {
		log.Fatal(err)
	}

	// run server
	server.Serve()

}

func connectToMongoDB(config Config) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: config.MongoUser,
		Password: config.MongoPassword,
	})

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting to MongoDB", err)
		return nil, err
	}

	return client, nil
}
