package client

import (
	"errors"

	ports "github.com/anglesson/simple-web-server/internal/application/common"
	client_domain "github.com/anglesson/simple-web-server/internal/domain/client"
)

type CreateClientUseCase struct {
	clientRepository      ports.ClientRepositoryInterface
	receitaFederalService ports.ReceitaFederalServiceInterface
}

func NewCreateClientUseCase(clientRepository ports.ClientRepositoryInterface, receitaFederalService ports.ReceitaFederalServiceInterface) *CreateClientUseCase {
	return &CreateClientUseCase{
		clientRepository:      clientRepository,
		receitaFederalService: receitaFederalService,
	}
}

func (cuc *CreateClientUseCase) Execute(input CreateClientInput) (*CreateClientOutput, error) {
	foundClient := cuc.clientRepository.FindByCPF(input.CPF)
	if foundClient != nil {
		return nil, errors.New("client already exists")
	}

	result, err := cuc.receitaFederalService.Search(input.CPF, input.BirthDay)
	if err != nil {
		return nil, errors.New("failed to validate CPF")
	}

	client, err := client_domain.NewClient(result.NomeDaPF, input.CPF, input.BirthDay, input.Email, input.Phone)
	if err != nil {
		return nil, err
	}

	cuc.clientRepository.Create(client)

	return &CreateClientOutput{}, nil
}
