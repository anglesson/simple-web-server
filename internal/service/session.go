package service

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/anglesson/simple-web-server/internal/config"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/pkg/database"
	"github.com/anglesson/simple-web-server/pkg/utils"
)

type SessionService interface {
	GenerateSessionToken() string
	GenerateCSRFToken() string
	SetSessionToken(w http.ResponseWriter)
	SetCSRFToken(w http.ResponseWriter)
	ClearSessionToken(w http.ResponseWriter)
	ClearCSRFToken(w http.ResponseWriter)
	GetSessionToken(r *http.Request) string
	GetCSRFToken(r *http.Request) string
	ClearSession(w http.ResponseWriter)
	SetSession(w http.ResponseWriter)
	GetSession(w http.ResponseWriter, r *http.Request) (string, string)
	InitSession(w http.ResponseWriter, email string)
}

type SessionServiceImpl struct {
	SessionToken string
	CSRFToken    string
	encrypter    utils.Encrypter
}

// maskEmail masks sensitive parts of email for logging
func maskEmail(email string) string {
	if email == "" {
		return "[EMPTY]"
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "[INVALID_EMAIL]"
	}
	username := parts[0]
	domain := parts[1]
	
	if len(username) <= 2 {
		return fmt.Sprintf("%s***@%s", username, domain)
	}
	
	maskedUsername := username[:2] + "***"
	return fmt.Sprintf("%s@%s", maskedUsername, domain)
}

func NewSessionService() SessionService {
	return &SessionServiceImpl{
		SessionToken: "",
		CSRFToken:    "",
		encrypter:    utils.NewEncrypter(),
	}
}

func (s *SessionServiceImpl) GenerateSessionToken() string {
	s.SessionToken = s.encrypter.GenerateToken(32)
	return s.SessionToken
}

func (s *SessionServiceImpl) GenerateCSRFToken() string {
	s.CSRFToken = s.encrypter.GenerateToken(32)
	return s.CSRFToken
}

func (s *SessionServiceImpl) SetSessionToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    s.SessionToken,
		Expires:  time.Now().Add(8 * time.Hour), // Reduzido de 24h para 8h por segurança
		HttpOnly: true,
		Secure:   config.AppConfig.IsProduction(),
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
}

func (s *SessionServiceImpl) SetCSRFToken(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "csrf_token",
		Value:    s.CSRFToken,
		Expires:  time.Now().Add(8 * time.Hour), // Reduzido de 24h para 8h por segurança
		HttpOnly: true,
		Secure:   config.AppConfig.IsProduction(),
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
	log.Printf("CSRF token definido no cookie: [REDACTED]")
}

func (s *SessionServiceImpl) ClearSessionToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		MaxAge: -1,
	})
}

func (s *SessionServiceImpl) ClearCSRFToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "csrf_token",
		MaxAge: -1,
	})
}

func (s *SessionServiceImpl) GetSessionToken(r *http.Request) string {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func (s *SessionServiceImpl) GetCSRFToken(r *http.Request) string {
	cookie, err := r.Cookie("csrf_token")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func (s *SessionServiceImpl) ClearSession(w http.ResponseWriter) {
	s.ClearSessionToken(w)
	s.ClearCSRFToken(w)
}

func (s *SessionServiceImpl) SetSession(w http.ResponseWriter) {
	s.SetSessionToken(w)
	s.SetCSRFToken(w)
}

func (s *SessionServiceImpl) GetSession(w http.ResponseWriter, r *http.Request) (string, string) {
	sessionToken := s.GetSessionToken(r)
	csrfToken := s.GetCSRFToken(r)
	return sessionToken, csrfToken
}

func (s *SessionServiceImpl) InitSession(w http.ResponseWriter, email string) {
	// Generate new tokens
	s.SessionToken = s.GenerateSessionToken()
	s.CSRFToken = s.GenerateCSRFToken()

	// Update the session token in the user data
	userRepository := repository.NewGormUserRepository(database.DB)
	user := userRepository.FindByEmail(email)
	if user == nil {
		log.Printf("Erro: Usuário não encontrado para o email: %s", email)
		return
	}

	log.Printf("Atualizando tokens para o usuário: %s", maskEmail(email))
	log.Printf("Session token generated: [REDACTED]")
	log.Printf("CSRF token generated: [REDACTED]")

	user.SessionToken = s.SessionToken
	user.CSRFToken = s.CSRFToken

	if err := userRepository.Save(user); err != nil {
		log.Printf("Erro ao salvar tokens do usuário: %v", err)
		return
	}

	// Set cookies after saving to database
	s.SetSessionToken(w)
	s.SetCSRFToken(w)

	log.Printf("Sessão inicializada com sucesso para o usuário: %s", maskEmail(email))
}
