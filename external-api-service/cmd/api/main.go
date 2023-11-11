package main

import (
	"external-api-service/event"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"time"
)

func main() {
	conn, err := connectToRabbitMQ()
	if err != nil {
		log.Fatal("cannot connect to RabbitMQ!")
	}
	defer conn.Close()

	log.Println("Listening for and consuming RabbitMQ messages...")
	// create consumer
	consumer := event.NewEventConsumer(conn)
	err = consumer.Listen([]string{"info"})
	if err != nil {
		log.Fatal("can't consume messages from queue")
	}

}

func connectToRabbitMQ() (*amqp.Connection, error) {
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
