package handlers

import "net/http"

func (app brokerHandler) GetNewFilmsCollection(w http.ResponseWriter, r *http.Request) {
	//TODO push message to queue

	//TODO get response from external-api-service and write back
}
