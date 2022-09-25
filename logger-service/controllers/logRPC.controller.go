package controllers

import (
	"context"
	"logger/dto"
	"logger/logs"
	"logger/services"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	logService services.LogService
}

func NewLogServer(logService services.LogService) LogServer {
	return LogServer{
		logService: logService,
	}
}

func (controller *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()
	logDTO := dto.LogDTO{
		Name: input.Name,
		Data: input.Data,
	}
	err := controller.logService.Insert(logDTO)
	if err != nil {
		res := &logs.LogResponse{Result: "failed"}
		return res, err
	}
	res := &logs.LogResponse{Result: "logged"}
	return res, nil
}
