package data

import (
	"errors"
	"log/slog"

	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

func (p *PaymentRepository) InsertAccount(accountModel *AccountModel) error {
	err := p.db.Save(accountModel).Error
	if err != nil {
		slog.Error("Erro ao salvar conta: %v", err)
		return err
	}
	return nil
}

func (p *PaymentRepository) FindAccountBySellerID(sellerID uint) (*AccountModel, error) {
	var sellerAccount AccountModel
	err := p.db.Where("creator_id = ?", sellerID).First(&sellerAccount).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Não encontrado - não é um erro, apenas ausência de dados
		}
		slog.Error("[PaymentRepository] Erro ao buscar conta: %v", err)
		return nil, err // Erro real de banco de dados
	}
	return &sellerAccount, nil
}
