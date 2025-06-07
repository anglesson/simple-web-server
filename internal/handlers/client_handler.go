package handlers

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repositories"
	"github.com/anglesson/simple-web-server/internal/services"
	cookies "github.com/anglesson/simple-web-server/internal/shared/cookie"
	"github.com/anglesson/simple-web-server/internal/shared/middlewares"
	"github.com/anglesson/simple-web-server/internal/shared/template"
	"github.com/anglesson/simple-web-server/internal/shared/utils"
	"github.com/go-chi/chi/v5"
)

type ClientHandler struct {
}

func NewClientHandler() *ClientHandler {
	return &ClientHandler{}
}

func (ch *ClientHandler) CreateView(w http.ResponseWriter, r *http.Request) {
	loggedUser := middlewares.Auth(r)
	if loggedUser.ID == 0 {
		http.Error(w, "Não foi possível prosseguir com a sua solicitação", http.StatusInternalServerError)
		return
	}

	template.View(w, r, "client/create", nil, "admin")
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
	pagination := repositories.NewPagination(page, perPage)

	log.Printf("User Logado: %v", loggedUser.Email)

	creatorRepository := repositories.NewCreatorRepository()
	creator, err := creatorRepository.FindCreatorByUserID(loggedUser.ID)
	if err != nil {
		redirectBackWithErrors(w, r, err.Error())
	}

	clients, err := repositories.NewClientRepository().FindClientsByCreator(creator, repositories.ClientQuery{
		Term:       term,
		Pagination: pagination,
	})
	if err != nil {
		redirectBackWithErrors(w, r, err.Error())
	}

	template.View(w, r, "client", map[string]any{
		"Clients":    clients,
		"Pagination": pagination,
	}, "admin")
}

func ClientCreateSubmit(w http.ResponseWriter, r *http.Request) {
	errors := make(map[string]string)

	form := models.ClientRequest{
		Name:           r.FormValue("name"),
		CPF:            r.FormValue("cpf"),
		DataNascimento: r.FormValue("data_nascimento"),
		Email:          r.FormValue("email"),
		Phone:          r.FormValue("phone"),
	}

	errForm := utils.ValidateForm(form)
	for key, value := range errForm {
		errors[key] = value
	}

	if len(errors) > 0 {
		formJSON, _ := json.Marshal(form)
		errorsJSON, _ := json.Marshal(errors)

		http.SetCookie(w, &http.Cookie{
			Name:  "form",
			Value: url.QueryEscape(string(formJSON)),
			Path:  "/",
		})
		http.SetCookie(w, &http.Cookie{
			Name:  "errors",
			Value: url.QueryEscape(string(errorsJSON)),
			Path:  "/",
		})
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	user_email, ok := r.Context().Value(middlewares.UserEmailKey).(string)
	if !ok {
		http.Error(w, "Invalid user email", http.StatusInternalServerError)
		return
	}

	creatorService := services.NewCreatorService()
	creator, err := creatorService.FindCreatorByEmail(user_email)
	if err != nil {
		redirectBackWithErrors(w, r, err.Error())
	}

	// TODO: Validar se o cliente existe
	clientService := services.NewClientService()
	_, err = clientService.CreateClient(form.Name, form.CPF, form.DataNascimento, form.Email, form.Phone, creator)
	if err != nil {
		redirectBackWithErrors(w, r, err.Error())
	}
	cookies.NotifySuccess(w, "Cliente foi cadastrado!")

	http.Redirect(w, r, "/client", http.StatusSeeOther)
}

func ClientUpdateSubmit(w http.ResponseWriter, r *http.Request) {
	user_email, ok := r.Context().Value(middlewares.UserEmailKey).(string)
	if !ok {
		http.Error(w, "Invalid user email", http.StatusInternalServerError)
		return
	}

	creatorService := services.NewCreatorService()
	creator, err := creatorService.FindCreatorByEmail(user_email)
	if err != nil {
		http.Redirect(w, r, r.Referer(), http.StatusUnauthorized)
		return
	}

	errors := make(map[string]string)

	clientID, _ := strconv.Atoi(chi.URLParam(r, "id"))
	log.Printf("ClientID: %v", clientID)
	if clientID == 0 {
		cookies.NotifyError(w, "O ID deve ser informado")
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	form := models.ClientRequest{
		ID:             uint(clientID),
		Name:           r.FormValue("name"),
		CPF:            r.FormValue("cpf"),
		DataNascimento: r.FormValue("data_nascimento"),
		Email:          r.FormValue("email"),
		Phone:          r.FormValue("phone"),
	}

	clientService := services.NewClientService()
	client, err := clientService.FindCreatorsClientByID(form.ID, creator.ID)
	log.Printf("Id do client encontrado: %v", client.ID)
	if err != nil {
		cookies.NotifyError(w, err.Error())
		http.Redirect(w, r, r.Referer(), http.StatusNotFound)
		return
	}

	// Move to a validatorService
	errForm := utils.ValidateForm(form)
	for key, value := range errForm {
		errors[key] = value
	}

	if len(errors) > 0 {
		formJSON, _ := json.Marshal(form)
		errorsJSON, _ := json.Marshal(errors)

		http.SetCookie(w, &http.Cookie{
			Name:  "form",
			Value: url.QueryEscape(string(formJSON)),
			Path:  "/",
		})
		http.SetCookie(w, &http.Cookie{
			Name:  "errors",
			Value: url.QueryEscape(string(errorsJSON)),
			Path:  "/",
		})
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	err = clientService.Update(client, form)
	if err != nil {
		cookies.NotifyError(w, err.Error())
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}
	cookies.NotifySuccess(w, "Cliente foi atualizado!")

	http.Redirect(w, r, "/client", http.StatusSeeOther)
}

func ClientImportSubmit(w http.ResponseWriter, r *http.Request) {
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
		redirectBackWithErrors(w, r, "Erro ao ler o arquivo")
	}
	defer file.Close()

	// Verifica a extensão do arquivo (opcional)
	if !strings.HasSuffix(handler.Filename, ".csv") {
		redirectBackWithErrors(w, r, "Arquivo não é CSV")
	}

	log.Println("Arquivo validado!")

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		log.Printf("Erro na leitura do CSV: %s", err.Error())
		redirectBackWithErrors(w, r, "Erro na leitura do CSV")
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

	clientService := services.NewClientService()
	if err = clientService.CreateBatchClient(clients); err != nil {
		redirectBackWithErrors(w, r, err.Error())
		return
	}

	cookies.NotifySuccess(w, "Clientes foram importados!")
	http.Redirect(w, r, "/client", http.StatusSeeOther)
}

func redirectBackWithErrors(w http.ResponseWriter, r *http.Request, erroMessage string) {
	cookies.NotifyError(w, erroMessage)
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}
