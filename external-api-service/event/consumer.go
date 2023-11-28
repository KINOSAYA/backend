package event

import (
	"encoding/json"
	"external-api-service/internal/services"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)

type Consumer struct {
	conn *amqp.Connection
}

type Payload struct {
	Name string `json:"name"`
	Data struct {
		Language   string `json:"language"`
		TimeWindow string `json:"timeWindow"`
	} `json:"data,omitempty"`
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
			err := json.Unmarshal(d.Body, &payload)
			if err != nil {
				log.Println("consumer:57 ", err)
				continue
			}
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
		log.Println(os.Getenv("Bearer to TMDB"))
		services.GetNewFilms(os.Getenv("Bearer to TMDB"), payload.Data.Language, payload.Data.TimeWindow)
	default:
		fmt.Printf("payload.Name %s\n", payload.Name)
	}
}

func NewEventConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}
	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, err
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs",  // name
		"topic", //type
		true,    // durable
		false,   // autoDelete
		false,   // internal
		false,   // no-wait
		nil,     //arguments
	)
}
