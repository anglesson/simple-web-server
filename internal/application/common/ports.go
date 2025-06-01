package common

import (
	"github.com/anglesson/simple-web-server/internal/domain/client"
)

type ReceitaFederalData struct {
	NumeroDeCPF            string `json:"numero_de_cpf"`
	NomeDaPF               string `json:"nome_da_pf"`
	DataNascimento         string `json:"data_nascimento"`
	SituacaoCadastral      string `json:"situacao_cadastral"`
	DataInscricao          string `json:"data_inscricao"`
	DigitoVerificador      string `json:"digito_verificador"`
	ComprovanteEmitido     string `json:"comprovante_emitido"`
	ComprovanteEmitidoData string `json:"comprovante_emitido_data"`
}

type ClientRepositoryInterface interface {
	Create(*client.Client)
	FindByCPF(cpf string) *client.Client
}

type ReceitaFederalServiceInterface interface {
	Search(cpf, birthDay string) (ReceitaFederalData, error)
}
