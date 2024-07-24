package main

import (
	"log"

	"github.com/digisata/invitation-service/config"
	"github.com/digisata/invitation-service/internal/app"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("fatal error in config file: %s", err.Error())
	}

	app.Run(cfg)
}
