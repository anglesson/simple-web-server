package cookies

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

// Nome do cookie da sessão
const SessionCookieName = "session_token"

// Define um cookie com o token de sessão
func SetSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Coloque true em produção com HTTPS
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})
}

// Remove o cookie de sessão
func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

// Lê o token de sessão a partir do cookie
func GetSessionToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

type FlashMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func NotifySuccess(w http.ResponseWriter, message string) {
	b, _ := json.Marshal(FlashMessage{
		Message: message,
		Type:    "success",
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "flash",
		Value: url.QueryEscape(string(b)),
		Path:  "/",
	})
}
