package routes

import (
	"keerja-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupJobRoutes configures job routes
// Routes: /api/v1/jobs/*
func SetupJobRoutes(api fiber.Router, deps *Dependencies, authMw *middleware.AuthMiddleware) {
	jobs := api.Group("/jobs")

	// Public routes
	jobs.Get("/", func(c *fiber.Ctx) error {
		// TODO: Implement ListJobs handler
		// deps.JobHandler.ListJobs(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Job listing endpoint - Coming soon",
		})
	})

	jobs.Get("/:id", func(c *fiber.Ctx) error {
		// TODO: Implement GetJob handler
		// deps.JobHandler.GetJob(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Get job details endpoint - Coming soon",
		})
	})

	jobs.Post("/search",
		middleware.SearchRateLimiter(),
		func(c *fiber.Ctx) error {
			// TODO: Implement SearchJobs handler
			// deps.JobHandler.SearchJobs(c)
			return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
				"message": "Job search endpoint - Coming soon",
			})
		},
	)

	// Protected routes - employer only
	protected := jobs.Group("")
	protected.Use(authMw.AuthRequired())
	protected.Use(authMw.EmployerOnly())

	protected.Post("/", func(c *fiber.Ctx) error {
		// TODO: Implement CreateJob handler
		// deps.JobHandler.CreateJob(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Create job endpoint - Coming soon",
		})
	})

	protected.Put("/:id", func(c *fiber.Ctx) error {
		// TODO: Implement UpdateJob handler
		// deps.JobHandler.UpdateJob(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Update job endpoint - Coming soon",
		})
	})

	protected.Delete("/:id", func(c *fiber.Ctx) error {
		// TODO: Implement DeleteJob handler
		// deps.JobHandler.DeleteJob(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Delete job endpoint - Coming soon",
		})
	})

	protected.Get("/my-jobs", func(c *fiber.Ctx) error {
		// TODO: Implement GetMyJobs handler
		// deps.JobHandler.GetMyJobs(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Get my jobs endpoint - Coming soon",
		})
	})

	protected.Get("/:id/applications", func(c *fiber.Ctx) error {
		// TODO: Implement GetJobApplications handler
		// deps.JobHandler.GetJobApplications(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Get job applications endpoint - Coming soon",
		})
	})
}
