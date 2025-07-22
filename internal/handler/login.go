package handler

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/internal/service"
	"github.com/anglesson/simple-web-server/pkg/database"
	"github.com/anglesson/simple-web-server/pkg/template"
	"github.com/anglesson/simple-web-server/pkg/utils"
)

var sessionService = service.NewSessionService()
var encrypter = utils.NewEncrypter()

func LoginView(w http.ResponseWriter, r *http.Request) {
	csrfToken := sessionService.GenerateCSRFToken()
	sessionService.SetCSRFToken(w)
	template.View(w, r, "login", map[string]interface{}{
		"csrf_token": csrfToken,
	}, "guest")
}

func LoginSubmit(w http.ResponseWriter, r *http.Request) {
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
	user := repository.NewGormUserRepository(database.DB).FindByEmail(form.Email)
	if user == nil || !encrypter.CheckPasswordHash(user.Password, form.Password) {
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
