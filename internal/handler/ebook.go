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

	"github.com/anglesson/simple-web-server/pkg/gov"

	"github.com/anglesson/simple-web-server/internal/handler/web"
	"github.com/anglesson/simple-web-server/internal/repository/gorm"

	"github.com/anglesson/simple-web-server/internal/handler/middleware"
	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/internal/service"
	cookies "github.com/anglesson/simple-web-server/pkg/cookie"
	"github.com/anglesson/simple-web-server/pkg/database"
	"github.com/anglesson/simple-web-server/pkg/storage"
	"github.com/anglesson/simple-web-server/pkg/template"
	"github.com/anglesson/simple-web-server/pkg/utils"
	"github.com/go-chi/chi/v5"
)

func EbookIndexView(w http.ResponseWriter, r *http.Request) {
	loggedUser := middleware.Auth(r)
	if loggedUser.ID == 0 {
		http.Error(w, "Não foi possível prosseguir com a sua solicitação", http.StatusInternalServerError)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	title := r.URL.Query().Get("title")

	pagination := models.NewPagination(page, perPage)

	ebookService := service.NewEbookService()
	ebooks, err := ebookService.ListEbooksForUser(loggedUser.ID, repository.EbookQuery{
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

	template.View(w, r, "ebook/index", map[string]any{
		"Ebooks":     ebooks,
		"Pagination": pagination,
	}, "admin")
}

func EbookCreateView(w http.ResponseWriter, r *http.Request) {
	loggedUser := middleware.Auth(r)
	if loggedUser.ID == 0 {
		http.Error(w, "Não foi possível prosseguir com a sua solicitação", http.StatusInternalServerError)
		return
	}

	// Buscar arquivos do criador para seleção
	creatorRepo := gorm.NewCreatorRepository(database.DB)
	rfService := gov.NewHubDevService()
	userRepository := repository.NewGormUserRepository(database.DB)
	encrypter := utils.NewEncrypter()
	userService := service.NewUserService(userRepository, encrypter)
	subscriptionRepository := gorm.NewSubscriptionGormRepository()
	stripeService := service.NewStripeService()
	creatorService := service.NewCreatorService(
		creatorRepo,
		rfService,
		userService,
		service.NewSubscriptionService(subscriptionRepository, rfService),
		service.NewStripePaymentGateway(stripeService),
	)

	creator, err := creatorService.FindCreatorByUserID(loggedUser.ID)
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
	fileRepository := repository.NewGormFileRepository(database.DB)
	s3Storage := storage.NewS3Storage()
	fileService := service.NewFileService(fileRepository, s3Storage)

	query := repository.FileQuery{
		Pagination: pagination,
	}

	files, total, err := fileService.GetFilesByCreatorPaginated(creator.ID, query)
	if err != nil {
		log.Printf("Erro ao buscar arquivos: %v", err)
		files = []*models.File{} // Lista vazia em caso de erro
		total = 0
	}

	// Configurar paginação com total
	pagination.SetTotal(total)

	template.View(w, r, "ebook/create", map[string]interface{}{
		"Files":      files,
		"Creator":    creator,
		"Pagination": pagination,
	}, "admin")
}

func EbookCreateSubmit(w http.ResponseWriter, r *http.Request) {
	loggedUser := middleware.Auth(r)
	if loggedUser.ID == 0 {
		http.Error(w, "Não foi possível prosseguir com a sua solicitação", http.StatusInternalServerError)
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
		return
	}

	fmt.Printf("Criando e-book para: %v", loggedUser)

	// Busca o criador
	creatorRepo := gorm.NewCreatorRepository(database.DB)
	rfService := gov.NewHubDevService()
	userRepository := repository.NewGormUserRepository(database.DB)
	encrypter := utils.NewEncrypter()
	userService := service.NewUserService(userRepository, encrypter)
	subscriptionRepository := gorm.NewSubscriptionGormRepository()
	stripeService := service.NewStripeService()
	creatorService := service.NewCreatorService(
		creatorRepo,
		rfService,
		userService,
		service.NewSubscriptionService(subscriptionRepository, rfService),
		service.NewStripePaymentGateway(stripeService),
	)
	creator, err := creatorService.FindCreatorByUserID(loggedUser.ID)
	if err != nil {
		log.Printf("Falha ao cadastrar e-book: %s", err)
		web.RedirectBackWithErrors(w, r, "Falha ao cadastrar e-book")
		return
	}

	fmt.Printf("Criando e-book para creator: %v", creator.ID)

	// Processar upload da imagem
	var imageURL string
	imageFile, imageHeader, imageErr := r.FormFile("image")
	if imageErr == nil && imageFile != nil && imageHeader != nil && imageHeader.Filename != "" {
		// Validar se é uma imagem
		contentType := imageHeader.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "image/") {
			errors["image"] = "O arquivo deve ser uma imagem"
		} else {
			// Upload da imagem usando o S3
			s3Storage := storage.NewS3Storage()

			// Gerar nome único para a imagem
			fileExt := filepath.Ext(imageHeader.Filename)
			uniqueID := fmt.Sprintf("%d-%d", time.Now().Unix(), creator.ID)
			imageName := fmt.Sprintf("ebook-covers/%s%s", uniqueID, fileExt)

			// Upload para S3
			imageURL, err = s3Storage.UploadFile(imageHeader, imageName)
			if err != nil {
				log.Printf("Erro ao fazer upload da imagem: %v", err)
				errors["image"] = "Erro ao fazer upload da imagem"
			}
		}
	}

	// Se houve erro no upload da imagem, retornar
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
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		return
	}

	// Criar ebook
	ebook := models.NewEbook(form.Title, form.Description, form.SalesPage, form.Value, *creator)

	// Definir a URL da imagem se foi enviada
	if imageURL != "" {
		ebook.Image = imageURL
	}

	// Adicionar arquivos selecionados ao ebook
	fileRepository := repository.NewGormFileRepository(database.DB)
	for _, fileIDStr := range selectedFiles {
		fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
		if err != nil {
			continue
		}

		file, err := fileRepository.FindByID(uint(fileID))
		if err != nil {
			continue
		}

		// Verificar se o arquivo pertence ao criador
		if file.CreatorID == creator.ID {
			ebook.AddFile(file)
		}
	}

	// Salvar ebook
	ebookRepository := repository.NewGormEbookRepository(database.DB)
	err = ebookRepository.Create(ebook)
	if err != nil {
		log.Printf("Falha ao salvar e-book: %s", err)
		web.RedirectBackWithErrors(w, r, "Falha ao salvar e-book")
		return
	}

	log.Println("E-book criado com sucesso")
	web.RedirectBackWithSuccess(w, r, "E-book criado com sucesso!")
}

// TODO: Move to a service
func validateFile(file multipart.File, expectedContentType string) map[string]string {
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

// TODO: Move to a repository
func GetEbookByID(w http.ResponseWriter, r *http.Request) *models.Ebook {
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

	ebookRepository := repository.NewGormEbookRepository(database.DB)
	ebook, err := ebookRepository.FindByID(uint(id))
	if err != nil {
		http.Error(w, "Erro ao buscar e-book", http.StatusInternalServerError)
		return nil
	}

	return ebook
}

func EbookUpdateView(w http.ResponseWriter, r *http.Request) {
	// Recupera o e-book
	loggedUser := GetSessionUser(r)

	ebook := GetEbookByID(w, r)
	if ebook == nil {
		http.Error(w, "Erro ao buscar e-book", http.StatusNotFound)
		return
	}

	if loggedUser.ID != ebook.Creator.UserID {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}

	// Buscar arquivos disponíveis da biblioteca para adicionar ao ebook
	creatorRepo := gorm.NewCreatorRepository(database.DB)
	rfService := gov.NewHubDevService()
	userRepository := repository.NewGormUserRepository(database.DB)
	encrypter := utils.NewEncrypter()
	userService := service.NewUserService(userRepository, encrypter)
	subscriptionRepository := gorm.NewSubscriptionGormRepository()
	stripeService := service.NewStripeService()
	creatorService := service.NewCreatorService(
		creatorRepo,
		rfService,
		userService,
		service.NewSubscriptionService(subscriptionRepository, rfService),
		service.NewStripePaymentGateway(stripeService),
	)

	creator, err := creatorService.FindCreatorByUserID(loggedUser.ID)
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
	fileRepository := repository.NewGormFileRepository(database.DB)
	s3Storage := storage.NewS3Storage()
	fileService := service.NewFileService(fileRepository, s3Storage)

	query := repository.FileQuery{
		Pagination: pagination,
	}

	allFiles, total, err := fileService.GetFilesByCreatorPaginated(creator.ID, query)
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

	template.View(w, r, "ebook/update", map[string]interface{}{
		"ebook":          ebook,
		"AvailableFiles": availableFiles,
		"Pagination":     pagination,
	}, "admin")
}

func EbookUpdateSubmit(w http.ResponseWriter, r *http.Request) {
	errors := make(map[string]string)

	value, err := utils.BRLToFloat(r.FormValue("value"))
	if err != nil {
		http.Error(w, "erro na conversão", http.StatusInternalServerError)
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
		errFile := validateFile(uploadFile, "application/pdf")
		for key, value := range errFile {
			errors[key] = value
		}
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
		referer := r.Header.Get("Referer")
		if referer == "" {
			referer = "/ebook/edit/" + chi.URLParam(r, "id")
		}
		http.Redirect(w, r, referer, http.StatusSeeOther)
		return
	}

	user := GetSessionUser(r)

	// Verificar se o usuário é um criador
	creatorRepository := gorm.NewCreatorRepository(database.DB)
	_, err = creatorRepository.FindCreatorByUserID(user.ID)
	if err != nil {
		log.Printf("Falha ao buscar criador: %s", err)
		http.Error(w, "Entre em contato", http.StatusInternalServerError)
		return
	}

	ebook := GetEbookByID(w, r)
	if ebook == nil {
		http.Error(w, "E-book não encontrado", http.StatusNotFound)
		return
	}

	// Processar upload da nova imagem
	imageFile, imageHeader, imageErr := r.FormFile("image")
	if imageErr == nil && imageFile != nil && imageHeader != nil && imageHeader.Filename != "" {
		// Validar se é uma imagem
		contentType := imageHeader.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "image/") {
			errors["image"] = "O arquivo deve ser uma imagem"
		} else {
			// Upload da nova imagem usando o S3
			s3Storage := storage.NewS3Storage()

			// Gerar nome único para a imagem
			fileExt := filepath.Ext(imageHeader.Filename)
			uniqueID := fmt.Sprintf("%d-%d", time.Now().Unix(), ebook.CreatorID)
			imageName := fmt.Sprintf("ebook-covers/%s%s", uniqueID, fileExt)

			// Upload para S3
			imageURL, err := s3Storage.UploadFile(imageHeader, imageName)
			if err != nil {
				log.Printf("Erro ao fazer upload da imagem: %v", err)
				errors["image"] = "Erro ao fazer upload da imagem"
			} else {
				// Se o upload foi bem-sucedido, atualizar a URL da imagem
				ebook.Image = imageURL
			}
		}
	}

	// Se houve erro no upload da imagem, retornar
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
		referer := r.Header.Get("Referer")
		if referer == "" {
			referer = "/ebook/edit/" + chi.URLParam(r, "id")
		}
		http.Redirect(w, r, referer, http.StatusSeeOther)
		return
	}

	ebook.Title = form.Title
	ebook.Description = form.Description
	ebook.SalesPage = form.SalesPage
	ebook.Value = form.Value
	ebook.Status = form.Status

	// Processar novos arquivos selecionados
	newFiles := r.Form["new_files"]
	if len(newFiles) > 0 {
		fileRepository := repository.NewGormFileRepository(database.DB)
		for _, fileIDStr := range newFiles {
			fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
			if err != nil {
				continue
			}

			file, err := fileRepository.FindByID(uint(fileID))
			if err != nil {
				continue
			}

			// Verificar se o arquivo pertence ao criador
			if file.CreatorID == ebook.CreatorID {
				ebook.AddFile(file)
			}
		}
	}

	// Salvar usando o repositório
	ebookRepository := repository.NewGormEbookRepository(database.DB)
	err = ebookRepository.Update(ebook)
	if err != nil {
		log.Printf("Falha ao atualizar e-book: %s", err)
		http.Error(w, "Erro ao atualizar e-book", http.StatusInternalServerError)
		return
	}

	cookies.NotifySuccess(w, "Dados do e-book foram atualizados!")
	http.Redirect(w, r, "/ebook", http.StatusSeeOther)
}

