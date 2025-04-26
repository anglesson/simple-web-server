package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/anglesson/simple-web-server/internal/auth/models"
	"github.com/anglesson/simple-web-server/internal/auth/repositories"
	"github.com/anglesson/simple-web-server/internal/auth/services"
	"github.com/anglesson/simple-web-server/internal/shared/template"
	"github.com/anglesson/simple-web-server/internal/shared/utils"
)

var sessionService = services.NewSessionService()

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

	sessionService.InitSession(w, form.Email)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
