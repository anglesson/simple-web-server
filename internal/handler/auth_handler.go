package handler

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/anglesson/simple-web-server/internal/service"
	"github.com/anglesson/simple-web-server/pkg/template"
)

type AuthHandler struct {
	userService    service.UserService
	sessionService service.SessionService
}

func NewAuthHandler(userService service.UserService, sessionService service.SessionService) *AuthHandler {
	return &AuthHandler{
		userService:    userService,
		sessionService: sessionService,
	}
}

// LoginView renders the login page with CSRF token
func (h *AuthHandler) LoginView(w http.ResponseWriter, r *http.Request) {
	csrfToken := h.sessionService.GenerateCSRFToken()
	h.sessionService.SetCSRFToken(w)
	template.View(w, r, "login", map[string]interface{}{
		"csrf_token": csrfToken,
	}, "guest")
}

// LoginSubmit handles user login authentication
func (h *AuthHandler) LoginSubmit(w http.ResponseWriter, r *http.Request) {
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
		h.redirectWithErrors(w, r, loginInput, errors)
		return
	}

	// Authenticate user using UserService
	user, err := h.userService.AuthenticateUser(loginInput)
	if err != nil {
		errors["password"] = "Email ou senha inválidos"
		h.redirectWithErrors(w, r, loginInput, errors)
		return
	}

	// Initialize session for authenticated user
	h.sessionService.InitSession(w, user.Email)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// LogoutSubmit handles user logout
func (h *AuthHandler) LogoutSubmit(w http.ResponseWriter, r *http.Request) {
	h.sessionService.ClearSession(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// redirectWithErrors is a helper function to redirect with form data and errors
func (h *AuthHandler) redirectWithErrors(w http.ResponseWriter, r *http.Request, form interface{}, errors map[string]string) {
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
