package template

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/anglesson/simple-web-server/internal/config"
	"github.com/anglesson/simple-web-server/internal/handlers/middleware"
	"github.com/anglesson/simple-web-server/internal/models"
	cookies "github.com/anglesson/simple-web-server/internal/shared/cookie"
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
			return middleware.Auth(r)
		},
		"json": func(data any) (template.JS, error) {
			jsonData, err := json.Marshal(data)
			if err != nil {
				return "", err // Or handle error appropriately
			}
			return template.JS(jsonData), nil
		},
	}
}

func View(w http.ResponseWriter, r *http.Request, page string, data map[string]interface{}, layout string) {
	if data == nil {
		data = make(map[string]interface{})
	}

	// Get form data from cookies if available
	formCookie, err := r.Cookie("form")
	if err == nil {
		formValue, _ := url.QueryUnescape(formCookie.Value)
		var formData map[string]interface{}
		if err := json.Unmarshal([]byte(formValue), &formData); err == nil {
			data["Form"] = formData
		}
		http.SetCookie(w, &http.Cookie{
			Name:   "form",
			MaxAge: -1,
		})
	}

	// Get error data from cookies if available
	errorsCookie, err := r.Cookie("errors")
	if err == nil {
		errorsValue, _ := url.QueryUnescape(errorsCookie.Value)
		var errorsData map[string]string
		if err := json.Unmarshal([]byte(errorsValue), &errorsData); err == nil {
			data["Errors"] = errorsData
		}
		http.SetCookie(w, &http.Cookie{
			Name:   "errors",
			MaxAge: -1,
		})
	}

	// Get error data from cookies if available
	var flash *cookies.FlashMessage
	flash = nil
	if c, err := r.Cookie("flash"); err == nil {
		decoded, _ := url.QueryUnescape(c.Value)
		_ = json.Unmarshal([]byte(decoded), &flash)
		http.SetCookie(w, &http.Cookie{Name: "flash", MaxAge: -1})
	}

	// Get CSRF token from context
	if csrfToken := middleware.GetCSRFToken(r); csrfToken != "" {
		log.Printf("CSRF token encontrado no contexto: %s", csrfToken)
		data["csrf_token"] = csrfToken
	} else {
		log.Printf("CSRF token não encontrado no contexto")
	}

	// Get user from context
	if user := middleware.Auth(r); user != nil {
		log.Printf("Usuário encontrado no contexto: %s", user.Email)
		data["user"] = user
		if user.CSRFToken != "" {
			log.Printf("Usando CSRF token do usuário: %s", user.CSRFToken)
			data["csrf_token"] = user.CSRFToken
		}
	} else {
		log.Printf("Usuário não encontrado no contexto")
	}

	// Parse the template
	tmpl, err := template.New("").Funcs(TemplateFunctions(r)).ParseGlob("web/layouts/*.html")
	if err != nil {
		log.Printf("Erro ao carregar layouts: %v", err)
		http.Error(w, "Erro ao carregar página", http.StatusInternalServerError)
		return
	}

	// Parse partial templates
	_, err = tmpl.ParseGlob("web/partials/*.html")
	if err != nil {
		log.Printf("Erro ao carregar parciais: %v", err)
		http.Error(w, "Erro ao carregar página", http.StatusInternalServerError)
		return
	}

	// Parse the page template
	_, err = tmpl.ParseFiles("web/pages/" + page + ".html")
	if err != nil {
		log.Printf("Erro ao carregar página: %v", err)
		http.Error(w, "Erro ao carregar página", http.StatusInternalServerError)
		return
	}

	// Execute the template
	err = tmpl.ExecuteTemplate(w, layout, map[string]interface{}{
		"Data":  data,
		"Flash": flash,
	})
	if err != nil {
		log.Printf("Erro ao renderizar template: %v", err)
		http.Error(w, "Erro ao renderizar página", http.StatusInternalServerError)
		return
	}
}
