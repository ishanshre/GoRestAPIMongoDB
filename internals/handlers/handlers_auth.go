package handlers

import (
	"context"
	"encoding/json"
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

func (h *handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken := &models.RefreshToken{}
	if err := json.NewDecoder(r.Body).Decode(&refreshToken); err != nil {
		helpers.StatusBadRequest(w, "could not parse json data")
		return
	}
	tokenDetails, err := utils.VerifyTokenWithClaims(refreshToken.RefreshToken, "refresh_token")
	if err != nil {
		helpers.StatusBadRequest(w, err.Error())
		return
	}
	exists, err := h.RedisClient.Exists(context.Background(), tokenDetails.Username).Result()
	if err != nil {
		helpers.InternalServerError(w, err.Error())
		return
	}
	if exists == 1 {
		if err := h.RedisClient.Del(context.Background(), tokenDetails.Username).Err(); err != nil {
			helpers.InternalServerError(w, "error in deleting previous token")
			return
		}
	} else {
		helpers.StatusBadRequest(w, "token already revoked")
	}
	loginResponse, token, err := utils.GenerateLoginResponse(tokenDetails.UserId, tokenDetails.Username)
	if err != nil {
		helpers.InternalServerError(w, "token creating error")
		return
	}
	tokenJson, _ := json.Marshal(token)
	if err := h.RedisClient.Set(context.Background(), loginResponse.Username, tokenJson, time.Until(token.RefreshToken.ExpiresAt)).Err(); err != nil {
		helpers.InternalServerError(w, "error storing tokens")
		return
	}
	helpers.StatusOk(w, token)
}
