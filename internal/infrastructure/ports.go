package infrastructure

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

var _ FlashMessagePort = (*GorillaFlashMessage)(nil)

var SessionStore = sessions.NewCookieStore([]byte("SESSION_KEY"))

type Message struct {
	Type    string // "success", "error", "info", "warning"
	Content string
}

type FlashMessagePort interface {
	Success(message string)
	Error(message string)
}

type FlashMessageFactory func(w http.ResponseWriter, r *http.Request) FlashMessagePort

type GorillaFlashMessage struct {
	sessionStore *sessions.CookieStore
	w            http.ResponseWriter
	r            *http.Request
}

func NewGorillaFlashMessage(w http.ResponseWriter, r *http.Request) FlashMessagePort {
	return &GorillaFlashMessage{
		w:            w,
		r:            r,
		sessionStore: SessionStore,
	}
}

func (fm *GorillaFlashMessage) Success(message string) {
	session, err := fm.sessionStore.Get(fm.r, "SESSION_KEY")
	if err != nil {
		return
	}
	session.AddFlash(message, "success")
	session.Save(fm.r, fm.w)
	log.Printf("[SUCCESS FLASH]: " + message)
}

func (fm *GorillaFlashMessage) Error(message string) {
	session, err := fm.sessionStore.Get(fm.r, "SESSION_KEY")
	if err != nil {
		return
	}
	session.AddFlash(message, "error")
	session.Save(fm.r, fm.w)
	log.Printf("[ERROR FLASH]: " + message)
}
