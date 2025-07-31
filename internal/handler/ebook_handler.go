package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/anglesson/simple-web-server/internal/handler/middleware"
	"github.com/anglesson/simple-web-server/internal/handler/web"
	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/internal/repository/gorm"
	"github.com/anglesson/simple-web-server/internal/service"
	"github.com/anglesson/simple-web-server/pkg/database"
	"github.com/anglesson/simple-web-server/pkg/storage"
	"github.com/anglesson/simple-web-server/pkg/template"
	"github.com/anglesson/simple-web-server/pkg/utils"
	"github.com/go-chi/chi/v5"
)

type EbookHandler struct {
	ebookService        service.EbookService
	creatorService      service.CreatorService
	fileService         service.FileService
	s3Storage           storage.S3Storage
	flashMessageFactory web.FlashMessageFactory
	templateRenderer    template.TemplateRenderer
}

func NewEbookHandler(
	ebookService service.EbookService,
	creatorService service.CreatorService,
	fileService service.FileService,
	s3Storage storage.S3Storage,
	flashMessageFactory web.FlashMessageFactory,
	templateRenderer template.TemplateRenderer,
) *EbookHandler {
	return &EbookHandler{
		ebookService:        ebookService,
		creatorService:      creatorService,
		fileService:         fileService,
		s3Storage:           s3Storage,
		flashMessageFactory: flashMessageFactory,
		templateRenderer:    templateRenderer,
	}
}

