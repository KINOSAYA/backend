package event

import (
	"encoding/json"
	"external-api-service/internal/services"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Consumer struct {
	conn *amqp.Connection
}

type Payload struct {
	Name string `json:"name"`
	Data struct {
		Page     int    `json:"page"`
		Limit    int    `json:"limit"`
		Category string `json:"category,omitempty"`
		Slug     string `json:"slug,omitempty"`
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
			go handlePayload(payload, c, d)
		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] [external_api, %s]\n", q.Name)
	<-forever

	return nil
}

func sendMsgBack(payload []byte, ch *amqp.Channel, delivery amqp.Delivery) {
	err := ch.Publish(
		"",               // Обменник (пусто для обмена по умолчанию)
		delivery.ReplyTo, // Ответная очередь
		false,            // Опубликованное сообщение не сохраняется в хранилище
		false,            // Не устанавливать подтверждение доставки
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
		})
	if err != nil {
		log.Printf("cannot respond back %v\n", err)
	}
}

type ErrorJson struct {
	Error string `json:"error"`
}

func sendErrorBack(err error, ch *amqp.Channel, delivery amqp.Delivery) {
	res := ErrorJson{
		Error: err.Error(),
	}
	jsonRes, _ := json.Marshal(res)
	sendMsgBack(jsonRes, ch, delivery)
}
func handlePayload(payload Payload, ch *amqp.Channel, delivery amqp.Delivery) {
	switch payload.Name {
	case "get-slugs":
		responseMessage, err := services.GetSlugsByCategory(payload.Data.Page, payload.Data.Limit, payload.Data.Category)

		if err != nil {
			sendErrorBack(err, ch, delivery)
			return
		}

		// Отправляем ответ в ту же очередь, откуда пришел запрос
		sendMsgBack(responseMessage, ch, delivery)

	case "get-films":
		log.Println("got request for get-films by slugs")
		responseMessage, err := services.GetFilmsBySlug(payload.Data.Page, payload.Data.Limit, payload.Data.Slug)

		if err != nil {
			sendErrorBack(err, ch, delivery)
			return
		}

		// Отправляем ответ в ту же очередь, откуда пришел запрос
		sendMsgBack(responseMessage, ch, delivery)
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
