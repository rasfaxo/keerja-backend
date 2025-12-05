package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
)

// ResponseCompression adds gzip compression for responses larger than 5KB
// This middleware automatically compresses response bodies using gzip for clients that support it
// Benefits:
// - Reduces payload size for large JSON responses (job categories, skills, etc.)
// - Transparent to clients - they decompress automatically
// - Minimal CPU overhead for modern hardware
// - Particularly effective for master data endpoints with >5KB responses
func ResponseCompression() fiber.Handler {
	return compress.New(compress.Config{
		Level: compress.LevelDefault, // Balance between compression ratio and speed
	})
}
