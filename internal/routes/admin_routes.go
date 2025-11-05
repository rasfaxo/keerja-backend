package routes

import (
	"keerja-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupAdminRoutes configures admin routes
// Routes: /api/v1/admin/*
func SetupAdminRoutes(api fiber.Router, deps *Dependencies, authMw *middleware.AuthMiddleware) {
	admin := api.Group("/admin")

	// All admin routes require authentication and admin role
	admin.Use(authMw.AuthRequired())
	admin.Use(authMw.AdminOnly())

	// Dashboard
	admin.Get("/dashboard", func(c *fiber.Ctx) error {
		// TODO: Implement GetDashboard handler
		// deps.AdminHandler.GetDashboard(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Admin dashboard endpoint - Coming soon",
		})
	})

	// User management
	admin.Get("/users", func(c *fiber.Ctx) error {
		// TODO: Implement GetUsers handler
		// deps.AdminHandler.GetUsers(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Get users endpoint - Coming soon",
		})
	})

	admin.Get("/users/:id", func(c *fiber.Ctx) error {
		// TODO: Implement GetUserDetail handler
		// deps.AdminHandler.GetUserDetail(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Get user detail endpoint - Coming soon",
		})
	})

	admin.Put("/users/:id/status", func(c *fiber.Ctx) error {
		// TODO: Implement UpdateUserStatus handler
		// deps.AdminHandler.UpdateUserStatus(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Update user status endpoint - Coming soon",
		})
	})

	// Company management
	admin.Get("/companies", func(c *fiber.Ctx) error {
		// TODO: Implement GetCompanies handler
		// deps.AdminHandler.GetCompanies(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Get companies endpoint - Coming soon",
		})
	})

	admin.Get("/companies/:id", func(c *fiber.Ctx) error {
		// TODO: Implement GetCompanyDetail handler
		// deps.AdminHandler.GetCompanyDetail(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Get company detail endpoint - Coming soon",
		})
	})

	admin.Put("/companies/:id/verify", func(c *fiber.Ctx) error {
		// TODO: Implement VerifyCompany handler
		// deps.AdminHandler.VerifyCompany(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Verify company endpoint - Coming soon",
		})
	})

	admin.Put("/companies/:id/status", func(c *fiber.Ctx) error {
		// TODO: Implement UpdateCompanyStatus handler
		// deps.AdminHandler.UpdateCompanyStatus(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Update company status endpoint - Coming soon",
		})
	})

	// Job management
	admin.Get("/jobs", func(c *fiber.Ctx) error {
		// TODO: Implement GetJobs handler to list pending jobs
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Get jobs endpoint - Coming soon",
		})
	})

	admin.Patch("/jobs/:id/approve",
		deps.AdminHandler.ApproveJob,
	)

	admin.Patch("/jobs/:id/reject",
		deps.AdminHandler.RejectJob,
	)

	admin.Put("/jobs/:id/status", func(c *fiber.Ctx) error {
		// TODO: Implement UpdateJobStatus handler for other status changes
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Update job status endpoint - Coming soon",
		})
	})

	// Application monitoring
	admin.Get("/applications", func(c *fiber.Ctx) error {
		// TODO: Implement GetApplications handler
		// deps.AdminHandler.GetApplications(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Get applications endpoint - Coming soon",
		})
	})

	// Reports & Analytics
	admin.Get("/reports/users", func(c *fiber.Ctx) error {
		// TODO: Implement GetUserReport handler
		// deps.AdminHandler.GetUserReport(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Get user report endpoint - Coming soon",
		})
	})

	admin.Get("/reports/jobs", func(c *fiber.Ctx) error {
		// TODO: Implement GetJobReport handler
		// deps.AdminHandler.GetJobReport(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Get job report endpoint - Coming soon",
		})
	})

	admin.Get("/reports/applications", func(c *fiber.Ctx) error {
		// TODO: Implement GetApplicationReport handler
		// deps.AdminHandler.GetApplicationReport(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Get application report endpoint - Coming soon",
		})
	})
}
