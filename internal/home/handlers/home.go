package home

import (
	"net/http"

	"github.com/anglesson/simple-web-server/internal/shared/template"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		notFoundHandler(w, r)
		return
	}
	template.View(w, "home", nil)
}
