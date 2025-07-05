package handler

import (
	"net/http"
	"strconv"

	"github.com/anglesson/simple-web-server/domain"
	"github.com/anglesson/simple-web-server/internal/handler/web"
	"github.com/anglesson/simple-web-server/internal/repositories/gorm"

	"github.com/anglesson/simple-web-server/internal/handler/middleware"
	"github.com/anglesson/simple-web-server/internal/repositories"
	"github.com/anglesson/simple-web-server/internal/services"
	cookies "github.com/anglesson/simple-web-server/pkg/cookie"
	"github.com/anglesson/simple-web-server/pkg/template"
)

func SendViewHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	loggedUser := middleware.Auth(r)
	if loggedUser.ID == 0 {
		http.Error(w, "Não foi possível prosseguir com a sua solicitação", http.StatusInternalServerError)
		return
	}

	creatorRepository := gorm.NewCreatorRepository()
	creator, err := creatorRepository.FindCreatorByUserID(loggedUser.ID)
	if err != nil {
		web.RedirectBackWithErrors(w, r, err.Error())
	}

	viewData := map[string]any{
		"Ebooks":         nil,
		"Clients":        nil,
		"Pagination":     domain.NewPagination(1, 1000),
		"EbookID":        nil,
		"ClientsCreator": len(creator.Clients),
	}

	ebookService := services.NewEbookService()
	ebooks, err := ebookService.ListEbooksForUser(loggedUser.ID, repositories.EbookQuery{
		Pagination: viewData["Pagination"].(*domain.Pagination),
	})
	if err != nil {
		cookies.NotifyError(w, "Ocorre um erro ao listar seus ebooks")
		return
	}
	if ebooks != nil && len(*ebooks) > 0 {
		viewData["Ebooks"] = ebooks
	}

	ebookID, _ := strconv.Atoi(r.URL.Query().Get("ebook_id"))
	if ebookID != 0 {
		viewData["EbookID"] = ebookID
	}
	clients, err := gorm.NewClientGormRepository().FindByClientsWhereEbookNotSend(creator, domain.ClientFilter{
		EbookID:    uint(ebookID),
		Pagination: viewData["Pagination"].(*domain.Pagination),
		Term:       r.URL.Query().Get("term"),
	})
	if clients != nil && len(*clients) > 0 {
		viewData["Clients"] = clients
	}

	if err != nil {
		web.RedirectBackWithErrors(w, r, err.Error())
	}

	template.View(w, r, "send_ebook", viewData, "admin")
}
