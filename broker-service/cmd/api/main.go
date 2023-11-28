package main

import (
	_ "broker-service/docs"
	"broker-service/event"
	"broker-service/internal/config"
	"broker-service/internal/handlers"
	"broker-service/internal/routes"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

// @title kinosaya API
// @version 1.0

// @host localhost:8080
// @BasePath /

var webPort = os.Getenv("webPort")
var authHost = os.Getenv("authHost")

var authGrpcPort = os.Getenv("authGrpcPort")

func main() {
	// connecting to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	// Dependency Injection
	app := config.Config{
		AmqpService: event.NewRabbitEventService(rabbitConn),
	}
	broker := handlers.NewBrokerHandler(authGrpcPort, authHost, app)
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
