package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	CORS     CORSConfig
	Redis    RedisConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	URL      string
}

type ServerConfig struct {
	Port    string
	Host    string
	ENV     string
	GinMode string
}

type JWTConfig struct {
	Secret string
}

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

type RedisConfig struct {
	URL  string
	Port string
}

func LoadConfig() *Config {
	// Load .env file if it exists (for local development)
	_ = godotenv.Load()
	
	port, _ := strconv.Atoi(getEnv("DB_PORT", "3307"))
	
	// Parse CORS origins
	corsOrigins := []string{
		"http://localhost:3000",
		"http://127.0.0.1:3000",
	}
	
	if envOrigins := getEnv("CORS_ORIGINS", ""); envOrigins != "" {
		// Split by comma and add to origins
		origins := strings.Split(envOrigins, ",")
		for _, origin := range origins {
			corsOrigins = append(corsOrigins, strings.TrimSpace(origin))
		}
	}
	
	// Set GinMode based on environment
	ginMode := "debug"
	if getEnv("ENV", "development") == "production" {
		ginMode = "release"
	}
	
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     port,
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "root"),
			DBName:   getEnv("DB_NAME", "stock_simulation"),
			URL:      getEnv("DATABASE_URL", ""),
		},
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			Host:    getEnv("HOST", "0.0.0.0"),
			ENV:     getEnv("ENV", "development"),
			GinMode: ginMode,
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-secret-key-change-this-in-production"),
		},
		CORS: CORSConfig{
			AllowedOrigins: corsOrigins,
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Content-Type", "Authorization", "X-Requested-With"},
		},
		Redis: RedisConfig{
			URL:  getEnv("REDIS_URL", "redis://localhost:6379"),
			Port: getEnv("REDIS_PORT", "6379"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue string) time.Duration {
	valueStr := getEnv(key, defaultValue)
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	// If parsing fails, parse the default value
	if defaultDuration, err := time.ParseDuration(defaultValue); err == nil {
		return defaultDuration
	}
	return time.Minute // fallback
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	// For Railway production, try internal connection first
	if c.IsProduction() {
		// If we have DATABASE_URL, try to use it but handle Railway-specific issues
		if c.Database.URL != "" {
			// For Railway internal connections, don't add TLS
			if strings.Contains(c.Database.URL, "railway.internal") {
				// Clean internal URL - remove TLS for internal connections
				cleanURL := strings.Replace(c.Database.URL, "&tls=true", "", -1)
				cleanURL = strings.Replace(cleanURL, "?tls=true", "", -1)
				
				// Add required parameters for Railway internal
				if strings.Contains(cleanURL, "?") {
					return cleanURL + "&charset=utf8mb4&parseTime=True&loc=Local"
				} else {
					return cleanURL + "?charset=utf8mb4&parseTime=True&loc=Local"
				}
			}
			
			// For external Railway connections, add TLS
			if strings.Contains(c.Database.URL, "proxy.rlwy.net") || strings.Contains(c.Database.URL, "railway.app") {
				if !strings.Contains(c.Database.URL, "tls=") {
					if strings.Contains(c.Database.URL, "?") {
						return c.Database.URL + "&tls=true&charset=utf8mb4&parseTime=True&loc=Local"
					} else {
						return c.Database.URL + "?tls=true&charset=utf8mb4&parseTime=True&loc=Local"
					}
				}
			}
			
			// Default URL handling
			if !strings.Contains(c.Database.URL, "charset=") {
				if strings.Contains(c.Database.URL, "?") {
					return c.Database.URL + "&charset=utf8mb4&parseTime=True&loc=Local"
				} else {
					return c.Database.URL + "?charset=utf8mb4&parseTime=True&loc=Local"
				}
			}
			return c.Database.URL
		}
		
		// If individual variables are set, use them
		if c.Database.Host != "" {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				c.Database.User,
				c.Database.Password,
				c.Database.Host,
				c.Database.Port,
				c.Database.DBName,
			)
			
			// Add TLS only for external Railway connections
			if strings.Contains(c.Database.Host, "proxy.rlwy.net") || strings.Contains(c.Database.Host, "railway.app") {
				dsn += "&tls=true"
			}
			// Do NOT add TLS for .railway.internal connections
			
			return dsn
		}
	}
	
	// Build DSN manually for local development
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
	)
	
	return dsn
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return c.Server.ENV == "production"
}

// IsDevelopment returns true if running in development environment
func (c *Config) IsDevelopment() bool {
	return c.Server.ENV == "development"
}

// GetServerAddress returns the full server address
func (c *Config) GetServerAddress() string {
	return c.Server.Host + ":" + c.Server.Port
}
