package handler

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/internal/service"
	"github.com/anglesson/simple-web-server/pkg/database"
	"github.com/anglesson/simple-web-server/pkg/template"
	"github.com/anglesson/simple-web-server/pkg/utils"
)

var sessionService = service.NewSessionService()
var userService = service.NewUserService(
	repository.NewGormUserRepository(database.DB),
	utils.NewEncrypter(),
)

// LoginView renders the login page with CSRF token
func LoginView(w http.ResponseWriter, r *http.Request) {
	csrfToken := sessionService.GenerateCSRFToken()
	sessionService.SetCSRFToken(w)
	template.View(w, r, "login", map[string]interface{}{
		"csrf_token": csrfToken,
	}, "guest")
}

// LoginSubmit handles user login authentication
func LoginSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	loginInput := service.InputLogin{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	errors := make(map[string]string)

	// Validate required fields
	if loginInput.Email == "" {
		errors["email"] = "Email é obrigatório."
	}
	if loginInput.Password == "" {
		errors["password"] = "Senha é obrigatória."
	}

	// If there are validation errors, redirect back with errors
	if len(errors) > 0 {
		redirectWithErrors(w, r, loginInput, errors)
		return
	}

	// Authenticate user using UserService
	user, err := userService.AuthenticateUser(loginInput)
	if err != nil {
		errors["password"] = "Email ou senha inválidos"
		redirectWithErrors(w, r, loginInput, errors)
		return
	}

	// Initialize session for authenticated user
	sessionService.InitSession(w, user.Email)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// LogoutSubmit handles user logout
func LogoutSubmit(w http.ResponseWriter, r *http.Request) {
	sessionService.ClearSession(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// redirectWithErrors is a helper function to redirect with form data and errors
func redirectWithErrors(w http.ResponseWriter, r *http.Request, form interface{}, errors map[string]string) {
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
}
