package service

import (
	"errors"
	"strings"

	"github.com/anglesson/simple-web-server/internal/models"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/pkg/utils"
)

var ErrUserAlreadyExists = errors.New("usuário já existe")
var ErrInvalidCredentials = errors.New("email ou senha inválidos")

type InputCreateUser struct {
	Username             string
	Email                string
	Password             string
	PasswordConfirmation string
}

type InputLogin struct {
	Email    string
	Password string
}

type UserService interface {
	CreateUser(input InputCreateUser) (*models.User, error)
	AuthenticateUser(input InputLogin) (*models.User, error)
}

type UserServiceImpl struct {
	userRepository repository.UserRepository
	encrypter      utils.Encrypter
}

func NewUserService(userRepository repository.UserRepository, encrypter utils.Encrypter) UserService {
	return &UserServiceImpl{
		userRepository: userRepository,
		encrypter:      encrypter,
	}
}

func (us *UserServiceImpl) CreateUser(input InputCreateUser) (*models.User, error) {
	// Validate input
	if err := validateUserInput(input); err != nil {
		return nil, err
	}

	// Clean username
	input.Username = strings.TrimSpace(input.Username)

	existingUser := us.userRepository.FindByUserEmail(input.Email)
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	hashedPassword := us.encrypter.HashPassword(input.Password)
	user := models.NewUser(input.Username, hashedPassword, input.Email)

	err := us.userRepository.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserServiceImpl) AuthenticateUser(input InputLogin) (*models.User, error) {
	// Validate input
	if input.Email == "" || input.Password == "" {
		return nil, ErrInvalidCredentials
	}

	// Find user by email
	user := us.userRepository.FindByUserEmail(input.Email)
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Check password
	if !us.encrypter.CheckPasswordHash(user.Password, input.Password) {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
