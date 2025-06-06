package application

import (
	"github.com/anglesson/simple-web-server/internal/domain"
)

type ClientRepositoryInterface interface {
	Create(client *domain.Client) error
	Update(client *domain.Client) error
	FindByCPF(cpf domain.CPF) *domain.Client
	FindByID(id uint) (*domain.Client, error)
	FindByCreatorID(creatorID uint, query ClientQuery) ([]*domain.Client, error)
	CreateBatch(clients []*domain.Client) error
	List(query ClientQuery) ([]*domain.Client, int, error)
}

type ClientUseCasePort interface {
	CreateClient(input CreateClientInput) (*CreateClientOutput, error)
	UpdateClient(input UpdateClientInput) (*UpdateClientOutput, error)
	ImportClients(input ImportClientsInput) (*ImportClientsOutput, error)
	ListClients(input ListClientsInput) (*ListClientsOutput, error)
}

type ClientQuery struct {
	Term       string
	Pagination *Pagination
	CreatorID  uint
}

type CPFServicePort interface {
	ConsultCPF(cpf domain.CPF, birthDay domain.BirthDate) (CPFOutput, error)
}

type CPFOutput struct {
	NumeroDeCPF            string `json:"numero_de_cpf"`
	NomeDaPF               string `json:"nome_da_pf"`
	DataNascimento         string `json:"data_nascimento"`
	SituacaoCadastral      string `json:"situacao_cadastral"`
	DataInscricao          string `json:"data_inscricao"`
	DigitoVerificador      string `json:"digito_verificador"`
	ComprovanteEmitido     string `json:"comprovante_emitido"`
	ComprovanteEmitidoData string `json:"comprovante_emitido_data"`
}
