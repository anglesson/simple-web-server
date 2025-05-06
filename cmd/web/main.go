package main

import (
	"log"
	"net/http"

	"github.com/anglesson/simple-web-server/config"
	"github.com/anglesson/simple-web-server/internal/handlers"
	"github.com/anglesson/simple-web-server/internal/shared/database"
	"github.com/anglesson/simple-web-server/internal/shared/middlewares"
	"github.com/go-chi/chi/v5"
)

func main() {
	config.LoadConfigs()
	database.Connect()

	r := chi.NewRouter()

	fs := http.FileServer(http.Dir("web/templates/assets"))
	r.Get("/assets/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/assets/", fs).ServeHTTP(w, r)
	})

	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)
		r.Get("/login", handlers.LoginHandler)
		r.Post("/login", handlers.LoginHandler)
		r.Get("/register", handlers.RegisterHandler)
		r.Post("/register", handlers.RegisterHandler)
		r.Get("/forget-password", handlers.ForgetPasswordHandler)
		r.Get("/dashboard", handlers.DashboardHandler)
		r.Get("/ebook", handlers.IndexHandler)
		r.Get("/ebook/create", handlers.CreateHandler)
		r.Post("/ebook/create", handlers.CreateHandler)
		r.Get("/ebook/edit/{id}", handlers.EditEbookHandler)
		r.Post("/ebook/update/{id}", handlers.EditEbookHandler)
	})

	r.Get("/", handlers.HomeHandler) // Home page deve ser a ultima rota

	port := config.AppConfig.Port

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	log.Println("Starting server on :" + port)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
