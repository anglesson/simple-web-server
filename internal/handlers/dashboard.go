package handlers

import (
	"net/http"

	"github.com/anglesson/simple-web-server/internal/shared/template"
)

func DashboardView(w http.ResponseWriter, r *http.Request) {
	template.View(w, r, "dashboard", nil, "admin")
}
