package http_server

import (
	"log"
	"net/http"

	"github.com/anglesson/simple-web-server/internal/application"
	common_infrastructure "github.com/anglesson/simple-web-server/internal/common/infrastructure"
	cookies "github.com/anglesson/simple-web-server/internal/shared/cookie"
)

type ClientHandler struct {
	clientUseCase application.ClientUseCasePort
}

func NewClientSSRHandler(useCase application.ClientUseCasePort) *ClientHandler {
	return &ClientHandler{
		clientUseCase: useCase,
	}
}

func (h *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Erro ao analisar formulário para SSR: %v", err)
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	input := application.CreateClientInput{
		Name:             r.FormValue("name"),
		CPF:              r.FormValue("cpf"),
		BirthDay:         r.FormValue("birth_day"),
		Email:            r.FormValue("email"),
		Phone:            r.FormValue("phone"),
		CreatorUserEmail: r.Context().Value(common_infrastructure.LoggedUserKey).(string),
	}

	_, err = h.clientUseCase.CreateClient(input)
	if err != nil {
		cookies.NotifyError(w, err.Error())
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	cookies.NotifySuccess(w, "Cliente foi atualizado!")

	http.Redirect(w, r, "/client", http.StatusSeeOther)
}
