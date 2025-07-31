package main

import (
	"log"
	"net/http"

	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/pkg/storage"
	"github.com/anglesson/simple-web-server/pkg/utils"

	handler "github.com/anglesson/simple-web-server/internal/handler"
	"github.com/anglesson/simple-web-server/internal/handler/web"
	"github.com/anglesson/simple-web-server/internal/repository/gorm"
	"github.com/anglesson/simple-web-server/internal/service"
	"github.com/anglesson/simple-web-server/pkg/gov"

	"github.com/anglesson/simple-web-server/internal/config"
	"github.com/anglesson/simple-web-server/internal/handler/middleware"
	"github.com/anglesson/simple-web-server/pkg/database"
	"github.com/anglesson/simple-web-server/pkg/template"
	"github.com/go-chi/chi/v5"
)

func main() {
	// ========== Infrastructure Initialization ==========
	config.LoadConfigs()
	database.Connect()

	flashServiceFactory := func(w http.ResponseWriter, r *http.Request) web.FlashMessagePort {
		return web.NewCookieFlashMessage(w, r)
	}

	// Template renderer
	templateRenderer := template.DefaultTemplateRenderer()

	// Utils
	encrypter := utils.NewEncrypter()

	// Repositories
	creatorRepository := gorm.NewCreatorRepository(database.DB)
	clientRepository := gorm.NewClientGormRepository()
	userRepository := repository.NewGormUserRepository(database.DB)
	fileRepository := repository.NewGormFileRepository(database.DB)

	// Services
	commonRFService := gov.NewHubDevService()
	userService := service.NewUserService(userRepository, encrypter)
	sessionService := service.NewSessionService()
	subscriptionRepository := gorm.NewSubscriptionGormRepository()
	subscriptionService := service.NewSubscriptionService(subscriptionRepository, commonRFService)
	stripeService := service.NewStripeService()
	paymentGateway := service.NewStripePaymentGateway(stripeService)
	creatorService := service.NewCreatorService(creatorRepository, commonRFService, userService, subscriptionService, paymentGateway)
	clientService := service.NewClientService(clientRepository, creatorRepository, commonRFService)
	s3Storage := storage.NewS3Storage()
	fileService := service.NewFileService(fileRepository, s3Storage)
	ebookService := service.NewEbookService(s3Storage)

	// Handlers
	authHandler := handler.NewAuthHandler(userService, sessionService, templateRenderer)
	clientHandler := handler.NewClientHandler(clientService, creatorService, flashServiceFactory, templateRenderer)
	creatorHandler := handler.NewCreatorHandler(creatorService, sessionService, templateRenderer)
	settingsHandler := handler.NewSettingsHandler(sessionService, templateRenderer)
	fileHandler := handler.NewFileHandler(fileService, sessionService, templateRenderer)
	ebookHandler := handler.NewEbookHandler(ebookService, creatorService, fileService, s3Storage, flashServiceFactory, templateRenderer)
	salesPageHandler := handler.NewSalesPageHandler(ebookService, creatorService, templateRenderer)
	dashboardHandler := handler.NewDashboardHandler(templateRenderer)
	errorHandler := handler.NewErrorHandler(templateRenderer)
	homeHandler := handler.NewHomeHandler(templateRenderer, errorHandler)
	forgetPasswordHandler := handler.NewForgetPasswordHandler(templateRenderer)
	sendHandler := handler.NewSendHandler(templateRenderer)
	purchaseHandler := handler.NewPurchaseHandler(templateRenderer)

	r := chi.NewRouter()

	fs := http.FileServer(http.Dir("web/assets"))
	r.Get("/assets/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/assets/", fs).ServeHTTP(w, r)
	})

	// Public routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthGuard)
		r.Get("/login", authHandler.LoginView)
		r.Post("/login", authHandler.LoginSubmit)
		r.Get("/register", creatorHandler.RegisterView)
		r.Post("/register", creatorHandler.RegisterCreatorSSR)
		r.Get("/forget-password", forgetPasswordHandler.ForgetPasswordView)
		r.Post("/forget-password", forgetPasswordHandler.ForgetPasswordSubmit)
		r.Get("/sales/{slug}", salesPageHandler.SalesPageView) // Página de vendas pública
	})

	// Completely public routes (no middleware)
	r.Get("/purchase/download/{id}", purchaseHandler.PurchaseDownloadHandler)

	// Stripe routes
	r.Post("/api/create-checkout-session", handler.CreateCheckoutSession)
	r.Post("/api/webhook", handler.HandleStripeWebhook)
	r.Post("/api/watermark", handler.WatermarkHandler)

	// Private routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Use(middleware.TrialMiddleware)

		r.Post("/logout", authHandler.LogoutSubmit)
		r.Get("/dashboard", dashboardHandler.DashboardView)
		r.Get("/settings", settingsHandler.SettingsView)

		// Ebook routes
		r.Get("/ebook", ebookHandler.IndexView)
		r.Get("/ebook/create", ebookHandler.CreateView)
		r.Post("/ebook/create", ebookHandler.CreateSubmit)
		r.Get("/ebook/edit/{id}", ebookHandler.UpdateView)
		r.Get("/ebook/view/{id}", ebookHandler.ShowView)
		r.Post("/ebook/update/{id}", ebookHandler.UpdateSubmit)
		r.Get("/ebook/preview/{id}", salesPageHandler.SalesPagePreviewView) // Preview da página de vendas

		// File routes
		r.Get("/file", fileHandler.FileIndexView)
		r.Get("/file/upload", fileHandler.FileUploadView)
		r.Post("/file/upload", fileHandler.FileUploadSubmit)
		r.Post("/file/{id}/update", fileHandler.FileUpdateSubmit)
		r.Post("/file/{id}/delete", fileHandler.FileDeleteSubmit)

		// Client routes
		r.Get("/client", clientHandler.ClientIndexView)
		r.Get("/client/new", clientHandler.CreateView)
		r.Post("/client", clientHandler.ClientCreateSubmit)
		r.Get("/client/update/{id}", clientHandler.UpdateView)
		r.Post("/client/update/{id}", clientHandler.ClientUpdateSubmit)
		r.Post("/client/import", clientHandler.ClientImportSubmit)

		// Purchase routes
		r.Post("/purchase/ebook/{id}", purchaseHandler.PurchaseCreateHandler)
		r.Get("/send", sendHandler.SendViewHandler)
	})

	r.Get("/", homeHandler.HomeView) // Home page deve ser a ultima rota

	port := config.AppConfig.Port

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	log.Println("Starting server on :" + port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
