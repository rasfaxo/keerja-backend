package middleware

import (
	"fmt"
	"time"

	"keerja-backend/internal/config"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RateLimiter creates a rate limiter middleware
func RateLimiter(cfg *config.Config) fiber.Handler {
	// Skip if rate limiting is disabled
	if !cfg.RateLimitEnabled {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	// Default values
	max := cfg.RateLimitMax
	if max == 0 {
		max = 100 // Default: 100 requests
	}

	window := cfg.RateLimitWindow
	if window == 0 {
		window = 1 * time.Minute // Default: per minute
	}

	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: window,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use IP address as key
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return utils.ErrorResponse(
				c,
				fiber.StatusTooManyRequests,
				"Rate limit exceeded",
				fmt.Sprintf("Too many requests. Please try again later. Limit: %d requests per %v", max, window),
			)
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
		LimiterMiddleware:      limiter.SlidingWindow{},
	})
}

// CustomRateLimiter creates a custom rate limiter with specific limits
type RateLimiterConfig struct {
	Max            int           // Maximum number of requests
	Window         time.Duration // Time window
	Message        string        // Custom error message
	KeyGenerator   func(*fiber.Ctx) string
	SkipSuccessful bool // Skip counting successful requests
	SkipFailed     bool // Skip counting failed requests
}

// NewCustomRateLimiter creates a rate limiter with custom config
func NewCustomRateLimiter(config RateLimiterConfig) fiber.Handler {
	// Set defaults
	if config.Max == 0 {
		config.Max = 100
	}

	if config.Window == 0 {
		config.Window = 1 * time.Minute
	}

	if config.KeyGenerator == nil {
		config.KeyGenerator = func(c *fiber.Ctx) string {
			return c.IP()
		}
	}

	if config.Message == "" {
		config.Message = "Rate limit exceeded. Please try again later."
	}

	return limiter.New(limiter.Config{
		Max:          config.Max,
		Expiration:   config.Window,
		KeyGenerator: config.KeyGenerator,
		LimitReached: func(c *fiber.Ctx) error {
			return utils.ErrorResponse(
				c,
				fiber.StatusTooManyRequests,
				"Rate limit exceeded",
				config.Message,
			)
		},
		SkipSuccessfulRequests: config.SkipSuccessful,
		SkipFailedRequests:     config.SkipFailed,
		LimiterMiddleware:      limiter.SlidingWindow{},
	})
}

// RateLimitByIP limits requests by IP address
func RateLimitByIP(max int, window time.Duration) fiber.Handler {
	return NewCustomRateLimiter(RateLimiterConfig{
		Max:    max,
		Window: window,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		Message: fmt.Sprintf("Too many requests from your IP. Limit: %d per %v", max, window),
	})
}

// RateLimitByUser limits requests by authenticated user
func RateLimitByUser(max int, window time.Duration) fiber.Handler {
	return NewCustomRateLimiter(RateLimiterConfig{
		Max:    max,
		Window: window,
		KeyGenerator: func(c *fiber.Ctx) string {
			userID := GetUserID(c)
			if userID > 0 {
				return fmt.Sprintf("user:%d", userID)
			}
			// Fallback to IP for unauthenticated users
			return c.IP()
		},
		Message: fmt.Sprintf("Too many requests. Limit: %d per %v", max, window),
	})
}

// AuthRateLimiter for login/authentication endpoints (stricter)
func AuthRateLimiter() fiber.Handler {
	return NewCustomRateLimiter(RateLimiterConfig{
		Max:    5,                // 5 attempts
		Window: 15 * time.Minute, // per 15 minutes
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		Message: "Too many login attempts. Please try again in 15 minutes.",
	})
}

// APIRateLimiter for general API endpoints
func APIRateLimiter() fiber.Handler {
	return NewCustomRateLimiter(RateLimiterConfig{
		Max:    100,             // 100 requests
		Window: 1 * time.Minute, // per minute
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		Message: "API rate limit exceeded. Please slow down.",
	})
}

// UploadRateLimiter for file upload endpoints
func UploadRateLimiter() fiber.Handler {
	return NewCustomRateLimiter(RateLimiterConfig{
		Max:    10,              // 10 uploads
		Window: 1 * time.Minute, // per minute
		KeyGenerator: func(c *fiber.Ctx) string {
			userID := GetUserID(c)
			if userID > 0 {
				return fmt.Sprintf("upload:user:%d", userID)
			}
			return fmt.Sprintf("upload:ip:%s", c.IP())
		},
		Message: "Too many upload requests. Please wait before uploading more files.",
	})
}

// SearchRateLimiter for search endpoints
func SearchRateLimiter() fiber.Handler {
	return NewCustomRateLimiter(RateLimiterConfig{
		Max:    30,              // 30 searches
		Window: 1 * time.Minute, // per minute
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		Message: "Too many search requests. Please slow down.",
	})
}

// EmailRateLimiter for email sending endpoints
func EmailRateLimiter() fiber.Handler {
	return NewCustomRateLimiter(RateLimiterConfig{
		Max:    3,             // 3 emails
		Window: 1 * time.Hour, // per hour
		KeyGenerator: func(c *fiber.Ctx) string {
			userID := GetUserID(c)
			if userID > 0 {
				return fmt.Sprintf("email:user:%d", userID)
			}
			return fmt.Sprintf("email:ip:%s", c.IP())
		},
		Message: "Too many email requests. Please try again later.",
	})
}

// ApplicationRateLimiter for job application endpoints
func ApplicationRateLimiter() fiber.Handler {
	return NewCustomRateLimiter(RateLimiterConfig{
		Max:    10,            // 10 applications
		Window: 1 * time.Hour, // per hour
		KeyGenerator: func(c *fiber.Ctx) string {
			userID := GetUserID(c)
			if userID > 0 {
				return fmt.Sprintf("apply:user:%d", userID)
			}
			return fmt.Sprintf("apply:ip:%s", c.IP())
		},
		Message: "Too many job applications. Please wait before applying to more jobs.",
	})
}

// RegistrationRateLimiter for user registration
func RegistrationRateLimiter() fiber.Handler {
	return NewCustomRateLimiter(RateLimiterConfig{
		Max:    3,             // 3 registrations
		Window: 1 * time.Hour, // per hour
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		Message: "Too many registration attempts. Please try again later.",
	})
}
