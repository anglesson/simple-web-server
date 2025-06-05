package application

import (
	"errors"

	common_application "github.com/anglesson/simple-web-server/internal/common/application"
	domain "github.com/anglesson/simple-web-server/internal/domain"
)

type ClientUseCase struct {
	clientRepository      ClientRepositoryInterface
	receitaFederalService common_application.CPFServicePort
}

func NewClientUseCase(clientRepository ClientRepositoryInterface, receitaFederalService common_application.CPFServicePort) *ClientUseCase {
	return &ClientUseCase{
		clientRepository:      clientRepository,
		receitaFederalService: receitaFederalService,
	}
}

func (cuc *ClientUseCase) CreateClient(input CreateClientInput) (*CreateClientOutput, error) {
	foundClient := cuc.clientRepository.FindByCPF(input.CPF)
	if foundClient != nil {
		return nil, errors.New("client already exists")
	}

	client, err := domain.NewClient(input.Name, input.CPF, input.BirthDay, input.Email, input.Phone)
	if err != nil {
		return nil, err
	}

	result, err := cuc.receitaFederalService.ConsultCPF(client.CPF, *client.BirthDay)
	if err != nil {
		return nil, errors.New("failed to validate CPF")
	}

	client.Name = result.NomeDaPF

	cuc.clientRepository.Create(client)

	return &CreateClientOutput{}, nil
}
