package utils

import (
	"github.com/gofiber/fiber/v2"
)

// Response represents a standard API response
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func SuccessResponse(c *fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessResponseWithMeta(c *fiber.Ctx, message string, data any, meta any) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func CreatedResponse(c *fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *fiber.Ctx, statusCode int, message string, details ...string) error {
	resp := Response{
		Success: false,
		Message: message,
	}

	// Add error details if provided
	if len(details) > 0 {
		resp.Errors = details[0]
	}

	return c.Status(statusCode).JSON(resp)
}

func ErrorResponseWithErrors(c *fiber.Ctx, statusCode int, message string, errors any) error {
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

func ValidationErrorResponse(c *fiber.Ctx, message string, errors any) error {
	if message == "" {
		message = "Validation failed"
	}
	return c.Status(fiber.StatusBadRequest).JSON(Response{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

func UnauthorizedResponse(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Unauthorized"
	}
	return ErrorResponse(c, fiber.StatusUnauthorized, message)
}

func ForbiddenResponse(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Forbidden"
	}
	return ErrorResponse(c, fiber.StatusForbidden, message)
}

func NotFoundResponse(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Resource not found"
	}
	return ErrorResponse(c, fiber.StatusNotFound, message)
}

func ConflictResponse(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Resource already exists"
	}
	return ErrorResponse(c, fiber.StatusConflict, message)
}

func InternalServerErrorResponse(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Internal server error"
	}
	return ErrorResponse(c, fiber.StatusInternalServerError, message)
}

func BadRequestResponse(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Bad request"
	}
	return ErrorResponse(c, fiber.StatusBadRequest, message)
}

func NoContentResponse(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func PaginatedResponse(c *fiber.Ctx, data any, page, limit int, totalRows int64) error {
	meta := GetPaginationMeta(page, limit, totalRows)
	return SuccessResponseWithMeta(c, "Data retrieved successfully", data, meta)
}
