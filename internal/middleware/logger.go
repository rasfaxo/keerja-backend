package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// RequestLogger logs incoming requests and responses
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Start timer
		start := time.Now()

		// Get request info
		method := c.Method()
		path := c.Path()
		ip := c.IP()
		userAgent := c.Get("User-Agent")

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get response info
		statusCode := c.Response().StatusCode()

		// Get user info if authenticated
		userID := GetUserID(c)
		userType := GetUserType(c)

		// Build log message
		logMsg := fmt.Sprintf(
			"[%s] %s %s | Status: %d | Duration: %v | IP: %s | UserID: %d | UserType: %s",
			time.Now().Format("2006-01-02 15:04:05"),
			method,
			path,
			statusCode,
			duration,
			ip,
			userID,
			userType,
		)

		// Log based on status code
		if statusCode >= 500 {
			log.Errorf("%s | UserAgent: %s | Error: %v", logMsg, userAgent, err)
		} else if statusCode >= 400 {
			log.Warnf("%s | UserAgent: %s", logMsg, userAgent)
		} else {
			log.Infof("%s", logMsg)
		}

		return err
	}
}

// DetailedLogger provides more detailed logging including request/response bodies
func DetailedLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Log request
		log.Infof("📥 Incoming Request:")
		log.Infof("  Method: %s", c.Method())
		log.Infof("  Path: %s", c.Path())
		log.Infof("  IP: %s", c.IP())
		log.Infof("  Headers: %v", c.GetReqHeaders())

		// Log body for non-GET requests (be careful with sensitive data)
		if c.Method() != "GET" {
			log.Infof("  Body: %s", string(c.Body()))
		}

		// Process request
		err := c.Next()

		// Log response
		duration := time.Since(start)
		statusCode := c.Response().StatusCode()

		log.Infof("📤 Outgoing Response:")
		log.Infof("  Status: %d", statusCode)
		log.Infof("  Duration: %v", duration)
		log.Infof("  Body Size: %d bytes", len(c.Response().Body()))

		if err != nil {
			log.Errorf("  Error: %v", err)
		}

		return err
	}
}

// ErrorLogger logs errors with stack trace in development mode
func ErrorLogger(isDevelopment bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		if err != nil {
			// Log error with context
			log.Errorf("❌ Error occurred:")
			log.Errorf("  Method: %s", c.Method())
			log.Errorf("  Path: %s", c.Path())
			log.Errorf("  IP: %s", c.IP())
			log.Errorf("  UserID: %d", GetUserID(c))
			log.Errorf("  Error: %v", err)

			// In development, log more details
			if isDevelopment {
				log.Errorf("  Request Headers: %v", c.GetReqHeaders())
				log.Errorf("  Request Body: %s", string(c.Body()))
			}
		}

		return err
	}
}

// LoggerConfig configures the logging middleware
type LoggerConfig struct {
	// Skip logging for specific paths
	SkipPaths []string

	// Log request body (be careful with sensitive data)
	LogRequestBody bool

	// Log response body
	LogResponseBody bool

	// Custom format function
	CustomFormat func(c *fiber.Ctx, duration time.Duration) string
}

// CustomLogger creates a customizable logger middleware
func CustomLogger(config LoggerConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if path should be skipped
		path := c.Path()
		for _, skipPath := range config.SkipPaths {
			if path == skipPath {
				return c.Next()
			}
		}

		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		// Use custom format if provided
		if config.CustomFormat != nil {
			logMsg := config.CustomFormat(c, duration)
			log.Info(logMsg)
		} else {
			// Default format
			statusCode := c.Response().StatusCode()
			logMsg := fmt.Sprintf(
				"%s | %s %s | Status: %d | Duration: %v | IP: %s",
				time.Now().Format("2006-01-02 15:04:05"),
				c.Method(),
				path,
				statusCode,
				duration,
				c.IP(),
			)
			log.Info(logMsg)
		}

		return err
	}
}

// PerformanceLogger logs slow requests
func PerformanceLogger(threshold time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		// Log if request took longer than threshold
		if duration > threshold {
			log.Warnf(
				"🐌 Slow Request Detected: %s %s | Duration: %v | Threshold: %v | IP: %s",
				c.Method(),
				c.Path(),
				duration,
				threshold,
				c.IP(),
			)
		}

		return err
	}
}

// AccessLogger logs access attempts to protected resources
func AccessLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := GetUserID(c)
		userType := GetUserType(c)

		// Log access attempt
		log.Infof(
			"🔐 Access Attempt: User %d (%s) accessing %s %s from IP %s",
			userID,
			userType,
			c.Method(),
			c.Path(),
			c.IP(),
		)

		return c.Next()
	}
}
