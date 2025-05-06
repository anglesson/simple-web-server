package handlers

import (
	"net/http"

	"github.com/anglesson/simple-web-server/internal/shared/template"
)

func HomeView(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		notFoundHandler(w, r)
		return
	}
	template.View(w, r, "home", nil, "guest")
}
