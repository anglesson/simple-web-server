package services

import (
	"log"
	"net/http"
	"time"

	"github.com/anglesson/simple-web-server/internal/repositories"
	"github.com/anglesson/simple-web-server/internal/shared/utils"
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
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    s.CSRFToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})
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
	s.SessionToken = s.GenerateSessionToken()
	s.SetSessionToken(w)
	s.CSRFToken = s.GenerateCSRFToken()
	s.SetCSRFToken(w)

	// Update the session token in the user data
	userFound := repositories.Users[email]
	userFound.SessionToken = s.SessionToken
	userFound.CSRFToken = s.CSRFToken
	repositories.Users[email] = userFound

	log.Printf("Initialized session with EMAIL: %s", email)
}
