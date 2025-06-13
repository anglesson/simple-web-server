package common_application

type ReceitaFederalServicePort interface {
	ConsultaCPF(cpf, dataNascimento string) *ReceitaFederalResponse
}
