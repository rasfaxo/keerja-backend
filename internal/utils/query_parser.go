package utils

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ParseBoolQuery parses boolean query parameter with various formats
// Accepts: true, false, 1, 0, yes, no, on, off (case insensitive)
func ParseBoolQuery(c *fiber.Ctx, key string, defaultValue bool) bool {
	value := strings.ToLower(strings.TrimSpace(c.Query(key)))

	if value == "" {
		return defaultValue
	}

	// True values
	if value == "true" || value == "1" || value == "yes" || value == "on" {
		return true
	}

	// False values
	if value == "false" || value == "0" || value == "no" || value == "off" {
		return false
	}

	// Default if invalid
	return defaultValue
}

// ParseIntQuery parses integer query parameter with validation
func ParseIntQuery(c *fiber.Ctx, key string, defaultValue int) int {
	value := strings.TrimSpace(c.Query(key))

	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

// ParseIntQueryWithRange parses integer query parameter with min/max validation
func ParseIntQueryWithRange(c *fiber.Ctx, key string, defaultValue, min, max int) int {
	value := ParseIntQuery(c, key, defaultValue)

	if value < min {
		return min
	}
	if value > max {
		return max
	}

	return value
}

// ParseInt64Query parses int64 query parameter
func ParseInt64Query(c *fiber.Ctx, key string, defaultValue int64) int64 {
	value := strings.TrimSpace(c.Query(key))

	if value == "" {
		return defaultValue
	}

	int64Value, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}

	return int64Value
}

// ParseFloatQuery parses float64 query parameter
func ParseFloatQuery(c *fiber.Ctx, key string, defaultValue float64) float64 {
	value := strings.TrimSpace(c.Query(key))

	if value == "" {
		return defaultValue
	}

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}

	return floatValue
}

// ParseStringQuery parses and sanitizes string query parameter
func ParseStringQuery(c *fiber.Ctx, key string, defaultValue string) string {
	value := strings.TrimSpace(c.Query(key))

	if value == "" {
		return defaultValue
	}

	return SanitizeString(value)
}

// ParseStringArrayQuery parses comma-separated query parameter into array
func ParseStringArrayQuery(c *fiber.Ctx, key string) []string {
	value := strings.TrimSpace(c.Query(key))

	if value == "" {
		return []string{}
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, SanitizeString(trimmed))
		}
	}

	return result
}

// QueryParams holds common query parameters
type QueryParams struct {
	Page     int
	Limit    int
	Sort     string
	Order    string
	Search   string
	Status   string
	FromDate string
	ToDate   string
}

// ParseCommonQueryParams parses common pagination and filtering query parameters
func ParseCommonQueryParams(c *fiber.Ctx) QueryParams {
	return QueryParams{
		Page:     ParseIntQueryWithRange(c, "page", 1, 1, 10000),
		Limit:    ParseIntQueryWithRange(c, "limit", 10, 1, 100),
		Sort:     ParseStringQuery(c, "sort", "created_at"),
		Order:    ParseStringQuery(c, "order", "desc"),
		Search:   ParseStringQuery(c, "search", ""),
		Status:   ParseStringQuery(c, "status", ""),
		FromDate: ParseStringQuery(c, "from_date", ""),
		ToDate:   ParseStringQuery(c, "to_date", ""),
	}
}
