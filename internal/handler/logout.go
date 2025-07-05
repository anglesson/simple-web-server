package handler

import "net/http"

func LogoutSubmit(w http.ResponseWriter, r *http.Request) {
	sessionService.ClearSession(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
