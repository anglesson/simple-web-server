package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/anglesson/simple-web-server/internal/application"
	"github.com/anglesson/simple-web-server/internal/infrastructure/http_server/utils"
	cookies "github.com/anglesson/simple-web-server/internal/shared/cookie"
	"github.com/go-chi/chi/v5"
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
	// 1. Validar e processar o formulário
	err := r.ParseForm()
	if err != nil {
		log.Printf("Erro ao analisar formulário: %v", err)
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	// 2. Extrair e validar dados do formulário
	name := r.FormValue("name")
	if name == "" {
		cookies.NotifyError(w, "Nome é obrigatório")
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	// 3. Construir input com validação básica
	input := application.CreateClientInput{
		Name:             name,
		CPF:              r.FormValue("cpf"),
		BirthDay:         r.FormValue("birth_day"),
		Email:            r.FormValue("email"),
		Phone:            r.FormValue("phone"),
		CreatorUserEmail: r.Context().Value(utils.LoggedUserKey).(string),
	}

	// 4. Chamar o caso de uso
	output, err := h.clientUseCase.CreateClient(input)
	if err != nil {
		cookies.NotifyError(w, err.Error())
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	// 5. Sucesso - redirecionar com mensagem
	log.Printf("Cliente criado com ID: %d", output.ID)
	cookies.NotifySuccess(w, "Cliente foi cadastrado com sucesso!")
	http.Redirect(w, r, "/client", http.StatusSeeOther)
}

func (h *ClientHandler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	// 1. Validar e processar o formulário
	err := r.ParseForm()
	if err != nil {
		log.Printf("Erro ao analisar formulário: %v", err)
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	// 2. Validar ID do cliente
	clientID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		cookies.NotifyError(w, "ID do cliente inválido")
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	// 3. Validar nome (campo obrigatório)
	name := r.FormValue("name")
	if name == "" {
		cookies.NotifyError(w, "Nome é obrigatório")
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	// 4. Construir input
	input := application.UpdateClientInput{
		ID:               uint(clientID),
		Name:             name,
		CPF:              r.FormValue("cpf"),
		BirthDay:         r.FormValue("birth_day"),
		Email:            r.FormValue("email"),
		Phone:            r.FormValue("phone"),
		CreatorUserEmail: r.Context().Value(utils.LoggedUserKey).(string),
	}

	// 5. Chamar o caso de uso
	output, err := h.clientUseCase.UpdateClient(input)
	if err != nil {
		cookies.NotifyError(w, err.Error())
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	// 6. Sucesso - redirecionar com mensagem
	log.Printf("Cliente atualizado com ID: %d", output.ID)
	cookies.NotifySuccess(w, "Cliente foi atualizado com sucesso!")
	http.Redirect(w, r, "/client", http.StatusSeeOther)
}

func (h *ClientHandler) ImportClients(w http.ResponseWriter, r *http.Request) {
	// 1. Validar e processar o formulário multipart
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		cookies.NotifyError(w, "Erro ao processar o formulário")
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	// 2. Obter e validar o arquivo
	file, handler, err := r.FormFile("file")
	if err != nil {
		cookies.NotifyError(w, "Erro ao ler o arquivo")
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}
	defer file.Close()

	// 3. Validar extensão do arquivo
	if !strings.HasSuffix(strings.ToLower(handler.Filename), ".csv") {
		cookies.NotifyError(w, "Apenas arquivos CSV são permitidos")
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	// 4. Ler conteúdo do arquivo
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		cookies.NotifyError(w, "Erro ao ler o arquivo")
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	// 5. Construir input
	input := application.ImportClientsInput{
		File:             fileBytes,
		FileName:         handler.Filename,
		CreatorUserEmail: r.Context().Value(utils.LoggedUserKey).(string),
	}

	// 6. Chamar o caso de uso
	output, err := h.clientUseCase.ImportClients(input)
	if err != nil {
		cookies.NotifyError(w, err.Error())
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	// 7. Sucesso - redirecionar com mensagem
	log.Printf("Importados %d clientes", output.ImportedCount)
	cookies.NotifySuccess(w, fmt.Sprintf("%d clientes foram importados com sucesso!", output.ImportedCount))
	http.Redirect(w, r, "/client", http.StatusSeeOther)
}

func (h *ClientHandler) ListClients(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 {
		pageSize = 10
	}
	term := r.URL.Query().Get("term")

	input := application.ListClientsInput{
		Term:             term,
		Page:             page,
		PageSize:         pageSize,
		CreatorUserEmail: r.Context().Value(utils.LoggedUserKey).(string),
	}

	_, err := h.clientUseCase.ListClients(input)
	if err != nil {
		cookies.NotifyError(w, err.Error())
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))
}
