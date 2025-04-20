package main

import (
	"log"
	"net/http"
	"os"

	"github.com/anglesson/simple-web-server/handlers"
	"github.com/anglesson/simple-web-server/middlewares"
	"github.com/anglesson/simple-web-server/templates"
)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	templates.View(w, "404-error", nil)
}

func main() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("templates/assets"))
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", fs))

	mux.HandleFunc("GET /login", handlers.LoginGetHandler)
	mux.HandleFunc("POST /login", handlers.LoginPostHandler)
	mux.HandleFunc("GET /register", handlers.RegisterGetHandler)
	mux.HandleFunc("GET /forget-password", handlers.ForgetPasswordGetHandler)
	mux.Handle("GET /dashboard", middlewares.AuthMiddleware(http.HandlerFunc(handlers.DashboardGetHandler)))

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			notFoundHandler(w, r)
			return
		}
		handlers.HomeGetHandler(w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Println("Starting server on :" + port)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
