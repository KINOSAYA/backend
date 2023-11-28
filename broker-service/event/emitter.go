package event

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type Emitter struct {
	connection *amqp.Connection
}

func (e Emitter) Push(conn *amqp.Connection, payload any) error {
	c, err := conn.Channel()
	if err != nil {
		log.Fatalf("channel.open: %s", err)
	}

	err = c.ExchangeDeclare("logs", "topic", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("exchange.declare: %v", err)
	}

	jsonBody, err := json.Marshal(payload)
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "application/json", // Change content type to "application/json"
		Body:         jsonBody,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = c.PublishWithContext(ctx, "logs", "info", false, false, msg)
	if err != nil {
		log.Fatalf("basic.publish: %v", err)
	}

	return nil
}

func NewEventEmitter(conn *amqp.Connection) Emitter {
	emitter := Emitter{
		connection: conn,
	}
	return emitter
}
