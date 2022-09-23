package controllers

import (
	"log"
	"mail/config"
	"mail/dto"
	"mail/helpers"
	"net/http"
)

type MailController interface {
	SendMail(w http.ResponseWriter, r *http.Request)
}

type mailController struct {
	Mailer config.Mail
}

func NewMailController(Mailer config.Mail) MailController {
	return &mailController{
		Mailer: Mailer,
	}
}

func (controller *mailController) SendMail(w http.ResponseWriter, r *http.Request) {

	var requestPayload dto.MailMessageDTO

	err := helpers.ReadJSON(w, r, &requestPayload)
	if err != nil {
		log.Println(err)
		helpers.ErrorJSON(w, err)
		return
	}

	msg := config.Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = controller.Mailer.SendSMTPMessage(msg)
	if err != nil {
		log.Println(err)
		helpers.ErrorJSON(w, err)
		return
	}

	payload := helpers.JsonResponse{
		Error:   false,
		Message: "sent to " + requestPayload.To,
	}

	helpers.WriteJSON(w, http.StatusAccepted, payload)
}
