package service

import (
	"errors"
	"fmt"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/pkg/mail"
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

	// TODO: Implementar com novo sistema de arquivos
	fileLocation := ""
	if len(purchase.Ebook.Files) > 0 {
		fileLocation = purchase.Ebook.Files[0].S3Key
	}
	if fileLocation == "" {
		return "", errors.New("arquivo não encontrado")
	}

	outputFilePath, err := ApplyWatermark(fileLocation, fmt.Sprintf("%s - %s - %s", purchase.Client.Name, purchase.Client.CPF, purchase.Client.Email))
	if err != nil {
		return "", err
	}

	purchase.UseDownload()
	ps.purchaseRepository.Update(purchase)

	return outputFilePath, nil
}
