package handler

import (
	"log"
	"net/http"

	"github.com/anglesson/simple-web-server/internal/handler/middleware"
	"github.com/anglesson/simple-web-server/pkg/template"
)

func SettingsView(w http.ResponseWriter, r *http.Request) {
	user := middleware.Auth(r)
	if user == nil {
		log.Printf("Usuário não autenticado ao acessar configurações")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Gerar novo token CSRF se necessário
	if user.CSRFToken == "" {
		user.CSRFToken = sessionService.GenerateCSRFToken()
		sessionService.SetCSRFToken(w)
	}

	log.Printf("Renderizando página de configurações para o usuário: %s", user.Email)
	log.Printf("Token CSRF: %s", user.CSRFToken)

	// Passar apenas o objeto user para o template
	template.View(w, r, "settings", map[string]interface{}{
		"user": user,
	}, "admin")
}
