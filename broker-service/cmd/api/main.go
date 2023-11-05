package main

import (
	_ "broker-service/cmd/api/docs"
	"fmt"
	"log"
	"net/http"
	"os"
)

const webPort = "8080"

var authGrpcPort = os.Getenv("authGrpcPort")

type Config struct {
}

// @title kinosaya API
// @version 1.0

// @host localhost:8080
// @BasePath /

func main() {
	// TODO: try to connect to rabbitmq

	app := Config{}

	log.Printf("Starting broker service on port %s\n", webPort)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
