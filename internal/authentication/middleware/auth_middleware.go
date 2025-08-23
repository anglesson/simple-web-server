package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/anglesson/simple-web-server/internal/authentication/session"
	"github.com/anglesson/simple-web-server/internal/config"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/pkg/database"
)

// First, define a custom type for context keys (typically at package level)
type contextKey string

// Define a constant for your key
const UserEmailKey contextKey = "user_email"
const CSRFTokenKey contextKey = "csrf_token"
const UserKey contextKey = "user"

var ErrUnauthorized = errors.New("unauthorized")

// GoogleSessionMiddleware middleware que usa apenas o Google OAuth SessionStore
type GoogleSessionMiddleware struct {
	sessionStore *session.SessionStore
}

// NewGoogleSessionMiddleware cria um novo middleware baseado apenas em sessões Google
func NewGoogleSessionMiddleware(sessionStore *session.SessionStore) *GoogleSessionMiddleware {
	return &GoogleSessionMiddleware{
		sessionStore: sessionStore,
	}
}

// AuthMiddleware verifica se o usuário está autenticado via Google OAuth
func (gsm *GoogleSessionMiddleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log para debugging
		log.Printf("AuthMiddleware: Verificando autenticação para %s %s", r.Method, r.URL.Path)

		// Listar todos os cookies para debugging
		log.Printf("AuthMiddleware: Cookies recebidos:")
		for _, cookie := range r.Cookies() {
			// Não logar o valor completo dos cookies por segurança
			log.Printf("  - %s: %s", cookie.Name, cookie.Value[:min(len(cookie.Value), 10)]+"...")
		}

		// Tenta obter o session ID do cookie
		sessionID, err := gsm.getSessionIDFromRequest(r)
		if err != nil {
			gsm.handleUnauthorized(w, r)
			return
		}

		// Obtém o usuário da sessão
		user, err := gsm.sessionStore.GetSessionUser(sessionID)
		if err != nil {
			log.Printf("Sessão inválida ou expirada: %v", err)
			gsm.handleUnauthorized(w, r)
			return
		}

		// Gera um CSRF token simples baseado na sessão (em produção, use algo mais robusto)
		csrfToken := "google-" + sessionID[:16]

		// Armazena as informações do usuário no contexto
		ctx := context.WithValue(r.Context(), UserKey, user)
		ctx = context.WithValue(ctx, UserEmailKey, user.Email)
		ctx = context.WithValue(ctx, CSRFTokenKey, csrfToken)
		*r = *r.WithContext(ctx)

		// Define o CSRF token no header para templates
		w.Header().Set("X-CSRF-Token", csrfToken)

		// Define o CSRF token em um cookie se necessário
		if _, err := r.Cookie("csrf_token"); err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:     "csrf_token",
				Value:    csrfToken,
				Path:     "/",
				HttpOnly: true,
				Secure:   config.AppConfig.IsProduction(),
				SameSite: http.SameSiteLaxMode,
			})
		}

		next.ServeHTTP(w, r)
	})
}

// getSessionIDFromRequest extrai o session ID do cookie
func (gsm *GoogleSessionMiddleware) getSessionIDFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// handleUnauthorized lida com requisições não autorizadas
func (gsm *GoogleSessionMiddleware) handleUnauthorized(w http.ResponseWriter, r *http.Request) {
	log.Printf("Tentativa de acesso não autorizado: %s %s", r.Method, r.URL.Path)

	// Para requisições de API, retorna JSON
	if strings.HasPrefix(r.URL.Path, "/api/") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Não autorizado - faça login com Google",
		})
		return
	}

	// Para outras requisições, redireciona para a página inicial com o botão de login do Google
	http.Redirect(w, r, "/?login_required=true", http.StatusSeeOther)
}

// GetAuthUser retorna o usuário autenticado via Google
func GetAuthUser(r *http.Request) *session.User {
	if user, ok := r.Context().Value(UserKey).(*session.User); ok {
		return user
	}
	log.Printf("User not found in context")
	return nil
}

// GetUserEmail retorna o email do usuário autenticado
func GetUserEmail(r *http.Request) string {
	if email, ok := r.Context().Value(UserEmailKey).(string); ok {
		return email
	}
	log.Printf("User email not found in context")
	return ""
}

// GetCSRFToken retrieves the CSRF token from the request context
func GetCSRFToken(r *http.Request) string {
	if token, ok := r.Context().Value(CSRFTokenKey).(string); ok {
		return token
	}
	log.Printf("CSRF token not found in context")
	return ""
}

func authorizer(r *http.Request) (string, error) {
	// Get session token from the cookie
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		log.Printf("Session token not found in cookie: %v", err)
		return "", ErrUnauthorized
	}

	// Find user by session token
	userRepository := repository.NewGormUserRepository(database.DB)
	user := userRepository.FindBySessionToken(cookie.Value)
	if user == nil {
		log.Printf("UserKey not found for session token")
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
		log.Printf("CSRF token mismatch for user: %s", user.Email)
		return "", ErrUnauthorized
	}

	// Store the email and CSRF token in request context
	ctx := context.WithValue(r.Context(), UserEmailKey, user.Email)
	ctx = context.WithValue(ctx, CSRFTokenKey, user.CSRFToken)
	ctx = context.WithValue(ctx, UserKey, user)
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
				HttpOnly: true,
				Secure:   config.AppConfig.IsProduction(),
				SameSite: http.SameSiteStrictMode,
			})
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func Auth(r *http.Request) *session.User {
	return GetAuthUser(r)
}

// GetUserIDForDB retorna um ID que pode ser usado para relacionamentos no banco de dados
// Como não temos mais models.User, vamos usar o email como identificador único
func GetUserIDForDB(r *http.Request) string {
	user := GetAuthUser(r)
	if user == nil {
		return ""
	}
	// Você pode usar o ID do Google ou o email como identificador
	// Por simplicidade, vou usar o email que é único
	return user.Email
}

// GetCurrentUserEmail retorna o email do usuário atual
func GetCurrentUserEmail(r *http.Request) string {
	return GetUserEmail(r)
}

// GetCurrentUserName retorna o nome do usuário atual
func GetCurrentUserName(r *http.Request) string {
	user := GetAuthUser(r)
	if user == nil {
		return ""
	}
	return user.Name
}

// GetCurrentUserID retorna o ID do usuário atual
func GetCurrentUserID(r *http.Request) string {
	user := GetAuthUser(r)
	if user == nil {
		return ""
	}
	return user.ID
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
