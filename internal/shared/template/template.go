package template

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/anglesson/simple-web-server/internal/config"
	"github.com/anglesson/simple-web-server/internal/models"
	cookies "github.com/anglesson/simple-web-server/internal/shared/cookie"
	"github.com/anglesson/simple-web-server/internal/shared/middlewares"
)

type PageData struct {
	ErrorMessage string
}

// TemplateFunctions returns a map of functions available to templates
func TemplateFunctions(r *http.Request) template.FuncMap {
	return template.FuncMap{
		"appName": func() string {
			return config.AppConfig.AppName
		},
		"user": func() *models.User {
			return middlewares.Auth(r)
		},
	}
}

func View(w http.ResponseWriter, r *http.Request, templateName string, data any, layout string) {
	var form map[string]interface{}
	var errors map[string]string
	var flash cookies.FlashMessage

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

	if c, err := r.Cookie("flash"); err == nil {
		decoded, _ := url.QueryUnescape(c.Value)
		_ = json.Unmarshal([]byte(decoded), &flash)
		http.SetCookie(w, &http.Cookie{Name: "flash", MaxAge: -1})
	}

	files := []string{
		fmt.Sprintf("internal/templates/layouts/%s.html", layout),
		fmt.Sprintf("internal/templates/pages/%s.html", templateName),
	}
	partials, _ := filepath.Glob("internal/templates/partials/*.html")
	files = append(files, partials...)

	t := template.New(layout + ".html").Funcs(TemplateFunctions(r))
	t, err := t.ParseFiles(files...)
	if err != nil {
		http.Error(w, "Erro ao carregar template", http.StatusInternalServerError)
		log.Println("Template parse error:", err)
		return
	}

	err = t.ExecuteTemplate(w, layout+".html", map[string]any{
		"Form":   form,
		"Errors": errors,
		"Flash":  flash,
		"Data":   data,
	})
	if err != nil {
		log.Printf("Erro ao renderizar template: %s", err)
	}
}
