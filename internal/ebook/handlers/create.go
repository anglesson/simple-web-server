package handlers

import (
	"net/http"

	"github.com/anglesson/simple-web-server/internal/shared/template"
)

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		renderCreatePage(w, r)
	case http.MethodPost:
		processCreateEbook(w, r)
	}
}

func renderCreatePage(w http.ResponseWriter, r *http.Request) {
	template.View(w, r, "create_ebook", nil, "base_logged")
}

func processCreateEbook(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/ebook", http.StatusSeeOther)
}
