package event

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Consumer struct {
	conn *amqp.Connection
}

type Payload struct {
	Name string `json:"name"`
}

func (cons Consumer) Listen(topics []string) error {
	c, err := cons.conn.Channel()
	if err != nil {
		log.Fatalf("channel.open: %s", err)
	}
	defer c.Close()

	q, err := c.QueueDeclare("", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("exchange.declare: %s", err)
	}

	for _, s := range topics {
		err = c.QueueBind(
			q.Name,
			s,
			"logs",
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}
	messages, err := c.Consume(q.Name, "", true, false, false, false, nil)

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] [external_api, %s]\n", q.Name)
	<-forever

	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "new-films":
		//TODO get from api some info
		fmt.Println("successfully consumed info from rabbitMQ!!!!!!!!")
	}
}

func NewEventConsumer(conn *amqp.Connection) Consumer {
	return Consumer{conn: conn}
}
