package handlers

import (
	"net/http"

	"github.com/anglesson/simple-web-server/internal/shared/template"
)

func ForgetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		renderForgetPasswordPage(w, r)
	case http.MethodPost:
		processForgetPasswordPage(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	template.View(w, "forget-password", nil)
}

func renderForgetPasswordPage(w http.ResponseWriter, r *http.Request) {
	template.View(w, "forget-password", nil)
}

func processForgetPasswordPage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/password-reset-success", http.StatusSeeOther)
}
