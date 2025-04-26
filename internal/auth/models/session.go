package models

import "net/http"

type Session interface {
	Start(name string)
	SetCookie(w http.ResponseWriter, cookie *http.Cookie)
	GetCookie(r *http.Request, name string) (*http.Cookie, error)
	DeleteCookie(w http.ResponseWriter, name string, path string)
}
