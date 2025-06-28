package client

import (
	"errors"
	"log"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/shared/database"
	"gorm.io/gorm"
)

type GormRepository struct {
}

func NewGormRepository() *GormRepository {
	return &GormRepository{}
}

func (cr *GormRepository) Save(client *models.Client) error {
	// Start a transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var existingClient models.Client
	if client.ID != 0 {
		// Check if client exists
		if err := tx.First(&existingClient, client.ID).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				tx.Rollback()
				log.Printf("Erro ao buscar cliente: %s", err)
				return errors.New("erro ao buscar cliente")
			}
		}
	}

	// Handle contact
	if existingClient.ContactID != 0 {
		// Update existing contact
		client.Contact.ID = existingClient.ContactID
		if err := tx.Model(&models.Contact{}).Where("id = ?", existingClient.ContactID).Updates(client.Contact).Error; err != nil {
			tx.Rollback()
			log.Printf("Erro ao atualizar contato: %s", err)
			return errors.New("erro ao atualizar contato")
		}
	} else {
		// Create new contact
		if err := tx.Save(&client.Contact).Error; err != nil {
			tx.Rollback()
			log.Printf("Erro ao salvar contato: %s", err)
			return errors.New("erro ao salvar contato")
		}
	}

	// Update the client's ContactID
	client.ContactID = client.Contact.ID

	// Save or update the client
	if existingClient.ID != 0 {
		if err := tx.Model(&models.Client{}).Where("id = ?", existingClient.ID).Updates(client).Error; err != nil {
			tx.Rollback()
			log.Printf("Erro ao atualizar cliente: %s", err)
			return errors.New("erro ao atualizar cliente")
		}
	} else {
		if err := tx.Save(client).Error; err != nil {
			tx.Rollback()
			log.Printf("Erro ao salvar cliente: %s", err)
			return errors.New("erro ao salvar cliente")
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Printf("Erro ao commitar transação: %s", err)
		return errors.New("erro ao salvar dados")
	}

	return nil
}

func (cr *GormRepository) FindClientsByCreator(creator *models.Creator, query ClientQuery) (*[]models.Client, error) {
	var clients []models.Client

	err := database.DB.
		Offset(query.Pagination.GetOffset()).
		Limit(query.Pagination.GetLimit()).
		Model(&models.Client{}).
		Joins("JOIN client_creators ON client_creators.client_id = clients.id").
		Where("client_creators.creator_id = ?", creator.ID).
		Preload("Contact").
		Preload("Creators").
		Scopes(ContainsNameCpfEmailOrPhoneWith(query.Term)).
		Find(&clients).
		Error
	if err != nil {
		log.Printf("Erro na busca de ebooks: %s", err)
		return nil, errors.New("erro na busca de ebooks")
	}

	query.Pagination.Total = int64(len(clients))

	return &clients, nil
}

func ContainsNameCpfEmailOrPhoneWith(term string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		searchTerm := "%" + term + "%"
		return db.Joins("LEFT JOIN contacts ON clients.contact_id = contacts.id").
			Where("clients.name LIKE ? OR clients.cpf LIKE ? OR contacts.email LIKE ? OR contacts.phone LIKE ?",
				searchTerm, searchTerm, searchTerm, searchTerm)
	}
}

func (cr *GormRepository) FindByIDAndCreators(client *models.Client, clientID uint, creator string) error {
	err := database.DB.
		Preload("Contact").
		Preload("Creators").
		First(&client).
		Joins("JOIN client_creators ON client_creators.client_id = clients.id AND client_creators.creator_id = creators.id").
		Joins("JOIN contacts ON contacts.id = creators.contact_id").
		Where("clients.id = ? AND contacts.email = ?", clientID, creator).
		Error
	if err != nil {
		log.Printf("Erro na busca do client: %s", err)
		return errors.New("não foi possível recuperar as informações do cliente")
	}
	return nil
}

func (cr *GormRepository) InsertBatch(clients []*models.Client) error {
	err := database.DB.CreateInBatches(clients, 1000).Error
	if err != nil {
		log.Printf("[CLIENT-REPOSITORY] ERROR: %s", err)
		return errors.New("falha no processamento dos clientes")
	}
	return nil
}

func (cr *GormRepository) FindByClientsWhereEbookNotSend(creator *models.Creator, query ClientQuery) (*[]models.Client, error) {
	var clients []models.Client
	err := database.DB.Debug().
		Offset(query.Pagination.GetOffset()).
		Limit(query.Pagination.GetLimit()).
		Model(&models.Client{}).
		Joins("JOIN client_creators ON client_creators.client_id = clients.id and client_creators.creator_id = ?", creator.ID).
		Where("clients.id NOT IN (SELECT client_id FROM purchases WHERE ebook_id = ?)", query.EbookID).
		Preload("Contact").
		Preload("Creators").
		Scopes(ContainsNameCpfEmailOrPhoneWith(query.Term)).
		Find(&clients).Error

	if err != nil {
		log.Printf("Erro na busca de clientes: %s", err)
		return nil, errors.New("erro na busca de clientes")
	}

	return &clients, nil
}

func (cr *GormRepository) FindByClientsWhereEbookWasSend(creator *models.Creator, query ClientQuery) (*[]models.Client, error) {
	var clients []models.Client
	err := database.DB.Debug().
		Offset(query.Pagination.GetOffset()).
		Limit(query.Pagination.GetLimit()).
		Model(&models.Client{}).
		Joins("JOIN client_creators ON client_creators.client_id = clients.id and client_creators.creator_id = ?", creator.ID).
		Joins("JOIN purchases ON purchases.client_id = clients.id").
		Where("clients.id IN (SELECT client_id FROM purchases WHERE ebook_id = ?)", query.EbookID).
		Preload("Contact").
		Preload("Creators").
		Preload("Purchases").
		Scopes(ContainsNameCpfEmailOrPhoneWith(query.Term)).
		Find(&clients).Error

	if err != nil {
		log.Printf("Erro na busca de clientes: %s", err)
		return nil, errors.New("erro na busca de clientes")
	}

	return &clients, nil
}

func (cr *GormRepository) FindByEmail(email string) (*models.Client, error) {
	var client models.Client
	err := database.DB.
		Model(&models.Client{}).
		Joins("JOIN contacts ON contacts.id = clients.contact_id").
		Preload("Contact").
		Preload("Creators").
		Where("contacts.email = ?", email).
		First(&client).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Printf("Erro na busca do cliente: %s", err)
		return nil, errors.New("erro na busca do cliente")
	}
	return &client, nil
}
