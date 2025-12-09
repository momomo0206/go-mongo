package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Response struct {
	Msg  string
	Code int
}

func CreateRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Get("/healthcheck", healthCheck)
			r.Get("/todos", getTodos)
			r.Get("/todos/{id}", getTodoById)
			r.Post("/todos/create", createTodo)
			r.Put("/todos/update/{id}", updateTodo)
			r.Delete("/todos/delete/{id}", deleteTodo)
		})
	})

	return r
}
