package repositories

import (
	"errors"
	"log"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/shared/database"
	"gorm.io/gorm"
)

type ClientQuery struct {
	Term       string
	Pagination *Pagination
}

type ClientRepository struct {
}

func NewClientRepository() *ClientRepository {
	return &ClientRepository{}
}

func (cr *ClientRepository) Save(client *models.Client) error {
	result := database.DB.Save(&client)
	if result.Error != nil {
		log.Panic("Erro ao salvar client")
		return errors.New("erro ao salvar cliente")
	}

	return nil
}

func (cr *ClientRepository) FindClientsByCreator(creator *models.Creator, query ClientQuery) (*[]models.Client, error) {
	var clients []models.Client

	err := database.DB.
		Offset(query.Pagination.GetOffset()).
		Limit(query.Pagination.GetLimit()).
		Preload("Contact").
		Preload("Creators").
		// Scopes(ContainsNameCpfEmailOrPhoneWith(query.Term)).
		Find(&clients).
		Error
	if err != nil {
		log.Printf("Erro na busca de ebooks: %s", err)
		return nil, errors.New("erro na busca de dados")
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
