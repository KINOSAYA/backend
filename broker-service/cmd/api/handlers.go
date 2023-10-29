package main

import (
	_ "broker-service/cmd/api/docs"
	"net/http"
)

// broker is a sample API endpoint that returns a JSON response.
// @Summary hit the broker
// @Description Returns a JSON response with success status.
// @ID get-sample-response
// @Produce json
// @Success 200 {object} jsonResponse
// @Router / [get]
func (app *Config) broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Successfully hit the broker",
	}
	app.writeJSON(w, http.StatusOK, payload)
}
