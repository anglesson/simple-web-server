package client

import (
	"encoding/csv"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/anglesson/simple-web-server/internal/application/dtos"
	"github.com/anglesson/simple-web-server/internal/common"
	"github.com/anglesson/simple-web-server/internal/infrastructure"
	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repositories"
	"github.com/anglesson/simple-web-server/internal/services"
	cookies "github.com/anglesson/simple-web-server/internal/shared/cookie"
	"github.com/anglesson/simple-web-server/internal/shared/middlewares"
	"github.com/anglesson/simple-web-server/internal/shared/template"
	"github.com/go-chi/chi/v5"
)

type ClientHandler struct {
	clientService       ClientServicePort
	flashMessageFactory infrastructure.FlashMessageFactory
}

func NewClientHandler(clientService ClientServicePort, flashMessageFactory infrastructure.FlashMessageFactory) *ClientHandler {
	return &ClientHandler{
		clientService:       clientService,
		flashMessageFactory: flashMessageFactory,
	}
}

func (ch *ClientHandler) CreateView(w http.ResponseWriter, r *http.Request) {
	loggedUser := middlewares.Auth(r)
	if loggedUser.ID == 0 {
		http.Error(w, "Não foi possível prosseguir com a sua solicitação", http.StatusInternalServerError)
		return
	}

	template.View(w, r, "client/create", nil, "admin")
}

func (ch *ClientHandler) UpdateView(w http.ResponseWriter, r *http.Request) {
	loggedUser := middlewares.Auth(r)
	if loggedUser.ID == 0 {
		http.Error(w, "Não foi possível prosseguir com a sua solicitação", http.StatusInternalServerError)
		return
	}

	clientID := chi.URLParam(r, "id")
	id, _ := strconv.ParseUint(clientID, 10, 32)
	client, err := ch.clientService.FindCreatorsClientByID(uint(id), loggedUser.Email)
	if err != nil {
		http.Redirect(w, r, r.Referer(), http.StatusNotFound)
	}

	template.View(w, r, "client/update", map[string]interface{}{"Client": client}, "admin")
}

func ClientIndexView(w http.ResponseWriter, r *http.Request) {
	loggedUser := middlewares.Auth(r)
	if loggedUser.ID == 0 {
		http.Error(w, "Não foi possível prosseguir com a sua solicitação", http.StatusInternalServerError)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	term := r.URL.Query().Get("term")
	pagination := dtos.NewPagination(page, perPage)

	log.Printf("User Logado: %v", loggedUser.Email)

	creatorRepository := repositories.NewCreatorRepository()
	creator, err := creatorRepository.FindCreatorByUserID(loggedUser.ID)
	if err != nil {
		common.RedirectBackWithErrors(w, r, err.Error())
	}

	clients, err := repositories.NewClientRepository().FindClientsByCreator(creator, dtos.ClientQuery{
		Term:       term,
		Pagination: pagination,
	})
	if err != nil {
		common.RedirectBackWithErrors(w, r, err.Error())
	}

	template.View(w, r, "client", map[string]any{
		"Clients":    clients,
		"Pagination": pagination,
	}, "admin")
}

func (ch *ClientHandler) ClientCreateSubmit(w http.ResponseWriter, r *http.Request) {
	flashMessage := ch.flashMessageFactory(w, r)

	user_email, ok := r.Context().Value(middlewares.UserEmailKey).(string)
	if !ok {
		flashMessage.Error("Unauthorized. Invalid user email")
		http.Error(w, "Invalid user email", http.StatusUnauthorized)
		return
	}

	input := dtos.CreateClientInput{
		Name:      r.FormValue("name"),
		CPF:       r.FormValue("cpf"),
		BirthDate: r.FormValue("birthdate"),
		Email:     r.FormValue("email"),
		Phone:     r.FormValue("phone"),
	}

	input.EmailCreator = user_email

	// TODO: Validar se o cliente existe
	_, err := ch.clientService.CreateClient(input)
	if err != nil {
		flashMessage.Error(err.Error())
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}
	flashMessage.Success("Cliente foi cadastrado!")

	http.Redirect(w, r, "/client", http.StatusSeeOther)
}

func (ch *ClientHandler) ClientUpdateSubmit(w http.ResponseWriter, r *http.Request) {
	flashMessage := ch.flashMessageFactory(w, r)
	user_email, ok := r.Context().Value(middlewares.UserEmailKey).(string)
	if !ok {
		http.Error(w, "Invalid user email", http.StatusInternalServerError)
		return
	}

	input := dtos.UpdateClientInput{
		CPF:          r.FormValue("cpf"),
		Email:        r.FormValue("email"),
		Phone:        r.FormValue("phone"),
		EmailCreator: user_email,
	}

	_, err := ch.clientService.Update(input)
	if err != nil {
		flashMessage.Error(err.Error())
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}
	flashMessage.Success("Cliente foi atualizado!")

	http.Redirect(w, r, "/client", http.StatusSeeOther)
}

func (ch *ClientHandler) ClientImportSubmit(w http.ResponseWriter, r *http.Request) {
	log.Println("Iniciando processamento de CSV")
	user_email, ok := r.Context().Value(middlewares.UserEmailKey).(string)
	if !ok {
		http.Error(w, "Invalid user email", http.StatusInternalServerError)
		return
	}

	creatorService := services.NewCreatorService()
	creator, err := creatorService.FindCreatorByEmail(user_email)
	if err != nil {
		log.Println("Nao autorizado")
		http.Redirect(w, r, r.Referer(), http.StatusUnauthorized)
		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Erro ao processar o formulário", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		common.RedirectBackWithErrors(w, r, "Erro ao ler o arquivo")
	}
	defer file.Close()

	// Verifica a extensão do arquivo (opcional)
	if !strings.HasSuffix(handler.Filename, ".csv") {
		common.RedirectBackWithErrors(w, r, "Arquivo não é CSV")
	}

	log.Println("Arquivo validado!")

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		log.Printf("Erro na leitura do CSV: %s", err.Error())
		common.RedirectBackWithErrors(w, r, "Erro na leitura do CSV")
	}

	// Validate header
	var clients []*models.Client

	for i, linha := range rows {
		log.Printf("linha: %s", linha)
		if i > 0 {
			client := models.NewClient(linha[0], linha[1], linha[2], linha[3], linha[4], creator)
			clients = append(clients, client)
		}
	}

	if err = ch.clientService.CreateBatchClient(clients); err != nil {
		common.RedirectBackWithErrors(w, r, err.Error())
		return
	}

	cookies.NotifySuccess(w, "Clientes foram importados!")
	http.Redirect(w, r, "/client", http.StatusSeeOther)
}
