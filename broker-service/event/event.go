package event

import amqp "github.com/rabbitmq/amqp091-go"

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      //type
		true,         // durable
		false,        // autoDelete
		false,        // internal
		false,        // no-wait
		nil,          //arguments
	)
}
