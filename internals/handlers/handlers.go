package handlers

import (
	"net/http"

	"github.com/ishanshre/GoRestAPIMongoDB/internals/database"
	"github.com/ishanshre/GoRestAPIMongoDB/internals/helpers"
)

type Handlers interface {
	GetUsers(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	MG *database.DB
}

func NewHandler(mg *database.DB) Handlers {
	return &handler{
		MG: mg,
	}
}

func (h *handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	helpers.WriteJSON(w, http.StatusOK, helpers.Message{
		MessageStatus: "success",
		Message:       "working fine",
	})
}
