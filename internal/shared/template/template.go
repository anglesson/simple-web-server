package template

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/anglesson/simple-web-server/internal/config"
)

type PageData struct {
	ErrorMessage string
}

// TemplateFunctions returns a map of functions available to templates
func TemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"appName": func() string {
			return config.AppConfig.AppName
		},
	}
}

func View(w http.ResponseWriter, r *http.Request, templateName string, data any, layout string) {
	var form map[string]interface{}
	var errors map[string]string

	if c, err := r.Cookie("form"); err == nil {
		decodedValue, decodeErr := url.QueryUnescape(c.Value)
		if decodeErr != nil {
			log.Println("Error decoding cookie value:", decodeErr)
		}
		_ = json.Unmarshal([]byte(decodedValue), &form)
		http.SetCookie(w, &http.Cookie{Name: "form", MaxAge: -1})
	}
	if c, err := r.Cookie("errors"); err == nil {
		decodedValue, decodeErr := url.QueryUnescape(c.Value)
		if decodeErr != nil {
			log.Println("Error decoding cookie value:", decodeErr)
		}
		_ = json.Unmarshal([]byte(decodedValue), &errors)
		http.SetCookie(w, &http.Cookie{Name: "errors", MaxAge: -1})
	}

	files := []string{
		fmt.Sprintf("internal/templates/layouts/%s.html", layout),
		fmt.Sprintf("internal/templates/pages/%s.html", templateName),
	}

	t := template.New(layout + ".html").Funcs(TemplateFunctions())

	t, err := t.ParseFiles(files...)
	if err != nil {
		http.Error(w, "Erro ao carregar template", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, map[string]any{
		"Form":   form,
		"Errors": errors,
		"Data":   data,
	})
	if err != nil {
		log.Printf("Erro ao renderizar template %s", err)
		return
	}
}
