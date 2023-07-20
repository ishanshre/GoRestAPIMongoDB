package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ishanshre/GoRestAPIMongoDB/internals/database"
	"github.com/ishanshre/GoRestAPIMongoDB/internals/helpers"
	"github.com/ishanshre/GoRestAPIMongoDB/internals/models"
	"github.com/ishanshre/GoRestAPIMongoDB/internals/repository"
	"github.com/ishanshre/GoRestAPIMongoDB/internals/repository/dbrepo"
	"github.com/ishanshre/GoRestAPIMongoDB/utils"
)

type Handlers interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUsers(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	GetUserByUsername(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	MG repository.MongoDbRepo
}

func NewHandler(mg database.DbInterface) Handlers {
	return &handler{
		MG: dbrepo.NewMongoDbRepo(mg, context.Background()),
	}
}

func (h *handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 0
	}

	users, err := h.MG.GetAllUsers(page, limit)
	if err != nil {
		helpers.InternalServerError(w, err.Error())
		return
	}
	helpers.StatusOkAll(w, limit, page, users)
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		helpers.StatusBadRequest(w, "error in parsing json")
		return
	}
	if err := h.MG.UsernameExists(user.Username); err != nil {
		helpers.StatusBadRequest(w, err.Error())
		return
	}
	hashedPassword, err := utils.GeneratePasswordHash(user.Password)
	if err != nil {
		helpers.InternalServerError(w, "cannot generate hash password")
		return
	}
	user.Password = hashedPassword
	user, err = h.MG.CreateUser(user)
	if err != nil {
		helpers.InternalServerError(w, err.Error())
		return
	}
	helpers.StatusCreated(w, user)
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if err := h.MG.DeleteUser(username); err != nil {
		helpers.InternalServerError(w, err.Error())
		return
	}
	helpers.StatusAccepted(w, fmt.Sprintf("%s deleted", username))
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	updateObj := &models.User{}
	if err := json.NewDecoder(r.Body).Decode(&updateObj); err != nil {
		helpers.StatusBadRequest(w, err.Error())
		return
	}
	user, err := h.MG.UpdateUser(username, updateObj)
	if err != nil {
		helpers.InternalServerError(w, err.Error())
		return
	}
	helpers.StatusAcceptedData(w, user)
}

func (h *handler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	user, err := h.MG.GetUserByUsername(username)
	if err != nil {
		helpers.InternalServerError(w, err.Error())
		return
	}
	helpers.StatusOk(w, user)
}
