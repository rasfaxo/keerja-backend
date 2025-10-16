package routes

import (
	"keerja-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupCompanyRoutes configures company routes
// Routes: /api/v1/companies/*
func SetupCompanyRoutes(api fiber.Router, deps *Dependencies, authMw *middleware.AuthMiddleware) {
	companies := api.Group("/companies")

	// Public routes
	companies.Get("/", func(c *fiber.Ctx) error {
		// TODO: Implement ListCompanies handler
		// deps.CompanyHandler.ListCompanies(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "List companies endpoint - Coming soon",
		})
	})

	companies.Get("/:id", func(c *fiber.Ctx) error {
		// TODO: Implement GetCompany handler
		// deps.CompanyHandler.GetCompany(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Get company details endpoint - Coming soon",
		})
	})

	// Protected routes
	protected := companies.Group("")
	protected.Use(authMw.AuthRequired())

	protected.Post("/", func(c *fiber.Ctx) error {
		// TODO: Implement RegisterCompany handler
		// deps.CompanyHandler.RegisterCompany(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Register company endpoint - Coming soon",
		})
	})

	protected.Put("/:id", func(c *fiber.Ctx) error {
		// TODO: Implement UpdateCompany handler
		// deps.CompanyHandler.UpdateCompany(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Update company endpoint - Coming soon",
		})
	})

	protected.Post("/:id/follow", func(c *fiber.Ctx) error {
		// TODO: Implement FollowCompany handler
		// deps.CompanyHandler.FollowCompany(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Follow company endpoint - Coming soon",
		})
	})

	protected.Delete("/:id/follow", func(c *fiber.Ctx) error {
		// TODO: Implement UnfollowCompany handler
		// deps.CompanyHandler.UnfollowCompany(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Unfollow company endpoint - Coming soon",
		})
	})

	protected.Post("/:id/review", func(c *fiber.Ctx) error {
		// TODO: Implement AddReview handler
		// deps.CompanyHandler.AddReview(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Add company review endpoint - Coming soon",
		})
	})

	protected.Post("/:id/invite-employee", func(c *fiber.Ctx) error {
		// TODO: Implement InviteEmployee handler
		// deps.CompanyHandler.InviteEmployee(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Invite employee endpoint - Coming soon",
		})
	})

	protected.Post("/:id/verify", func(c *fiber.Ctx) error {
		// TODO: Implement RequestVerification handler
		// deps.CompanyHandler.RequestVerification(c)
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"message": "Request company verification endpoint - Coming soon",
		})
	})
}
