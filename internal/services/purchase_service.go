package services

import (
	"errors"

	"github.com/anglesson/simple-web-server/internal/mail"
	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repositories"
)

type PurchaseService struct {
	purchaseRepository *repositories.PurchaseRepository
	mailService        *mail.EmailService
}

func NewPurchaseService(purchaseRepository *repositories.PurchaseRepository, mailService *mail.EmailService) *PurchaseService {
	return &PurchaseService{
		purchaseRepository: purchaseRepository,
		mailService:        mailService,
	}
}

func (ps *PurchaseService) CreatePurchase(ebookId uint, clients []uint) error {
	var purchases []*models.Purchase

	for _, clientId := range clients {
		if clientId != 0 && ebookId != 0 {
			purchases = append(purchases, models.NewPurchase(ebookId, uint(clientId)))
		}
	}

	err := ps.purchaseRepository.CreateManyPurchases(purchases)
	if err != nil {
		return errors.New(err.Error())
	}

	go ps.mailService.SendLinkToDownload(purchases)
	return nil
}
