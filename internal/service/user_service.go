package service

import (
	"errors"
	"github.com/anglesson/simple-web-server/domain"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/pkg/utils"
)

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

	existingUser := us.userRepository.FindByEmail(email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	user, err := domain.NewUser(username, email, us.encrypter.HashPassword(password))
	if err != nil {
		return nil, err
	}

	err = us.userRepository.Save(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
