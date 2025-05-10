package handlers

import (
	"net/http"
	"strconv"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repositories"
	"github.com/anglesson/simple-web-server/internal/shared/middlewares"
	"github.com/anglesson/simple-web-server/internal/shared/template"
)

func ClientIndexView(w http.ResponseWriter, r *http.Request) {
	loggedUser := middlewares.Auth(r)
	if loggedUser.ID == 0 {
		http.Error(w, "Não foi possível prosseguir com a sua solicitação", http.StatusInternalServerError)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	// term := r.URL.Query().Get("term")
	pagination := repositories.NewPagination(page, perPage)

	client1 := models.NewClient("Jose Arimatéia", "000.000.000-00", "jose.arimateia@test.com", "+55996265197")
	clients := []*models.Client{
		client1,
	}

	template.View(w, r, "client", map[string]any{
		"Clients":    clients,
		"Pagination": pagination,
	}, "admin")
}
