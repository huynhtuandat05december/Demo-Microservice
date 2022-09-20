package services

import (
	"authentication/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	PasswordMatches(userPassword string, plainText string) (bool, error)
}

type authService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRepository repository.UserRepository) AuthService {
	return &authService{
		userRepository: userRepository,
	}

}

func (service *authService) PasswordMatches(userPassword string, plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
