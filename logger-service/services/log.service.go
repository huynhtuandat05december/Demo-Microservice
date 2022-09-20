package services

import (
	"logger/dto"
	"logger/models"
	"logger/repository"
	"time"
)

type LogService interface {
	Insert(logDTO dto.LogDTO) error
}

type logService struct {
	logRepository repository.LogRepository
}

func NewLogService(logRepository repository.LogRepository) LogService {
	return &logService{
		logRepository: logRepository,
	}
}

func (service *logService) Insert(logDTO dto.LogDTO) error {
	logPayload := models.LogEntry{
		Name:      logDTO.Name,
		Data:      logDTO.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := service.logRepository.Insert(logPayload)
	if err != nil {
		return err
	}

	return nil

}
