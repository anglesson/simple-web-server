package handlers

import (
	"net/http"

	"github.com/anglesson/simple-web-server/internal/repositories"
	"github.com/anglesson/simple-web-server/internal/shared/template"
)

func DashboardView(w http.ResponseWriter, r *http.Request) {
	loggedUser := GetSessionUser(r)
	dashRepository := repositories.NewDashboardRepository(loggedUser.ID)

	totalEbooks := dashRepository.GetTotalEbooks()
	totalSendEbooks := dashRepository.GetTotalSendEbooks()
	totalClients := dashRepository.GetTotalClients()
	ebookStats, _ := dashRepository.GetEbookStats()

	template.View(w, r, "dashboard", map[string]any{
		"TotalEbooks":        totalEbooks,
		"GetTotalSendEbooks": totalSendEbooks,
		"GetTotalClients":    totalClients,
		"EbookStats":         ebookStats,
	}, "admin")
}
