package home

import (
	"net/http"

	"github.com/anglesson/simple-web-server/internal/shared/template"
)

func HomeGetHandler(w http.ResponseWriter, r *http.Request) {
	template.View(w, "home", nil)
}
