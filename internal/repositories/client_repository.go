package repositories

import (
	"errors"
	"log"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/shared/database"
)

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
