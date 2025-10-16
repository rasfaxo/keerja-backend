package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	// Server Configuration
	ServerPort string
	ServerHost string
	AppEnv     string
	AppName    string

	// Database Configuration
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// JWT Configuration
	JWTSecret                string
	JWTExpirationHours       int
	JWTRefreshExpirationDays int

	// Email Configuration
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string

	// Storage Configuration
	StorageProvider string // "local" or "s3" or "cloudinary"
	AWSRegion       string
	AWSBucket       string
	AWSAccessKey    string
	AWSSecretKey    string
	CloudinaryURL   string
	UploadPath      string

	// Redis Configuration (optional)
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	// Rate Limiting
	RateLimitEnabled bool
	RateLimitMax     int
	RateLimitWindow  time.Duration

	// CORS Configuration
	AllowedOrigins []string

	// Pagination
	DefaultPageSize int
	MaxPageSize     int

	// Frontend URLs
	FrontendURL      string
	VerifyEmailURL   string
	ResetPasswordURL string
	DashboardURL     string
	SupportEmail     string
}

var globalConfig *Config

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	config := &Config{
		// Server Configuration
		ServerPort: getEnv("SERVER_PORT", "8080"),
		ServerHost: getEnv("SERVER_HOST", "0.0.0.0"),
		AppEnv:     getEnv("APP_ENV", "development"),
		AppName:    getEnv("APP_NAME", "Keerja API"),

		// Database Configuration
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "bekerja"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "keerja"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		// JWT Configuration
		JWTSecret:                getEnv("JWT_SECRET", "your-secret-key-change-this"),
		JWTExpirationHours:       getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
		JWTRefreshExpirationDays: getEnvAsInt("JWT_REFRESH_EXPIRATION_DAYS", 7),

		// Email Configuration
		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("SMTP_FROM", "noreply@keerja.com"),

		// Storage Configuration
		StorageProvider: getEnv("STORAGE_PROVIDER", "local"),
		AWSRegion:       getEnv("AWS_REGION", ""),
		AWSBucket:       getEnv("AWS_BUCKET", ""),
		AWSAccessKey:    getEnv("AWS_ACCESS_KEY", ""),
		AWSSecretKey:    getEnv("AWS_SECRET_KEY", ""),
		CloudinaryURL:   getEnv("CLOUDINARY_URL", ""),
		UploadPath:      getEnv("UPLOAD_PATH", "./uploads"),

		// Redis Configuration
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),

		// Rate Limiting
		RateLimitEnabled: getEnvAsBool("RATE_LIMIT_ENABLED", true),
		RateLimitMax:     getEnvAsInt("RATE_LIMIT_MAX", 100),
		RateLimitWindow:  time.Duration(getEnvAsInt("RATE_LIMIT_WINDOW_SECONDS", 60)) * time.Second,

		// CORS Configuration
		AllowedOrigins: getEnvAsSlice("ALLOWED_ORIGINS", []string{"http://localhost:3000", "http://localhost:5173"}),

		// Pagination
		DefaultPageSize: getEnvAsInt("DEFAULT_PAGE_SIZE", 10),
		MaxPageSize:     getEnvAsInt("MAX_PAGE_SIZE", 100),

		// Frontend URLs
		FrontendURL:      getEnv("FRONTEND_URL", "http://localhost:3000"),
		VerifyEmailURL:   getEnv("VERIFY_EMAIL_URL", "http://localhost:3000/verify-email"),
		ResetPasswordURL: getEnv("RESET_PASSWORD_URL", "http://localhost:3000/reset-password"),
		DashboardURL:     getEnv("DASHBOARD_URL", "http://localhost:3000/dashboard"),
		SupportEmail:     getEnv("SUPPORT_EMAIL", "support@keerja.com"),
	}

	// Validate required configurations
	if err := config.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	globalConfig = config
	return config
}

// GetConfig returns the global configuration instance
func GetConfig() *Config {
	if globalConfig == nil {
		log.Fatal("Configuration not initialized. Call LoadConfig() first")
	}
	return globalConfig
}

// Validate checks if all required configurations are set
func (c *Config) Validate() error {
	if c.DBPassword == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}

	if c.JWTSecret == "your-secret-key-change-this" {
		log.Println("WARNING: Using default JWT secret. Please change it in production!")
	}

	if c.AppEnv == "production" {
		if c.JWTSecret == "your-secret-key-change-this" {
			return fmt.Errorf("JWT_SECRET must be changed in production")
		}
		if c.StorageProvider == "local" {
			log.Println("WARNING: Using local storage in production. Consider using S3 or Cloudinary")
		}
	}

	return nil
}

// IsDevelopment returns true if the app is running in development mode
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

// IsProduction returns true if the app is running in production mode
func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort, c.DBSSLMode,
	)
}

// Helper functions

func getEnv(key string, defaultValue string) string {
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

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	// Simple implementation - split by comma
	// kalo mau complex parsing, bisa pake proper CSV parser
	var result []string
	for _, v := range defaultValue {
		if valueStr != "" {
			result = append(result, valueStr)
			break
		}
		result = append(result, v)
	}
	return result
}
