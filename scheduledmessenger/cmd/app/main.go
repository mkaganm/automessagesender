package main

import (
	"github.com/joho/godotenv"
	"log"
	"scheduledmessenger/internal/container"
)

func main() {

	// Load env
	err := godotenv.Load()

	if err != nil {
		// fixme :
		log.Fatal("error when read env")
	}

	c := container.Create()
	c.Initialize()

	c.MessageService.Run()
}
