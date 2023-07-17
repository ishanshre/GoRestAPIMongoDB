package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/ishanshre/GoRestAPIMongoDB/internals/handlers"
)

func Router(h handlers.Handlers) http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Get("/", h.GetUsers)
	return mux
}
