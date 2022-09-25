package main

import (
	"context"
	"fmt"
	"log"
	"logger/config"
	"logger/controllers"
	"logger/repository"
	"logger/routes"
	"logger/services"
	"net/http"
	"time"
)

const webPort = "8002"

const dbTimeout = 15 * time.Second

func main() {
	log.Println("Starting logger service")

	mongoClient, errConnect := config.ConnectDB()
	if errConnect != nil {
		log.Panic(errConnect)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	defer func() {
		if errDisconnect := mongoClient.Disconnect(ctx); errDisconnect != nil {
			panic(errDisconnect)
		}
	}()

	var (
		//repository
		logRepository repository.LogRepository = repository.NewLogRepository(mongoClient, dbTimeout)
		//service
		logService services.LogService = services.NewLogService(logRepository)
		//controller
		logController controllers.LogController = controllers.NewLogController(logService)
		//gRPC
		logServer controllers.LogServer = controllers.NewLogServer(logService)
	)

	log.Printf("Starting logger service on port %s\n", webPort)

	go config.GRPCListen(logServer)

	server := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: routes.Routes(&controllers.Controller{
			LogController: logController,
		}),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}
