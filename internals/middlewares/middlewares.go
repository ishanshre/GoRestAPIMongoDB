package middlewares

import (
	"net/http"

	"github.com/redis/go-redis/v9"
)

type Middlewares interface {
	JwtAuth(next http.Handler) http.Handler
}

type middlewares struct {
	redisClient *redis.Client
}

func NewMiddleware(r *redis.Client) Middlewares {
	return &middlewares{
		redisClient: r,
	}
}
