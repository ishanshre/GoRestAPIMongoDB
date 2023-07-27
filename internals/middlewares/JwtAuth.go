package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ishanshre/GoRestAPIMongoDB/internals/helpers"
	"github.com/ishanshre/GoRestAPIMongoDB/utils"
)

const (
	tokenDetailKey helpers.ContextKey = "tokenDetail"
)

func (m *middlewares) JwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")
		if bearerToken == "" {
			helpers.StatusUnauthorized(w, "access token not provided")
			return
		}
		tokenString := strings.Split(bearerToken, " ")
		if len(tokenString) != 2 && tokenString[0] != "Bearer" {
			helpers.StatusUnauthorized(w, "invalid token format")
			return
		}

		// verify the token
		tokenDetail, err := utils.VerifyTokenWithClaims(tokenString[1], "access_token")
		if err != nil {
			helpers.StatusUnauthorized(w, err.Error())
			return
		}
		exists, err := m.redisClient.Exists(context.Background(), tokenDetail.Username).Result()
		if err != nil {
			helpers.InternalServerError(w, err.Error())
			return
		}
		if exists == 0 {
			helpers.StatusBadRequest(w, "please use the latest token")
			return
		}
		data, err := m.redisClient.Get(context.Background(), tokenDetail.Username).Result()
		if err != nil {
			helpers.InternalServerError(w, err.Error())
			return
		}
		token := &utils.Token{}
		if err := json.Unmarshal([]byte(data), token); err != nil {
			helpers.InternalServerError(w, err.Error())
			return
		}
		if tokenDetail.TokenID != token.AccessToken.TokenID {
			helpers.StatusBadRequest(w, "please use the latest token")
			return
		}
		ctx := context.WithValue(r.Context(), tokenDetailKey, tokenDetail)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

	})
}
