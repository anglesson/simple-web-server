package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repositories"
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
		log.Printf("Session token not found in cookie: %v", err)
		return "", ErrUnauthorized
	}

	// Find user by session token
	userRepository := repositories.NewUserRepository()
	user := userRepository.FindBySessionToken(cookie.Value)
	if user == nil {
		log.Printf("User not found for session token: %s", cookie.Value)
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

	if csrfToken == "" {
		log.Printf("CSRF token is empty for user: %s", user.Email)
		return "", ErrUnauthorized
	}

	if csrfToken != user.CSRFToken {
		log.Printf("CSRF token mismatch for user: %s. Received: %s, Expected: %s",
			user.Email, csrfToken, user.CSRFToken)
		return "", ErrUnauthorized
	}

	// Store the email and CSRF token in request context
	ctx := context.WithValue(r.Context(), UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, CSRFTokenKey, user.CSRFToken)
	ctx = context.WithValue(ctx, User, user)
	*r = *r.WithContext(ctx)

	return user.CSRFToken, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Authentication logic
		csrfToken, err := authorizer(r)
		if err != nil {
			log.Printf("Unauthorized access attempt: %v", err)

			// Check if it's an API request
			if strings.HasPrefix(r.URL.Path, "/api/") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Não autorizado",
				})
				return
			}

			// For regular page requests, redirect to login
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Store CSRF token in a header that your templates can access
		w.Header().Set("X-CSRF-Token", csrfToken)

		// Set CSRF token in cookie if not present
		if _, err := r.Cookie("csrf_token"); err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:     "csrf_token",
				Value:    csrfToken,
				Path:     "/",
				HttpOnly: false,
				Secure:   false,
				SameSite: http.SameSiteStrictMode,
			})
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// GetCSRFToken retrieves the CSRF token from the request context
func GetCSRFToken(r *http.Request) string {
	if token, ok := r.Context().Value(CSRFTokenKey).(string); ok {
		return token
	}
	log.Printf("CSRF token not found in context")
	return ""
}

func Auth(r *http.Request) *models.User {
	var user *models.User

	user_email, ok := r.Context().Value(UserEmailKey).(string)
	if !ok {
		log.Printf("User email not found in context")
		return &models.User{}
	}

	userRepository := repositories.NewUserRepository()
	user = userRepository.FindByEmail(user_email)

	return user
}
