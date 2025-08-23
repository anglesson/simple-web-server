package handler

import (
	"log"
	"net/http"

	"github.com/anglesson/simple-web-server/internal/authentication/middleware"
	"github.com/anglesson/simple-web-server/pkg/template"
)

type SettingsHandler struct {
	templateRenderer template.TemplateRenderer
}

func NewSettingsHandler(templateRenderer template.TemplateRenderer) *SettingsHandler {
	return &SettingsHandler{
		templateRenderer: templateRenderer,
	}
}

func (h *SettingsHandler) SettingsView(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey)
	if user == nil {
		r.Header.Set("Location", "/login")
		return
	}

	csrfToken := r.Context().Value(middleware.CSRFTokenKey).(string)
	if csrfToken == "" {
		r.Header.Set("Location", "/login")
		return
	}

	log.Printf("Renderizando página de configurações para o usuário: %s", user)
	log.Printf("Token CSRF: %s", csrfToken)

	// Passar apenas o objeto user para o template
	h.templateRenderer.View(w, r, "settings", map[string]interface{}{
		"user": user,
	}, "admin")
}
