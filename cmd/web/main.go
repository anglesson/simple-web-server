package main

import (
	"log"
	"net/http"

	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/pkg/utils"

	handler "github.com/anglesson/simple-web-server/internal/handler"
	"github.com/anglesson/simple-web-server/internal/handler/web"
	"github.com/anglesson/simple-web-server/internal/repository/gorm"
	"github.com/anglesson/simple-web-server/internal/service"
	"github.com/anglesson/simple-web-server/pkg/gov"

	"github.com/anglesson/simple-web-server/internal/config"
	"github.com/anglesson/simple-web-server/internal/handler/middleware"
	"github.com/anglesson/simple-web-server/pkg/database"
	"github.com/go-chi/chi/v5"
)

func main() {
	// ========== Infrastructure Initialization ==========
	config.LoadConfigs()
	database.Connect()

	flashServiceFactory := func(w http.ResponseWriter, r *http.Request) web.FlashMessagePort {
		return web.NewCookieFlashMessage(w, r)
	}

	// Utils
	encrypter := utils.NewEncrypter()

	// Repositories
	creatorRepository := gorm.NewCreatorRepository(database.DB)
	clientRepository := gorm.NewClientGormRepository()
	userRepository := repository.NewGormUserRepository(database.DB)

	// Services
	commonRFService := gov.NewHubDevService()
	userService := service.NewUserService(userRepository, encrypter)
	subscriptionRepository := gorm.NewSubscriptionGormRepository()
	subscriptionService := service.NewSubscriptionService(subscriptionRepository, commonRFService)
	stripeService := service.NewStripeService()
	paymentGateway := service.NewStripePaymentGateway(stripeService)
	creatorService := service.NewCreatorService(creatorRepository, commonRFService, userService, subscriptionService, paymentGateway)
	clientService := service.NewClientService(clientRepository, creatorRepository, commonRFService)

	// Handlers
	clientHandler := handler.NewClientHandler(clientService, creatorService, flashServiceFactory)
	creatorHandler := handler.NewCreatorHandler(creatorService)

	r := chi.NewRouter()

	fs := http.FileServer(http.Dir("web/assets"))
	r.Get("/assets/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/assets/", fs).ServeHTTP(w, r)
	})

	// Public routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthGuard)
		r.Get("/login", handler.LoginView)
		r.Post("/login", handler.LoginSubmit)
		r.Get("/register", handler.RegisterView)
		r.Post("/register", creatorHandler.RegisterCreatorSSR)
		r.Get("/forget-password", handler.ForgetPasswordView)
		r.Post("/forget-password", handler.ForgetPasswordSubmit)
		r.Get("/purchase/download/{id}", handler.PurchaseDownloadHandler)
	})

	// Stripe routes
	r.Post("/api/create-checkout-session", handler.CreateCheckoutSession)
	r.Post("/api/webhook", handler.HandleStripeWebhook)
	r.Post("/api/watermark", handler.WatermarkHandler)

	// Private routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Use(middleware.TrialMiddleware)

		r.Post("/logout", handler.LogoutSubmit)
		r.Get("/dashboard", handler.DashboardView)
		r.Get("/settings", handler.SettingsView)

		// Ebook routes
		r.Get("/ebook", handler.EbookIndexView)
		r.Get("/ebook/create", handler.EbookCreateView)
		r.Post("/ebook/create", handler.EbookCreateSubmit)
		r.Get("/ebook/edit/{id}", handler.EbookUpdateView)
		r.Get("/ebook/view/{id}", handler.EbookShowView)
		r.Post("/ebook/update/{id}", handler.EbookUpdateSubmit)

		// Client routes
		r.Get("/client", clientHandler.ClientIndexView)
		r.Get("/client/new", clientHandler.CreateView)
		r.Post("/client", clientHandler.ClientCreateSubmit)
		r.Get("/client/update/{id}", clientHandler.UpdateView)
		r.Post("/client/update/{id}", clientHandler.ClientUpdateSubmit)
		r.Post("/client/import", clientHandler.ClientImportSubmit)

		// Purchase routes
		r.Post("/purchase/ebook/{id}", handler.PurchaseCreateHandler)
		r.Get("/purchase/download/{id}", handler.PurchaseDownloadHandler)
		r.Get("/send", handler.SendViewHandler)
	})

	r.Get("/", handler.HomeView) // Home page deve ser a ultima rota

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
