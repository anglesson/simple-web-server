package common_application

import (
	common_domain "github.com/anglesson/simple-web-server/internal/common/domain"
)

type ReceitaFederalService struct{}

func NewReceitaFederalService() *ReceitaFederalService {
	return &ReceitaFederalService{}
}

func (s *ReceitaFederalService) ConsultCPF(cpf common_domain.CPF, birthDay common_domain.BirthDate) (CPFOutput, error) {
	// TODO: Implement real CPF validation
	return CPFOutput{
		NomeDaPF: "Nome da Pessoa Física",
	}, nil
}
