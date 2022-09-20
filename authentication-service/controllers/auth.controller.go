package controllers

import (
	"authentication/dto"
	"authentication/helpers"
	"authentication/services"
	"errors"
	"fmt"
	"net/http"
)

type AuthController interface {
	Authenticate(w http.ResponseWriter, r *http.Request)
}

type authController struct {
	userService services.UserService
	authService services.AuthService
}

func NewAuthController(userService services.UserService, authService services.AuthService) AuthController {
	return &authController{
		userService: userService,
		authService: authService,
	}
}

func (controller *authController) Authenticate(w http.ResponseWriter, r *http.Request) {
	var loginDTO dto.LoginDTO

	err := helpers.ReadJSON(w, r, &loginDTO)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate the user against the database
	user, err := controller.userService.GetByEmail(loginDTO.Email)
	if err != nil {
		helpers.ErrorJSON(w, errors.New("wrong email"), http.StatusUnauthorized)
		return
	}

	valid, err := controller.authService.PasswordMatches(user.Password, loginDTO.Password)
	if err != nil || !valid {
		helpers.ErrorJSON(w, errors.New("wrong password"), http.StatusUnauthorized)
		return
	}

	payload := helpers.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	helpers.WriteJSON(w, http.StatusAccepted, payload)
}
