package main

import (
	db "authenticate/cmd/data/sqlc"
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog/log"
)

var counts int64

const webPort = "80"

func main() {
	config, err := LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	conn := connectToDB(config.DSN)
	store := db.NewStore(conn)

	server, err := NewServer(config, store)
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

func connectToDB(dsn string) *sql.DB {

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Log().Msg("Postgres not yet ready ...")
			counts++
		} else {
			log.Log().Msg("Connected to Postgres!")
			return connection
		}

		if counts > 10 {
			log.Log().Msg(err.Error())
			return nil
		}

		log.Log().Msg("Backing off for two seconds....")
		time.Sleep(2 * time.Second)
		continue
	}
}
