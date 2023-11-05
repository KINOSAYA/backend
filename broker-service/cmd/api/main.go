package main

import (
	_ "broker-service/docs"
	"broker-service/internal/config"
	"broker-service/internal/handlers"
	"broker-service/internal/routes"
	"fmt"
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

	// Dependency Injection
	broker := handlers.NewBrokerHandler(authGrpcPort, authHost)
	chiRouter := routes.NewChiRouter(broker, webPort)
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
