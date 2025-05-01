package handlers

import (
	"net/http"

	"github.com/anglesson/simple-web-server/internal/shared/template"
)

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		renderCreatePage(w, r)
	}
}

func renderCreatePage(w http.ResponseWriter, r *http.Request) {
	template.View(w, "create_ebook", nil, "base_logged")
}
