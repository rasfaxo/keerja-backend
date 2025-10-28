package routes

import (
	"keerja-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupUserRoutes configures user routes
// Routes: /api/v1/users/*
func SetupUserRoutes(api fiber.Router, deps *Dependencies, authMw *middleware.AuthMiddleware) {
	users := api.Group("/users")

	// Protected routes - require authentication
	users.Use(authMw.AuthRequired())

	// Profile routes
	users.Get("/me", deps.UserHandler.GetProfile) // Support ?include=all or ?include=educations,skills
	users.Put("/me", deps.UserHandler.UpdateProfile)

	// Dedicated section GET routes (for granular access)
	users.Get("/me/educations", deps.UserHandler.GetEducations)
	users.Get("/me/experiences", deps.UserHandler.GetExperiences)
	users.Get("/me/skills", deps.UserHandler.GetSkills)
	users.Get("/me/certifications", deps.UserHandler.GetCertifications)
	users.Get("/me/languages", deps.UserHandler.GetLanguages)
	users.Get("/me/projects", deps.UserHandler.GetProjects)
	users.Get("/me/documents", deps.UserHandler.GetDocuments)

	// Education routes
	users.Post("/me/education", deps.UserHandler.AddEducation)
	users.Put("/me/education/:id", deps.UserHandler.UpdateEducation)
	users.Delete("/me/education/:id", deps.UserHandler.DeleteEducation)

	// Experience routes
	users.Post("/me/experience", deps.UserHandler.AddExperience)
	users.Put("/me/experience/:id", deps.UserHandler.UpdateExperience)
	users.Delete("/me/experience/:id", deps.UserHandler.DeleteExperience)

	// Skills routes
	users.Post("/me/skills", deps.UserHandler.AddSkill)        // Single skill
	users.Post("/me/skills/batch", deps.UserHandler.AddSkills) // Multiple skills
	users.Delete("/me/skills/:id", deps.UserHandler.DeleteSkill)

	// Document upload routes
	users.Post("/me/documents",
		middleware.UploadRateLimiter(),
		middleware.ValidateFileUpload(middleware.FileUploadConfig{
			MaxFileSize: 10 * 1024 * 1024, // 10MB
			AllowedMimeTypes: []string{
				"application/pdf",
				"image/jpeg",
				"image/png",
				"application/msword",
				"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			},
			AllowedExtensions: []string{".pdf", ".jpg", ".jpeg", ".png", ".doc", ".docx"},
			Required:          true,
			FieldName:         "file",
		}),
		deps.UserHandler.UploadDocument,
	)
}
