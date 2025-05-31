package application

import (
	"errors"

	"github.com/anglesson/simple-web-server/internal/client/domain"
)

type ClientUseCase struct {
	clientRepository ClientRepositoryInterface
}

func NewClientUseCase() *ClientUseCase {
	return &ClientUseCase{}
}

func (cuc *ClientUseCase) Create(input CreateClientInput) error {
	foundClient := cuc.clientRepository.FindByCPF(input.CPF)
	if foundClient != nil {
		return errors.New("client already exists")
	}

	client, err := domain.NewClient(input.Name, input.CPF, input.BirthDay, input.Email, input.Phone)
	if err != nil {
		return err
	}

	cuc.clientRepository.Create(client)

	return nil
}
