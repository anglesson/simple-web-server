package repository

import (
	"errors"
	"log"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/pkg/database"
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

func (pr *PurchaseRepository) FindByID(id uint) (*models.Purchase, error) {
	var purchase models.Purchase
	log.Printf("Buscando a compra: %v", id)
	err := database.DB.Preload("Client.Contact").
		Preload("Ebook").First(&purchase, id).Error
	if err != nil {
		log.Printf("Erro na busca da compra: %s", err)
		return nil, errors.New("erro na busca da compra")
	}

	return &purchase, nil
}

func (pr *PurchaseRepository) Update(purchase *models.Purchase) error {
	if purchase.ID == 0 {
		log.Printf("error to update purchase: %v", purchase)
		return errors.New("erro ao atualizar downloads")
	}

	err := database.DB.Save(purchase).Error
	if err != nil {
		log.Printf("Erro na busca da compra: %s", err)
		return errors.New("erro na busca da compra")
	}

	return nil
}
