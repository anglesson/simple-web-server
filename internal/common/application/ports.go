package common_application

type CPFServicePort interface {
	ConsultCPF(cpf, birthDay string) (CPFOutput, error)
}
