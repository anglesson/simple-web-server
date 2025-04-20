package handlers

import (
	"net/http"

	"github.com/anglesson/simple-web-server/templates"
)

func DashboardGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.View(w, "dashboard", nil)
}
