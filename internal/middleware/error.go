package middleware

import (
	"errors"
	"fmt"

	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

// ErrorHandler is a global error handler for the application
func ErrorHandler(isDevelopment bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Process request
		err := c.Next()

		// If no error, return
		if err == nil {
			return nil
		}

		// Log the error
		log.Errorf("Error occurred: %v", err)

		// Handle Fiber errors
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			return handleFiberError(c, fiberErr, isDevelopment)
		}

		// Handle GORM errors
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Resource not found", "The requested resource does not exist")
		}

		if errors.Is(err, gorm.ErrInvalidData) {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid data", err.Error())
		}

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return utils.ErrorResponse(c, fiber.StatusConflict, "Duplicate entry", "A record with this data already exists")
		}

		// Handle JWT errors
		if errors.Is(err, utils.ErrInvalidToken) {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token", "Your authentication token is invalid")
		}

		if errors.Is(err, utils.ErrExpiredToken) {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Token expired", "Your authentication token has expired")
		}

		// Default internal server error
		return handleInternalError(c, err, isDevelopment)
	}
}

// handleFiberError handles Fiber-specific errors
func handleFiberError(c *fiber.Ctx, err *fiber.Error, isDevelopment bool) error {
	code := err.Code
	message := err.Message

	// Map common Fiber errors to user-friendly messages
	switch code {
	case fiber.StatusNotFound:
		message = "The requested endpoint does not exist"
	case fiber.StatusMethodNotAllowed:
		message = "HTTP method not allowed for this endpoint"
	case fiber.StatusRequestEntityTooLarge:
		message = "Request body too large"
	case fiber.StatusUnsupportedMediaType:
		message = "Unsupported media type"
	}

	response := fiber.Map{
		"success": false,
		"message": message,
		"code":    code,
	}

	// Add error details in development
	if isDevelopment {
		response["error"] = err.Error()
	}

	return c.Status(code).JSON(response)
}

// handleInternalError handles internal server errors
func handleInternalError(c *fiber.Ctx, err error, isDevelopment bool) error {
	response := fiber.Map{
		"success": false,
		"message": "An internal server error occurred",
		"code":    fiber.StatusInternalServerError,
	}

	// Add error details in development
	if isDevelopment {
		response["error"] = err.Error()
		response["path"] = c.Path()
		response["method"] = c.Method()
	}

	return c.Status(fiber.StatusInternalServerError).JSON(response)
}

// RecoverPanic recovers from panics and returns 500 error
func RecoverPanic(isDevelopment bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				// Log the panic
				log.Errorf("❗ PANIC RECOVERED: %v", r)
				log.Errorf("  Method: %s", c.Method())
				log.Errorf("  Path: %s", c.Path())
				log.Errorf("  IP: %s", c.IP())
				log.Errorf("  UserID: %d", GetUserID(c))

				// Build error response
				response := fiber.Map{
					"success": false,
					"message": "Internal server error",
					"code":    fiber.StatusInternalServerError,
				}

				// Add panic details in development
				if isDevelopment {
					response["panic"] = fmt.Sprintf("%v", r)
					response["path"] = c.Path()
					response["method"] = c.Method()
				}

				// Send error response
				c.Status(fiber.StatusInternalServerError).JSON(response)
			}
		}()

		return c.Next()
	}
}

// NotFoundHandler handles 404 errors
func NotFoundHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return utils.ErrorResponse(
			c,
			fiber.StatusNotFound,
			"Endpoint not found",
			fmt.Sprintf("The endpoint '%s %s' does not exist", c.Method(), c.Path()),
		)
	}
}

// MethodNotAllowedHandler handles 405 errors
func MethodNotAllowedHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return utils.ErrorResponse(
			c,
			fiber.StatusMethodNotAllowed,
			"Method not allowed",
			fmt.Sprintf("The HTTP method '%s' is not allowed for this endpoint", c.Method()),
		)
	}
}

// ValidationErrorHandler handles validation errors
func ValidationErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		if err != nil {
			// Check if it's a validation error
			validationErrors := utils.FormatValidationErrors(err)
			if len(validationErrors) > 0 {
				return utils.ValidationErrorResponse(c, "Validation failed", validationErrors)
			}
		}

		return err
	}
}

// DatabaseErrorHandler handles database errors
func DatabaseErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		if err != nil {
			// Handle specific database errors
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return utils.ErrorResponse(c, fiber.StatusNotFound, "Record not found", "The requested record does not exist")
			}

			if errors.Is(err, gorm.ErrInvalidData) {
				return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid data", "The provided data is invalid")
			}

			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return utils.ErrorResponse(c, fiber.StatusConflict, "Duplicate entry", "A record with this data already exists")
			}

			// Check for foreign key constraint violations
			if isForeignKeyError(err) {
				return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid reference", "Referenced record does not exist")
			}

			// Check for unique constraint violations
			if isUniqueConstraintError(err) {
				return utils.ErrorResponse(c, fiber.StatusConflict, "Duplicate entry", "This value already exists and must be unique")
			}
		}

		return err
	}
}

// Helper functions to check error types

func isForeignKeyError(err error) bool {
	errMsg := err.Error()
	return contains(errMsg, "foreign key constraint") || contains(errMsg, "violates foreign key")
}

func isUniqueConstraintError(err error) bool {
	errMsg := err.Error()
	return contains(errMsg, "unique constraint") || contains(errMsg, "duplicate key")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TimeoutHandler handles request timeouts
func TimeoutHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		if err != nil && errors.Is(err, fiber.ErrRequestTimeout) {
			return utils.ErrorResponse(
				c,
				fiber.StatusRequestTimeout,
				"Request timeout",
				"The request took too long to process. Please try again.",
			)
		}

		return err
	}
}
