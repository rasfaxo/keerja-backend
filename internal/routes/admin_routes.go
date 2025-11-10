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

	// Master Data Management
	setupAdminMasterDataRoutes(admin, deps)
}

// setupAdminMasterDataRoutes configures admin master data CRUD routes
func setupAdminMasterDataRoutes(admin fiber.Router, deps *Dependencies) {
	// Provinces CRUD
	provinces := admin.Group("/master/provinces")
	provinces.Post("/", deps.AdminMasterDataHandler.CreateProvince)
	provinces.Get("/", deps.AdminMasterDataHandler.GetProvinces)
	provinces.Get("/:id", deps.AdminMasterDataHandler.GetProvinceByID)
	provinces.Put("/:id", deps.AdminMasterDataHandler.UpdateProvince)
	provinces.Delete("/:id", deps.AdminMasterDataHandler.DeleteProvince)

	// Cities CRUD
	cities := admin.Group("/master/cities")
	cities.Post("/", deps.AdminMasterDataHandler.CreateCity)
	cities.Get("/", deps.AdminMasterDataHandler.GetCities)
	cities.Get("/:id", deps.AdminMasterDataHandler.GetCityByID)
	cities.Put("/:id", deps.AdminMasterDataHandler.UpdateCity)
	cities.Delete("/:id", deps.AdminMasterDataHandler.DeleteCity)

	// Districts CRUD
	districts := admin.Group("/master/districts")
	districts.Post("/", deps.AdminMasterDataHandler.CreateDistrict)
	districts.Get("/", deps.AdminMasterDataHandler.GetDistricts)
	districts.Get("/:id", deps.AdminMasterDataHandler.GetDistrictByID)
	districts.Put("/:id", deps.AdminMasterDataHandler.UpdateDistrict)
	districts.Delete("/:id", deps.AdminMasterDataHandler.DeleteDistrict)

	// Industries CRUD
	industries := admin.Group("/master/industries")
	industries.Post("/", deps.AdminMasterDataHandler.CreateIndustry)
	industries.Get("/", deps.AdminMasterDataHandler.GetIndustries)
	industries.Get("/:id", deps.AdminMasterDataHandler.GetIndustryByID)
	industries.Put("/:id", deps.AdminMasterDataHandler.UpdateIndustry)
	industries.Delete("/:id", deps.AdminMasterDataHandler.DeleteIndustry)

	// Job Types CRUD
	jobTypes := admin.Group("/master/job-types")
	jobTypes.Post("/", deps.AdminMasterDataHandler.CreateJobType)
	jobTypes.Get("/", deps.AdminMasterDataHandler.GetJobTypes)
	jobTypes.Get("/:id", deps.AdminMasterDataHandler.GetJobTypeByID)
	jobTypes.Put("/:id", deps.AdminMasterDataHandler.UpdateJobType)
	jobTypes.Delete("/:id", deps.AdminMasterDataHandler.DeleteJobType)

	// Company Sizes CRUD (note: endpoint is /admin/meta/company-sizes as per requirement)
	companySizes := admin.Group("/meta/company-sizes")
	companySizes.Post("/", deps.AdminMasterDataHandler.CreateCompanySize)
	companySizes.Get("/", deps.AdminMasterDataHandler.GetCompanySizes)
	companySizes.Get("/:id", deps.AdminMasterDataHandler.GetCompanySizeByID)
	companySizes.Put("/:id", deps.AdminMasterDataHandler.UpdateCompanySize)
	companySizes.Delete("/:id", deps.AdminMasterDataHandler.DeleteCompanySize)
}
