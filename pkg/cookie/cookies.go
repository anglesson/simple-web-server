package cookies

import (
	"encoding/json"
	"net/http"
	"net/url"
)

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

func NotifyError(w http.ResponseWriter, message string) {
	b, _ := json.Marshal(FlashMessage{
		Message: message,
		Type:    "danger",
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "flash",
		Value: url.QueryEscape(string(b)),
		Path:  "/",
	})
}
