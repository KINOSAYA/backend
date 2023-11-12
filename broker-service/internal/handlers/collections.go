package handlers

import (
	"broker-service/event"
	"broker-service/internal/helpers"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/http"
)

// GetNewFilmsCollection is an API endpoint that gets a new collection of films from kinopoisk API.
// @Tags Collections
// @Summary New Collections
// @Description this endpoint goes to kinopoisk API and gets a JSON data with new films
// @Produce json
// @Success 200 {object} jsonResponse "Successful registration"
// @Router /collections/new-films [get]
func (app *brokerHandler) GetNewFilmsCollection(w http.ResponseWriter, r *http.Request) {
	//TODO push message to queue
	err := pushToQueue("get-new-films", nil, app.Rabbit)
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
