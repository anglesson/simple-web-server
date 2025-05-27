package handlers

import (
	"net/http"

	"github.com/anglesson/simple-web-server/internal/shared/template"
)

func HomeView(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorView(w, r, 404)
		return
	}

	// Generate CSRF token for the home page
	csrfToken := sessionService.GenerateCSRFToken()
	sessionService.SetCSRFToken(w)

	template.View(w, r, "home", map[string]interface{}{
		"csrf_token": csrfToken,
	}, "guest")
}
