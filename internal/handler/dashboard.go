package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/anglesson/simple-web-server/internal/handler/middleware"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/pkg/template"
)

type DashboardHandler struct {
	templateRenderer template.TemplateRenderer
}

func NewDashboardHandler(templateRenderer template.TemplateRenderer) *DashboardHandler {
	return &DashboardHandler{
		templateRenderer: templateRenderer,
	}
}

func (h *DashboardHandler) DashboardView(w http.ResponseWriter, r *http.Request) {
	loggedUser := middleware.Auth(r)
	dashRepository := repository.NewDashboardRepository(loggedUser.ID)

	totalEbooks := dashRepository.GetTotalEbooks()
	totalSendEbooks := dashRepository.GetTotalSendEbooks()
	totalClients := dashRepository.GetTotalClients()
	ebookStats, _ := dashRepository.GetEbookStats()

	// Get data for charts
	dailyPurchases, _ := dashRepository.GetDailyPurchases()
	dailyDownloads, _ := dashRepository.GetDailyDownloads()
	topEbooks, _ := dashRepository.GetTopEbooks()
	topClients, _ := dashRepository.GetTopClients()
	topDownloadedEbooks, _ := dashRepository.GetTopDownloadedEbooks()

	// Marshal data to JSON strings
	dailyPurchasesJSON, err := json.Marshal(dailyPurchases)
	if err != nil {
		log.Printf("Error marshaling daily purchases data to JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	dailyDownloadsJSON, err := json.Marshal(dailyDownloads)
	if err != nil {
		log.Printf("Error marshaling daily downloads data to JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	topEbooksJSON, err := json.Marshal(topEbooks)
	if err != nil {
		log.Printf("Error marshaling top ebooks data to JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	topDownloadedEbooksJSON, err := json.Marshal(topDownloadedEbooks)
	if err != nil {
		log.Printf("Error marshaling top downloaded ebooks data to JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.templateRenderer.View(w, r, "dashboard", map[string]any{
		"TotalEbooks":             totalEbooks,
		"GetTotalSendEbooks":      totalSendEbooks,
		"GetTotalClients":         totalClients,
		"EbookStats":              ebookStats,
		"DailyPurchasesJSON":      string(dailyPurchasesJSON),
		"DailyDownloadsJSON":      string(dailyDownloadsJSON),
		"TopEbooksJSON":           string(topEbooksJSON),
		"TopClients":              topClients,
		"TopDownloadedEbooksJSON": string(topDownloadedEbooksJSON),
	}, "admin")
}
