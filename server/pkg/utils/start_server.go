package utils

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"os/signal"
)

// StartServerWithGracefulShutdown  function for starting server with a graceful shutdown.
func StartServerWithGracefulShudown(a *fiber.App) {
	//Create channel for idle connections
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		//Recived an interrupt signal,shutdown
		if err := a.Shutdown(); err != nil {
			//Error from closing listeners, or context timeout:
			log.Printf("Oops..Server is not shutting down! Reason: %v", err)
		}
		close(idleConnsClosed)
	}()
	// Build Fiber connection URL.
	fiberConnURL, _ := ConnectionURLBuilder("fiber")

	// Run server.
	if err := a.Listen(fiberConnURL); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}

	<-idleConnsClosed
}

// StartServer func for starting a simple server.
func StartServer(a *fiber.App) {
	// Build Fiber connection URL.
	fiberConnURL, _ := ConnectionURLBuilder("fiber")
	fmt.Println(fiberConnURL)

	// Run server.
	if err := a.Listen(fiberConnURL); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}
}
