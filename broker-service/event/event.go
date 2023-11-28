package event

import (
	"encoding/json"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"external_api", // name
		"topic",        //type
		true,           // durable
		false,          // autoDelete
		false,          // internal
		false,          // no-wait
		nil,            //arguments
	)
}

type Service interface {
	SendRequestAndWaitForResponse(payload any) (any, error)
}
type RabbitEventService struct {
	conn *amqp.Connection
}

func NewRabbitEventService(connection *amqp.Connection) Service {
	return &RabbitEventService{
		conn: connection,
	}
}

func (rabbit RabbitEventService) SendRequestAndWaitForResponse(payload any) (any, error) {
	c, err := rabbit.conn.Channel()
	if err != nil {
		log.Fatalf("channel.open: %s", err)
	}
	defer c.Close()

	// Создание временной очереди для ответов
	responseQueue, err := c.QueueDeclare("", false, true, true, false, nil)
	if err != nil {
		return nil, fmt.Errorf("queue.declare: %w", err)
	}

	// Привязка временной очереди ответов к обмену
	err = c.QueueBind(
		responseQueue.Name,
		responseQueue.Name,
		"logs",
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("queue.bind: %w", err)
	}

	// Определение сообщения для отправки
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("json.marshal: %w", err)
	}

	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "application/json",
		Body:         jsonBody,
		ReplyTo:      responseQueue.Name,
	}

	// Публикация сообщения
	err = c.Publish("logs", "info", false, false, msg)
	if err != nil {
		return nil, fmt.Errorf("basic.publish: %w", err)
	}

	log.Println("successfully sent to queue")
	// Получение ответа из временной очереди
	messages, err := c.Consume(
		responseQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("queue.consume: %w", err)
	}

	// Ожидание и обработка ответа
	for d := range messages {
		return d.Body, nil
	}

	// Если не удалось получить ответ
	return nil, errors.New("Timeout: No response received")

}
