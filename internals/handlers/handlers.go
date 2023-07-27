package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/ishanshre/GoRestAPIMongoDB/internals/database"
	"github.com/ishanshre/GoRestAPIMongoDB/internals/helpers"
	"github.com/ishanshre/GoRestAPIMongoDB/internals/models"
	"github.com/ishanshre/GoRestAPIMongoDB/internals/repository"
	"github.com/ishanshre/GoRestAPIMongoDB/internals/repository/dbrepo"
	"github.com/ishanshre/GoRestAPIMongoDB/internals/validators"
	"github.com/ishanshre/GoRestAPIMongoDB/utils"
	"github.com/redis/go-redis/v9"
)

type Handlers interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUsers(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	GetUserByUsername(w http.ResponseWriter, r *http.Request)

	UserLogin(w http.ResponseWriter, r *http.Request)
	UserLogout(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
}

var validate *validator.Validate

type handler struct {
	MG          repository.MongoDbRepo
	RedisClient *redis.Client
}

func NewHandler(mg database.DbInterface, r *redis.Client) Handlers {
	validate = validator.New()
	validate.RegisterValidation("uppercase", validators.UpperCase)
	validate.RegisterValidation("lowercase", validators.LowerCase)
	validate.RegisterValidation("number", validators.Number)
	return &handler{
		MG:          dbrepo.NewMongoDbRepo(mg, context.Background()),
		RedisClient: r,
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
	newUser := &models.CreateUser{}
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		helpers.StatusBadRequest(w, "error in parsing json")
		return
	}
	if err := validate.Struct(newUser); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.Field()
			if err.Tag() == "containsany" {
				helpers.StatusBadRequest(w, fmt.Sprintf("%s must have at least one special characters from: %v", fieldName, err.Param()))
				return
			}
			helpers.StatusBadRequest(w, fmt.Sprintf("%s must have at least one %s %v characters", fieldName, err.Tag(), err.Param()))
		}
		return
	}
	if err := h.MG.UsernameExists(newUser.Username); err != nil {
		helpers.StatusBadRequest(w, err.Error())
		return
	}
	hashedPassword, err := utils.GeneratePasswordHash(newUser.Password)
	if err != nil {
		helpers.InternalServerError(w, "cannot generate hash password")
		return
	}
	user := &models.User{
		Username: newUser.Username,
		Password: hashedPassword,
	}
	getUser, err := h.MG.CreateUser(user)
	if err != nil {
		helpers.InternalServerError(w, err.Error())
		return
	}
	helpers.StatusCreated(w, getUser)
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	tokenDetail := r.Context().Value(tokenDetailKey).(*utils.TokenDetail)
	username := chi.URLParam(r, "username")

	if tokenDetail.Username != username {
		helpers.StatusUnauthorized(w, "You are not authorized to delete others")
		return
	}
	if err := h.MG.DeleteUser(username); err != nil {
		helpers.InternalServerError(w, err.Error())
		return
	}
	helpers.StatusAccepted(w, fmt.Sprintf("%s deleted", username))
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	tokenDetail := r.Context().Value(tokenDetailKey).(*utils.TokenDetail)
	username := chi.URLParam(r, "username")

	if tokenDetail.Username != username {
		helpers.StatusUnauthorized(w, "You are not authorized to delete others")
		return
	}
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
	tokenDetail := r.Context().Value(tokenDetailKey).(*utils.TokenDetail)
	username := chi.URLParam(r, "username")

	if tokenDetail.Username != username {
		helpers.StatusUnauthorized(w, "You are not authorized to delete others")
		return
	}
	user, err := h.MG.GetUserByUsername(username)
	if err != nil {
		helpers.InternalServerError(w, err.Error())
		return
	}
	helpers.StatusOk(w, user)
}
