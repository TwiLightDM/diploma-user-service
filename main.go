package main

import (
	"log"

	"github.com/TwiLightDM/diploma-user-service/internal/app"
	"github.com/TwiLightDM/diploma-user-service/internal/config"
)

func main() {
	cfg := config.Load()

	log.Printf("Starting user-service on %s", cfg.GRPCPort)

	if err := app.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
