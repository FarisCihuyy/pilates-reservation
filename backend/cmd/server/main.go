package main

import (
	"log"
	"reservation-api/api/routes"
	"reservation-api/internal/config"
	"reservation-api/internal/database"

	"github.com/gin-gonic/gin"
)

// @title Pilates Reservation API
// @version 1.0
// @description API untuk sistem reservasi Pilates dengan payment gateway
// @contact.name API Support
// @contact.email support@pilates.com
// @host localhost:8080
// @BasePath /api
// @schemes http https

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db := database.InitDB(cfg)

	// Run database seeding
	database.SeedData(db)

	// Setup Gin mode
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Setup router
	router := gin.Default()

	// Initialize routes
	routes.SetupRoutes(router, db, cfg)

	// Start server
	log.Printf("🚀 Server running on port %s", cfg.Port)
	log.Printf("📝 API Documentation: http://localhost:%s/api/docs", cfg.Port)
	log.Printf("🌍 Environment: %s", cfg.AppEnv)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}