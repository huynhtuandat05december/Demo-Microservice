package controllers

import (
	"logger/dto"
	"logger/helpers"
	"logger/services"
	"net/http"
)

type LogController interface {
	Insert(w http.ResponseWriter, r *http.Request)
}

type logController struct {
	logService services.LogService
}

func NewLogController(logService services.LogService) LogController {
	return &logController{
		logService: logService,
	}
}

func (controller *logController) Insert(w http.ResponseWriter, r *http.Request) {
	var logDTO dto.LogDTO

	err := helpers.ReadJSON(w, r, &logDTO)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
	}

	err = controller.logService.Insert(logDTO)

	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
	}

	response := helpers.JsonResponse{
		Error:   false,
		Message: "logged",
	}

	helpers.WriteJSON(w, http.StatusAccepted, response)

}
