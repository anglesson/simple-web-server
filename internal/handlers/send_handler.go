package handlers

import (
	"net/http"
	"strconv"

	client_dtos "github.com/anglesson/simple-web-server/internal/client/dtos"
	client_repo "github.com/anglesson/simple-web-server/internal/client/repositories"
	common_application "github.com/anglesson/simple-web-server/internal/common/application"
	common_http "github.com/anglesson/simple-web-server/internal/common/infrastructure/http_serve"
	"github.com/anglesson/simple-web-server/internal/repositories"
	"github.com/anglesson/simple-web-server/internal/services"
	cookies "github.com/anglesson/simple-web-server/internal/shared/cookie"
	"github.com/anglesson/simple-web-server/internal/shared/middlewares"
	"github.com/anglesson/simple-web-server/internal/shared/template"
)

func SendViewHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	loggedUser := middlewares.Auth(r)
	if loggedUser.ID == 0 {
		http.Error(w, "Não foi possível prosseguir com a sua solicitação", http.StatusInternalServerError)
		return
	}

	creatorRepository := repositories.NewCreatorRepository()
	creator, err := creatorRepository.FindCreatorByUserID(loggedUser.ID)
	if err != nil {
		common_http.RedirectBackWithErrors(w, r, err.Error())
	}

	viewData := map[string]any{
		"Ebooks":         nil,
		"Clients":        nil,
		"Pagination":     common_application.NewPagination(1, 1000),
		"EbookID":        nil,
		"ClientsCreator": len(creator.Clients),
	}

	ebookService := services.NewEbookService()
	ebooks, err := ebookService.ListEbooksForUser(loggedUser.ID, repositories.EbookQuery{
		Pagination: viewData["Pagination"].(*common_application.Pagination),
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
	clients, err := client_repo.NewClientRepository().FindByClientsWhereEbookNotSend(creator, client_dtos.ClientQuery{
		EbookID:    uint(ebookID),
		Pagination: viewData["Pagination"].(*common_application.Pagination),
		Term:       r.URL.Query().Get("term"),
	})
	if clients != nil && len(*clients) > 0 {
		viewData["Clients"] = clients
	}

	if err != nil {
		common_http.RedirectBackWithErrors(w, r, err.Error())
	}

	template.View(w, r, "send_ebook", viewData, "admin")
}
