package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
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
	URL    string
	Port   string
	Client *redis.Client
}

func LoadConfig() *Config {
	// Load .env file if it exists (for local development)
	_ = godotenv.Load()

	port, _ := strconv.Atoi(getEnv("DB_PORT", "3307"))

	// Parse CORS origins
	corsOrigins := []string{
		"http://localhost:3000",
		"http://127.0.0.1:3000",
		"https://stock-simulation-frontend.vercel.app",
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

	// Initialize Redis client
	redisURL := getEnv("REDIS_URL", "redis://localhost:6379")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     strings.TrimPrefix(redisURL, "redis://"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Test Redis connection
	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("⚠️ Redis connection failed: %v\n", err)
		fmt.Println("📝 Redis features will be disabled")
		redisClient = nil
	} else {
		fmt.Println("✅ Redis connected successfully")
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
			URL:    redisURL,
			Port:   getEnv("REDIS_PORT", "6379"),
			Client: redisClient,
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	// For Railway production, try internal connection first
	if c.IsProduction() {
		// If we have DATABASE_URL, convert mysql:// format to Go MySQL driver format
		if c.Database.URL != "" {
			// Convert mysql://user:password@host:port/database format to Go MySQL driver format
			if strings.HasPrefix(c.Database.URL, "mysql://") {
				// Parse the mysql:// URL and convert to Go format
				url := strings.TrimPrefix(c.Database.URL, "mysql://")

				// Split into credentials@host/database
				parts := strings.Split(url, "@")
				if len(parts) == 2 {
					credentials := parts[0]
					hostAndDB := parts[1]

					// Split host:port/database
					hostParts := strings.Split(hostAndDB, "/")
					if len(hostParts) == 2 {
						hostPort := hostParts[0]
						database := hostParts[1]

						// Build Go MySQL driver DSN format: user:password@tcp(host:port)/database
						dsn := fmt.Sprintf("%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", credentials, hostPort, database)

						// Add TLS configuration based on host
						if strings.Contains(hostPort, "railway.internal") {
							dsn += "&tls=false"
						} else if strings.Contains(hostPort, "proxy.rlwy.net") {
							dsn += "&tls=false"
						} else if strings.Contains(hostPort, "railway.app") {
							dsn += "&tls=true"
						}

						return dsn
					}
				}
			}

			// For non-mysql:// URLs, handle as before
			cleanURL := c.Database.URL

			// For Railway internal connections, don't add TLS
			if strings.Contains(c.Database.URL, "railway.internal") {
				// Clean internal URL - remove TLS for internal connections
				cleanURL = strings.Replace(c.Database.URL, "&tls=true", "", -1)
				cleanURL = strings.Replace(cleanURL, "?tls=true", "", -1)

				// Add required parameters for Railway internal
				if strings.Contains(cleanURL, "?") {
					return cleanURL + "&charset=utf8mb4&parseTime=True&loc=Local"
				} else {
					return cleanURL + "?charset=utf8mb4&parseTime=True&loc=Local"
				}
			}

			// For external Railway connections, disable TLS for proxy
			if strings.Contains(c.Database.URL, "proxy.rlwy.net") {
				if !strings.Contains(c.Database.URL, "tls=") {
					if strings.Contains(c.Database.URL, "?") {
						return c.Database.URL + "&tls=false&charset=utf8mb4&parseTime=True&loc=Local"
					} else {
						return c.Database.URL + "?tls=false&charset=utf8mb4&parseTime=True&loc=Local"
					}
				}
			}

			// For railway.app connections, add TLS
			if strings.Contains(c.Database.URL, "railway.app") {
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

			// Add TLS configuration for Railway connections
			if strings.Contains(c.Database.Host, "proxy.rlwy.net") {
				dsn += "&tls=false"
			} else if strings.Contains(c.Database.Host, "railway.app") {
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
