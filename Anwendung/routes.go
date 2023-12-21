package Anwendung

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func loadRoutes(userHandler *User) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.Route("/product", func(r chi.Router) {
		loadUserRoutes(r, userHandler)
	})

	return router
}

func loadUserRoutes(router chi.Router, userHandler *User) {
	router.Post("/", userHandler.Create)
	router.Get("/", userHandler.List)
	router.Get("/{id}", userHandler.GetByID)
	router.Put("/{id}", userHandler.UpdateByID)
	router.Delete("/{id}", userHandler.DeleteByID)
}
