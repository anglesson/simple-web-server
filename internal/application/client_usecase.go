package application

import (
	"encoding/csv"
	"errors"
	"strings"

	common_application "github.com/anglesson/simple-web-server/internal/common/application"
	common_domain "github.com/anglesson/simple-web-server/internal/common/domain"
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
	cpf, err := common_domain.NewCPF(input.CPF)
	if err != nil {
		return nil, err
	}
	foundClient := cuc.clientRepository.FindByCPF(cpf)
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

	err = cuc.clientRepository.Create(client)
	if err != nil {
		return nil, err
	}

	return &CreateClientOutput{ID: client.ID}, nil
}

func (cuc *ClientUseCase) UpdateClient(input UpdateClientInput) (*UpdateClientOutput, error) {
	client, err := cuc.clientRepository.FindByID(input.ID)
	if err != nil {
		return nil, errors.New("client not found")
	}

	validCPF, err := common_domain.NewCPF(input.CPF)
	if err != nil {
		return nil, err
	}

	validBirthDay, err := common_domain.NewBirthDate(input.BirthDay)
	if err != nil {
		return nil, err
	}

	client.Name = input.Name
	client.CPF = validCPF
	client.BirthDay = validBirthDay
	client.Email = input.Email
	client.Phone = input.Phone

	result, err := cuc.receitaFederalService.ConsultCPF(client.CPF, *client.BirthDay)
	if err != nil {
		return nil, errors.New("failed to validate CPF")
	}

	client.Name = result.NomeDaPF

	err = cuc.clientRepository.Update(client)
	if err != nil {
		return nil, err
	}

	return &UpdateClientOutput{ID: client.ID}, nil
}

func (cuc *ClientUseCase) ImportClients(input ImportClientsInput) (*ImportClientsOutput, error) {
	if !strings.HasSuffix(input.FileName, ".csv") {
		return nil, errors.New("file must be a CSV")
	}

	reader := csv.NewReader(strings.NewReader(string(input.File)))
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, errors.New("failed to read CSV file")
	}

	var clients []*domain.Client
	for i, row := range rows {
		if i == 0 { // Skip header
			continue
		}

		if len(row) < 5 {
			return nil, errors.New("invalid CSV format")
		}

		client, err := domain.NewClient(row[0], row[1], row[2], row[3], row[4])
		if err != nil {
			return nil, err
		}

		result, err := cuc.receitaFederalService.ConsultCPF(client.CPF, *client.BirthDay)
		if err != nil {
			return nil, errors.New("failed to validate CPF")
		}

		client.Name = result.NomeDaPF
		clients = append(clients, client)
	}

	err = cuc.clientRepository.CreateBatch(clients)
	if err != nil {
		return nil, err
	}

	return &ImportClientsOutput{ImportedCount: len(clients)}, nil
}

func (cuc *ClientUseCase) ListClients(input ListClientsInput) (*ListClientsOutput, error) {
	query := ClientQuery{
		Term: input.Term,
		Pagination: &Pagination{
			Page:     input.Page,
			PageSize: input.PageSize,
		},
	}

	clients, total, err := cuc.clientRepository.List(query)
	if err != nil {
		return nil, err
	}

	return &ListClientsOutput{
		Clients:    clients,
		TotalCount: total,
		Page:       input.Page,
		PageSize:   input.PageSize,
	}, nil
}
