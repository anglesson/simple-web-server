package repositories

import (
	"errors"
	"log"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/shared/database"
)

type PurchaseRepository struct {
}

func NewPurchaseRepository() *PurchaseRepository {
	return &PurchaseRepository{}
}

func (pr *PurchaseRepository) CreateManyPurchases(purchases []*models.Purchase) error {
	err := database.DB.CreateInBatches(purchases, 1000).Error
	if err != nil {
		log.Printf("[PURCHASE-REPOSITORY] ERROR: %s", err)
		return errors.New("falha no processamento do envio")
	}

	// Recarrega as purchases com os dados relacionados
	ids := make([]uint, len(purchases))
	for i, p := range purchases {
		ids[i] = p.ID
	}

	err = database.DB.
		Preload("Client.Contact").
		Preload("Ebook").
		Find(&purchases, "id IN ?", ids).Error
	if err != nil {
		log.Printf("[PURCHASE-REPOSITORY] LOAD ERROR: %s", err)
		return errors.New("falha ao carregar dados relacionados")
	}

	return nil
}
