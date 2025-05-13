package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"net/http"
	"os"
	"time"
)

const webPort = "80"

var counts int64

type Config struct {
	//DB     *sql.DB
	//Models data.Models
	Repo   data.Repository
	Client *http.Client
}

func main() {
	log.Println("Starting authentication service at: " + webPort)

	// Connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Could not connect to Postgres")
	}
	// Set up config
	app := Config{
		//DB:     conn,
		//Models: data.New(conn),
		Client: &http.Client{},
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not ready yet ...")
			counts++
		} else {
			log.Println("Connected to Postgres")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for 2 seconds ...")
		time.Sleep(2 * time.Second)
		continue
	}
}

// Added for test
func (app *Config) setupRepo(conn *sql.DB) {
	db := data.NewPostgresRepository(conn)
	app.Repo = db
}
