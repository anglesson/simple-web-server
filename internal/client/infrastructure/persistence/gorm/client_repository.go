package client_persistence

import (
	client_domain "github.com/anglesson/simple-web-server/internal/client/domain"
	"gorm.io/gorm"
)

type GormClientRepository struct {
	db *gorm.DB
}

func NewGormClientRepository(db *gorm.DB) *GormClientRepository {
	return &GormClientRepository{db: db}
}

func (r *GormClientRepository) Create(client *client_domain.Client) {
	r.db.Create(client)
}

func (r *GormClientRepository) FindByCPF(cpf string) *client_domain.Client {
	var client client_domain.Client
	r.db.Where("cpf = ?", cpf).First(&client)
	return &client
}
