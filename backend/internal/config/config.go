package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	// Database
	DatabaseURL string

	// Server
	Port   string
	AppEnv string

	// JWT
	JWTSecret string

	// Midtrans
	MidtransServerKey string
	MidtransClientKey string
	MidtransBaseURL   string

	// CORS
	AllowedOrigins []string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using system environment variables")
	}

	// Determine Midtrans environment
	appEnv := getEnv("APP_ENV", "development")
	midtransBaseURL := "https://app.sandbox.midtrans.com/snap/v1"
	if appEnv == "production" {
		midtransBaseURL = "https://app.midtrans.com/snap/v1"
	}

	config := &Config{
		// Database
		DatabaseURL: getEnv("DATABASE_URL", "host=localhost user=postgres password=postgres dbname=pilates_db port=5432 sslmode=disable"),

		// Server
		Port:   getEnv("PORT", "8080"),
		AppEnv: appEnv,

		// JWT
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key-change-this-in-production"),

		// Midtrans
		MidtransServerKey: getEnv("MIDTRANS_SERVER_KEY", ""),
		MidtransClientKey: getEnv("MIDTRANS_CLIENT_KEY", ""),
		MidtransBaseURL:   midtransBaseURL,

		// CORS
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://localhost:3001",
			getEnv("FRONTEND_URL", "http://localhost:3000"),
		},
	}

	// Validate required configs
	config.validate()

	return config
}

// validate checks if required configuration values are set
func (c *Config) validate() {
	if c.JWTSecret == "your-secret-key-change-this-in-production" && c.AppEnv == "production" {
		log.Fatal("❌ JWT_SECRET must be set in production environment")
	}

	if c.AppEnv == "production" {
		if c.MidtransServerKey == "" || c.MidtransClientKey == "" {
			log.Println("⚠️  Warning: Midtrans credentials not set. Payment features will not work.")
		}
	}
}

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// IsDevelopment checks if app is in development mode
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

// IsProduction checks if app is in production mode
func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}