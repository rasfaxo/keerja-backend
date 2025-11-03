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
// - Invitations: CompanyInviteHandler (5 endpoints)
// Total: 31 endpoints
func SetupCompanyRoutes(api fiber.Router, deps *Dependencies, authMw *middleware.AuthMiddleware, permMw *middleware.PermissionMiddleware) {
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

	// Update company details (admin only)
	protected.Put("/:id",
		permMw.RequireAdmin(),
		deps.CompanyBasicHandler.UpdateCompany,
	)

	// Delete company (owner or admin only)
	protected.Delete("/:id",
		permMw.RequireOwnerOrAdmin(),
		deps.CompanyBasicHandler.DeleteCompany,
	)

	// Upload company logo (admin only)
	protected.Post("/:id/logo",
		permMw.RequireAdmin(),
		deps.CompanyBasicHandler.UploadLogo,
	)

	// Delete company logo (admin only)
	protected.Delete("/:id/logo",
		permMw.RequireAdmin(),
		deps.CompanyBasicHandler.DeleteLogo,
	)

	// Upload company banner (admin only)
	protected.Post("/:id/banner",
		permMw.RequireAdmin(),
		deps.CompanyBasicHandler.UploadBanner,
	)

	// Delete company banner (admin only)
	protected.Delete("/:id/banner",
		permMw.RequireAdmin(),
		deps.CompanyBasicHandler.DeleteBanner,
	)

	// ------------------------------------------
	// Profile Management (CompanyProfileHandler)
	// ------------------------------------------

	// Update company profile (admin only)
	protected.Put("/:id/profile",
		permMw.RequireAdmin(),
		deps.CompanyProfileHandler.UpdateProfile,
	)

	// Publish company profile (admin only)
	protected.Post("/:id/profile/publish",
		permMw.RequireAdmin(),
		deps.CompanyProfileHandler.PublishProfile,
	)

	// Unpublish company profile (admin only)
	protected.Post("/:id/profile/unpublish",
		permMw.RequireAdmin(),
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
	// Additional Protected Routes (Employee Invitations)
	// ------------------------------------------

	// Invite employee to company (admin only)
	protected.Post("/:id/invite-employee",
		middleware.EmailRateLimiter(), // Rate limit invitations - 3/hour
		permMw.CanManageEmployees(),
		deps.CompanyInviteHandler.InviteEmployee,
	)

	// Accept invitation (global endpoint - not tied to specific company)
	protected.Post("/invitations/accept",
		deps.CompanyInviteHandler.AcceptInvitation,
	)

	// Get pending invitations for a company (admin only)
	protected.Get("/:id/invitations",
		permMw.CanManageEmployees(),
		deps.CompanyInviteHandler.GetPendingInvitations,
	)

	// Resend invitation (admin only)
	protected.Post("/:id/invitations/:invitationId/resend",
		middleware.EmailRateLimiter(), // Rate limit resends
		permMw.CanManageEmployees(),
		deps.CompanyInviteHandler.ResendInvitation,
	)

	// Cancel invitation (admin only)
	protected.Delete("/:id/invitations/:invitationId",
		permMw.CanManageEmployees(),
		deps.CompanyInviteHandler.CancelInvitation,
	)

	// TODO: Implement RequestVerification handler
	protected.Post("/:id/verify",
		middleware.APIRateLimiter(), // Rate limit verification requests
		permMw.RequireAdmin(),
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
				"message": "Request company verification endpoint - Coming soon",
			})
		},
	)
}
