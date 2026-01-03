package main

import (
	"log"

	"task-management-platform/backend/internal/config"
	"task-management-platform/backend/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	r := server.New(cfg)

	log.Printf("api starting on port %s", cfg.Port)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
