package client

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes registra todas as rotas relacionadas ao módulo client
func RegisterRoutes(r chi.Router, handler *ClientHandler, authMiddleware, trialMiddleware, subscriptionMiddleware func(next http.Handler) http.Handler) {
	// Client routes - todas protegidas por autenticação
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Use(trialMiddleware)
		r.Use(subscriptionMiddleware)

		r.Get("/client", handler.ClientIndexView)
		r.Get("/client/new", handler.CreateView)
		r.Post("/client", handler.ClientCreateSubmit)
		r.Get("/client/update/{id}", handler.UpdateView)
		r.Post("/client/update/{id}", handler.ClientUpdateSubmit)
		r.Post("/client/import", handler.ClientImportSubmit)
	})
}
