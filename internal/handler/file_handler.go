package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/anglesson/simple-web-server/internal/authentication/middleware"
	"github.com/anglesson/simple-web-server/internal/handler/web"
	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/internal/repository/gorm"
	"github.com/anglesson/simple-web-server/internal/service"
	"github.com/anglesson/simple-web-server/pkg/database"
	"github.com/anglesson/simple-web-server/pkg/template"
	"github.com/go-chi/chi/v5"
)

type FileHandler struct {
	fileService         service.FileService
	templateRenderer    template.TemplateRenderer
	flashMessageFactory web.FlashMessageFactory
}

func NewFileHandler(fileService service.FileService, templateRenderer template.TemplateRenderer, flashMessageFactory web.FlashMessageFactory) *FileHandler {
	return &FileHandler{
		fileService:         fileService,
		templateRenderer:    templateRenderer,
		flashMessageFactory: flashMessageFactory,
	}
}

// FileIndexView exibe a lista de arquivos do criador
func (h *FileHandler) FileIndexView(w http.ResponseWriter, r *http.Request) {
	creatorID := h.getCreatorIDFromSession(r)
	if creatorID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Obter parâmetros de paginação e busca
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	searchTerm := r.URL.Query().Get("search")
	fileType := r.URL.Query().Get("type")

	// Criar paginação
	pagination := models.NewPagination(page, perPage)

	// Log para debug
	log.Printf("Buscando arquivos para creator ID: %d, página: %d, por página: %d", creatorID, page, perPage)

	// Criar query para busca paginada
	query := repository.FileQuery{
		CreatorID:  creatorID,
		FileType:   fileType,
		SearchTerm: searchTerm,
		Pagination: pagination,
	}

	// Buscar arquivos com paginação
	files, total, err := h.fileService.GetFilesByCreatorPaginated(creatorID, query)
	if err != nil {
		log.Printf("Erro ao buscar arquivos: %v", err)
		http.Error(w, "Erro ao carregar arquivos", http.StatusInternalServerError)
		return
	}

	// Configurar paginação com total
	pagination.SetTotal(total)

	// Adicionar parâmetros de busca à paginação
	pagination.SearchTerm = searchTerm
	pagination.FileType = fileType

	// Log para debug
	log.Printf("Arquivos encontrados: %d de %d total", len(files), total)

	data := map[string]interface{}{
		"Files":      files,
		"Pagination": pagination,
		"Title":      "Minha Biblioteca de Arquivos",
	}

	h.templateRenderer.View(w, r, "file/index", data, "admin")
}

// FileUploadView exibe o formulário de upload
func (h *FileHandler) FileUploadView(w http.ResponseWriter, r *http.Request) {
	creatorID := h.getCreatorIDFromSession(r)
	if creatorID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := map[string]interface{}{
		"Title": "Upload de Arquivo",
	}

	h.templateRenderer.View(w, r, "file/upload", data, "admin")
}

// FileUploadSubmit processa o upload de arquivo
func (h *FileHandler) FileUploadSubmit(w http.ResponseWriter, r *http.Request) {
	creatorID := h.getCreatorIDFromSession(r)
	if creatorID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Parse multipart form (máximo 50MB)
	err := r.ParseMultipartForm(50 << 20)
	if err != nil {
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Arquivo não encontrado", http.StatusBadRequest)
		return
	}
	defer file.Close()

	description := r.FormValue("description")

	_, err = h.fileService.UploadFile(header, description, creatorID)
	if err != nil {
		http.Error(w, "Erro ao fazer upload: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirecionar com mensagem de sucesso
	http.Redirect(w, r, "/file?success=upload", http.StatusSeeOther)
}

// FileDeleteSubmit deleta um arquivo
func (h *FileHandler) FileDeleteSubmit(w http.ResponseWriter, r *http.Request) {
	creatorID := h.getCreatorIDFromSession(r)
	if creatorID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	fileIDStr := chi.URLParam(r, "id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Verificar se o arquivo pertence ao criador antes de deletar
	file, err := h.fileService.GetFileByID(uint(fileID))
	if err != nil {
		http.Error(w, "Arquivo não encontrado", http.StatusNotFound)
		return
	}

	if file.CreatorID != creatorID {
		http.Error(w, "Acesso negado", http.StatusForbidden)
		return
	}

	err = h.fileService.DeleteFile(uint(fileID))
	if err != nil {
		http.Error(w, "Erro ao deletar arquivo", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/file?success=delete", http.StatusSeeOther)
}

// FileUpdateSubmit atualiza nome e descrição do arquivo
func (h *FileHandler) FileUpdateSubmit(w http.ResponseWriter, r *http.Request) {
	creatorID := h.getCreatorIDFromSession(r)
	if creatorID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	fileIDStr := chi.URLParam(r, "id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Verificar se o arquivo pertence ao criador antes de atualizar
	file, err := h.fileService.GetFileByID(uint(fileID))
	if err != nil {
		http.Error(w, "Arquivo não encontrado", http.StatusNotFound)
		return
	}

	if file.CreatorID != creatorID {
		http.Error(w, "Acesso negado", http.StatusForbidden)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")

	// Validar se o nome não está vazio
	if name == "" {
		http.Error(w, "Nome do arquivo é obrigatório", http.StatusBadRequest)
		return
	}

	err = h.fileService.UpdateFile(uint(fileID), name, description)
	if err != nil {
		flashMessage := h.flashMessageFactory(w, r)
		flashMessage.Error("Erro ao atualizar arquivo")
		http.Redirect(w, r, "/file", http.StatusSeeOther)
		return
	}

	flashMessage := h.flashMessageFactory(w, r)
	flashMessage.Success("Arquivo atualizado com sucesso!")
	http.Redirect(w, r, "/file", http.StatusSeeOther)
}

// getCreatorIDFromSession extrai o ID do criador da sessão usando o SessionService injetado
func (h *FileHandler) getCreatorIDFromSession(r *http.Request) uint {
	// Obter usuário da sessão usando o middleware Auth
	userID := middleware.GetCurrentUserID(r)
	if userID == "" {
		log.Printf("Usuário não encontrado na sessão")
		return 0
	}

	log.Printf("Usuário encontrado: ID=%d", userID)

	// Buscar o creator associado ao usuário
	creatorRepository := gorm.NewCreatorRepository(database.DB)
	creator, err := creatorRepository.FindCreatorByUserID(userID)
	if err != nil || creator == nil {
		log.Printf("Erro ao buscar creator para usuário %d: %v", userID, err)
		return 0
	}

	log.Printf("Creator encontrado: ID=%d, Nome=%s", creator.ID, creator.Name)
	return creator.ID
}
