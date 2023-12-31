package helpers

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func ToBSON(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return nil, err
	}
	err = bson.Unmarshal(data, &doc)
	return
}

type Message struct {
	MessageStatus string `json:"status,omitempty"`
	Message       string `json:"message,omitempty"`
	Limit         int    `json:"limit,omitempty"`
	Page          int    `json:"page,omitempty"`
	Data          any    `json:"data,omitempty"`
}

type ContextKey string

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func InternalServerError(w http.ResponseWriter, err string) {
	WriteJSON(w, http.StatusInternalServerError, Message{
		MessageStatus: "error",
		Message:       err,
	})
}

func StatusUnauthorized(w http.ResponseWriter, err string) {
	WriteJSON(w, http.StatusUnauthorized, Message{
		MessageStatus: "error",
		Message:       err,
	})
}

func StatusBadRequest(w http.ResponseWriter, err string) {
	WriteJSON(w, http.StatusBadRequest, Message{
		MessageStatus: "error",
		Message:       err,
	})
}

func StatusOk(w http.ResponseWriter, data any) {
	WriteJSON(w, http.StatusOK, Message{
		MessageStatus: "success",
		Data:          data,
	})
}
func StatusOkAll(w http.ResponseWriter, limit, page int, data any) {
	WriteJSON(w, http.StatusOK, Message{
		MessageStatus: "success",
		Limit:         limit,
		Page:          page,
		Data:          data,
	})
}
func StatusCreated(w http.ResponseWriter, data any) {
	WriteJSON(w, http.StatusCreated, Message{
		MessageStatus: "success",
		Data:          data,
	})
}
func StatusAccepted(w http.ResponseWriter, message string) {
	WriteJSON(w, http.StatusAccepted, Message{
		MessageStatus: "success",
		Message:       message,
	})
}
func StatusAcceptedData(w http.ResponseWriter, data any) {
	WriteJSON(w, http.StatusAccepted, Message{
		MessageStatus: "success",
		Data:          data,
	})
}
