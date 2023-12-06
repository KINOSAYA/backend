package handlers

import (
	"broker-service/event"
	"broker-service/internal/helpers"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"strconv"
)

// GetSlugsByCollection is an API endpoint that triggers the retrieval of a slugs of films from the Kinopoisk API.
// @Tags Kinopoisk
// @Summary Get Collection
// @Description This endpoint communicates with the Kinopoisk API to fetch JSON data containing information about slugs by chosen category.
// @Produce json
// @Param page query string day "Page of results"
// @Param limit query string day "Limit records for one page"
// @Param category query string day "Category e.g. ('Фильмы', 'Сериалы')"
// @Success 200 {object} jsonResponse "Successful request"
// @Router /kinopoisk/slugs [get]
func (app *brokerHandler) GetSlugsByCollection(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	page, _ := strconv.Atoi(queryParams.Get("page"))
	limit, _ := strconv.Atoi(queryParams.Get("limit"))
	payload := struct {
		Name string `json:"name"`
		Data any
	}{
		Name: "get-slugs",
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

	response, err := app.EventService.SendRequestAndWaitForResponse(payload)
	if err != nil {
		_ = helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		_ = helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	//_ = helpers.WriteJSON(w, http.StatusOK, rsp)

}

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}
type RequestForRabbitPayload struct {
	Name       string     `json:"name"`
	Pagination Pagination `json:"pagination,omitempty"`
	Data       any        `json:"data"`
}

// GetFilmsBySlug is an API endpoint that retrieves films based on the provided slug.
// @Tags Kinopoisk
// @Summary Get Films by Slug
// @Description Retrieve films based on the provided slug.
// @Produce json
// @Param slug query string true "The slug of the film"
// @Param page query string false "Page number for pagination"
// @Param limit query string false "Number of items per page"
// @Success 200 {string} json "Successful response"
// @Failure 400 {string} json "Bad request"
// @Router /kinopoisk/films [get]
func (app *brokerHandler) GetFilmsBySlug(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	page, err := strconv.Atoi(queryParams.Get("page"))
	if err != nil {
		_ = helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil {
		_ = helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	payload := RequestForRabbitPayload{
		Name: "get-films",
		Pagination: Pagination{
			Page:  page,
			Limit: limit,
		},
		Data: struct {
			Slug string `json:"slug"`
		}{
			queryParams.Get("slug"),
		},
	}
	log.Println(payload)
	log.Println(queryParams.Get("slug"))
	response, err := app.EventService.SendRequestAndWaitForResponse(payload)
	if err != nil {
		_ = helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		_ = helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
}

func pushToQueue(conn *amqp.Connection, payload any) error {
	emitter := event.NewEventEmitter(conn)

	err := emitter.Push(conn, payload)
	if err != nil {
		return err
	}

	return nil
}
