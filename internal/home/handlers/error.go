package home

import (
	"net/http"

	"github.com/anglesson/simple-web-server/internal/shared/template"
)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	template.View(w, r, "404-error", nil, "base_logged")
}
