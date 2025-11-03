package middleware

import (
	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// PermissionMiddleware handles role-based permission checks for employer users
type PermissionMiddleware struct {
	companyService company.CompanyService
}

// NewPermissionMiddleware creates a new permission middleware
func NewPermissionMiddleware(companyService company.CompanyService) *PermissionMiddleware {
	return &PermissionMiddleware{
		companyService: companyService,
	}
}

// RequirePermission checks if the user has a specific permission for a company
func (pm *PermissionMiddleware) RequirePermission(permission company.Permission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		// Get user ID from context (set by auth middleware)
		userID := GetUserID(c)
		if userID == 0 {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated", "")
		}

		// Get company ID from path parameter
		companyID, err := c.ParamsInt("id")
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", err.Error())
		}

		// Get user's role in the company
		employerUser, err := pm.companyService.GetEmployerUser(ctx, userID, int64(companyID))
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "You are not an employer of this company", err.Error())
		}

		// Check if user has the required permission
		if !company.HasPermission(employerUser.Role, permission) {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "You don't have permission to perform this action", "")
		}

		// Store employer user in context for handlers to use
		c.Locals("employer_user", employerUser)
		c.Locals("company_id", int64(companyID))

		return c.Next()
	}
}

// RequireRole checks if the user has a specific role or higher
func (pm *PermissionMiddleware) RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		// Get user ID from context
		userID := GetUserID(c)
		if userID == 0 {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated", "")
		}

		// Get company ID from path parameter
		companyID, err := c.ParamsInt("id")
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", err.Error())
		}

		// Get user's role in the company
		employerUser, err := pm.companyService.GetEmployerUser(ctx, userID, int64(companyID))
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "You are not an employer of this company", err.Error())
		}

		// Check if user has required role or higher
		if !company.HasHigherRole(employerUser.Role, requiredRole) {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Insufficient role. Required: "+requiredRole+", Your role: "+employerUser.Role, "")
		}

		// Store employer user in context
		c.Locals("employer_user", employerUser)
		c.Locals("company_id", int64(companyID))

		return c.Next()
	}
}

// RequireAdmin checks if the user is an admin
func (pm *PermissionMiddleware) RequireAdmin() fiber.Handler {
	return pm.RequireRole("admin")
}

// RequireOwnerOrAdmin checks if the user is an owner or admin
func (pm *PermissionMiddleware) RequireOwnerOrAdmin() fiber.Handler {
	return pm.RequireRole("admin") // This will pass for both admin (level 3) and owner (level 4)
}

// RequireRecruiterOrAbove checks if the user is a recruiter or admin
func (pm *PermissionMiddleware) RequireRecruiterOrAbove() fiber.Handler {
	return pm.RequireRole("recruiter")
}

// CanManageJobs checks if user can manage jobs
func (pm *PermissionMiddleware) CanManageJobs() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		userID := GetUserID(c)
		if userID == 0 {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated", "")
		}

		companyID, err := c.ParamsInt("id")
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", err.Error())
		}

		employerUser, err := pm.companyService.GetEmployerUser(ctx, userID, int64(companyID))
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "You are not an employer of this company", err.Error())
		}

		if !company.CanManageJobs(employerUser.Role) {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "You don't have permission to manage jobs", "")
		}

		c.Locals("employer_user", employerUser)
		c.Locals("company_id", int64(companyID))

		return c.Next()
	}
}

// CanManageEmployees checks if user can manage employees
func (pm *PermissionMiddleware) CanManageEmployees() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		userID := GetUserID(c)
		if userID == 0 {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated", "")
		}

		companyID, err := c.ParamsInt("id")
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", err.Error())
		}

		employerUser, err := pm.companyService.GetEmployerUser(ctx, userID, int64(companyID))
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "You are not an employer of this company", err.Error())
		}

		if !company.CanManageEmployees(employerUser.Role) {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "You don't have permission to manage employees", "")
		}

		c.Locals("employer_user", employerUser)
		c.Locals("company_id", int64(companyID))

		return c.Next()
	}
}

// CanManageApplications checks if user can manage applications
func (pm *PermissionMiddleware) CanManageApplications() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		userID := GetUserID(c)
		if userID == 0 {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated", "")
		}

		companyID, err := c.ParamsInt("id")
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", err.Error())
		}

		employerUser, err := pm.companyService.GetEmployerUser(ctx, userID, int64(companyID))
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "You are not an employer of this company", err.Error())
		}

		if !company.CanManageApplications(employerUser.Role) {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "You don't have permission to manage applications", "")
		}

		c.Locals("employer_user", employerUser)
		c.Locals("company_id", int64(companyID))

		return c.Next()
	}
}

// GetEmployerUser retrieves employer user from context
func GetEmployerUser(c *fiber.Ctx) *company.EmployerUser {
	if employerUser, ok := c.Locals("employer_user").(*company.EmployerUser); ok {
		return employerUser
	}
	return nil
}

// GetCompanyIDFromContext retrieves company ID from context
func GetCompanyIDFromContext(c *fiber.Ctx) int64 {
	if companyID, ok := c.Locals("company_id").(int64); ok {
		return companyID
	}
	return 0
}
