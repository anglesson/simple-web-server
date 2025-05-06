package handlers

import (
	"net/http"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/shared/database"
	"github.com/anglesson/simple-web-server/internal/shared/middlewares"
	"github.com/anglesson/simple-web-server/internal/shared/template"
)

func EbookIndexView(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		renderIndexPage(w, r)
	}
}

func renderIndexPage(w http.ResponseWriter, r *http.Request) {
	loggedUser := middlewares.Auth(r)
	if loggedUser.ID == 0 {
		http.Error(w, "Não foi possível prosseguir com a sua solicitação", http.StatusInternalServerError)
		return
	}

	var creator models.Creator

	// Busca o criador com os ebooks associados
	err := database.DB.
		Preload("Ebooks").
		Where("user_id = ?", loggedUser.ID).
		First(&creator).Error
	if err != nil {
		http.Error(w, "Erro ao buscar dados", http.StatusInternalServerError)
		return
	}

	// Renderiza a página com os ebooks do criador
	template.View(w, r, "ebook", map[string]any{
		"Ebooks": creator.Ebooks,
	}, "admin")
}
