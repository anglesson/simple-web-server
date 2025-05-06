package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/anglesson/simple-web-server/internal/config"
	"github.com/anglesson/simple-web-server/internal/mail"
	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repositories"
	"github.com/anglesson/simple-web-server/internal/shared/database"
	"github.com/anglesson/simple-web-server/internal/shared/template"
	"github.com/anglesson/simple-web-server/internal/shared/utils"
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
	template.View(w, r, "register", nil, "base_guest")
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
	foundedUser := repositories.FindByEmail(form.Email)
	if foundedUser == nil {
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

	hashedPassword := utils.HashPassword(form.Password)

	user := models.NewUser(form.Username, hashedPassword, form.Email)
	repositories.Save(user)
	creator := models.NewCreator(user.Username, user.Email, "", user.ID)

	database.DB.Save(&creator) // TODO: Create middleware

	sessionService.InitSession(w, form.Email)

	log.Printf("User registered ID: %d, EMAIL: %s", user.ID, user.Email)
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
