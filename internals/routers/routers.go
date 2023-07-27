package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/ishanshre/GoRestAPIMongoDB/internals/handlers"
	"github.com/ishanshre/GoRestAPIMongoDB/internals/middlewares"
)

func Router(h handlers.Handlers, m middlewares.Middlewares) http.Handler {
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

	mux.Post("/", h.CreateUser)
	mux.Post("/login", h.UserLogin)
	mux.Group(func(mux chi.Router) {
		mux.Use(m.JwtAuth)
		mux.Get("/", h.GetUsers)
		mux.Delete("/users/{username}", h.DeleteUser)
		mux.Post("/logout", h.UserLogout)
		mux.Put("/users/{username}", h.UpdateUser)
		mux.Get("/users/{username}", h.GetUserByUsername)
	})
	return mux
}
