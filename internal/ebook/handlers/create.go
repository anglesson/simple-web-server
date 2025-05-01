package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	auth "github.com/anglesson/simple-web-server/internal/auth/models"
	"github.com/anglesson/simple-web-server/internal/ebook/models"
	"github.com/anglesson/simple-web-server/internal/shared/template"
	"github.com/go-playground/validator/v10"
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
	var form models.EbookRequest
	var errors auth.FormErrors

	if c, err := r.Cookie("form"); err == nil {
		decodedValue, decodeErr := url.QueryUnescape(c.Value) // Decodifica o valor do cookie
		if decodeErr != nil {
			log.Println("Error decoding cookie value:", decodeErr)
		}
		_ = json.Unmarshal([]byte(decodedValue), &form)
		http.SetCookie(w, &http.Cookie{Name: "form", MaxAge: -1})
	}
	if c, err := r.Cookie("errors"); err == nil {
		decodedValue, decodeErr := url.QueryUnescape(c.Value) // Decodifica o valor do cookie
		if decodeErr != nil {
			log.Println("Error decoding cookie value:", decodeErr)
		}
		_ = json.Unmarshal([]byte(decodedValue), &errors)
		http.SetCookie(w, &http.Cookie{Name: "errors", MaxAge: -1})
	}
	template.View(w, "create_ebook", map[string]any{
		"Form":   form,
		"Errors": errors,
	}, "base_logged")
}

func processCreateEbook(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	fmt.Println(r.FormValue("title"))

	form := models.EbookRequest{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Value:       r.FormValue("value"),
	}

	errors := validateForm(form)

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
		http.Redirect(w, r, "/ebook/create", http.StatusSeeOther)
		return
	}
}

func validateForm(form interface{}) map[string]string {
	validate := validator.New()
	err := validate.Struct(form)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			// Create a map to store field-specific error messages
			errors := make(map[string]string)
			for _, e := range validationErrors {
				// Add the field name and its error message to the map
				switch e.Tag() {
				case "required":
					errors[e.Field()] = "Preenchimento obrigatório"
				case "min":
					errors[e.Field()] = fmt.Sprintf("Digite no mínimo %s caracteres", e.Param())
				case "max":
					errors[e.Field()] = fmt.Sprintf("Digite no máximo %s caracteres", e.Param())
				default:
					errors[e.Field()] = "Revise este campo"
				}

			}
			return errors
		}
	}
	return nil
}
