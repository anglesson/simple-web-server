package handler

import (
	"fmt"
	"net/http"

	"github.com/anglesson/simple-web-server/internal/service"
	"github.com/anglesson/simple-web-server/pkg/template"
)

type CreatorHandler struct {
	userService service.UserService
}

func NewCreatorHandler(userService service.UserService) *CreatorHandler {
	return &CreatorHandler{
		userService,
	}
}

func (ch *CreatorHandler) RegisterCreatorSSR(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	input := service.InputCreateCreator{
		Name:                 r.FormValue("name"),
		BirthDate:            r.FormValue("birthDate"),
		PhoneNumber:          r.FormValue("phoneNumber"),
		Email:                r.FormValue("email"),
		CPF:                  r.FormValue("cpf"),
		Password:             r.FormValue("password"),
		PasswordConfirmation: r.FormValue("password_confirmation"),
	}

	inputCreateUser := service.InputCreateUser{
		Username:             input.Name,
		Email:                input.Email,
		Password:             input.Password,
		PasswordConfirmation: input.PasswordConfirmation,
	}

	_, err := ch.userService.CreateUser(inputCreateUser)
	if err != nil {
		fmt.Printf("[ERROR]: %s", err.Error())
		template.View(w, r, "creator/register", map[string]interface{}{
			"Error": err.Error(),
			"Form":  input,
		}, "guest")
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
