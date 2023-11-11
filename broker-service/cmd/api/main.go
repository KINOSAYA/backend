package main

import (
	_ "broker-service/docs"
	"broker-service/internal/config"
	"broker-service/internal/handlers"
	"broker-service/internal/routes"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"os"
)

// @title kinosaya API
// @version 1.0

// @host localhost:8080
// @BasePath /

var webPort = os.Getenv("webPort")
var authHost = os.Getenv("authHost")

var authGrpcPort = os.Getenv("authGrpcPort")

func main() {
	// TODO: try to connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	// Dependency Injection
	app := config.Config{Rabbit: rabbitConn}
	broker := handlers.NewBrokerHandler(authGrpcPort, authHost, app.Rabbit)
	router := routes.NewChiRouter(broker, webPort)

	log.Printf("Starting broker service on port %s\n", webPort)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: router.GetRoutes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp.Connection, error) {

	return nil, nil
}
