package repositories

import (
	"log"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/shared/database"
)

type DashboardRepository struct {
	UserID uint
}

func NewDashboardRepository(userID uint) *DashboardRepository {
	return &DashboardRepository{
		UserID: userID,
	}
}

// TODO: Create a service and use EbookRepository
func (dr *DashboardRepository) GetTotalEbooks() int64 {
	var count int64
	err := database.DB.
		Model(models.Ebook{}).
		InnerJoins("Creator").
		Where("user_id = ?", dr.UserID).Count(&count).Error
	if err != nil {
		log.Panicf("Erro na busca de totais: %s", err)
		return 0
	}

	return count
}

// TODO: Create a service and use EbookRepository
func (dr *DashboardRepository) GetTotalSendEbooks() int64 {
	var total int64
	err := database.DB.
		Model(models.Purchase{}).
		Joins("JOIN ebooks ON ebooks.id = purchases.ebook_id").
		Joins("JOIN creators ON creators.id = ebooks.creator_id").
		Joins("JOIN users ON users.id = creators.user_id").
		Where("users.id = ?", dr.UserID).
		Count(&total).Error
	if err != nil {
		log.Panicf("Erro na busca de totais: %s", err)
		return 0
	}

	return total
}

func (dr *DashboardRepository) GetTotalClients() int64 {
	var total int64
	err := database.DB.
		Model(models.Client{}).
		Joins("JOIN client_creators ON client_creators.client_id = clients.id").
		Joins("JOIN creators ON creators.id = client_creators.creator_id").
		Joins("JOIN users ON users.id = creators.user_id").
		Where("users.id = ?", dr.UserID).
		Count(&total).Error
	if err != nil {
		log.Panicf("Erro na busca de totais: %s", err)
		return 0
	}

	return total
}

func (dr *DashboardRepository) GetEngagementMetric() float64 {
	var metric float64
	err := database.DB.Model(&models.Purchase{}).
		Select("ROUND(AVG(CAST(downloads_used AS float) / NULLIF(download_limit, 0)) * 100, 2) AS average_engagement").
		Joins("INNER JOIN ebooks ON ebooks.id = purchases.ebook_id").
		Joins("INNER JOIN creators ON creators.id = ebooks.creator_id").
		Where("creators.user_id = ?", dr.UserID).
		Scan(&metric).Error

	if err != nil {
		log.Panicf("Erro na busca de engajamento: %s", err)
		return 0
	}

	return metric
}

func (dr *DashboardRepository) GetLastPurchases() []models.Purchase {
	var purchases []models.Purchase

	err := database.DB.
		Model(&models.Purchase{}).
		Preload("Ebook").
		Joins("INNER JOIN ebooks ON ebooks.id = purchases.ebook_id").
		Joins("INNER JOIN creators ON creators.id = ebooks.creator_id").
		Where("creators.user_id = ?", dr.UserID).
		Group("ebooks.id").
		Order("purchases.created_at DESC").
		Limit(10).
		Find(&purchases).Error

	if err != nil {
		log.Panicf("Erro ao buscar últimos envios/compras: %s", err)
	}

	return purchases
}

type EbookStats struct {
	ID              uint   `json:"id"`
	Title           string `json:"title"`
	TotalPurchases  int64  `json:"total_purchases"`
	TotalDownloads  int64  `json:"total_downloads"`
	TotalClients    int64  `json:"total_clients"`
	UniqueDownloads int64  `json:"unique_downloads"`
}

func (es *EbookStats) PercentDownloads() float64 {
	return float64((es.UniqueDownloads * 100) / es.TotalClients)
}

func (dr *DashboardRepository) GetEbookStats() ([]EbookStats, error) {
	var stats []EbookStats

	err := database.DB.
		Table("ebooks").
		Select(`
			ebooks.id,
			ebooks.title,
			COUNT(purchases.id) AS total_purchases,
			COALESCE(SUM(purchases.downloads_used), 0) AS total_downloads,
			COUNT(DISTINCT CASE WHEN purchases.downloads_used > 0 THEN purchases.client_id END) AS unique_downloads,
			COUNT(DISTINCT purchases.client_id) AS total_clients
		`).
		Joins("INNER JOIN purchases ON purchases.ebook_id = ebooks.id").
		Joins("INNER JOIN creators ON creators.id = ebooks.creator_id").
		Where("creators.user_id = ?", dr.UserID).
		Group("ebooks.id").
		Scan(&stats).Error

	if err != nil {
		log.Printf("Erro ao buscar estatísticas dos ebooks: %v", err)
		return nil, err
	}

	return stats, nil
}
