package middleware

import (
	"strings"

	"keerja-backend/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORSConfig returns CORS middleware with configuration
func CORSConfig(cfg *config.Config) fiber.Handler {
	// Determine allowed origins
	allowedOrigins := cfg.AllowedOrigins
	if len(allowedOrigins) == 0 {
		// Default to localhost in development
		if cfg.AppEnv == "development" {
			allowedOrigins = []string{
				"http://localhost:3000",
				"http://localhost:3001",
				"http://localhost:5173", // Vite default
				"http://127.0.0.1:3000",
			}
		} else {
			// In production, must specify allowed origins
			allowedOrigins = []string{"https://yourdomain.com"}
		}
	}

	return cors.New(cors.Config{
		AllowOrigins: strings.Join(allowedOrigins, ","),
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPut,
			fiber.MethodPatch,
			fiber.MethodDelete,
			fiber.MethodOptions,
		}, ","),
		AllowHeaders: strings.Join([]string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
			"X-CSRF-Token",
		}, ","),
		AllowCredentials: true,
		ExposeHeaders: strings.Join([]string{
			"Content-Length",
			"Content-Type",
			"Authorization",
		}, ", "),
		MaxAge: 3600, // 1 hour
	})
}

// CustomCORS creates a custom CORS middleware with more control
type CustomCORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// NewCustomCORS creates a custom CORS middleware
func NewCustomCORS(config CustomCORSConfig) fiber.Handler {
	// Set defaults
	if len(config.AllowedMethods) == 0 {
		config.AllowedMethods = []string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPut,
			fiber.MethodPatch,
			fiber.MethodDelete,
			fiber.MethodOptions,
		}
	}

	if len(config.AllowedHeaders) == 0 {
		config.AllowedHeaders = []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		}
	}

	if config.MaxAge == 0 {
		config.MaxAge = 3600
	}

	return cors.New(cors.Config{
		AllowOrigins:     strings.Join(config.AllowedOrigins, ", "),
		AllowMethods:     strings.Join(config.AllowedMethods, ", "),
		AllowHeaders:     strings.Join(config.AllowedHeaders, ", "),
		ExposeHeaders:    strings.Join(config.ExposedHeaders, ", "),
		AllowCredentials: config.AllowCredentials,
		MaxAge:           config.MaxAge,
	})
}

// DynamicCORS validates origins dynamically
func DynamicCORS(allowedOrigins []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		origin := c.Get("Origin")

		// Check if origin is allowed
		if isOriginAllowed(origin, allowedOrigins) {
			c.Set("Access-Control-Allow-Origin", origin)
			c.Set("Access-Control-Allow-Credentials", "true")
			c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
			c.Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, Authorization")
			c.Set("Access-Control-Max-Age", "3600")
		}

		// Handle preflight requests
		if c.Method() == fiber.MethodOptions {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	}
}

// isOriginAllowed checks if origin is in allowed list
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	// Allow all origins if list is empty (development only)
	if len(allowedOrigins) == 0 {
		return true
	}

	// Check if origin matches any allowed origin
	for _, allowed := range allowedOrigins {
		// Support wildcard (*) in allowed origins
		if allowed == "*" {
			return true
		}

		// Exact match
		if origin == allowed {
			return true
		}

		// Support subdomain wildcard (*.example.com)
		if strings.HasPrefix(allowed, "*.") {
			domain := strings.TrimPrefix(allowed, "*.")
			if strings.HasSuffix(origin, domain) {
				return true
			}
		}
	}

	return false
}

// SecurityHeaders adds security headers
func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Prevent clickjacking
		c.Set("X-Frame-Options", "DENY")

		// Prevent MIME type sniffing
		c.Set("X-Content-Type-Options", "nosniff")

		// Enable XSS protection
		c.Set("X-XSS-Protection", "1; mode=block")

		// Enforce HTTPS
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Referrer policy
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy (adjust as needed)
		c.Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")

		// Permissions policy (formerly Feature-Policy)
		c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		return c.Next()
	}
}
