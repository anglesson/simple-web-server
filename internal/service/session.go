package service

import (
	"log"
	"net/http"
	"time"

	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/pkg/utils"
)

type SessionService struct {
	SessionToken string
	CSRFToken    string
}

func NewSessionService() *SessionService {
	return &SessionService{
		SessionToken: "",
		CSRFToken:    "",
	}
}

func (s *SessionService) GenerateSessionToken() string {
	s.SessionToken = utils.GenerateToken(32)
	return s.SessionToken
}

func (s *SessionService) GenerateCSRFToken() string {
	s.CSRFToken = utils.GenerateToken(32)
	return s.CSRFToken
}

func (s *SessionService) SetSessionToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    s.SessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})
}

func (s *SessionService) SetCSRFToken(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "csrf_token",
		Value:    s.CSRFToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
	http.SetCookie(w, cookie)
	log.Printf("CSRF token definido no cookie: %s", s.CSRFToken)
}

func (s *SessionService) ClearSessionToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		MaxAge: -1,
	})
}

func (s *SessionService) ClearCSRFToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "csrf_token",
		MaxAge: -1,
	})
}

func (s *SessionService) GetSessionToken(r *http.Request) string {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func (s *SessionService) GetCSRFToken(r *http.Request) string {
	cookie, err := r.Cookie("csrf_token")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func (s *SessionService) ClearSession(w http.ResponseWriter) {
	s.ClearSessionToken(w)
	s.ClearCSRFToken(w)
}

func (s *SessionService) SetSession(w http.ResponseWriter) {
	s.SetSessionToken(w)
	s.SetCSRFToken(w)
}

func (s *SessionService) GetSession(w http.ResponseWriter, r *http.Request) (string, string) {
	sessionToken := s.GetSessionToken(r)
	csrfToken := s.GetCSRFToken(r)
	return sessionToken, csrfToken
}

func (s *SessionService) InitSession(w http.ResponseWriter, email string) {
	// Generate new tokens
	s.SessionToken = s.GenerateSessionToken()
	s.CSRFToken = s.GenerateCSRFToken()

	// Update the session token in the user data
	userRepository := repository.NewUserRepository()
	user := userRepository.FindByEmail(email)
	if user == nil {
		log.Printf("Erro: Usuário não encontrado para o email: %s", email)
		return
	}

	log.Printf("Atualizando tokens para o usuário: %s", email)
	log.Printf("Session Token: %s", s.SessionToken)
	log.Printf("CSRF Token: %s", s.CSRFToken)

	user.SessionToken = s.SessionToken
	user.CSRFToken = s.CSRFToken

	if err := userRepository.Save(user); err != nil {
		log.Printf("Erro ao salvar tokens do usuário: %v", err)
		return
	}

	// Set cookies after saving to database
	s.SetSessionToken(w)
	s.SetCSRFToken(w)

	log.Printf("Sessão inicializada com sucesso para o usuário: %s", email)
}
