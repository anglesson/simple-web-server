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
		"Form":   form,
		"Errors": errors,
	})
}

func processRegisterPage(w http.ResponseWriter, r *http.Request) {
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
	if foundedUser.ID != 0 {
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

	sessionService.InitSession(w, form.Email)

	log.Printf("User registered ID: %d, EMAIL: %s", user.ID, user.Email)

	// Redirect to the protected area
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
