package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ishanshre/GoRestAPIMongoDB/internals/helpers"
	"github.com/ishanshre/GoRestAPIMongoDB/internals/models"
	"github.com/ishanshre/GoRestAPIMongoDB/utils"
)

func (h *handler) UserLogin(w http.ResponseWriter, r *http.Request) {
	loginUser := &models.CreateUser{}
	if err := json.NewDecoder(r.Body).Decode(&loginUser); err != nil {
		helpers.StatusBadRequest(w, err.Error())
		return
	}
	if err := h.MG.UsernameExists(loginUser.Username); err == nil {
		helpers.StatusBadRequest(w, "username does not exists")
		return
	}
	user, err := h.MG.GetUserByUsername(loginUser.Username)
	if err != nil {
		helpers.InternalServerError(w, err.Error())
		return
	}
	if err := utils.CompareHashPassword(user.Password, loginUser.Password); err != nil {
		helpers.StatusBadRequest(w, "invalid username/password")
		return
	}
	exists, err := h.RedisClient.Exists(context.Background(), user.Username).Result()
	if err != nil {
		helpers.InternalServerError(w, err.Error())
		return
	}
	if exists == 1 {
		if err := h.RedisClient.Del(context.Background(), user.Username).Err(); err != nil {
			helpers.InternalServerError(w, err.Error())
			return
		}
	}
	loginResponse, token, err := utils.GenerateLoginResponse(user.ID.Hex(), user.Username)
	if err != nil {
		helpers.InternalServerError(w, err.Error())
		return
	}
	tokenJson, err := json.Marshal(token)
	if err != nil {
		helpers.InternalServerError(w, err.Error())
		return
	}
	if err := h.RedisClient.Set(context.Background(), token.AccessToken.Username, tokenJson, time.Until(token.RefreshToken.ExpiresAt)).Err(); err != nil {
		helpers.InternalServerError(w, err.Error())
		return
	}
	helpers.StatusAcceptedData(w, loginResponse)

}

const (
	tokenDetailKey helpers.ContextKey = "tokenDetail"
)

func (h *handler) UserLogout(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Context().Value(tokenDetailKey))
	tokenDetail := r.Context().Value(tokenDetailKey).(*utils.TokenDetail)
	username := tokenDetail.Username
	if username == "" {
		helpers.StatusBadRequest(w, "Not Authorized")
		return
	}
	exists, err := h.RedisClient.Exists(context.Background(), username).Result()
	if err != nil {
		helpers.InternalServerError(w, err.Error())
		return
	}
	if exists == 1 {
		if err := h.RedisClient.Del(context.Background(), username).Err(); err != nil {
			helpers.InternalServerError(w, err.Error())
			return
		}
	}
	helpers.StatusOk(w, "Logout Success")
}
