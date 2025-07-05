package handler

import (
	"net/http"

	"github.com/anglesson/simple-web-server/pkg/template"
)

func ForgetPasswordView(w http.ResponseWriter, r *http.Request) {
	template.View(w, r, "forget-password", nil, "guest")
}

func ForgetPasswordSubmit(w http.ResponseWriter, r *http.Request) {
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