// TODO: Move GetSession to a service
func GetSessionUser(r *http.Request) *models.User {
	user_email, ok := r.Context().Value(middleware.UserEmailKey).(string)
	if !ok {
		log.Fatalf("Erro ao recuperar usuário da sessão: %s", user_email)
		return nil
	}
	return repository.NewGormUserRepository(database.DB).FindByEmail(user_email)
}

func EbookShowView(w http.ResponseWriter, r *http.Request) {
	loggedUser := middleware.Auth(r)
	if loggedUser.ID == 0 {
		http.Error(w, "Não foi possível prosseguir com a sua solicitação", http.StatusInternalServerError)
		return
	}

	ebook := GetEbookByID(w, r)
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

	creatorRepository := gorm.NewCreatorRepository(database.DB)
	creator, err := creatorRepository.FindCreatorByUserID(loggedUser.ID)
	if err != nil {
		web.RedirectBackWithErrors(w, r, err.Error())
	}

	clients, err := gorm.NewClientGormRepository().FindByClientsWhereEbookWasSend(creator, models.ClientFilter{
		Term:       term,
		EbookID:    ebook.ID,
		Pagination: pagination,
	})
	if err != nil {
		web.RedirectBackWithErrors(w, r, err.Error())
	}

	// Set total count for pagination
	if clients != nil {
		pagination.SetTotal(int64(len(*clients)))
	}

	template.View(w, r, "ebook/view", map[string]any{
		"Ebook":      ebook,
		"Clients":    clients,
		"Pagination": pagination,
	}, "admin")
}
