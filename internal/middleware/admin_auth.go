package middleware

import (
	"strings"

	"keerja-backend/internal/config"
	"keerja-backend/internal/domain/admin"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// AdminAuthMiddleware handles admin-specific authentication
type AdminAuthMiddleware struct {
	config        *config.Config
	adminUserRepo admin.AdminUserRepository
}

// NewAdminAuthMiddleware creates new admin auth middleware
func NewAdminAuthMiddleware(cfg *config.Config, adminUserRepo admin.AdminUserRepository) *AdminAuthMiddleware {
	return &AdminAuthMiddleware{
		config:        cfg,
		adminUserRepo: adminUserRepo,
	}
}

// AdminAuthRequired middleware requires valid admin JWT token
func (m *AdminAuthMiddleware) AdminAuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from Authorization header
		token, err := m.extractToken(c)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authentication required", err.Error())
		}

		// Validate admin token
		claims, err := utils.ValidateAdminToken(token, m.config.JWTSecret)
		if err != nil {
			if err == utils.ErrExpiredToken {
				return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Token expired", err.Error())
			}
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid admin token", err.Error())
		}

		// Verify admin exists and is active
		adminUser, err := m.adminUserRepo.FindByID(c.Context(), claims.AdminID)
		if err != nil || adminUser == nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Admin not found", "Admin account does not exist")
		}

		// Check admin status
		if adminUser.Status != "active" {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Admin account inactive", "Your admin account is not active")
		}

		// Store admin info in context
		c.Locals("admin_id", claims.AdminID)
		c.Locals("admin_email", claims.Email)
		c.Locals("admin_role_id", claims.RoleID)
		c.Locals("admin_access_level", claims.AccessLevel)
		c.Locals("admin_claims", claims)

		return c.Next()
	}
}

// AdminRoleRequired middleware checks if admin has required access level
func (m *AdminAuthMiddleware) AdminRoleRequired(minAccessLevel int16) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get access level from context
		accessLevel, ok := c.Locals("admin_access_level").(int16)
		if !ok {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authentication required", "Access level not found in context")
		}

		// Check if admin has required access level
		if accessLevel < minAccessLevel {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Access denied", "Insufficient access level for this operation")
		}

		return c.Next()
	}
}

// SuperAdminOnly middleware allows only super admins (access level >= 9)
func (m *AdminAuthMiddleware) SuperAdminOnly() fiber.Handler {
	return m.AdminRoleRequired(9)
}

// AdminOnly middleware allows admins (access level >= 7)
func (m *AdminAuthMiddleware) AdminOnly() fiber.Handler {
	return m.AdminRoleRequired(7)
}

// ModeratorOnly middleware allows moderators (access level >= 5)
func (m *AdminAuthMiddleware) ModeratorOnly() fiber.Handler {
	return m.AdminRoleRequired(5)
}

// extractToken extracts JWT token from Authorization header
func (m *AdminAuthMiddleware) extractToken(c *fiber.Ctx) (string, error) {
	// Get Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Authorization header missing")
	}

	// Check if it's Bearer token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization header format")
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Token is empty")
	}

	return token, nil
}

// GetAdminAccessLevel extracts admin access level from context
func GetAdminAccessLevel(c *fiber.Ctx) int16 {
	accessLevel, ok := c.Locals("admin_access_level").(int16)
	if !ok {
		return 0
	}
	return accessLevel
}

// GetAdminRoleID extracts admin role ID from context
func GetAdminRoleID(c *fiber.Ctx) *int64 {
	roleID, ok := c.Locals("admin_role_id").(*int64)
	if !ok {
		return nil
	}
	return roleID
}

// GetAdminEmail extracts admin email from context
func GetAdminEmail(c *fiber.Ctx) string {
	email, ok := c.Locals("admin_email").(string)
	if !ok {
		return ""
	}
	return email
}
