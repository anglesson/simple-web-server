package handler

import (
	"net/http"

	"github.com/anglesson/simple-web-server/pkg/template"
)

type HomeHandler struct {
	templateRenderer template.TemplateRenderer
	errorHandler     *ErrorHandler
}

func NewHomeHandler(templateRenderer template.TemplateRenderer, errorHandler *ErrorHandler) *HomeHandler {
	return &HomeHandler{
		templateRenderer: templateRenderer,
		errorHandler:     errorHandler,
	}
}

func (h *HomeHandler) HomeView(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.errorHandler.ErrorView(w, r, 404)
		return
	}

	h.templateRenderer.View(w, r, "home", nil, "guest")
}
