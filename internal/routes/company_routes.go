package routes

import (
	"keerja-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupCompanyRoutes configures company routes
// Routes: /api/v1/companies/*
//
// Route Organization:
// - Basic CRUD: CompanyBasicHandler (10 endpoints)
// - Profile & Social: CompanyProfileHandler (8 endpoints)
// - Reviews & Ratings: CompanyReviewHandler (5 endpoints)
// - Statistics & Queries: CompanyStatsHandler (3 endpoints)
// Total: 26 endpoints
func SetupCompanyRoutes(api fiber.Router, deps *Dependencies, authMw *middleware.AuthMiddleware) {
	companies := api.Group("/companies")

	// ==========================================
	// PROTECTED ROUTES - User-Specific
	// ==========================================

	// Get my companies (where user is a member)
	companies.Get("/my-companies",
		authMw.AuthRequired(),
		deps.CompanyBasicHandler.GetMyCompanies,
	)

	// Get user's followed companies
	companies.Get("/followed",
		authMw.AuthRequired(),
		deps.CompanyProfileHandler.GetFollowedCompanies,
	)

	// ==========================================
	// PUBLIC ROUTES - Basic CRUD (CompanyBasicHandler)
	// ==========================================

	// List companies with filters/search/pagination
	companies.Get("/",
		middleware.SearchRateLimiter(), // Rate limit search/listing - 30 req/min
		deps.CompanyBasicHandler.ListCompanies,
	)

	// Get company by ID
	companies.Get("/:id",
		deps.CompanyBasicHandler.GetCompany,
	)

	// Get company by slug (SEO-friendly)
	companies.Get("/slug/:slug",
		deps.CompanyBasicHandler.GetCompanyBySlug,
	)

	// ==========================================
	// PUBLIC ROUTES - Statistics (CompanyStatsHandler)
	// ==========================================

	// Get verified companies
	companies.Get("/verified",
		deps.CompanyStatsHandler.GetVerifiedCompanies,
	)

	// Get top-rated companies
	companies.Get("/top-rated",
		deps.CompanyStatsHandler.GetTopRatedCompanies,
	)

	// ==========================================
	// PUBLIC ROUTES - Reviews (CompanyReviewHandler)
	// ==========================================

	// Get company reviews (public)
	companies.Get("/:id/reviews",
		deps.CompanyReviewHandler.GetCompanyReviews,
	)

	// Get company average ratings (public)
	companies.Get("/:id/ratings",
		deps.CompanyReviewHandler.GetAverageRatings,
	)

	// ==========================================
	// PUBLIC ROUTES - Profile (CompanyProfileHandler)
	// ==========================================

	// Get company profile (public)
	companies.Get("/:id/profile",
		deps.CompanyProfileHandler.GetProfile,
	)

	// Get company followers (public)
	companies.Get("/:id/followers",
		deps.CompanyProfileHandler.GetFollowers,
	)

	// Get company statistics (public)
	companies.Get("/:id/stats",
		deps.CompanyStatsHandler.GetCompanyStats,
	)

	// ==========================================
	// PROTECTED ROUTES - Authentication Required
	// ==========================================
	protected := companies.Group("")
	protected.Use(authMw.AuthRequired())

	// ------------------------------------------
	// Basic CRUD Operations (CompanyBasicHandler)
	// ------------------------------------------

	// Create company (register)
	protected.Post("/",
		deps.CompanyBasicHandler.CreateCompany,
	)

	// Update company details
	protected.Put("/:id",
		deps.CompanyBasicHandler.UpdateCompany,
	)

	// Delete company
	protected.Delete("/:id",
		deps.CompanyBasicHandler.DeleteCompany,
	)

	// Upload company logo
	protected.Post("/:id/logo",
		deps.CompanyBasicHandler.UploadLogo,
	)

	// Delete company logo
	protected.Delete("/:id/logo",
		deps.CompanyBasicHandler.DeleteLogo,
	)

	// Upload company banner
	protected.Post("/:id/banner",
		deps.CompanyBasicHandler.UploadBanner,
	)

	// Delete company banner
	protected.Delete("/:id/banner",
		deps.CompanyBasicHandler.DeleteBanner,
	)

	// ------------------------------------------
	// Profile Management (CompanyProfileHandler)
	// ------------------------------------------

	// Update company profile
	protected.Put("/:id/profile",
		deps.CompanyProfileHandler.UpdateProfile,
	)

	// Publish company profile (make public)
	protected.Post("/:id/profile/publish",
		deps.CompanyProfileHandler.PublishProfile,
	)

	// Unpublish company profile (make private)
	protected.Post("/:id/profile/unpublish",
		deps.CompanyProfileHandler.UnpublishProfile,
	)

	// ------------------------------------------
	// Social Features (CompanyProfileHandler)
	// ------------------------------------------

	// Follow company
	protected.Post("/:id/follow",
		deps.CompanyProfileHandler.FollowCompany,
	)

	// Unfollow company
	protected.Delete("/:id/follow",
		deps.CompanyProfileHandler.UnfollowCompany,
	)

	// ------------------------------------------
	// Review Management (CompanyReviewHandler)
	// ------------------------------------------

	// Add company review with rate limiting (prevent spam)
	protected.Post("/:id/review",
		middleware.APIRateLimiter(), // Rate limit reviews - 100 req/min
		deps.CompanyReviewHandler.AddReview,
	)

	// Update company review (own review only)
	protected.Put("/:id/review/:reviewId",
		deps.CompanyReviewHandler.UpdateReview,
	)

	// Delete company review (own review only)
	protected.Delete("/:id/review/:reviewId",
		deps.CompanyReviewHandler.DeleteReview,
	)

	// ------------------------------------------
	// Additional Protected Routes (Future Implementation)
	// ------------------------------------------

	// Invite employee to company
	protected.Post("/:id/invite-employee",
		middleware.EmailRateLimiter(), // Rate limit invitations - 3/hour
		deps.CompanyInviteHandler.InviteEmployee,
	)

	// TODO: Implement RequestVerification handler
	protected.Post("/:id/verify",
		middleware.APIRateLimiter(), // Rate limit verification requests
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
				"message": "Request company verification endpoint - Coming soon",
			})
		},
	)
}
