package handlers

import (
	"net/http"

	"github.com/anglesson/simple-web-server/templates"
)

func HomeGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.View(w, "home", nil)
}
