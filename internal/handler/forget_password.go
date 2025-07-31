package handler

import (
	"net/http"

	"github.com/anglesson/simple-web-server/pkg/template"
)

type ForgetPasswordHandler struct {
	templateRenderer template.TemplateRenderer
}

func NewForgetPasswordHandler(templateRenderer template.TemplateRenderer) *ForgetPasswordHandler {
	return &ForgetPasswordHandler{
		templateRenderer: templateRenderer,
	}
}

func (h *ForgetPasswordHandler) ForgetPasswordView(w http.ResponseWriter, r *http.Request) {
	h.templateRenderer.View(w, r, "forget-password", nil, "guest")
}

func (h *ForgetPasswordHandler) ForgetPasswordSubmit(w http.ResponseWriter, r *http.Request) {
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