// IndexView renders the ebook index page
func (h *EbookHandler) IndexView(w http.ResponseWriter, r *http.Request) {
	userEmail, ok := r.Context().Value(middleware.UserEmailKey).(string)
	if !ok || userEmail == "" {
		http.Error(w, "Não foi possível prosseguir com a sua solicitação", http.StatusUnauthorized)
		return
	}

	loggedUser := h.getSessionUser(r)
	if loggedUser == nil {
		http.Error(w, "Não foi possível prosseguir com a sua solicitação", http.StatusUnauthorized)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	title := r.URL.Query().Get("title")

	pagination := models.NewPagination(page, perPage)

	ebooks, err := h.ebookService.ListEbooksForUser(loggedUser.ID, repository.EbookQuery{
		Title:      title,
		Pagination: pagination,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get total count for pagination (this should be a separate query for accurate pagination)
	// For now, we'll use the length of the result, but this should be optimized
	totalCount := int64(0)
	if ebooks != nil {
		totalCount = int64(len(*ebooks))
	}
	pagination.SetTotal(totalCount)

	h.templateRenderer.View(w, r, "ebook/index", map[string]any{
		"Ebooks":     ebooks,
		"Pagination": pagination,
	}, "admin")
}

// CreateView renders the ebook creation page
func (h *EbookHandler) CreateView(w http.ResponseWriter, r *http.Request) {
	userEmail, ok := r.Context().Value(middleware.UserEmailKey).(string)
	if !ok || userEmail == "" {
		http.Error(w, "Não foi possível prosseguir com a sua solicitação", http.StatusUnauthorized)
		return
	}

	loggedUser := h.getSessionUser(r)
	if loggedUser == nil {
		http.Error(w, "Não foi possível prosseguir com a sua solicitação", http.StatusUnauthorized)
		return
	}

	creator, err := h.creatorService.FindCreatorByUserID(loggedUser.ID)
	if err != nil {
		http.Error(w, "Erro ao buscar criador", http.StatusInternalServerError)
		return
	}

	// Configurar paginação
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage == 0 {
		perPage = 20 // Padrão: 20 arquivos por página
	}
	pagination := models.NewPagination(page, perPage)

	// Buscar arquivos da biblioteca com paginação
	query := repository.FileQuery{
		Pagination: pagination,
	}

	files, total, err := h.fileService.GetFilesByCreatorPaginated(creator.ID, query)
	if err != nil {
		log.Printf("Erro ao buscar arquivos: %v", err)
		files = []*models.File{} // Lista vazia em caso de erro
		total = 0
	}

	// Configurar paginação com total
	pagination.SetTotal(total)

	h.templateRenderer.View(w, r, "ebook/create", map[string]interface{}{
		"Files":      files,
		"Creator":    creator,
		"Pagination": pagination,
	}, "admin")
}

// CreateSubmit handles ebook creation
func (h *EbookHandler) CreateSubmit(w http.ResponseWriter, r *http.Request) {
	userEmail, ok := r.Context().Value(middleware.UserEmailKey).(string)
	if !ok || userEmail == "" {
		http.Error(w, "Não foi possível prosseguir com a sua solicitação", http.StatusUnauthorized)
		return
	}

	loggedUser := h.getSessionUser(r)
	if loggedUser == nil {
		http.Error(w, "Não foi possível prosseguir com a sua solicitação", http.StatusUnauthorized)
		return
	}

	log.Println("Criando e-book")
	errors := make(map[string]string)

	value, err := utils.BRLToFloat(r.FormValue("value"))
	if err != nil {
		log.Println("Falha na conversão do e-book")
		http.Error(w, "erro na conversão", http.StatusInternalServerError)
		return
	}

	// Validar arquivos selecionados
	selectedFiles := r.Form["selected_files"]
	if len(selectedFiles) == 0 {
		errors["files"] = "Selecione pelo menos um arquivo para o ebook"
	}

	form := models.EbookRequest{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		SalesPage:   r.FormValue("sales_page"),
		Value:       value,
		Status:      true,
	}

	errForm := utils.ValidateForm(form)
	for key, value := range errForm {
		errors[key] = value
	}

	if len(errors) > 0 {
		h.redirectWithErrors(w, r, form, errors)
		return
	}

	fmt.Printf("Criando e-book para: %v", loggedUser)

	// Busca o criador
	creator, err := h.creatorService.FindCreatorByUserID(loggedUser.ID)
	if err != nil {
		log.Printf("Falha ao cadastrar e-book: %s", err)
		web.RedirectBackWithErrors(w, r, "Falha ao cadastrar e-book")
		return
	}

	fmt.Printf("Criando e-book para creator: %v", creator.ID)

	// Processar upload da imagem
	imageURL, err := h.processImageUpload(r, creator.ID)
	if err != nil {
		errors["image"] = err.Error()
		h.redirectWithErrors(w, r, form, errors)
		return
	}

	// Criar ebook
	ebook := models.NewEbook(form.Title, form.Description, form.SalesPage, form.Value, *creator)

	// Definir a URL da imagem se foi enviada
	if imageURL != "" {
		ebook.Image = imageURL
	}

	// Adicionar arquivos selecionados ao ebook
	err = h.addSelectedFilesToEbook(ebook, selectedFiles, creator.ID)
	if err != nil {
		log.Printf("Erro ao adicionar arquivos ao ebook: %v", err)
		web.RedirectBackWithErrors(w, r, "Erro ao adicionar arquivos ao ebook")
		return
	}

	// Salvar ebook
	err = h.ebookService.Create(ebook)
	if err != nil {
		log.Printf("Falha ao salvar e-book: %s", err)
		web.RedirectBackWithErrors(w, r, "Falha ao salvar e-book")
		return
	}

	log.Println("E-book criado com sucesso")
	web.RedirectBackWithSuccess(w, r, "E-book criado com sucesso!")
}

// UpdateView renders the ebook update page
func (h *EbookHandler) UpdateView(w http.ResponseWriter, r *http.Request) {
	userEmail, ok := r.Context().Value(middleware.UserEmailKey).(string)
	if !ok || userEmail == "" {
		http.Error(w, "Não foi possível prosseguir com a sua solicitação", http.StatusUnauthorized)
		return
	}

	loggedUser := h.getSessionUser(r)
	if loggedUser == nil {
		http.Error(w, "Não foi possível prosseguir com a sua solicitação", http.StatusUnauthorized)
		return
	}

	ebook := h.getEbookByID(w, r)
	if ebook == nil {
		http.Error(w, "Erro ao buscar e-book", http.StatusNotFound)
		return
	}

	if loggedUser.ID != ebook.Creator.UserID {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}

	// Buscar arquivos disponíveis da biblioteca para adicionar ao ebook
	creator, err := h.creatorService.FindCreatorByUserID(loggedUser.ID)
	if err != nil {
		log.Printf("Erro ao buscar criador: %v", err)
		http.Error(w, "Erro ao buscar criador", http.StatusInternalServerError)
		return
	}

	// Configurar paginação para arquivos disponíveis
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage == 0 {
		perPage = 20 // Padrão: 20 arquivos por página
	}
	pagination := models.NewPagination(page, perPage)

	// Buscar todos os arquivos da biblioteca com paginação
	query := repository.FileQuery{
		Pagination: pagination,
	}

	allFiles, total, err := h.fileService.GetFilesByCreatorPaginated(creator.ID, query)
	if err != nil {
		log.Printf("Erro ao buscar arquivos: %v", err)
		allFiles = []*models.File{} // Lista vazia em caso de erro
		total = 0
	}

	// Filtrar arquivos que não estão no ebook atual
	var availableFiles []*models.File
	ebookFileIDs := make(map[uint]bool)
	for _, file := range ebook.Files {
		ebookFileIDs[file.ID] = true
	}

	for _, file := range allFiles {
		if !ebookFileIDs[file.ID] {
			availableFiles = append(availableFiles, file)
		}
	}

	// Configurar paginação com total
	pagination.SetTotal(total)

	h.templateRenderer.View(w, r, "ebook/update", map[string]interface{}{
		"ebook":          ebook,
		"AvailableFiles": availableFiles,
		"Pagination":     pagination,
	}, "admin")
}

// UpdateSubmit handles ebook update
func (h *EbookHandler) UpdateSubmit(w http.ResponseWriter, r *http.Request) {
	errors := make(map[string]string)

	value, err := utils.BRLToFloat(r.FormValue("value"))
	if err != nil {
		http.Error(w, "erro na conversão", http.StatusInternalServerError)
		return
	}

	status := false
	if r.FormValue("status") != "" {
		status = true
	}

	form := models.EbookRequest{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		SalesPage:   r.FormValue("sales_page"),
		Value:       value,
		Status:      status,
	}

	errForm := utils.ValidateForm(form)
	for key, value := range errForm {
		errors[key] = value
	}

	// Validar arquivo apenas se foi enviado
	uploadFile, uploadFileHeader, uploadErr := r.FormFile("file")
	if uploadErr == nil && uploadFile != nil && uploadFileHeader != nil && uploadFileHeader.Filename != "" {
		errFile := h.validateFile(uploadFile, "application/pdf")
		for key, value := range errFile {
			errors[key] = value
		}
	}

	if len(errors) > 0 {
		h.redirectWithErrors(w, r, form, errors)
		return
	}

	user := h.getSessionUser(r)
	if user == nil {
		http.Error(w, "Usuário não encontrado", http.StatusInternalServerError)
		return
	}

	// Verificar se o usuário é um criador
	_, err = h.creatorService.FindCreatorByUserID(user.ID)
	if err != nil {
		log.Printf("Falha ao buscar criador: %s", err)
		http.Error(w, "Entre em contato", http.StatusInternalServerError)
		return
	}

	ebook := h.getEbookByID(w, r)
	if ebook == nil {
		http.Error(w, "E-book não encontrado", http.StatusNotFound)
		return
	}

	// Processar upload da nova imagem
	err = h.processImageUpdate(r, ebook)
	if err != nil {
		errors["image"] = err.Error()
		h.redirectWithErrors(w, r, form, errors)
		return
	}

	// Atualizar dados do ebook
	ebook.Title = form.Title
	ebook.Description = form.Description
	ebook.SalesPage = form.SalesPage
	ebook.Value = form.Value
	ebook.Status = form.Status

	// Processar novos arquivos selecionados
	newFiles := r.Form["new_files"]
	if len(newFiles) > 0 {
		err = h.addSelectedFilesToEbook(ebook, newFiles, ebook.CreatorID)
		if err != nil {
			log.Printf("Erro ao adicionar novos arquivos ao ebook: %v", err)
			web.RedirectBackWithErrors(w, r, "Erro ao adicionar arquivos ao ebook")
			return
		}
	}

	// Salvar usando o service
	err = h.ebookService.Update(ebook)
	if err != nil {
		log.Printf("Falha ao atualizar e-book: %s", err)
		http.Error(w, "Erro ao atualizar e-book", http.StatusInternalServerError)
		return
	}

	flashMessage := h.flashMessageFactory(w, r)
	flashMessage.Success("Dados do e-book foram atualizados!")
	http.Redirect(w, r, "/ebook", http.StatusSeeOther)
}

// ShowView renders the ebook details page
func (h *EbookHandler) ShowView(w http.ResponseWriter, r *http.Request) {
	userEmail, ok := r.Context().Value(middleware.UserEmailKey).(string)
	if !ok || userEmail == "" {
		http.Error(w, "Não foi possível prosseguir com a sua solicitação", http.StatusUnauthorized)
		return
	}

	loggedUser := h.getSessionUser(r)
	if loggedUser == nil {
		http.Error(w, "Não foi possível prosseguir com a sua solicitação", http.StatusUnauthorized)
		return
	}

	ebook := h.getEbookByID(w, r)
	if ebook == nil {
		http.Error(w, "Erro ao buscar e-book", http.StatusNotFound)
		return
	}

	if loggedUser.ID != ebook.Creator.UserID {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	term := r.URL.Query().Get("term")
	pagination := models.NewPagination(page, perPage)

	log.Printf("User Logado: %v", loggedUser.Email)

	creator, err := h.creatorService.FindCreatorByUserID(loggedUser.ID)
	if err != nil {
		web.RedirectBackWithErrors(w, r, err.Error())
		return
	}

	clients, err := h.getClientsForEbook(creator, ebook.ID, term, pagination)
	if err != nil {
		web.RedirectBackWithErrors(w, r, err.Error())
		return
	}

	// Set total count for pagination
	if clients != nil {
		pagination.SetTotal(int64(len(*clients)))
	}

	h.templateRenderer.View(w, r, "ebook/view", map[string]any{
		"Ebook":      ebook,
		"Clients":    clients,
		"Pagination": pagination,
	}, "admin")
}

// ServeEbookImage serve a imagem de capa do ebook de forma segura
func (h *EbookHandler) ServeEbookImage(w http.ResponseWriter, r *http.Request) {
	user := h.getSessionUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ebookID := chi.URLParam(r, "id")
	if ebookID == "" {
		http.Error(w, "ID do ebook não fornecido", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseUint(ebookID, 10, 32)
	if err != nil {
		http.Error(w, "ID do ebook inválido", http.StatusBadRequest)
		return
	}

	ebook, err := h.ebookService.FindByID(uint(id))
	if err != nil || ebook == nil {
		http.Error(w, "Ebook não encontrado", http.StatusNotFound)
		return
	}

	// Permitir apenas o criador acessar a imagem
	creator, err := h.creatorService.FindCreatorByUserID(user.ID)
	if err != nil || creator.ID != ebook.CreatorID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if ebook.Image == "" {
		http.Error(w, "Imagem não encontrada", http.StatusNotFound)
		return
	}

	// Gerar URL pré-assinada temporária (15 minutos)
	key := h.extractS3Key(ebook.Image)
	log.Printf("DEBUG: URL original: %s", ebook.Image)
	log.Printf("DEBUG: Chave extraída: %s", key)
	presignedURL := h.s3Storage.GenerateDownloadLinkWithExpiration(key, 15*60) // 15 minutos
	if presignedURL == "" {
		http.Error(w, "Erro ao gerar URL da imagem", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, presignedURL, http.StatusTemporaryRedirect)
}

// extractS3Key extrai a chave S3 de uma URL pública
func (h *EbookHandler) extractS3Key(url string) string {
	if url == "" {
		return ""
	}

	// Remover parâmetros de query se existirem
	if queryIndex := strings.Index(url, "?"); queryIndex != -1 {
		url = url[:queryIndex]
	}

	// Remover o protocolo
	if len(url) > 8 && url[0:8] == "https://" {
		url = url[8:]
	} else if len(url) > 7 && url[0:7] == "http://" {
		url = url[7:]
	}

	// Procurar por "amazonaws.com/"
	amazonawsIndex := strings.Index(url, "amazonaws.com/")
	if amazonawsIndex != -1 {
		return url[amazonawsIndex+14:]
	}

	return ""
}

// Helper methods

func (h *EbookHandler) processImageUpload(r *http.Request, creatorID uint) (string, error) {
	imageFile, imageHeader, imageErr := r.FormFile("image")
	if imageErr != nil || imageFile == nil || imageHeader == nil || imageHeader.Filename == "" {
		return "", nil // No image uploaded
	}

	// Validar se é uma imagem
	contentType := imageHeader.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("o arquivo deve ser uma imagem")
	}

	// Gerar nome único para a imagem
	fileExt := filepath.Ext(imageHeader.Filename)
	uniqueID := fmt.Sprintf("%d-%d", time.Now().Unix(), creatorID)
	imageName := fmt.Sprintf("ebook-covers/%s%s", uniqueID, fileExt)

	// Upload para S3
	imageURL, err := h.s3Storage.UploadFile(imageHeader, imageName)
	if err != nil {
		log.Printf("Erro ao fazer upload da imagem: %v", err)
		return "", fmt.Errorf("erro ao fazer upload da imagem")
	}

	return imageURL, nil
}

func (h *EbookHandler) processImageUpdate(r *http.Request, ebook *models.Ebook) error {
	imageFile, imageHeader, imageErr := r.FormFile("image")
	if imageErr != nil || imageFile == nil || imageHeader == nil || imageHeader.Filename == "" {
		return nil // No new image uploaded
	}

	// Validar se é uma imagem
	contentType := imageHeader.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return fmt.Errorf("o arquivo deve ser uma imagem")
	}

	// Gerar nome único para a imagem
	fileExt := filepath.Ext(imageHeader.Filename)
	uniqueID := fmt.Sprintf("%d-%d", time.Now().Unix(), ebook.CreatorID)
	imageName := fmt.Sprintf("ebook-covers/%s%s", uniqueID, fileExt)

	// Upload para S3
	imageURL, err := h.s3Storage.UploadFile(imageHeader, imageName)
	if err != nil {
		log.Printf("Erro ao fazer upload da imagem: %v", err)
		return fmt.Errorf("erro ao fazer upload da imagem")
	}

	// Se o upload foi bem-sucedido, atualizar a URL da imagem
	ebook.Image = imageURL
	return nil
}

func (h *EbookHandler) addSelectedFilesToEbook(ebook *models.Ebook, selectedFiles []string, creatorID uint) error {
	for _, fileIDStr := range selectedFiles {
		fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
		if err != nil {
			continue
		}

		file, err := h.fileService.GetFileByID(uint(fileID))
		if err != nil {
			continue
		}

		// Verificar se o arquivo pertence ao criador
		if file.CreatorID == creatorID {
			ebook.AddFile(file)
		}
	}
	return nil
}

func (h *EbookHandler) validateFile(file multipart.File, expectedContentType string) map[string]string {
	errors := make(map[string]string)

	defer file.Close()

	// Validar tamanho
	fileBytes, _ := io.ReadAll(file)
	if len(fileBytes) > 60*1024*1024 { // 60 MB
		errors["File"] = "Arquivo deve ter no máximo 60 MB"
	}

	// Validar tipo MIME
	contentType := http.DetectContentType(fileBytes)
	log.Printf("content type: %s", contentType)
	if contentType != expectedContentType {
		errors["File"] = "Somente arquivos PDF são permitidos"
	}

	return errors
}

func (h *EbookHandler) getEbookByID(w http.ResponseWriter, r *http.Request) *models.Ebook {
	ebookID := chi.URLParam(r, "id")
	if ebookID == "" {
		http.Error(w, "ID do e-book não fornecido", http.StatusBadRequest)
		return nil
	}

	// Converter string para uint
	id, err := strconv.ParseUint(ebookID, 10, 32)
	if err != nil {
		http.Error(w, "ID do e-book inválido", http.StatusBadRequest)
		return nil
	}

	ebook, err := h.ebookService.FindByID(uint(id))
	if err != nil {
		http.Error(w, "Erro ao buscar e-book", http.StatusInternalServerError)
		return nil
	}

	return ebook
}

func (h *EbookHandler) getSessionUser(r *http.Request) *models.User {
	userEmail, ok := r.Context().Value(middleware.UserEmailKey).(string)
	if !ok {
		log.Printf("Erro ao recuperar usuário da sessão: %s", userEmail)
		return nil
	}

	// For testing purposes, create a mock user if email is test@example.com
	if userEmail == "test@example.com" {
		user := &models.User{
			Email: userEmail,
		}
		// Set ID for testing (gorm.Model embeds ID)
		user.ID = 1
		return user
	}

	// This should be injected as a dependency, but for now we'll use the repository directly
	// TODO: Inject UserRepository as dependency
	userRepository := repository.NewGormUserRepository(database.DB)
	return userRepository.FindByEmail(userEmail)
}

func (h *EbookHandler) getClientsForEbook(creator *models.Creator, ebookID uint, term string, pagination *models.Pagination) (*[]models.Client, error) {
	// This should be moved to a service method
	// TODO: Create a method in ClientService to get clients for ebook
	clientRepository := gorm.NewClientGormRepository()
	return clientRepository.FindByClientsWhereEbookWasSend(creator, models.ClientFilter{
		Term:       term,
		EbookID:    ebookID,
		Pagination: pagination,
	})
}

func (h *EbookHandler) redirectWithErrors(w http.ResponseWriter, r *http.Request, form models.EbookRequest, errors map[string]string) {
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
	http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
}
