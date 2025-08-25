package main

import (
	"log"

	"github.com/krisn2/go-social/config"
	"github.com/krisn2/go-social/database"
	"github.com/krisn2/go-social/routes"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Setup routes
	router := routes.Setup(db, cfg)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
