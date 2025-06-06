package main

import (
	"log"
	"net/http"

	"github.com/anglesson/simple-web-server/internal/application"
	common_infrastructure "github.com/anglesson/simple-web-server/internal/common/infrastructure"
	"github.com/anglesson/simple-web-server/internal/config"
	"github.com/anglesson/simple-web-server/internal/infrastructure/http_server"
	"github.com/anglesson/simple-web-server/internal/infrastructure/persistence"
	"github.com/anglesson/simple-web-server/internal/shared/database"
	"github.com/anglesson/simple-web-server/internal/shared/middlewares"
	"github.com/go-chi/chi/v5"
)

func main() {
	// Load configurations
	config.LoadConfigs()

	// Initialize database
	database.Connect()

	// Create router
	r := chi.NewRouter()

	// Serve static files
	fs := http.FileServer(http.Dir("web/assets"))
	r.Get("/assets/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/assets/", fs).ServeHTTP(w, r)
	})

	// Public routes
	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthGuard)
		r.Get("/login", http_server.LoginView)
		r.Post("/login", http_server.LoginSubmit)
		r.Get("/register", http_server.RegisterView)
		r.Post("/register", http_server.RegisterSubmit)
		r.Get("/forget-password", http_server.ForgetPasswordView)
		r.Post("/forget-password", http_server.ForgetPasswordSubmit)
		r.Get("/purchase/download/{id}", http_server.PurchaseDownloadHandler)
	})

	// Stripe routes
	r.Post("/api/create-checkout-session", http_server.CreateCheckoutSession)
	r.Post("/api/webhook", http_server.HandleStripeWebhook)

	// Private routes
	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)
		r.Use(middlewares.TrialMiddleware)

		r.Post("/logout", http_server.LogoutSubmit)
		r.Get("/dashboard", http_server.DashboardView)
		r.Get("/settings", http_server.SettingsView)

		// Ebook routes
		r.Get("/ebook", http_server.EbookIndexView)
		r.Get("/ebook/create", http_server.EbookCreateView)
		r.Post("/ebook/create", http_server.EbookCreateSubmit)
		r.Get("/ebook/edit/{id}", http_server.EbookUpdateView)
		r.Get("/ebook/view/{id}", http_server.EbookShowView)
		r.Post("/ebook/update/{id}", http_server.EbookUpdateSubmit)

		// Client routes
		clientHandler := http_server.NewClientSSRHandler(application.NewClientUseCase(
			persistence.NewClientRepository(),
			common_infrastructure.NewHubDevService(config.AppConfig.HubDesenvolvedorApi, config.AppConfig.HubDesenvolvedorToken),
		))
		r.Get("/client", clientHandler.ListClients)
		r.Post("/client", clientHandler.CreateClient)
		r.Post("/client/update/{id}", clientHandler.UpdateClient)
		r.Post("/client/import", clientHandler.ImportClients)

		// Purchase routes
		r.Post("/purchase/ebook/{id}", http_server.PurchaseCreateHandler)
		r.Get("/purchase/download/{id}", http_server.PurchaseDownloadHandler)
		r.Get("/send", http_server.SendViewHandler)
	})

	r.Get("/", http_server.HomeView) // Home page deve ser a ultima rota

	// Start server
	port := config.AppConfig.Port
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	log.Printf("Starting server on :%s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
