package main

import (
	"log"
	"net/http"

	"github.com/anglesson/simple-web-server/config"
	auth "github.com/anglesson/simple-web-server/internal/auth/handlers"
	dashboard "github.com/anglesson/simple-web-server/internal/dashboard/handlers"
	ebook "github.com/anglesson/simple-web-server/internal/ebook/handlers"
	home "github.com/anglesson/simple-web-server/internal/home/handlers"
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
		r.Get("/login", auth.LoginHandler)
		r.Post("/login", auth.LoginHandler)
		r.Get("/register", auth.RegisterHandler)
		r.Post("/register", auth.RegisterHandler)
		r.Get("/forget-password", auth.ForgetPasswordHandler)
		r.Get("/dashboard", dashboard.DashboardHandler)
		r.Get("/ebooks", ebook.IndexHandler)
	})

	r.Get("/", home.HomeHandler) // Home page deve ser a ultima rota

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
