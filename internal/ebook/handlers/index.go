package handlers

import (
	"net/http"

	"github.com/anglesson/simple-web-server/internal/shared/template"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		renderIndexPage(w, r)
	}
}

func renderIndexPage(w http.ResponseWriter, r *http.Request) {
	template.View(w, "ebook", nil, "base_logged")
}
