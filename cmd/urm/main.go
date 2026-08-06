package main

import (
	"log"

	"urm/internal/config"
	"urm/internal/server"
)

func startAllServices() {
	cfg := config.Load()

	srv := server.New(cfg.HTTP())
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	startAllServices()
}
