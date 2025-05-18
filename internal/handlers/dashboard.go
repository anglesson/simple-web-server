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
	engagementMetric := dashRepository.GetEngagementMetric()
	ebookStats, _ := dashRepository.GetEbookStats()

	template.View(w, r, "dashboard", map[string]any{
		"TotalEbooks":         totalEbooks,
		"GetTotalSendEbooks":  totalSendEbooks,
		"GetTotalClients":     totalClients,
		"GetEngagementMetric": engagementMetric,
		"EbookStats":          ebookStats,
	}, "admin")
}
