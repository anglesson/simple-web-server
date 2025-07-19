package handler

import (
	"fmt"
	"net/http"

	"github.com/anglesson/simple-web-server/internal/service"
	"github.com/anglesson/simple-web-server/pkg/template"
)

type CreatorHandler struct {
	creatorService service.CreatorService
}

func NewCreatorHandler(creatorService service.CreatorService) *CreatorHandler {
	return &CreatorHandler{
		creatorService,
	}
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
	}

	_, err := ch.creatorService.CreateCreator(input)
	if err != nil {
		fmt.Printf("[ERROR]: %s\n", err.Error())
		template.View(w, r, "creator/register", map[string]interface{}{
			"Error": err.Error(),
			"Form":  input,
		}, "guest")
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
