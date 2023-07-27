package handlers

import (
	"encoding/json"
	"net/http"

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
	_, token, err := utils.GenerateLoginResponse(user.ID.Hex(), user.Username)
	if err != nil {
		helpers.InternalServerError(w, err.Error())
		return
	}
	helpers.StatusAcceptedData(w, token)

}
