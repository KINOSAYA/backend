package handlers

import (
	"broker-service/event"
	"broker-service/internal/helpers"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/http"
	"strconv"
)

// GetCollections is an API endpoint that triggers the retrieval of a new collection of films from the Kinopoisk API.
// @Tags Collections
// @Summary Get Collection
// @Description This endpoint communicates with the Kinopoisk API to fetch JSON data containing information about new films.
// @Produce json
// @Param page query string day "Page of results"
// @Param limit query string day "Limit records for one page"
// @Param category query string day "Category e.g. ('Фильмы', 'Сериалы')"
// @Success 200 {object} jsonResponse "Successful request"
// @Router /collections/new-films [get]
func (app *brokerHandler) GetCollections(w http.ResponseWriter, r *http.Request) {
	//TODO push message to queue
	queryParams := r.URL.Query()
	page, _ := strconv.Atoi(queryParams.Get("page"))
	limit, _ := strconv.Atoi(queryParams.Get("limit"))
	payload := struct {
		Name string `json:"name"`
		Data any
	}{
		Name: "new-films",
		Data: struct {
			Page     int    `json:"page"`
			Limit    int    `json:"limit"`
			Category string `json:"category"`
		}{
			Page:     page,
			Limit:    limit,
			Category: queryParams.Get("category"),
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
