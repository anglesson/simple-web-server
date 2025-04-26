package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/anglesson/simple-web-server/internal/auth/models"
	"github.com/anglesson/simple-web-server/internal/auth/repositories"
	"github.com/anglesson/simple-web-server/internal/shared/template"
	"github.com/anglesson/simple-web-server/internal/shared/utils"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		renderLoginPage(w, r)
	case http.MethodPost:
		processLogin(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func renderLoginPage(w http.ResponseWriter, r *http.Request) {
	var form models.LoginForm
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

	template.View(w, "login", map[string]interface{}{
		"Form":   form,
		"Errors": errors,
	})
}

func processLogin(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	form := models.LoginForm{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	errors := make(map[string]string)

	if form.Email == "" {
		errors["email"] = "Email é obrigatório."
	}
	if form.Password == "" {
		errors["password"] = "Senha é obrigatória."
	}

	errors["email"] = "Email inválido"
	errors["password"] = "Senha inválida"

	// Check if the user exists
	user, exists := repositories.Users[form.Email]
	if !exists || !utils.CheckPasswordHash(user.HashedPassword, form.Password) {
		errors["password"] = "Email ou senha inválidos"
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
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	sessionLogin(w, form.Email)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

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

	// Login and redirect to the protected area
	sessionLogin(w, form.Email)

	log.Printf("User registered: %s", form.Email)

	// Redirect to the protected area
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func ForgetPasswordGetHandler(w http.ResponseWriter, r *http.Request) {
	template.View(w, "forget-password", nil)
}

func sessionLogin(w http.ResponseWriter, email string) {
	sessionToken := utils.GenerateToken(32)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})
	csrfToken := utils.GenerateToken(32)
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})
	// Update the session token in the user data
	userFound := repositories.Users[email]
	userFound.SessionToken = sessionToken
	userFound.CSRFToken = csrfToken
	repositories.Users[email] = userFound
}
