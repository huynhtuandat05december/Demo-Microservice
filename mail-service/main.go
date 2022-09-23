package main

import (
	"fmt"
	"log"
	"mail/config"
	"mail/controllers"
	"mail/routes"
	"net/http"
)

const webPort = "8003"

func main() {

	log.Println("Starting mail service on port", webPort)

	mailer := config.CreateMail()

	//controller
	mailController := controllers.NewMailController(mailer)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: routes.Routes(&controllers.Controller{
			MailController: mailController,
		}),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
