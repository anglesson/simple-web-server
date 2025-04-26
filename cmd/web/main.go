package main

import (
	"log"
	"net/http"
	"os"

	auth "github.com/anglesson/simple-web-server/internal/auth/handlers"
	dashboard "github.com/anglesson/simple-web-server/internal/dashboard/handlers"
	home "github.com/anglesson/simple-web-server/internal/home/handlers"
)

func main() {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("web/templates/assets"))
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", fs))

	mux.HandleFunc("GET /login", auth.LoginHandler)
	mux.HandleFunc("POST /login", auth.LoginHandler)
	mux.HandleFunc("GET /register", auth.RegisterHandler)
	mux.HandleFunc("POST /register", auth.RegisterHandler)
	mux.HandleFunc("GET /forget-password", auth.ForgetPasswordHandler)
	mux.Handle("GET /dashboard", http.HandlerFunc(dashboard.DashboardHandler))

	mux.HandleFunc("GET /", home.HomeHandler) // Home page deve ser a ultima rota

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
