package controllers

import (
	"authentication/dto"
	"authentication/helpers"
	"authentication/services"
	"bytes"
	"encoding/json"
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

	// log authentication
	err = logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}

	payload := helpers.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	helpers.WriteJSON(w, http.StatusAccepted, payload)
}

func logRequest(name string, data string) error {
	entry := dto.LogDTO{
		Name: name,
		Data: data,
	}

	jsonData, _ := json.MarshalIndent(entry, "", " ")
	logServiceURL := "http://logger:8002/log"
	method := "POST"
	request, err := http.NewRequest(method, logServiceURL, bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
