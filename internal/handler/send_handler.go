package handler

import (
	"net/http"
	"strconv"

	"github.com/anglesson/simple-web-server/internal/handler/web"
	"github.com/anglesson/simple-web-server/internal/repository/gorm"

	"github.com/anglesson/simple-web-server/internal/handler/middleware"
	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/internal/service"
	cookies "github.com/anglesson/simple-web-server/pkg/cookie"
	"github.com/anglesson/simple-web-server/pkg/database"
	"github.com/anglesson/simple-web-server/pkg/storage"
	"github.com/anglesson/simple-web-server/pkg/template"
)

type SendHandler struct {
	templateRenderer template.TemplateRenderer
}

func NewSendHandler(templateRenderer template.TemplateRenderer) *SendHandler {
	return &SendHandler{
		templateRenderer: templateRenderer,
	}
}

func (h *SendHandler) SendViewHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	loggedUser := middleware.Auth(r)
	if loggedUser.ID == 0 {
		http.Error(w, "Não foi possível prosseguir com a sua solicitação", http.StatusInternalServerError)
		return
	}

	creatorRepository := gorm.NewCreatorRepository(database.DB)
	creator, err := creatorRepository.FindCreatorByUserID(loggedUser.ID)
	if err != nil {
		web.RedirectBackWithErrors(w, r, err.Error())
	}

	pagination := models.NewPagination(1, 1000)

	viewData := map[string]any{
		"Ebooks":         nil,
		"Clients":        nil,
		"Pagination":     pagination,
		"EbookID":        nil,
		"ClientsCreator": len(creator.Clients),
	}

	// TODO: This should be injected as dependency
	s3Storage := storage.NewS3Storage()
	ebookService := service.NewEbookService(s3Storage)
	ebooks, err := ebookService.ListEbooksForUser(loggedUser.ID, repository.EbookQuery{
		Pagination: pagination,
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
	clients, err := gorm.NewClientGormRepository().FindByClientsWhereEbookNotSend(creator, models.ClientFilter{
		EbookID:    uint(ebookID),
		Pagination: pagination,
		Term:       r.URL.Query().Get("term"),
	})
	if clients != nil && len(*clients) > 0 {
		viewData["Clients"] = clients
	}

	if err != nil {
		web.RedirectBackWithErrors(w, r, err.Error())
	}

	// Set total count for pagination
	if clients != nil {
		pagination.SetTotal(int64(len(*clients)))
	}

	h.templateRenderer.View(w, r, "ebook/send", viewData, "admin")
}
