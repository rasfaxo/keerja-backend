package routes

import (
	"keerja-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupAdminRoutes configures admin routes
// Routes: /api/v1/admin/*
func SetupAdminRoutes(api fiber.Router, deps *Dependencies, adminAuthMw *middleware.AdminAuthMiddleware) {
	admin := api.Group("/admin")

	// All admin routes require admin authentication
	admin.Use(adminAuthMw.AdminAuthRequired())

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
	// Task 2.1: List companies with filters, pagination, and search
	admin.Get("/companies", deps.AdminCompanyHandler.ListCompanies)

	// Task 2.2: Get company detail for review
	admin.Get("/companies/:id", deps.AdminCompanyHandler.GetCompanyDetail)

	// Task 2.3: Update company status (approve/reject/suspend)
	admin.Patch("/companies/:id/status", deps.AdminCompanyHandler.UpdateCompanyStatus)

	// Task 2.4: Edit company details (admin support)
	admin.Put("/companies/:id", deps.AdminCompanyHandler.UpdateCompany)

	// Task 2.5: Delete company with validation
	admin.Delete("/companies/:id", deps.AdminCompanyHandler.DeleteCompany)

	// Additional company endpoints
	admin.Get("/companies/:id/stats", deps.AdminCompanyHandler.GetCompanyStats)
	admin.Get("/companies/:id/audit-logs", deps.AdminCompanyHandler.GetAuditLogs)

	// Dashboard stats
	admin.Get("/dashboard/stats", deps.AdminCompanyHandler.GetDashboardStats)

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
