package main

import (
	"log"
	"net/http"

	"github.com/anglesson/simple-web-server/internal/config"
	"github.com/anglesson/simple-web-server/internal/handlers"
	"github.com/anglesson/simple-web-server/internal/shared/database"
	"github.com/anglesson/simple-web-server/internal/shared/middlewares"
	"github.com/go-chi/chi/v5"
)

func main() {
	config.LoadConfigs()
	database.Connect()

	r := chi.NewRouter()

	fs := http.FileServer(http.Dir("web/assets"))
	r.Get("/assets/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/assets/", fs).ServeHTTP(w, r)
	})

	// Public routes
	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthGuard)
		r.Get("/login", handlers.LoginView)
		r.Post("/login", handlers.LoginSubmit)
		r.Get("/register", handlers.RegisterView)
		r.Post("/register", handlers.RegisterSubmit)
		r.Get("/forget-password", handlers.ForgetPasswordView)
		r.Post("/forget-password", handlers.ForgetPasswordSubmit)
	})

	// Private routes
	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)
		r.Post("/logout", handlers.LogoutSubmit)
		r.Get("/dashboard", handlers.DashboardView)
		r.Get("/ebook", handlers.EbookIndexView)
		r.Get("/ebook/create", handlers.EbookCreateView)
		r.Post("/ebook/create", handlers.EbookCreateSubmit)
		r.Get("/ebook/edit/{id}", handlers.EbookUpdateView)
		r.Post("/ebook/update/{id}", handlers.EbookUpdateSubmit)
		r.Get("/client", handlers.ClientIndexView)
		r.Post("/client", handlers.ClientCreateSubmit)
	})

	r.Get("/", handlers.HomeView) // Home page deve ser a ultima rota

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
