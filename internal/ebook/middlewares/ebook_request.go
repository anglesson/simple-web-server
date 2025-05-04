package middlewares

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/anglesson/simple-web-server/internal/ebook/models"
	"github.com/anglesson/simple-web-server/internal/shared/utils"
)

func EbookRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errors := make(map[string]string)
		form := models.EbookRequest{
			Title:       r.FormValue("title"),
			Description: r.FormValue("description"),
			Value:       r.FormValue("value"),
			Status:      true,
		}

		log.Println("Iniciando validacao do form")
		errForm := utils.ValidateForm(form)
		for key, value := range errForm {
			errors[key] = value
		}

		errFile := validateFile(r)
		for key, value := range errFile {
			errors[key] = value
		}

		log.Println("Validação do arquivo finalizada")

		fmt.Print(errors)

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

		next.ServeHTTP(w, r)
	})
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
		if len(fileBytes) > 5*1024*1024 { // 5 MB
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
