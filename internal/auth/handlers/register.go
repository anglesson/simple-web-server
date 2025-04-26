package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/anglesson/simple-web-server/internal/auth/models"
	"github.com/anglesson/simple-web-server/internal/auth/repositories"
	"github.com/anglesson/simple-web-server/internal/shared/template"
	"github.com/anglesson/simple-web-server/internal/shared/utils"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		renderRegisterPage(w, r)
	case http.MethodPost:
		processRegisterPage(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func renderRegisterPage(w http.ResponseWriter, r *http.Request) {
	var form models.RegisterForm
	var errors models.FormErrors

	if c, err := r.Cookie("form"); err == nil {
		decodedValue, decodeErr := url.QueryUnescape(c.Value) // Decodifica o valor do cookie
		if decodeErr != nil {
			log.Println("Error decoding cookie value:", decodeErr)
		}
		_ = json.Unmarshal([]byte(decodedValue), &form)
		http.SetCookie(w, &http.Cookie{Name: "form", MaxAge: -1})
	}
	if c, err := r.Cookie("errors"); err == nil {
		decodedValue, decodeErr := url.QueryUnescape(c.Value) // Decodifica o valor do cookie
		if decodeErr != nil {
			log.Println("Error decoding cookie value:", decodeErr)
		}
		_ = json.Unmarshal([]byte(decodedValue), &errors)
		http.SetCookie(w, &http.Cookie{Name: "errors", MaxAge: -1})
	}
	template.View(w, "register", map[string]any{
		"username": form.Name,
		"email":    form.Email,
		"errors":   errors,
	})
}

func processRegisterPage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	form := models.RegisterForm{
		Name:                 r.FormValue("username"),
		Email:                r.FormValue("email"),
		Password:             r.FormValue("password"),
		PasswordConfirmation: r.FormValue("password_confirmation"),
	}

	errors := make(map[string]string)

	// Validate the input
	if form.Name == "" {
		errors["name"] = "Username é obrigatório."
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
	if _, exists := repositories.Users[form.Email]; exists {
		errors["email"] = "User already exists"
	}

	if len(errors) > 0 {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	hashedPassword := utils.HashPassword(form.Password)

	// Create a new user
	repositories.Users[form.Email] = models.Login{
		HashedPassword: hashedPassword,
		SessionToken:   "", // This should be generated securely
		CSRFToken:      "", // This should be generated securely
	}

	sessionService.InitSession(w, form.Email)

	log.Printf("User registered: %s", form.Email)

	// Redirect to the protected area
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
