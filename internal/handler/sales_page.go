package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/anglesson/simple-web-server/internal/authentication/middleware"
	"github.com/anglesson/simple-web-server/internal/service"
	"github.com/anglesson/simple-web-server/pkg/template"
	"github.com/go-chi/chi/v5"
)

// SalesPageHandler gerencia as páginas de vendas
type SalesPageHandler struct {
	ebookService     service.EbookService
	creatorService   service.CreatorService
	templateRenderer template.TemplateRenderer
}

// NewSalesPageHandler cria uma nova instância do SalesPageHandler
func NewSalesPageHandler(ebookService service.EbookService, creatorService service.CreatorService, templateRenderer template.TemplateRenderer) *SalesPageHandler {
	return &SalesPageHandler{
		ebookService:     ebookService,
		creatorService:   creatorService,
		templateRenderer: templateRenderer,
	}
}

// SalesPageView exibe a página de vendas pública do ebook
func (h *SalesPageHandler) SalesPageView(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.Error(w, "Ebook não encontrado", http.StatusNotFound)
		return
	}

	// Buscar o ebook pelo slug
	ebook, err := h.ebookService.FindBySlug(slug)
	if err != nil {
		log.Printf("Erro ao buscar ebook por slug %s: %v", slug, err)
		http.Error(w, "Ebook não encontrado", http.StatusNotFound)
		return
	}

	if ebook == nil {
		http.Error(w, "Ebook não encontrado", http.StatusNotFound)
		return
	}

	// Verificar se o ebook está ativo
	if !ebook.Status {
		http.Error(w, "Ebook não disponível", http.StatusNotFound)
		return
	}

	// Buscar o criador do ebook
	creator, err := h.creatorService.FindByID(ebook.CreatorID)
	if err != nil {
		log.Printf("Erro ao buscar criador do ebook: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	// Atualizar o ebook com os dados do criador
	ebook.Creator = *creator

	// Incrementar visualizações
	ebook.IncrementViews()
	if err := h.ebookService.Update(ebook); err != nil {
		log.Printf("Erro ao incrementar visualizações: %v", err)
	}

	// Calcular economia (preço original vs preço atual)
	originalPrice := ebook.Value * 1.5 // Simular preço original 50% maior
	savings := originalPrice - ebook.Value

	// Preparar dados para o template
	data := map[string]any{
		"Ebook":         ebook,
		"OriginalPrice": originalPrice,
		"Savings":       savings,
		"Creator":       creator,
	}

	h.templateRenderer.View(w, r, "sales_page", data, "guest")
}

// SalesPagePreviewView exibe a página de vendas em modo preview para o criador
func (h *SalesPageHandler) SalesPagePreviewView(w http.ResponseWriter, r *http.Request) {
	// Verificar se o usuário está logado
	userID := middleware.GetCurrentUserID(r)
	if userID == "" {
		http.Error(w, "Não autorizado", http.StatusUnauthorized)
		return
	}

	ebookIDStr := chi.URLParam(r, "id")
	if ebookIDStr == "" {
		http.Error(w, "ID do ebook não fornecido", http.StatusBadRequest)
		return
	}

	ebookID, err := strconv.ParseUint(ebookIDStr, 10, 32)
	if err != nil {
		http.Error(w, "ID do ebook inválido", http.StatusBadRequest)
		return
	}

	// Buscar o ebook
	ebook, err := h.ebookService.FindByID(uint(ebookID))
	if err != nil {
		log.Printf("Erro ao buscar ebook: %v", err)
		http.Error(w, "Ebook não encontrado", http.StatusNotFound)
		return
	}

	// Verificar se o usuário logado é o criador do ebook
	creator, err := h.creatorService.FindCreatorByUserID(userID)
	if err != nil || creator.ID != ebook.CreatorID {
		http.Error(w, "Não autorizado", http.StatusUnauthorized)
		return
	}

	// Atualizar o ebook com os dados do criador
	ebook.Creator = *creator

	// Calcular economia
	originalPrice := ebook.Value * 1.5
	savings := originalPrice - ebook.Value

	// Preparar dados para o template
	data := map[string]any{
		"Ebook":         ebook,
		"OriginalPrice": originalPrice,
		"Savings":       savings,
		"Creator":       creator,
		"IsPreview":     true,
	}

	h.templateRenderer.View(w, r, "sales_page", data, "guest")
}
