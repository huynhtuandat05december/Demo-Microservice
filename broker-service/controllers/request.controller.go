package controllers

import (
	"broken/dto"
	"broken/event"
	"broken/helpers"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type requestController struct {
	Rabbit *amqp.Connection
}

type RequestController interface {
	HandleSubmission(w http.ResponseWriter, r *http.Request)
	PushToQueue(name, msg string) error
	logEventViaRabbit(w http.ResponseWriter, logDTO dto.LogDTO)
}

func NewRequestController(Rabbit *amqp.Connection) RequestController {
	return &requestController{
		Rabbit: Rabbit,
	}

}

func (controller *requestController) HandleSubmission(w http.ResponseWriter, r *http.Request) {
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
	case "mail":
		sendMail(w, requestDTO.Mail)
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

func sendMail(w http.ResponseWriter, payload dto.SendMailDTO) {
	jsonData, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}
	url := "http://mailer:8003/send"
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
		helpers.ErrorJSON(w, errors.New("error calling mailer service"))
		return
	}

	// create a variable we'll read response.Body into

	var response helpers.JsonResponse
	response.Error = false
	response.Message = "Message sent to " + payload.To

	helpers.WriteJSON(w, http.StatusAccepted, response)
}

// logEventViaRabbit logs an event using the logger-service. It makes the call by pushing the data to RabbitMQ.
func (controller *requestController) logEventViaRabbit(w http.ResponseWriter, logDTO dto.LogDTO) {
	err := controller.PushToQueue(logDTO.Name, logDTO.Data)
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}

	var payload helpers.JsonResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	helpers.WriteJSON(w, http.StatusAccepted, payload)
}

// pushToQueue pushes a message into RabbitMQ
func (controller *requestController) PushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(controller.Rabbit)
	if err != nil {
		return err
	}

	payload := dto.LogDTO{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}
