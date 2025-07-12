package service

import (
	"github.com/anglesson/simple-web-server/domain"
	"github.com/anglesson/simple-web-server/internal/repository"
	"github.com/anglesson/simple-web-server/pkg/utils"
)

type UserService interface {
	CreateUser(username, email, password string) (*domain.User, error)
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

func (us *UserServiceImpl) CreateUser(username, email, password string) (*domain.User, error) {
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
