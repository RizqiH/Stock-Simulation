package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server Configuration
	Port    string
	GinMode string
	LogLevel string

	// Database Configuration
	DatabaseURL      string
	MySQLHost        string
	MySQLPort        string
	MySQLUser        string
	MySQLPassword    string
	MySQLDatabase    string
	MySQLRootPassword string

	// Redis Configuration
	RedisURL  string
	RedisPort string

	// JWT Configuration
	JWTSecret string

	// File Upload Configuration
	MaxUploadSize string
	UploadPath    string

	// Rate Limiting Configuration
	RateLimitRequests int
	RateLimitWindow   time.Duration

	// CORS Configuration
	CORSAllowedOrigins []string
	CORSAllowedMethods []string
	CORSAllowedHeaders []string

	// External APIs
	StockAPIKey string
	StockAPIURL string

	// Email Configuration
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		// Server Configuration
		Port:     getEnv("PORT", "8080"),
		GinMode:  getEnv("GIN_MODE", "debug"),
		LogLevel: getEnv("LOG_LEVEL", "info"),

		// Database Configuration
		DatabaseURL:       getEnvRequired("DATABASE_URL"),
		MySQLHost:         getEnv("MYSQL_HOST", "localhost"),
		MySQLPort:         getEnv("MYSQL_PORT", "3306"),
		MySQLUser:         getEnv("MYSQL_USER", "user"),
		MySQLPassword:     getEnvRequired("MYSQL_PASSWORD"),
		MySQLDatabase:     getEnv("MYSQL_DATABASE", "stock_simulation"),
		MySQLRootPassword: getEnvRequired("MYSQL_ROOT_PASSWORD"),

		// Redis Configuration
		RedisURL:  getEnv("REDIS_URL", "redis://localhost:6379"),
		RedisPort: getEnv("REDIS_PORT", "6379"),

		// JWT Configuration
		JWTSecret: getEnvRequired("JWT_SECRET"),

		// File Upload Configuration
		MaxUploadSize: getEnv("MAX_UPLOAD_SIZE", "10MB"),
		UploadPath:    getEnv("UPLOAD_PATH", "./uploads"),

		// Rate Limiting Configuration
		RateLimitRequests: getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow:   getEnvAsDuration("RATE_LIMIT_WINDOW", "1m"),

		// CORS Configuration
		CORSAllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
		CORSAllowedMethods: getEnvAsSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		CORSAllowedHeaders: getEnvAsSlice("CORS_ALLOWED_HEADERS", []string{"Content-Type", "Authorization"}),

		// External APIs
		StockAPIKey: getEnv("STOCK_API_KEY", ""),
		StockAPIURL: getEnv("STOCK_API_URL", ""),

		// Email Configuration
		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	if duration, err := time.ParseDuration(defaultValue); err == nil {
		return duration
	}
	return time.Minute
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}

// getEnvRequired returns the value of an environment variable or exits if not set
func getEnvRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Required environment variable %s is not set", key)
	}
	return value
}
