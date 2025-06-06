package http_server

import (
	"net/http"

	"github.com/anglesson/simple-web-server/internal/application"
	"github.com/anglesson/simple-web-server/internal/infrastructure/http_server/handlers"
	"github.com/go-chi/chi/v5"
)

func NewRouter(useCase application.ClientUseCasePort) *chi.Mux {
	r := chi.NewRouter()

	fs := http.FileServer(http.Dir("web/assets"))
	r.Get("/assets/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/assets/", fs).ServeHTTP(w, r)
	})

	handler := handlers.NewClientSSRHandler(useCase)
	r.Post("/client", handler.CreateClient)

	return r
}
