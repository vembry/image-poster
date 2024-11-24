package handlers

import (
	"encoding/json"
	"net/http"
)

// respondJson is a generic handler to return json response to api requester
func respondJson[T any](w http.ResponseWriter, httpstatusCode int, body T) {
	// construct response
	raw, _ := json.Marshal(response[T]{
		Data: body,
	})

	// write to response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpstatusCode)
	w.Write(raw)
}

// respondErrorJson is a generic handler to return error json response to api requester
func respondErrorJson(w http.ResponseWriter, httpstatusCode int, message string) {
	// construct body to json's raw message
	raw, _ := json.Marshal(response[struct{}]{
		Error: message,
	})

	// write to response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpstatusCode)
	w.Write(raw)
}

// response define the structure of http api response
type response[T any] struct {
	Error string `json:"error"`
	Data  T      `json:"data"`
}
