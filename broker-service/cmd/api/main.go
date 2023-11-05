package main

import (
	_ "broker-service/cmd/api/docs"
	"broker-service/internal/config"
	"broker-service/internal/handlers"
	"broker-service/internal/routers"
	"fmt"
	"log"
	"net/http"
	"os"
)

const webPort = "8080"

var authGrpcPort = os.Getenv("authGrpcPort")

// @title kinosaya API
// @version 1.0

// @host localhost:8080
// @BasePath /

func main() {
	// TODO: try to connect to rabbitmq

	// Dependency Injection
	broker := handlers.NewBrokerHandler(authGrpcPort)
	chiRouter := routers.NewChiRouters(broker, webPort)
	app := config.NewConfig(chiRouter)

	log.Printf("Starting broker service on port %s\n", webPort)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.Router.GetRoutes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
