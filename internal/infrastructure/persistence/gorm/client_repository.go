package persistence

import (
	"github.com/anglesson/simple-web-server/internal/domain"
	"gorm.io/gorm"
)

type GormClientRepository struct {
	db *gorm.DB
}

func NewGormClientRepository(db *gorm.DB) *GormClientRepository {
	return &GormClientRepository{db: db}
}

func (r *GormClientRepository) Create(client *domain.Client) {
	r.db.Create(client)
}

func (r *GormClientRepository) FindByCPF(cpf string) *domain.Client {
	var client domain.Client
	r.db.Where("cpf = ?", cpf).First(&client)
	return &client
}
