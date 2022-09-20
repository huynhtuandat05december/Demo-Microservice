package controllers

import (
	"broken/dto"
	"broken/helpers"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestDTO dto.RequestDTO

	err := helpers.ReadJSON(w, r, &requestDTO)

	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}

	switch requestDTO.Action {
	case "log":
		log(w, requestDTO.Log)
	case "auth":
		authenticate(w, requestDTO.Auth)
	default:
		helpers.ErrorJSON(w, errors.New("unknown action"))
	}
}

func authenticate(w http.ResponseWriter, payload dto.LoginDTO) {
	jsonData, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}
	//call service
	url := "http://authentication:8001/authenticate"
	method := "POST"
	request, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}
	client := &http.Client{}
	responseFromService, err := client.Do(request)
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}
	defer responseFromService.Body.Close()
	fmt.Printf("%v", responseFromService)
	// make sure we get back the correct status code
	if responseFromService.StatusCode == http.StatusUnauthorized {
		helpers.ErrorJSON(w, errors.New("invalid credentials"))
		return
	}
	if responseFromService.StatusCode != http.StatusAccepted {
		helpers.ErrorJSON(w, errors.New("error calling auth service"))
		return
	}

	// create a variable we'll read response.Body into
	var jsonFromService helpers.JsonResponse

	// decode the json from the auth service
	err = json.NewDecoder(responseFromService.Body).Decode(&jsonFromService)
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		helpers.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var response helpers.JsonResponse
	response.Error = false
	response.Message = "Authenticated!"
	response.Data = jsonFromService.Data

	helpers.WriteJSON(w, http.StatusAccepted, response)

}

func log(w http.ResponseWriter, payload dto.LogDTO) {
	jsonData, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}
	//call service
	url := "http://logger:8002/log"
	method := "POST"
	request, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}

	client := &http.Client{}
	responseFromService, err := client.Do(request)
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}
	defer responseFromService.Body.Close()
	if responseFromService.StatusCode != http.StatusAccepted {
		helpers.ErrorJSON(w, errors.New("error calling logger service"))
		return
	}

	// create a variable we'll read response.Body into

	var response helpers.JsonResponse
	response.Error = false
	response.Message = "Logged!"

	helpers.WriteJSON(w, http.StatusAccepted, response)

}
