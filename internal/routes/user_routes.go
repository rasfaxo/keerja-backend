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

	// Profile routes (UserProfileHandler)
	users.Get("/me", deps.UserProfileHandler.GetProfile) // Support ?include=all or ?include=educations,skills
	users.Put("/me", deps.UserProfileHandler.UpdateProfile)
	users.Put("/me/preferences", deps.UserProfileHandler.UpdatePreferences)
	users.Get("/me/preferences", deps.UserProfileHandler.GetPreferences)

	// Education routes (UserEducationHandler)
	users.Get("/me/educations", deps.UserEducationHandler.GetEducations)
	users.Post("/me/education", deps.UserEducationHandler.AddEducation)
	users.Put("/me/education/:id", deps.UserEducationHandler.UpdateEducation)
	users.Delete("/me/education/:id", deps.UserEducationHandler.DeleteEducation)

	// Experience routes (UserExperienceHandler)
	users.Get("/me/experiences", deps.UserExperienceHandler.GetExperiences)
	users.Post("/me/experience", deps.UserExperienceHandler.AddExperience)
	users.Put("/me/experience/:id", deps.UserExperienceHandler.UpdateExperience)
	users.Delete("/me/experience/:id", deps.UserExperienceHandler.DeleteExperience)

	// Skills routes (UserSkillHandler)
	users.Get("/me/skills", deps.UserSkillHandler.GetSkills)
	users.Post("/me/skills/batch", deps.UserSkillHandler.AddSkills) // Multiple skills
	users.Delete("/me/skills/:id", deps.UserSkillHandler.DeleteSkill)

	// Document routes (UserDocumentHandler)
	users.Get("/me/documents", deps.UserDocumentHandler.GetDocuments)
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
		deps.UserDocumentHandler.UploadDocument,
	)
	users.Delete("/me/documents/:id", deps.UserDocumentHandler.DeleteDocument)

	// Misc routes - certifications, languages, projects (UserMiscHandler)
	users.Get("/me/certifications", deps.UserMiscHandler.GetCertifications)
	users.Get("/me/languages", deps.UserMiscHandler.GetLanguages)
	users.Get("/me/projects", deps.UserMiscHandler.GetProjects)

	// Profile photo upload route (UserProfileHandler)
	users.Post("/profile-photo",
		middleware.UploadRateLimiter(),
		middleware.ValidateFileUpload(middleware.FileUploadConfig{
			MaxFileSize:       5 * 1024 * 1024, // 5MB for profile photo
			AllowedMimeTypes:  []string{"image/jpeg", "image/png", "image/webp"},
			AllowedExtensions: []string{".jpg", ".jpeg", ".png", ".webp"},
			Required:          true,
			FieldName:         "file",
		}),
		deps.UserProfileHandler.UploadProfilePhoto,
	)
}
