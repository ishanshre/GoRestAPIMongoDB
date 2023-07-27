package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	mux.Use(middleware.Logger)

	mux.Get("/", h.GetUsers)
	mux.Post("/", h.CreateUser)
	mux.Post("/login", h.UserLogin)
	mux.Delete("/users/{username}", h.DeleteUser)
	mux.Put("/users/{username}", h.UpdateUser)
	mux.Get("/users/{username}", h.GetUserByUsername)
	return mux
}
