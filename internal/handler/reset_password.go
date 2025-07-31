package handler

import (
	"net/http"

	"github.com/anglesson/simple-web-server/internal/service"
	cookies "github.com/anglesson/simple-web-server/pkg/cookie"
	"github.com/anglesson/simple-web-server/pkg/template"
)

type ResetPasswordHandler struct {
	templateRenderer template.TemplateRenderer
	userService      service.UserService
}

func NewResetPasswordHandler(templateRenderer template.TemplateRenderer, userService service.UserService) *ResetPasswordHandler {
	return &ResetPasswordHandler{
		templateRenderer: templateRenderer,
		userService:      userService,
	}
}

func (h *ResetPasswordHandler) ResetPasswordView(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		cookies.NotifyError(w, "Token de reset inválido")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := map[string]interface{}{
		"Token": token,
	}

	h.templateRenderer.View(w, r, "reset-password", data, "guest")
}

func (h *ResetPasswordHandler) ResetPasswordSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	token := r.FormValue("token")
	password := r.FormValue("password")
	passwordConfirmation := r.FormValue("password_confirmation")

	if token == "" {
		cookies.NotifyError(w, "Token de reset inválido")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if password == "" {
		cookies.NotifyError(w, "Senha é obrigatória")
		http.Redirect(w, r, "/reset-password?token="+token, http.StatusSeeOther)
		return
	}

	if password != passwordConfirmation {
		cookies.NotifyError(w, "As senhas não coincidem")
		http.Redirect(w, r, "/reset-password?token="+token, http.StatusSeeOther)
		return
	}

	if len(password) < 6 {
		cookies.NotifyError(w, "A senha deve ter pelo menos 6 caracteres")
		http.Redirect(w, r, "/reset-password?token="+token, http.StatusSeeOther)
		return
	}

	// Resetar senha
	err := h.userService.ResetPassword(token, password)
	if err != nil {
		if err == service.ErrInvalidResetToken {
			cookies.NotifyError(w, "Token de reset inválido ou expirado. Solicite um novo link.")
		} else {
			cookies.NotifyError(w, "Erro ao resetar senha. Tente novamente.")
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	cookies.NotifySuccess(w, "Senha alterada com sucesso! Você pode fazer login agora.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
