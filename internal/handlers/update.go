package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repositories"
	"github.com/anglesson/simple-web-server/internal/shared/database"
	"github.com/anglesson/simple-web-server/internal/shared/middlewares"
	"github.com/anglesson/simple-web-server/internal/shared/storage"
	"github.com/anglesson/simple-web-server/internal/shared/template"
	"github.com/anglesson/simple-web-server/internal/shared/utils"
	"github.com/go-chi/chi/v5"
)

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
		errFile := validateFile(file)
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

	http.Redirect(w, r, "/ebook", http.StatusSeeOther)
}

func GetSessionUser(r *http.Request) *models.User {
	user_email, ok := r.Context().Value(middlewares.UserEmailKey).(string)
	if !ok {
		log.Fatalf("Erro ao recuperar usuário da sessão: %s", user_email)
		return nil
	}

	return repositories.FindByEmail(user_email)
}
