package handlers

import (
	"net/http"

	cookies "github.com/anglesson/simple-web-server/internal/shared/cookie"
)

func LoginView(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement login view
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func LoginSubmit(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement login submit
	cookies.NotifySuccess(w, "Login realizado com sucesso!")
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func RegisterView(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement register view
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func RegisterSubmit(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement register submit
	cookies.NotifySuccess(w, "Cadastro realizado com sucesso!")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func ForgetPasswordView(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement forget password view
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func ForgetPasswordSubmit(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement forget password submit
	cookies.NotifySuccess(w, "Email de recuperação enviado!")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func LogoutSubmit(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement logout
	cookies.NotifySuccess(w, "Logout realizado com sucesso!")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
