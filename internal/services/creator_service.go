package services

import (
	"errors"
	"github.com/anglesson/simple-web-server/domain"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/pkg/gov"
	"log"

	"github.com/anglesson/simple-web-server/internal/models"
)

type CreatorService interface {
	CreateCreator(input InputCreateCreator) (*domain.Creator, error)
	FindCreatorByEmail(email string) (*models.Creator, error)
	FindCreatorByUserID(userID uint) (*models.Creator, error)
}

type InputCreateCreator struct {
	Name        string `json:"name"`
	CPF         string `json:"cpf"`
	BirthDate   string `json:"birthDate"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
}

type creatorServiceImpl struct {
	creatorRepo repository.CreatorRepository
	rfService   gov.ReceitaFederalService
}

func NewCreatorService(
	creatorRepo repository.CreatorRepository,
	receitaFederalService gov.ReceitaFederalService,
) CreatorService {
	return &creatorServiceImpl{
		creatorRepo: creatorRepo,
		rfService:   receitaFederalService,
	}
}

func (cs *creatorServiceImpl) CreateCreator(input InputCreateCreator) (*domain.Creator, error) {
	creator, err := domain.NewCreator(
		input.Name,
		input.Email,
		input.CPF,
		input.PhoneNumber,
		input.BirthDate,
	)
	if err != nil {
		return nil, err
	}

	creatorExists, err := cs.creatorRepo.FindByFilter(domain.CreatorFilter{
		CPF: creator.CPF.Value(),
	})
	if err != nil {
		return nil, err
	}

	if creatorExists != nil {
		return nil, errors.New("creator already exists")
	}

	err = cs.validateReceita(creator)
	if err != nil {
		return nil, err
	}

	err = cs.creatorRepo.Save(creator)
	if err != nil {
		return nil, err
	}
	return creator, nil
}

func (cs *creatorServiceImpl) FindCreatorByUserID(userID uint) (*models.Creator, error) {
	creator, err := cs.creatorRepo.FindCreatorByUserID(userID)
	if err != nil {
		log.Printf("Erro ao buscar creator: %s", err)
		return nil, errors.New("creator not found")
	}

	log.Printf("Usuário encontrado! ID: %v", creator.Name)

	return creator, nil
}

func (cs *creatorServiceImpl) FindCreatorByEmail(email string) (*models.Creator, error) {
	creator, err := cs.creatorRepo.FindCreatorByUserEmail(email)
	if err != nil {
		log.Printf("Erro ao buscar creator: %s", err)
		return nil, errors.New("creator not found")
	}

	log.Printf("Usuário encontrado! ID: %v", creator.Name)

	return creator, nil
}

// TODO: Mover para classe de servico e retornar apenas o nome
func (cs *creatorServiceImpl) validateReceita(creator *domain.Creator) error {
	if cs.rfService == nil {
		return errors.New("serviço da receita federal não está disponível")
	}

	response, err := cs.rfService.ConsultaCPF(creator.CPF.String(), creator.Birthdate.Format("02/01/2006"))
	if err != nil {
		return err
	}

	if !response.Status {
		return errors.New("CPF inválido ou não encontrado na receita federal")
	}

	if response.Result.NomeDaPF == "" {
		return errors.New("nome não encontrado na receita federal")
	}

	creator.Name = response.Result.NomeDaPF
	return nil
}
