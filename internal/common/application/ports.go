package common_application

type ReceitaFederalService interface {
	ConsultaCPF(cpf, dataNascimento string) (*ReceitaFederalResponse, error)
}
