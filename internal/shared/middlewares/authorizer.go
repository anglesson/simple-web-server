package middlewares

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/anglesson/simple-web-server/internal/auth/models"
	"github.com/anglesson/simple-web-server/internal/auth/repositories"
)

// First, define a custom type for context keys (typically at package level)
type contextKey string

// Define a constant for your key
const UserEmailKey contextKey = "user_email"
const CSRFTokenKey contextKey = "csrf_token"
const User contextKey = "user"

var ErrUnauthorized = errors.New("Unauthorized")

func authorizer(r *http.Request) (string, error) {
	// Get session token from the cookie
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		log.Println("Session token not found in cookie:", err)
		return "", ErrUnauthorized
	}

	// Find user by session token
	var foundUser models.Login
	var foundEmail string
	var userFound bool
	var user models.User

	for email, user := range repositories.Users {
		if user.SessionToken == cookie.Value {
			foundUser = user
			foundEmail = email
			userFound = true
			break
		}
	}

	if !userFound {
		log.Println("User not found for session token:", cookie.Value)
		return "", ErrUnauthorized
	}

	// Get CSRF token from the cookie or header
	csrfCookie, _ := r.Cookie("csrf_token")
	csrfHeader := r.Header.Get("X-CSRF-Token")

	// Try both cookie and header for CSRF
	csrfToken := csrfCookie.Value
	if csrfHeader != "" {
		csrfToken = csrfHeader
	}

	if csrfToken != foundUser.CSRFToken || csrfToken == "" {
		log.Println("CSRF token mismatch or empty for user:", foundEmail)
		return "", ErrUnauthorized
	}

	// Store the email and CSRF token in request context
	ctx := context.WithValue(r.Context(), UserEmailKey, foundEmail)
	ctx = context.WithValue(ctx, CSRFTokenKey, foundUser.CSRFToken)
	ctx = context.WithValue(ctx, User, user)
	*r = *r.WithContext(ctx)

	return foundUser.CSRFToken, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Authentication logic
		csrfToken, err := authorizer(r)
		if err != nil {
			if r.URL.Path == "/login" || r.URL.Path == "/register" || r.URL.Path == "/forget-password" {
				next.ServeHTTP(w, r)
				return
			}
			log.Println("Unauthorized access attempt:", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Fluxo de autenticação logado
		if r.URL.Path == "/login" || r.URL.Path == "/register" || r.URL.Path == "/forget-password" {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}

		// Store CSRF token in a header that your templates can access
		w.Header().Set("X-CSRF-Token", csrfToken)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// GetCSRFToken retrieves the CSRF token from the request context
func GetCSRFToken(r *http.Request) string {
	if token, ok := r.Context().Value(CSRFTokenKey).(string); ok {
		return token
	}
	return ""
}

func Auth(r *http.Request) *models.User {
	var user *models.User

	user_email, ok := r.Context().Value(UserEmailKey).(string)
	if !ok {
		log.Panic("Ocorreu erro ao recuperar as informações do usuário")
		return &models.User{}
	}

	user = repositories.FindByEmail(user_email)

	return user
}
