package handlers

import (
	"encoding/json"
	"net/http"
)

// respondJson is a generic handler to return json response to api requester
func respondJson(w http.ResponseWriter, httpstatusCode int, body interface{}) {
	// transform body to json's raw message
	raw, _ := json.Marshal(body)

	// write to response

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpstatusCode)
	w.Write(raw)
}
