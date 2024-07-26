package main

import (
	"log"

	"github.com/digisata/todo-service/config"
	"github.com/digisata/todo-service/internal/app"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("fatal error in config file: %s", err.Error())
	}

	app.Run(cfg)
}
