package main

import (
	"log"
	"users/infrastructure/server"
)

func main() {
	log.Println("Starting service")

	config, err := NewConfig()
	if err != nil {
		log.Fatalf("Configuration error: %s", err.Error())
	}

	app := server.Setup(config.Server)
	if err = app.ListenAndServe(); err != nil {
		log.Fatalf("Application error: %s", err.Error())
	}
}
