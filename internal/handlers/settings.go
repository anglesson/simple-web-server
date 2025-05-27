package handlers

import (
	"log"
	"net/http"

	"github.com/anglesson/simple-web-server/internal/shared/middlewares"
	"github.com/anglesson/simple-web-server/internal/shared/template"
)

func SettingsView(w http.ResponseWriter, r *http.Request) {
	user := middlewares.Auth(r)
	if user == nil {
		log.Printf("Usuário não autenticado ao acessar configurações")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Usar o token CSRF do usuário diretamente
	csrfToken := user.CSRFToken
	if csrfToken == "" {
		log.Printf("Token CSRF não encontrado para o usuário: %s", user.Email)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	log.Printf("Renderizando página de configurações para o usuário: %s", user.Email)
	log.Printf("Token CSRF: %s", csrfToken)

	template.View(w, r, "settings", map[string]interface{}{
		"user":       user,
		"csrf_token": csrfToken,
	}, "admin")
}
