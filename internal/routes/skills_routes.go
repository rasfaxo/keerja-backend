package routes

import (
	"github.com/gofiber/fiber/v2"

	"keerja-backend/internal/handler/http"
)

// SetupSkillsRoutes sets up routes for skills master data
func SetupSkillsRoutes(api fiber.Router, handler *http.SkillsMasterHandler) {
	skills := api.Group("/skills")

	// Public endpoints - no authentication required for mobile app to fetch skills
	skills.Get("/", handler.GetAllSkills)              // GET /api/v1/skills - Get all skills with filters
	skills.Get("/search", handler.SearchSkills)        // GET /api/v1/skills/search?q=java - Search skills
	skills.Get("/type/:type", handler.GetSkillsByType) // GET /api/v1/skills/type/technical - Get by type
}
