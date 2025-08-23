package presentation

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/anglesson/simple-web-server/internal/authentication/business"
	"github.com/anglesson/simple-web-server/internal/authentication/data"
	"github.com/anglesson/simple-web-server/internal/authentication/session"
	"github.com/anglesson/simple-web-server/internal/config"
)

// GoogleAuthHandlers contém apenas os handlers para autenticação Google
type GoogleAuthHandlers struct {
	authService    *business.AuthService
	googleAuthRepo *data.GoogleAuthRepository
	sessionStore   *session.SessionStore
}

// NewGoogleAuthHandlers cria uma nova instância focada apenas no Google OAuth
func NewGoogleAuthHandlers(sessionStore *session.SessionStore) *GoogleAuthHandlers {
	// Inicializa as dependências do Google OAuth
	googleAuthRepo := data.NewGoogleAuthRepository()
	authService := business.NewAuthService(googleAuthRepo, sessionStore)

	log.Printf("Google OAuth configured: %v", googleAuthRepo != nil)

	return &GoogleAuthHandlers{
		authService:    authService,
		googleAuthRepo: googleAuthRepo,
		sessionStore:   sessionStore,
	}
}

// HandleGoogleLogin redireciona para o Google OAuth
func (h *GoogleAuthHandlers) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	state, err := generateRandomState()
	if err != nil {
		fmt.Printf("Erro ao gerar state: %v\n", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	authURL := h.googleAuthRepo.GetAuthURL(state)
	fmt.Printf("Redirecionando para Google OAuth: %s\n", authURL)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// HandleGoogleCallback processa o retorno do Google
func (h *GoogleAuthHandlers) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	errorParam := r.URL.Query().Get("error")

	if errorParam != "" {
		fmt.Printf("Erro na autorização Google: %s\n", errorParam)
		http.Redirect(w, r, "/?error=google_auth_failed", http.StatusTemporaryRedirect)
		return
	}

	if code == "" {
		fmt.Printf("Código de autorização não fornecido\n")
		http.Redirect(w, r, "/?error=missing_code", http.StatusTemporaryRedirect)
		return
	}

	fmt.Printf("Processando callback - State: %s, Code: %s\n", state, code)

	// Processa o login via Google
	sessionID, err := h.authService.HandleGoogleLogin(code)
	if err != nil {
		fmt.Printf("Erro no login Google: %v\n", err)
		http.Redirect(w, r, "/?error=auth_failed", http.StatusTemporaryRedirect)
		return
	}

	// Define o cookie de sessão
	fmt.Printf("Definindo cookie de sessão com ID: %s\n", sessionID)
	// MUDANÇA: Definir o cookie com configurações mais permissivas para debugging
	fmt.Printf("Definindo cookie de sessão com ID: %s\n", sessionID)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,                // TEMPORÁRIO: false para debugging
		SameSite: http.SameSiteLaxMode, // Lax em vez de Strict
		MaxAge:   3600 * 24 * 7,
	})

	// Log dos cookies que estão sendo definidos
	fmt.Printf("Cookies definidos na resposta:\n")
	for _, cookie := range w.Header()["Set-Cookie"] {
		fmt.Printf("  - %s\n", cookie)
	}

	fmt.Printf("Login Google realizado com sucesso\n")
	http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
}

// HandleLogout remove a sessão e redireciona
func (h *GoogleAuthHandlers) HandleLogout(w http.ResponseWriter, r *http.Request) {
	// Obtém o session ID do cookie
	cookie, err := r.Cookie("session_id")
	if err == nil && cookie.Value != "" {
		// Remove a sessão do store
		if err := h.sessionStore.DeleteSession(cookie.Value); err != nil {
			fmt.Printf("Erro ao remover sessão: %v\n", err)
		} else {
			fmt.Printf("Sessão %s removida do store\n", cookie.Value)
		}
	}

	// Remove os cookies
	h.clearSessionCookies(w)

	fmt.Printf("Logout realizado com sucesso\n")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// setSessionCookie define um cookie seguro
func (h *GoogleAuthHandlers) setSessionCookie(w http.ResponseWriter, sessionID string) {
	secure := config.AppConfig.IsProduction()

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure, // false em desenvolvimento
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600 * 24 * 7, // 7 dias
	}

	fmt.Printf("Definindo cookie: Name=%s, Value=%s, Path=%s, HttpOnly=%v, Secure=%v, SameSite=%v, MaxAge=%d\n",
		cookie.Name, sessionID[:8]+"...", cookie.Path, cookie.HttpOnly, cookie.Secure, cookie.SameSite, cookie.MaxAge)

	http.SetCookie(w, cookie)

}

// clearSessionCookies remove os cookies de sessão
func (h *GoogleAuthHandlers) clearSessionCookies(w http.ResponseWriter) {
	// Remove o cookie de sessão
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   config.AppConfig.IsProduction(),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	// Remove o cookie CSRF
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   config.AppConfig.IsProduction(),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

// generateRandomState gera um state aleatório para segurança
func generateRandomState() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
