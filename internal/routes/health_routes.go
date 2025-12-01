package routes

import (
	"keerja-backend/internal/handler/http/health"

	"github.com/gofiber/fiber/v2"
)

// SetupHealthRoutes configures health check endpoints
// These endpoints are used for monitoring and orchestration (Kubernetes, Docker, etc.)
//
// Endpoints:
//   - GET /health         - Full health check with all components status
//   - GET /health/live    - Liveness probe (is the app running?)
//   - GET /health/ready   - Readiness probe (is the app ready to serve traffic?)
//   - GET /health/system  - System information (memory, goroutines, etc.)
func SetupHealthRoutes(app *fiber.App, handler *health.HealthHandler) {
	// Health check group - outside of /api/v1 for easier access
	healthGroup := app.Group("/health")

	// Full health check with component status
	// Returns 200 if healthy, 503 if any component is down
	healthGroup.Get("/", handler.Health)

	// Liveness probe - always returns 200 if app is running
	// Used by Kubernetes to know when to restart the container
	healthGroup.Get("/live", handler.Liveness)

	// Readiness probe - returns 200 if ready, 503 if not ready
	// Used by Kubernetes to know when to send traffic
	healthGroup.Get("/ready", handler.Readiness)

	// System info - memory, goroutines, GC stats
	// Useful for debugging and monitoring
	healthGroup.Get("/system", handler.SystemInfo)
}
