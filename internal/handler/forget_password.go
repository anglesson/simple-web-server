package handler

import (
	"net/http"

	"github.com/anglesson/simple-web-server/internal/service"
	cookies "github.com/anglesson/simple-web-server/pkg/cookie"
	"github.com/anglesson/simple-web-server/pkg/template"
)

type ForgetPasswordHandler struct {
	templateRenderer template.TemplateRenderer
	userService      service.UserService
	emailService     *service.EmailService
}

func NewForgetPasswordHandler(templateRenderer template.TemplateRenderer, userService service.UserService, emailService *service.EmailService) *ForgetPasswordHandler {
	return &ForgetPasswordHandler{
		templateRenderer: templateRenderer,
		userService:      userService,
		emailService:     emailService,
	}
}

func (h *ForgetPasswordHandler) ForgetPasswordView(w http.ResponseWriter, r *http.Request) {
	h.templateRenderer.View(w, r, "forget-password", nil, "guest")
}

func (h *ForgetPasswordHandler) ForgetPasswordSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	if email == "" {
		cookies.NotifyError(w, "Email é obrigatório")
		http.Redirect(w, r, "/forget-password", http.StatusSeeOther)
		return
	}

	// Solicitar reset de senha
	err := h.userService.RequestPasswordReset(email)
	if err != nil {
		cookies.NotifyError(w, "Erro ao processar solicitação de reset de senha")
		http.Redirect(w, r, "/forget-password", http.StatusSeeOther)
		return
	}

	// Buscar usuário para enviar e-mail
	user := h.userService.FindByEmail(email)
	if user != nil {
		// Enviar e-mail de reset
		err = h.emailService.SendPasswordResetEmail(user.Username, user.Email, user.PasswordResetToken)
		if err != nil {
			cookies.NotifyError(w, "Erro ao enviar e-mail de reset de senha")
			http.Redirect(w, r, "/forget-password", http.StatusSeeOther)
			return
		}
	}

	// Sempre redirecionar para sucesso (não revelar se o email existe ou não)
	http.Redirect(w, r, "/password-reset-success", http.StatusSeeOther)
}
