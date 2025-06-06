package repository

import (
	"github.com/anglesson/simple-web-server/internal/application"
	"github.com/anglesson/simple-web-server/internal/domain"
	"gorm.io/gorm"
)

type ClientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) *ClientRepository {
	return &ClientRepository{
		db: db,
	}
}

func (r *ClientRepository) Create(client *domain.Client) error {
	return r.db.Create(client).Error
}

func (r *ClientRepository) Update(client *domain.Client) error {
	return r.db.Save(client).Error
}

func (r *ClientRepository) FindByID(id uint) (*domain.Client, error) {
	var client domain.Client
	if err := r.db.First(&client, id).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *ClientRepository) FindByCreatorID(creatorID uint, query application.ClientQuery) ([]*domain.Client, error) {
	var clients []*domain.Client
	db := r.db.Where("creator_id = ?", creatorID)

	if query.Term != "" {
		term := "%" + query.Term + "%"
		db = db.Where("name LIKE ? OR cpf LIKE ?", term, term)
	}

	if query.Pagination != nil {
		offset := (query.Pagination.Page - 1) * query.Pagination.PageSize
		db = db.Offset(offset).Limit(query.Pagination.PageSize)
	}

	if err := db.Find(&clients).Error; err != nil {
		return nil, err
	}

	return clients, nil
}

func (r *ClientRepository) List(query application.ClientQuery) ([]*domain.Client, int, error) {
	var clients []*domain.Client
	var total int64

	db := r.db.Model(&domain.Client{})

	if query.Term != "" {
		term := "%" + query.Term + "%"
		db = db.Where("name LIKE ? OR cpf LIKE ?", term, term)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if query.Pagination != nil {
		offset := (query.Pagination.Page - 1) * query.Pagination.PageSize
		db = db.Offset(offset).Limit(query.Pagination.PageSize)
	}

	if err := db.Find(&clients).Error; err != nil {
		return nil, 0, err
	}

	return clients, int(total), nil
}
