package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/anglesson/simple-web-server/internal/auth/repositories"
	"github.com/anglesson/simple-web-server/internal/ebook/models"
	"github.com/anglesson/simple-web-server/internal/shared/database"
	"github.com/anglesson/simple-web-server/internal/shared/middlewares"
	"github.com/anglesson/simple-web-server/internal/shared/storage"
	"github.com/anglesson/simple-web-server/internal/shared/template"
	"github.com/anglesson/simple-web-server/internal/shared/utils"
)

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		renderCreatePage(w, r)
	case http.MethodPost:
		processCreateEbook(w, r)
	}
}

func renderCreatePage(w http.ResponseWriter, r *http.Request) {
	template.View(w, r, "create_ebook", nil, "base_logged")
}

func processCreateEbook(w http.ResponseWriter, r *http.Request) {
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

	errFile := validateFile(r)
	for key, value := range errFile {
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

	user_email, ok := r.Context().Value(middlewares.UserEmailKey).(string)
	if !ok {
		http.Error(w, "Invalid user email", http.StatusInternalServerError)
		return
	}

	user := repositories.FindByEmail(user_email)
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

func validateFile(r *http.Request) map[string]string {
	errors := make(map[string]string)
	file, _, err := r.FormFile("file")
	if err != nil {
		errors["File"] = "Arquivo é obrigatório"
	} else {
		defer file.Close()

		// Validar tamanho
		fileBytes, _ := io.ReadAll(file)
		if len(fileBytes) > 60*1024*1024 { // 60 MB
			errors["File"] = "Arquivo deve ter no máximo 5 MB"
		}

		// Validar tipo MIME
		contentType := http.DetectContentType(fileBytes)
		if contentType != "application/pdf" {
			errors["File"] = "Somente arquivos PDF são permitidos"
		}
	}

	return errors
}
