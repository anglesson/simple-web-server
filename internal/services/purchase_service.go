package services

import (
	"errors"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/pkg/mail"
	"github.com/anglesson/simple-web-server/pkg/storage"
)

type PurchaseService struct {
	purchaseRepository *repository.PurchaseRepository
	mailService        *mail.EmailService
}

func NewPurchaseService(purchaseRepository *repository.PurchaseRepository, mailService *mail.EmailService) *PurchaseService {
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

func (ps *PurchaseService) GetEbookFile(purchaseID int) (string, error) {
	purchase, err := ps.purchaseRepository.FindByID(uint(purchaseID))
	if err != nil {
		return "", errors.New(err.Error())
	}

	if !purchase.AvailableDownloads() {
		return "", errors.New("não é possível realizar o download, limite de downloads atingido")
	}

	if purchase.IsExpired() {
		return "", errors.New("não é possível realizar o download, o pedido está expirado")
	}

	fileLocation, err := storage.GetFile(purchase.Ebook.File)
	if err != nil {
		return "", errors.New("erro no download do objeto")
	}

	outputFilePath, err := ApplyWatermark(fileLocation, purchase.Client.Name, purchase.Client.CPF, purchase.Client.Contact.Email)
	if err != nil {
		return "", err
	}

	purchase.UseDownload()
	ps.purchaseRepository.Update(purchase)

	return outputFilePath, nil
}
