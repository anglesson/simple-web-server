package handlers

import (
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repositories"
	"github.com/anglesson/simple-web-server/internal/services"
	cookies "github.com/anglesson/simple-web-server/internal/shared/cookie"
	"github.com/anglesson/simple-web-server/internal/shared/database"
	"github.com/anglesson/simple-web-server/internal/shared/middlewares"
	"github.com/anglesson/simple-web-server/internal/shared/storage"
	"github.com/anglesson/simple-web-server/internal/shared/template"
	"github.com/anglesson/simple-web-server/internal/shared/utils"
	"github.com/go-chi/chi/v5"
)

func EbookIndexView(w http.ResponseWriter, r *http.Request) {
	loggedUser := middlewares.Auth(r)
	if loggedUser.ID == 0 {
		http.Error(w, "Não foi possível prosseguir com a sua solicitação", http.StatusInternalServerError)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	title := r.URL.Query().Get("title")
	pagination := repositories.NewPagination(page, perPage)

	ebookService := services.NewEbookService()
	ebooks, err := ebookService.ListEbooksForUser(loggedUser.ID, repositories.EbookQuery{
		Title:      title,
		Pagination: pagination,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	template.View(w, r, "ebook", map[string]any{
		"Ebooks":     ebooks,
		"Pagination": pagination,
	}, "admin")
}

func EbookCreateView(w http.ResponseWriter, r *http.Request) {
	template.View(w, r, "create_ebook", nil, "admin")
}

func EbookCreateSubmit(w http.ResponseWriter, r *http.Request) {
	errors := make(map[string]string)

	value, err := utils.BRLToFloat(r.FormValue("value"))
	if err != nil {
		http.Error(w, "erro na conversão", http.StatusInternalServerError)
	}
	form := models.EbookRequest{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Value:       value,
		Status:      true,
	}

	errForm := utils.ValidateForm(form)
	for key, value := range errForm {
		errors[key] = value
	}

	file, _, err := r.FormFile("file")
	if err == nil {
		errors["file"] = "Arquivo é obrigatório"
	} else {
		errFile := validateFile(file, "application/pdf")
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
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		return
	}

	user_email, ok := r.Context().Value(middlewares.UserEmailKey).(string)
	if !ok {
		http.Error(w, "Invalid user email", http.StatusInternalServerError)
		return
	}

	user := repositories.NewUserRepository().FindByEmail(user_email)
	creator := models.Creator{
		UserID: user.ID,
	}
	result := database.DB.First(&creator)

	if result.Error != nil {
		log.Printf("Falha ao cadastrar ebook: %s", result.Error)
		http.Error(w, "Entre em contato", http.StatusInternalServerError)
		return
	}

	// Obtenha o arquivo do formulário
	file, fileHeader, err := r.FormFile("file") // "file" deve ser o nome do campo no HTML
	if err != nil {
		log.Printf("Erro ao obter arquivo: %v", err)
		http.Error(w, "Erro ao obter arquivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	locationFile, _ := storage.Upload(file, fileHeader.Filename)
	ebook := models.NewEbook(form.Title, form.Description, locationFile, form.Value, creator)

	database.DB.Create(&ebook)

	http.Redirect(w, r, "/ebook", http.StatusSeeOther)
}

// TODO: Move to a service
func validateFile(file multipart.File, expectedContentType string) map[string]string {
	errors := make(map[string]string)

	defer file.Close()

	// Validar tamanho
	fileBytes, _ := io.ReadAll(file)
	if len(fileBytes) > 60*1024*1024 { // 60 MB
		errors["File"] = "Arquivo deve ter no máximo 5 MB"
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
	var ebook models.Ebook

	// Busca o criador com os ebooks associados
	err := database.DB.
		Preload("Creator").
		Where("id = ?", chi.URLParam(r, "id")).
		First(&ebook).Error
	if err != nil {
		http.Error(w, "Erro ao buscar ebook", http.StatusInternalServerError)
		return nil
	}

	return &ebook
}

func EbookUpdateView(w http.ResponseWriter, r *http.Request) {
	// Recupera o ebook
	loggedUser := GetSessionUser(r)

	ebook := GetEbookByID(w, r)
	if ebook == nil {
		http.Error(w, "Erro ao buscar ebook", http.StatusNotFound)
		return
	}

	if loggedUser.ID != ebook.Creator.UserID {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}

	ebook.FileURL = storage.GenerateDownloadLink(ebook.File)

	template.View(w, r, "update_ebook", ebook, "admin")
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
		Value:       value,
		Status:      status,
	}

	errForm := utils.ValidateForm(form)
	for key, value := range errForm {
		errors[key] = value
	}

	file, _, err := r.FormFile("file")
	if err == nil {
		errFile := validateFile(file, "application/pdf")
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

	creator := models.Creator{
		UserID: user.ID,
	}
	result := database.DB.First(&creator)

	if result.Error != nil {
		log.Printf("Falha ao cadastrar ebook: %s", result.Error)
		http.Error(w, "Entre em contato", http.StatusInternalServerError)
		return
	}

	ebook := GetEbookByID(w, r)

	// Obtenha o arquivo do formulário
	file, fileHeader, err := r.FormFile("file") // "file" deve ser o nome do campo no HTML
	if err == nil {
		defer file.Close()
		storage.Upload(file, fileHeader.Filename)
		ebook.FileURL = storage.GenerateDownloadLink(fileHeader.Filename)
		ebook.File = fileHeader.Filename
	}

	ebook.Title = form.Title
	ebook.Description = form.Description
	ebook.Value = form.Value
	ebook.Status = form.Status

	database.DB.Save(&ebook)

	cookies.NotifySuccess(w, "Dados do ebook foram atualizados!")
	http.Redirect(w, r, "/ebook", http.StatusSeeOther)
}

// TODO: Move GetSession to a service
func GetSessionUser(r *http.Request) *models.User {
	user_email, ok := r.Context().Value(middlewares.UserEmailKey).(string)
	if !ok {
		log.Fatalf("Erro ao recuperar usuário da sessão: %s", user_email)
		return nil
	}
	return repositories.NewUserRepository().FindByEmail(user_email)
}
