package service

import (
	"errors"
	"fmt"
	"github.com/anglesson/simple-web-server/domain"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/pkg/utils"
	"strings"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type UserService interface {
	CreateUser(username, email, password, passwordConfirmation string) (*domain.User, error)
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

func (us *UserServiceImpl) CreateUser(username, email, password, passwordConfirmation string) (*domain.User, error) {
	if password != passwordConfirmation {
		return nil, errors.New("passwords do not match")
	}

	username = strings.TrimSpace(username)

	existingUser := us.userRepository.FindByEmail(email)
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	user, err := domain.NewUser(username, email, password)
	if err != nil {
		return nil, fmt.Errorf("invalid user data: %w", err)
	}

	hashedPassword := us.encrypter.HashPassword(password)
	passwordVO, err := domain.NewPassword(hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = passwordVO

	err = us.userRepository.Save(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
