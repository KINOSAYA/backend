package handlers

import (
	"broker-service/event"
	"broker-service/internal/helpers"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/http"
)

func (app brokerHandler) GetNewFilmsCollection(w http.ResponseWriter, r *http.Request) {
	var requestPayload requestPayload
	err := helpers.ReadJSON(w, r, &requestPayload)
	if err != nil {
		return
	}
	//TODO push message to queue
	err = pushToQueue("get-new-films", nil, app.Rabbit)
	if err != nil {
		_ = helpers.ErrorJSON(w, err)
		return
	}

	//TODO get response from external-api-service and write back
}

func pushToQueue(name string, data any, conn *amqp.Connection) error {
	emitter := event.NewEventEmitter(conn)

	err := emitter.Push(conn)
	if err != nil {
		return err
	}

	return nil
}
