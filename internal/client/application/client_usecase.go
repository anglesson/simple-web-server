package application

import (
	"errors"

	"github.com/anglesson/simple-web-server/internal/client/domain"
)

type CreateClientUseCase struct {
	clientRepository      ClientRepositoryInterface
	receitaFederalService ReceitaFederalServiceInterface
}

func NewCreateClientUseCase(clientRepository ClientRepositoryInterface, receitaFederalService ReceitaFederalServiceInterface) *CreateClientUseCase {
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
		return nil, err
	}

	client, err := domain.NewClient(result.NomeDaPF, input.CPF, input.BirthDay, input.Email, input.Phone)
	if err != nil {
		return nil, err
	}

	cuc.clientRepository.Create(client)

	return &CreateClientOutput{}, nil
}
