package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/anglesson/simple-web-server/internal/handler/web"
	"github.com/anglesson/simple-web-server/internal/repository/gorm"

	"github.com/anglesson/simple-web-server/internal/config"
	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/internal/service"
	"github.com/anglesson/simple-web-server/pkg/database"
	"github.com/anglesson/simple-web-server/pkg/gov"
	"github.com/anglesson/simple-web-server/pkg/mail"
	"github.com/anglesson/simple-web-server/pkg/template"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		RegisterView(w, r)
	case http.MethodPost:
		RegisterSubmit(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func RegisterView(w http.ResponseWriter, r *http.Request) {
	template.View(w, r, "creator/register", nil, "guest")
}

func RegisterSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	form := models.RegisterForm{
		Username:             r.FormValue("username"),
		Email:                r.FormValue("email"),
		Password:             r.FormValue("password"),
		PasswordConfirmation: r.FormValue("password_confirmation"),
	}

	errors := make(map[string]string)

	// Validate the input
	if form.Username == "" {
		errors["username"] = "Username é obrigatório."
	}

	if form.Email == "" {
		errors["email"] = "Email é obrigatório."
	}
	if form.Password == "" {
		errors["password"] = "Password é obrigatório."
	}
	if form.PasswordConfirmation == "" {
		errors["password_confirmation"] = "Confirmação de senha é obrigatório."
	}

	if form.Password != form.PasswordConfirmation {
		errors["password_confirmation"] = "Senhas não coincidem."
	}

	// Check if the user already exists
	foundedUser := repository.NewGormUserRepository(database.DB).FindByEmail(form.Email)
	if foundedUser != nil {
		errors["email"] = "Email já cadastrado"
	}

	if len(errors) > 0 {
		formJSON, _ := json.Marshal(form)
		errorsJSON, _ := json.Marshal(errors)

		http.SetCookie(w, &http.Cookie{
			Name:  "form",
			Value: url.QueryEscape(string(formJSON)),
			Path:  "/",
		})
		http.SetCookie(w, &http.Cookie{
			Name:  "errors",
			Value: url.QueryEscape(string(errorsJSON)),
			Path:  "/",
		})
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	hashedPassword := encrypter.HashPassword(form.Password)

	user := models.NewUser(form.Username, hashedPassword, form.Email)
	if err := repository.NewGormUserRepository(database.DB).Save(user); err != nil {
		web.RedirectBackWithErrors(w, r, err.Error())
		return
	}

	// Create creator with subscription and payment gateway integration
	creatorRepository := gorm.NewCreatorRepository(database.DB)
	userRepository := repository.NewGormUserRepository(database.DB)
	subscriptionRepository := gorm.NewSubscriptionGormRepository()
	stripeService := service.NewStripeService()

	creatorService := service.NewCreatorService(
		creatorRepository,
		gov.NewHubDevService(),
		service.NewUserService(userRepository, encrypter),
		service.NewSubscriptionService(subscriptionRepository, gov.NewHubDevService()),
		service.NewStripePaymentGateway(stripeService),
	)

	creatorInput := service.InputCreateCreator{
		Name:                 user.Username,
		CPF:                  "",           // Will be filled later
		BirthDate:            "1990-01-01", // Default date
		PhoneNumber:          "",
		Email:                user.Email,
		Password:             form.Password,
		PasswordConfirmation: form.PasswordConfirmation,
	}

	_, err := creatorService.CreateCreator(creatorInput)
	if err != nil {
		log.Printf("Error creating creator: %v", err)
		web.RedirectBackWithErrors(w, r, "Erro ao criar conta")
		return
	}

	sessionService.InitSession(w, user.Email)

	mailPort, _ := strconv.Atoi(config.AppConfig.MailPort)
	s := mail.NewEmailService(mail.NewGoMailer(
		config.AppConfig.MailHost,
		mailPort,
		config.AppConfig.MailUsername,
		config.AppConfig.MailPassword))

	go s.SendAccountConfirmation(user.Username, user.Email, "any")

	// Redirect to the protected area
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
