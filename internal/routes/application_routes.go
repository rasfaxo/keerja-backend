package routes

import (
	"keerja-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupApplicationRoutes configures application routes
// Routes: /api/v1/applications/*
func SetupApplicationRoutes(api fiber.Router, deps *Dependencies, authMw *middleware.AuthMiddleware) {
	applications := api.Group("/applications")
	applications.Use(authMw.AuthRequired())

	// Job seeker routes
	applications.Post("/jobs/:id/apply",
		authMw.JobSeekerOnly(),
		middleware.ApplicationRateLimiter(),
		func(c *fiber.Ctx) error {
			// TODO: Implement ApplyJob handler
			// deps.ApplicationHandler.ApplyJob(c)
			return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
				"message": "Apply to job endpoint - Coming soon",
			})
		},
	)

	applications.Get("/my-applications", func(c *fiber.Ctx) error {
		// TODO: Implement GetMyApplications handler
		// deps.ApplicationHandler.GetMyApplications(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Get my applications endpoint - Coming soon",
		})
	})

	applications.Get("/:id", func(c *fiber.Ctx) error {
		// TODO: Implement GetApplicationDetail handler
		// deps.ApplicationHandler.GetApplicationDetail(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Get application details endpoint - Coming soon",
		})
	})

	applications.Delete("/:id/withdraw", func(c *fiber.Ctx) error {
		// TODO: Implement WithdrawApplication handler
		// deps.ApplicationHandler.WithdrawApplication(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Withdraw application endpoint - Coming soon",
		})
	})

	// Employer routes
	employer := applications.Group("/:id")
	employer.Use(authMw.EmployerOnly())

	employer.Put("/stage", func(c *fiber.Ctx) error {
		// TODO: Implement UpdateApplicationStage handler
		// deps.ApplicationHandler.UpdateApplicationStage(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Update application stage endpoint - Coming soon",
		})
	})

	employer.Post("/notes", func(c *fiber.Ctx) error {
		// TODO: Implement AddNote handler
		// deps.ApplicationHandler.AddNote(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Add application note endpoint - Coming soon",
		})
	})

	employer.Post("/schedule-interview", func(c *fiber.Ctx) error {
		// TODO: Implement ScheduleInterview handler
		// deps.ApplicationHandler.ScheduleInterview(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Schedule interview endpoint - Coming soon",
		})
	})
}
