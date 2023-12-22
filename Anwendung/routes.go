package Anwendung

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func loadRoutes(kundeHandler *Kunde) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.Route("/kunde", func(r chi.Router) {
		loadKundeRoutes(r, kundeHandler)
	})

	return router
}

func loadKundeRoutes(router chi.Router, kundeHandler *Kunde) {
	router.Post("/", kundeHandler.Create)
	router.Get("/", kundeHandler.List)
	router.Get("/{id}", kundeHandler.GetByID)
	router.Put("/{id}", kundeHandler.UpdateByID)
	router.Delete("/{id}", kundeHandler.DeleteByID)
}
