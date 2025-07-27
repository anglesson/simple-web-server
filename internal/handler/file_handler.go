package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/anglesson/simple-web-server/internal/handler/middleware"
	"github.com/anglesson/simple-web-server/internal/repository/gorm"
	"github.com/anglesson/simple-web-server/internal/service"
	"github.com/anglesson/simple-web-server/pkg/database"
	"github.com/anglesson/simple-web-server/pkg/template"
	"github.com/go-chi/chi/v5"
)

type FileHandler struct {
	fileService    service.FileService
	sessionService service.SessionService
}

func NewFileHandler(fileService service.FileService, sessionService service.SessionService) *FileHandler {
	return &FileHandler{
		fileService:    fileService,
		sessionService: sessionService,
	}
}

// FileIndexView exibe a lista de arquivos do criador
func (h *FileHandler) FileIndexView(w http.ResponseWriter, r *http.Request) {
	creatorID := h.getCreatorIDFromSession(r)
	if creatorID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Log para debug
	log.Printf("Buscando arquivos para creator ID: %d", creatorID)

	// Usar GetFilesByCreator em vez de GetActiveByCreator para mostrar todos os arquivos
	files, err := h.fileService.GetFilesByCreator(creatorID)
	if err != nil {
		log.Printf("Erro ao buscar arquivos: %v", err)
		http.Error(w, "Erro ao carregar arquivos", http.StatusInternalServerError)
		return
	}

	// Log para debug
	log.Printf("Arquivos encontrados: %d", len(files))
	for i, file := range files {
		log.Printf("Arquivo %d: ID=%d, Nome=%s, Tipo=%s, CreatorID=%d",
			i+1, file.Model.ID, file.OriginalName, file.FileType, file.CreatorID)
	}

	data := map[string]interface{}{
		"Files": files,
		"Title": "Minha Biblioteca de Arquivos",
	}

	template.View(w, r, "file/index", data, "admin")
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

	template.View(w, r, "file/upload", data, "admin")
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

// FileUpdateSubmit atualiza descrição do arquivo
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

	description := r.FormValue("description")

	err = h.fileService.UpdateFile(uint(fileID), description)
	if err != nil {
		http.Error(w, "Erro ao atualizar arquivo", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/file?success=update", http.StatusSeeOther)
}

// getCreatorIDFromSession extrai o ID do criador da sessão usando o SessionService injetado
func (h *FileHandler) getCreatorIDFromSession(r *http.Request) uint {
	// Obter usuário da sessão usando o middleware Auth
	user := middleware.Auth(r)
	if user == nil || user.ID == 0 {
		log.Printf("Usuário não encontrado na sessão")
		return 0
	}

	log.Printf("Usuário encontrado: ID=%d, Email=%s", user.ID, user.Email)

	// Buscar o creator associado ao usuário
	creatorRepository := gorm.NewCreatorRepository(database.DB)
	creator, err := creatorRepository.FindCreatorByUserID(user.ID)
	if err != nil || creator == nil {
		log.Printf("Erro ao buscar creator para usuário %d: %v", user.ID, err)
		return 0
	}

	log.Printf("Creator encontrado: ID=%d, Nome=%s", creator.ID, creator.Name)
	return creator.ID
}
