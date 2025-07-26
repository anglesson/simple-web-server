package service

import (
	"errors"
	"log"
	"time"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/pkg/gov"
)

type CreatorService interface {
	CreateCreator(input InputCreateCreator) (*models.Creator, error)
	FindCreatorByEmail(email string) (*models.Creator, error)
	FindCreatorByUserID(userID uint) (*models.Creator, error)
}

type InputCreateCreator struct {
	Name                 string `json:"name"`
	CPF                  string `json:"cpf"`
	BirthDate            string `json:"birthDate"`
	PhoneNumber          string `json:"phoneNumber"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

type creatorServiceImpl struct {
	creatorRepo repository.CreatorRepository
	rfService   gov.ReceitaFederalService
	userService UserService
}

func NewCreatorService(
	creatorRepo repository.CreatorRepository,
	receitaFederalService gov.ReceitaFederalService,
	userService UserService,
) CreatorService {
	return &creatorServiceImpl{
		creatorRepo: creatorRepo,
		rfService:   receitaFederalService,
		userService: userService,
	}
}

func (cs *creatorServiceImpl) CreateCreator(input InputCreateCreator) (*models.Creator, error) {
	// Validate input
	if err := validateCreatorInput(input); err != nil {
		return nil, err
	}

	// Parse birth date
	birthDate, err := time.Parse("2006-01-02", input.BirthDate)
	if err != nil {
		return nil, err
	}

	// Clean CPF (remove non-digits)
	cleanCPF := cleanCPF(input.CPF)

	// Check if creator already exists
	creatorExists, err := cs.creatorRepo.FindByCPF(cleanCPF)
	if err != nil {
		return nil, err
	}

	if creatorExists != nil {
		return nil, errors.New("creator already exists")
	}

	// Validate with Receita Federal
	validatedName, err := cs.validateReceita(cleanCPF, birthDate)
	if err != nil {
		return nil, err
	}

	// Create user
	inputCreateUser := InputCreateUser{
		Username:             validatedName,
		Email:                input.Email,
		Password:             input.Password,
		PasswordConfirmation: input.PasswordConfirmation,
	}

	user, err := cs.userService.CreateUser(inputCreateUser)
	if err != nil {
		return nil, err
	}

	// Create creator
	creator := models.NewCreator(
		validatedName,
		input.Email,
		cleanPhone(input.PhoneNumber),
		cleanCPF,
		birthDate,
		user.ID,
	)

	err = cs.creatorRepo.Create(creator)
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

// validateReceita validates CPF with Receita Federal and returns the validated name
func (cs *creatorServiceImpl) validateReceita(cpf string, birthDate time.Time) (string, error) {
	if cs.rfService == nil {
		return "", errors.New("serviço da receita federal não está disponível")
	}

	response, err := cs.rfService.ConsultaCPF(cpf, birthDate.Format("02/01/2006"))
	if err != nil {
		return "", err
	}

	if !response.Status {
		return "", errors.New("CPF inválido ou não encontrado na receita federal")
	}

	if response.Result.NomeDaPF == "" {
		return "", errors.New("nome não encontrado na receita federal")
	}

	return response.Result.NomeDaPF, nil
}
