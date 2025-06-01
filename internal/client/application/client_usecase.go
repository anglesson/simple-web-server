package client_application

import (
	"errors"

	client_domain "github.com/anglesson/simple-web-server/internal/client/domain"
	common_application "github.com/anglesson/simple-web-server/internal/common/application"
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

	result, err := cuc.receitaFederalService.ConsultCPF(input.CPF, input.BirthDay)
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
