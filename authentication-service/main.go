package main

import (
	"authentication/config"
	"authentication/controllers"
	"authentication/repository"
	"authentication/routes"
	"authentication/services"
	"fmt"
	"log"
	"net/http"
	"time"
)

const webPort = "8001"

const dbTimeout = time.Second * 3

func main() {

	log.Println("Starting authentication service")

	conn := config.ConnectDB()

	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	var (
		//repository
		userRepository repository.UserRepository = repository.NewUserRepository(conn, dbTimeout)
		//service
		authService services.AuthService = services.NewAuthService(userRepository)
		userService services.UserService = services.NewUserService(userRepository)
		//controller
		authController controllers.AuthController = controllers.NewAuthController(userService, authService)
	)

	log.Printf("Starting authentication service on port %s\n", webPort)

	server := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: routes.Routes(&controllers.Controller{
			AuthController: authController,
		}),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}
