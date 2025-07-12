package service

import (
	"github.com/anglesson/simple-web-server/domain"
	"github.com/anglesson/simple-web-server/internal/repository"
)

type UserService interface {
	CreateUser(username, email, password string) (any, error)
}

type UserServiceImpl struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &UserServiceImpl{
		userRepository: userRepository,
	}
}

func (us *UserServiceImpl) CreateUser(username, email, password string) (any, error) {
	user, err := domain.NewUser(username, email, password)
	if err != nil {
		return nil, err
	}

	err = us.userRepository.Save(user)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
