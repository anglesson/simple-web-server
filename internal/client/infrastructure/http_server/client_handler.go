package http_server

import (
	"log"
	"net/http"

	client_application "github.com/anglesson/simple-web-server/internal/client/application"
	cookies "github.com/anglesson/simple-web-server/internal/shared/cookie"
)

type ClientHandler struct {
	createClientUseCase client_application.CreateClientUseCase
}

func NewClientHandler(useCase client_application.CreateClientUseCase) *ClientHandler {
	return &ClientHandler{
		createClientUseCase: useCase,
	}
}

func (h *ClientHandler) CreateClientSubmit(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Erro ao analisar formulário para SSR: %v", err)
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	input := client_application.CreateClientInput{
		Name:     r.FormValue("name"),
		CPF:      r.FormValue("cpf"),
		BirthDay: r.FormValue("birth_day"),
		Email:    r.FormValue("email"),
		Phone:    r.FormValue("phone"),
	}

	_, err = h.createClientUseCase.Execute(input)
	if err != nil {
		cookies.NotifyError(w, err.Error())
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	cookies.NotifySuccess(w, "Cliente foi atualizado!")

	http.Redirect(w, r, "/client", http.StatusSeeOther)
}
