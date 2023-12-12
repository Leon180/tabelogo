package main

import (
	db "authenticate/cmd/data/sqlc"
	"context"
	"database/sql"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

func main() {
	config, err := LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	// connect to postgres
	conn, err := connectToDB(config.DSN)
	if err != nil {
		log.Fatal(err)
	}
	store := db.NewStore(conn)

	// connect to rabbitmq
	rabbitConn, err := connectToRabbitMQ(config.RabbitMQConnect)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	// connect to redis
	redisSession := connectToRedis(config.RedisConnectSession)
	redisPlace := connectToRedis(config.RedisConnectPlace)

	server, err := NewServer(config, store, rabbitConn, CacheInstance{
		Session: redisSession,
		Place:   redisPlace,
	})
	if err != nil {
		panic(err)
	}
	err = server.Run(":" + webPort)
	if err != nil {
		panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func connectToDB(dsn string) (*sql.DB, error) {
	var counts int64
	var backOff = 2 * time.Second
	var connection *sql.DB
	for {
		c, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			connection = c
			break
		}

		if counts > 10 {
			log.Fatal(err)
			return nil, err
		}

		log.Println("Backing off for two seconds....")
		time.Sleep(backOff)
		continue
	}
	return connection, nil
}

func connectToRabbitMQ(rabbit_conn string) (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		c, err := amqp.Dial(rabbit_conn)
		if err != nil {
			log.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			break
		}

		if counts > 5 {
			log.Fatal(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}

func connectToRedis(redis_conn string) *redis.Client {
	opts, err := redis.ParseURL(redis_conn)
	if err != nil {
		log.Panic(err)
	}
	rdb := redis.NewClient(opts)
	err = rdb.Ping(context.Background()).Err()
	if err != nil {
		log.Panic(err)
	}
	return rdb
}
