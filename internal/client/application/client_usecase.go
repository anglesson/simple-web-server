package application

import (
	"errors"

	"github.com/anglesson/simple-web-server/internal/client/domain"
)

type CreateClientUseCase struct {
	clientRepository ClientRepositoryInterface
}

func NewCreateClientUseCase(clientRepository ClientRepositoryInterface) *CreateClientUseCase {
	return &CreateClientUseCase{
		clientRepository: clientRepository,
	}
}

func (cuc *CreateClientUseCase) Execute(input CreateClientInput) (*CreateClientOutput, error) {
	foundClient := cuc.clientRepository.FindByCPF(input.CPF)
	if foundClient != nil {
		return nil, errors.New("client already exists")
	}

	client, err := domain.NewClient(input.Name, input.CPF, input.BirthDay, input.Email, input.Phone)
	if err != nil {
		return nil, err
	}

	cuc.clientRepository.Create(client)

	return &CreateClientOutput{}, nil
}
