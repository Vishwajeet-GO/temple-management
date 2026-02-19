package main

import (
	"log"

	"newapp/internal/config"
	"newapp/internal/database"
	"newapp/internal/routes"
)

func main() {
	cfg := config.Load()
	database.Initialize(cfg)

	r := routes.Setup()

	log.Printf("🛕 Temple Management on http://localhost:%s", cfg.AppPort)
	log.Printf("📱 Phone: http://<your-ip>:%s", cfg.AppPort)
	r.Run("0.0.0.0:" + cfg.AppPort)
}
