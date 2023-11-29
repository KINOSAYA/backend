package handlers

import (
	"broker-service/event"
	"broker-service/internal/helpers"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/http"
)

// GetNewFilmsCollection is an API endpoint that triggers the retrieval of a new collection of films from the TMDB API.
// @Tags Collections
// @Summary Get New Films Collection
// @Description This endpoint communicates with the TMDB API to fetch JSON data containing information about new films.
// @Produce json
// @Param lan query string ru "Language code for the films (e.g., 'en' for English, 'ru' for Russian)"
// @Param time-window query string day "Time window for the new films (e.g., 'day', 'week')"
// @Success 200 {object} jsonResponse "Successful request"
// @Router /collections/new-films [get]
func (app *brokerHandler) GetNewFilmsCollection(w http.ResponseWriter, r *http.Request) {
	//TODO push message to queue
	queryParams := r.URL.Query()
	payload := struct {
		Name string `json:"name"`
		Data any
	}{
		Name: "new-films",
		Data: struct {
			Language   string `json:"language"`
			TimeWindow string `json:"timeWindow"`
		}{
			Language:   queryParams.Get("lan"),
			TimeWindow: queryParams.Get("time-window"),
		},
	}
	//err := pushToQueue(app.Rabbit, payload)
	//if err != nil {
	//	_ = helpers.ErrorJSON(w, err)
	//	return
	//}

	response, err := app.EventService.SendRequestAndWaitForResponse(payload)
	if err != nil {
		_ = helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		_ = helpers.ErrorJSON(w, err, http.StatusBadRequest)
	}
	//_ = helpers.WriteJSON(w, http.StatusOK, rsp)

	//TODO get response from external-api-service and write back
}

func pushToQueue(conn *amqp.Connection, payload any) error {
	emitter := event.NewEventEmitter(conn)

	err := emitter.Push(conn, payload)
	if err != nil {
		return err
	}

	return nil
}
