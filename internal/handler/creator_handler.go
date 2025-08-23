package handler

import (
	"fmt"
	"net/http"

	"github.com/anglesson/simple-web-server/internal/authentication/business"
	"github.com/anglesson/simple-web-server/internal/service"
	"github.com/anglesson/simple-web-server/pkg/template"
)

type CreatorHandler struct {
	creatorService   service.CreatorService
	sessionService   business.AuthService
	templateRenderer template.TemplateRenderer
}

func NewCreatorHandler(creatorService service.CreatorService, sessionService business.AuthService, templateRenderer template.TemplateRenderer) *CreatorHandler {
	return &CreatorHandler{
		creatorService:   creatorService,
		sessionService:   sessionService,
		templateRenderer: templateRenderer,
	}
}

func (ch *CreatorHandler) RegisterView(w http.ResponseWriter, r *http.Request) {
	ch.templateRenderer.View(w, r, "creator/register", nil, "guest")
}

func (ch *CreatorHandler) RegisterCreatorSSR(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	input := service.InputCreateCreator{
		Name:                 r.FormValue("name"),
		BirthDate:            r.FormValue("birthdate"),
		PhoneNumber:          r.FormValue("phone"),
		Email:                r.FormValue("email"),
		CPF:                  r.FormValue("cpf"),
		Password:             r.FormValue("password"),
		PasswordConfirmation: r.FormValue("password_confirmation"),
		TermsAccepted:        r.FormValue("terms_accepted"),
	}

	_, err := ch.creatorService.CreateCreator(input)
	if err != nil {
		fmt.Printf("[ERROR]: %s\n", err.Error())
		ch.templateRenderer.View(w, r, "creator/register", map[string]interface{}{
			"Error": err.Error(),
			"Form":  input,
		}, "guest")
		return
	}

	// TODO: Fazer login com o Google

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
