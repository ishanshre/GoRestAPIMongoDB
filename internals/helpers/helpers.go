package helpers

import (
	"encoding/json"
	"net/http"
)

type Message struct {
	MessageStatus string `json:"status,omitempty"`
	Message       string `json:"message,omitempty"`
	Data          string `json:"data,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
