package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/anglesson/simple-web-server/domain"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/pkg/utils"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type InputCreateUser struct {
	Username             string
	Email                string
	Password             string
	PasswordConfirmation string
}

type UserService interface {
	CreateUser(input InputCreateUser) (*domain.User, error)
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

func (us *UserServiceImpl) CreateUser(input InputCreateUser) (*domain.User, error) {
	if input.Password != input.PasswordConfirmation {
		return nil, errors.New("passwords do not match")
	}

	input.Username = strings.TrimSpace(input.Username)

	existingUser := us.userRepository.FindByUserEmail(input.Email)
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	user, err := domain.NewUser(input.Username, input.Email, input.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid user data: %w", err)
	}

	hashedPassword := us.encrypter.HashPassword(input.Password)
	passwordVO, err := domain.NewPassword(hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = passwordVO

	err = us.userRepository.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
