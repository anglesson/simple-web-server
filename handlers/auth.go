package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/anglesson/simple-web-server/templates"
)

func LoginGetHandler(w http.ResponseWriter, r *http.Request) {
	errors := map[string]string{}
	if cookie, err := r.Cookie("form_error"); err == nil {
		decodedValue, decodeErr := url.QueryUnescape(cookie.Value) // Decodifica o valor do cookie
		if decodeErr != nil {
			log.Println("Error decoding cookie value:", decodeErr)
		} else {
			_ = json.Unmarshal([]byte(decodedValue), &errors)
			log.Printf("Decoded cookie value: %s", decodedValue)
		}

		log.Println("Cookie value:", cookie.Value)

		// Apaga o cookie após usar
		http.SetCookie(w, &http.Cookie{
			Name:   "form_error",
			MaxAge: -1,
			Path:   "/",
		})
	} else {
		log.Println("Error getting cookie:", err)
	}

	log.Printf("Errors: %v", errors)

	templates.View(w, "login", map[string]any{
		"errors": errors,
	})
}

func LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("LoginPostHandler called with method: %s", r.Method)
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	validationErrors := map[string]string{}

	if email == "" {
		validationErrors["email"] = "Email é obrigatório."
	}
	if password == "" {
		validationErrors["password"] = "Senha é obrigatória."
	}

	if email != "admin@example.com" || password != "123456" {
		validationErrors["password"] = "Credenciais invalidas"
	}

	if len(validationErrors) > 0 {
		log.Println("Validation errors:", validationErrors)
		data, err := json.Marshal(validationErrors)
		if err != nil {
			log.Println("Error marshaling validation errors:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		setFormError(w, string(data))
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func RegisterGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.View(w, "register", nil)
}

func ForgetPasswordGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.View(w, "forget-password", nil)
}

func setFormError(w http.ResponseWriter, message string) {
	log.Println("Setting form error cookie with message:", message)
	http.SetCookie(w, &http.Cookie{
		Name:     "form_error",
		Value:    url.QueryEscape(message),
		Path:     "/",
		MaxAge:   5,
		HttpOnly: true,
		Expires:  time.Now().Add(5 * time.Second),
	})
}
