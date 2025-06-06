package http_server

import (
	"net/http"
)

func DashboardView(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement dashboard view
	http.Redirect(w, r, "/client", http.StatusSeeOther)
}

func SettingsView(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement settings view
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func HomeView(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement home view
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
