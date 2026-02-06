package main

import (
	"log"
	"newapp/internal/config"
	"newapp/internal/database"
	"newapp/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	database.InitDatabase(cfg)

	// Set Gin mode
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r)

	// Start server
	log.Printf("Server starting on port %s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
