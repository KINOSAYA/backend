package event

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type Emitter struct {
	connection *amqp.Connection
}

//
//func (e *Emitter) setup() error {
//	channel, err := e.connection.Channel()
//	if err != nil {
//		return err
//	}
//
//	defer channel.Close()
//	return declareExchange(channel)
//
//}

func (e Emitter) Push(conn *amqp.Connection) error {
	// This waits for a server acknowledgment which means the sockets will have
	//flushed all outbound publishings prior to returning.  It's important to
	//block on Close to not lose any publishings.
	c, err := conn.Channel()
	if err != nil {
		log.Fatalf("channel.open: %s", err)
	}

	// We declare our topology on both the publisher and consumer to ensure they
	//are the same.  This is part of AMQP being a programmable messaging model.
	//// See the Channel.Consume example for the complimentary declare.
	err = c.ExchangeDeclare("logs", "topic", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("exchange.declare: %v", err)
	}

	// Prepare this message to be persistent.
	//Your publishing requirements may
	//be different.
	msg := amqp.Publishing{DeliveryMode: amqp.Persistent, Timestamp: time.Now(), ContentType: "text/plain", Body: []byte("Go Go AMQP!")}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// This is not a mandatory delivery, so it will be dropped if there are no
	//queues bound to the logs exchange.
	err = c.PublishWithContext(ctx, "logs", "info", false, false, msg)
	if err != nil {
		// Since publish is asynchronous this can happen if the network connection
		//is reset or if the server has run out of resources.
		log.Fatalf("basic.publish: %v", err)
	}

	return nil
}

func NewEventEmitter(conn *amqp.Connection) Emitter {
	emitter := Emitter{
		connection: conn,
	}
	//err := emitter.setup()
	//if err != nil {
	//	return Emitter{}, err
	//}

	return emitter
}
