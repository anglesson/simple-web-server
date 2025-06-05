package common_application

import common_domain "github.com/anglesson/simple-web-server/internal/common/domain"

type CPFServicePort interface {
	ConsultCPF(cpf common_domain.CPF, birthDay common_domain.BirthDate) (CPFOutput, error)
}
