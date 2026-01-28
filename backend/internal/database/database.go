package database

import (
	"log"
	"reservation-api/internal/config"
	"reservation-api/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB initializes database connection
func InitDB(cfg *config.Config) *gorm.DB {
	// Configure GORM logger
	logLevel := logger.Silent
	if cfg.IsDevelopment() {
		logLevel = logger.Info
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	log.Println("✅ Database connected successfully")

	// Auto migrate all models
	if err := autoMigrate(db); err != nil {
		log.Fatal("❌ Failed to migrate database:", err)
	}

	return db
}

// autoMigrate runs database migrations
func autoMigrate(db *gorm.DB) error {
	log.Println("🔄 Running database migrations...")

	err := db.AutoMigrate(
		&models.User{},
		&models.Court{},
		&models.Timeslot{},
		&models.Reservation{},
		&models.Payment{},
	)

	if err != nil {
		return err
	}

	log.Println("✅ Database migrations completed")
	return nil
}

// GetDB returns database instance (for testing purposes)
func GetDB(cfg *config.Config) *gorm.DB {
	return InitDB(cfg)
}