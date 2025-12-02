package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ParseIDParam parses a path parameter to int64 (e.g. c.Params("id"))
func ParseIDParam(c *fiber.Ctx, name string) (int64, error) {
	return strconv.ParseInt(c.Params(name), 10, 64)
}

// SanitizePtr sanitizes a *string using utils.SanitizeString, preserving nil
func SanitizePtr(s *string) *string {
	if s == nil {
		return nil
	}
	sanitized := SanitizeString(*s)
	return &sanitized
}

// SanitizeIfNonEmpty sanitizes a string only when non-empty
func SanitizeIfNonEmpty(s string) string {
	if s == "" {
		return s
	}
	return SanitizeString(s)
}
