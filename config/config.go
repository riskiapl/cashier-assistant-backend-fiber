package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// GetEnvironment returns the current environment (development, production, etc.)
func GetEnvironment() string {
	env := os.Getenv("GO_ENV")
	if env == "" {
		// Default to development if not specified
		env = "development"
	}
	return env
}

// LoadEnv loads the appropriate .env file based on the environment
func LoadEnv() {
	env := GetEnvironment()

	// Try to load environment-specific .env file first
	envFile := ".env." + env
	err := godotenv.Load(envFile)
	if err != nil {
		log.Printf("No %s file found, trying default .env\n", envFile)

		// Fall back to default .env
		err = godotenv.Load()
		if err != nil {
			log.Println("No .env file found, using system environment variables")
		} else {
			log.Println("Loaded default .env file")
		}
	} else {
		log.Printf("Loaded %s environment configuration\n", envFile)
	}
}

func GetEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// GetCorsOrigins returns the allowed CORS origins based on the environment
func GetCorsOrigins() string {
	env := GetEnvironment()

	// Get origins from environment variable
	origins := GetEnv("CORS_ORIGINS", "")

	// If not specified in environment, use defaults based on environment
	if origins == "" {
		if env == "production" {
			// In production, you might want to be more restrictive
			origins = GetEnv("FRONTEND_URL", "https://riskiapl.duckdns.org")
		} else {
			// In development, default to localhost:3000
			origins = "http://localhost:3000"
		}
	}

	return origins
}

// GetCorsConfig returns configuration for CORS based on environment
func GetCorsConfig() map[string]string {
	return map[string]string{
		"origins":     GetCorsOrigins(),
		"headers":     "Origin, Content-Type, Accept, Authorization",
		"methods":     "GET, POST, PUT, DELETE",
		"credentials": "true",
	}
}
