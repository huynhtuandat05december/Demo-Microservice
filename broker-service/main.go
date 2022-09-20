package main

import (
	"broken/routes"
	"fmt"
	"log"
	"net/http"
)

const webPort = "8000"

func main() {

	log.Printf("Starting broker service on port %s\n", webPort)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: routes.Routes(),
	}

	// start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
