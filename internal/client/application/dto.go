package application

type CreateClientInput struct {
	Name     string
	CPF      string
	BirthDay string
	Email    string
	Phone    string
}

type CreateClientOutput struct {
	Name     string
	CPF      string
	BirthDay string
	Email    string
	Phone    string
}

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
