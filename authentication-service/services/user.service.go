package services

import (
	"authentication/models"
	"authentication/repository"
)

type UserService interface {
	GetByEmail(email string) (*models.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (service *userService) GetByEmail(email string) (*models.User, error) {
	existUser, err := service.userRepository.GetByEmail(email)
	return existUser, err
}
