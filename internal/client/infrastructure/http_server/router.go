package client_http

import (
	"net/http"

	client_application "github.com/anglesson/simple-web-server/internal/client/application"
	"github.com/go-chi/chi/v5"
)

func NewRouter(useCase client_application.ClientUseCasePort) *chi.Mux {
	r := chi.NewRouter()

	fs := http.FileServer(http.Dir("web/assets"))
	r.Get("/assets/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/assets/", fs).ServeHTTP(w, r)
	})

	handler := NewClientSSRHandler(useCase)
	r.Post("/client", handler.CreateClient)

	return r
}
