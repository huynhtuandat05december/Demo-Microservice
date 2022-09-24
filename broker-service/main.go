package main

import (
	"broken/config"
	"broken/controllers"
	"broken/routes"
	"fmt"
	"log"
	"net/http"
	"os"
)

const webPort = "8000"

func main() {
	rabbitConn, err := config.ConnectRabbitMQ()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	log.Printf("Starting broker service on port %s\n", webPort)

	//controller
	requestController := controllers.NewRequestController(rabbitConn)

	// define http server
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: routes.Routes(&controllers.Controller{
			RequestController: requestController,
		}),
	}

	// start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
