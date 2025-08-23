package business

import (
	"errors"
	"log/slog"

	"github.com/anglesson/simple-web-server/internal/payment/data"
)

var ErrGateway = errors.New("erro ao criar conta no gateway")
var ErrDatabase = errors.New("erro ao salvar conta no banco de dados")
var ErrAlreadyExists = errors.New("vendedor já possui uma conta")

type PaymentServiceImpl struct {
	paymentGateway          data.PaymentGateway
	sellerAccountRepository *data.PaymentRepository
}

func NewPaymentService(paymentGateway data.PaymentGateway, sellerAccountRepository *data.PaymentRepository) *PaymentServiceImpl {
	return &PaymentServiceImpl{
		paymentGateway,
		sellerAccountRepository,
	}
}

func (s *PaymentServiceImpl) CreateAccount(name string, sellerID uint) (*data.AccountModel, error) {
	hasAccount, err := s.sellerAccountRepository.FindAccountBySellerID(sellerID)
	if err != nil {
		return nil, ErrDatabase
	}

	if hasAccount != nil {
		return nil, ErrAlreadyExists
	}

	stripeAccountID, err := s.paymentGateway.CreateAccount(name, sellerID)
	if err != nil {
		slog.Error("Erro ao criar conta no Stripe: %v", err)
		return nil, ErrGateway
	}
	account := &data.AccountModel{
		Origin:    "stripe",
		AccountID: stripeAccountID,
		SellerID:  sellerID,
	}
	err = s.sellerAccountRepository.InsertAccount(account)
	if err != nil {
		slog.Error("Erro ao salvar conta no banco de dados: %v", err)
		return nil, ErrDatabase
	}

	return account, nil
}
