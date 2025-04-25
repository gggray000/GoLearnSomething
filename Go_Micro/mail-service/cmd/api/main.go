package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	Mailer Mailer
}

const webPort = "80"

func main() {
	app := Config{
		Mailer: createMailer(),
	}

	log.Println("Starting mail service on " + webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func createMailer() Mailer {
	port, _ := strconv.Atoi(os.Getenv("MAILER_PORT"))

	m := Mailer{
		Domain:      os.Getenv("MAILER_DOMAIN"),
		Host:        os.Getenv("MAILER_HOST"),
		Port:        port,
		Username:    os.Getenv("MAILER_USERNAME"),
		Password:    os.Getenv("MAILER_PASSWORD"),
		Encryption:  os.Getenv("MAILER_ENCRYPTION"),
		FromName:    os.Getenv("MAILER_FROMNAME"),
		FromAddress: os.Getenv("MAILER_FROMADDRESS"),
	}

	return m
}
