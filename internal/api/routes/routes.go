package routes

import (
	"net/http"

	"vietnamese-converter/internal/api/handlers"
	"github.com/go-chi/chi/v5"
)

func SetupConvertRoutes(r *chi.Mux, convertHandler *handlers.ConvertHandler) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/convert", convertHandler.ConvertNumber)
		r.Get("/convert", convertHandler.ConvertFromURL)
	})
	
	r.Get("/health", convertHandler.HealthCheck)
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("pong"))
	})
}
