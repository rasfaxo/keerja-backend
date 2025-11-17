package middleware

import (
	"strings"

	"keerja-backend/internal/config"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// ContextKey types for storing user info in context
const (
	ContextKeyUserID   = "user_id"
	ContextKeyEmail    = "email"
	ContextKeyUserType = "user_type"
	ContextKeyClaims   = "claims"
)

// AuthMiddleware creates authentication middleware
type AuthMiddleware struct {
	config *config.Config
}

// NewAuthMiddleware creates a new auth middleware instance
func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		config: cfg,
	}
}

// AuthRequired middleware requires valid JWT token
func (m *AuthMiddleware) AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from Authorization header
		token, err := m.extractToken(c)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authentication required", err.Error())
		}

		// Validate token
		claims, err := utils.ValidateToken(token, m.config.JWTSecret)
		if err != nil {
			if err == utils.ErrExpiredToken {
				return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Token expired", err.Error())
			}
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid token", err.Error())
		}

		// Store user info in context
		c.Locals(ContextKeyUserID, claims.UserID)
		c.Locals(ContextKeyEmail, claims.Email)
		c.Locals(ContextKeyUserType, claims.UserType)
		c.Locals(ContextKeyClaims, claims)

		return c.Next()
	}
}

// OptionalAuth middleware parses token if exists but doesn't require it
func (m *AuthMiddleware) OptionalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Try to extract token
		token, err := m.extractToken(c)
		if err != nil {
			// No token found, continue without authentication
			return c.Next()
		}

		// Validate token
		claims, err := utils.ValidateToken(token, m.config.JWTSecret)
		if err != nil {
			// Invalid token, continue without authentication
			return c.Next()
		}

		// Store user info in context
		c.Locals(ContextKeyUserID, claims.UserID)
		c.Locals(ContextKeyEmail, claims.Email)
		c.Locals(ContextKeyUserType, claims.UserType)
		c.Locals(ContextKeyClaims, claims)

		return c.Next()
	}
}

// RoleRequired middleware checks if user has required role
func (m *AuthMiddleware) RoleRequired(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user type from context
		userType, ok := c.Locals(ContextKeyUserType).(string)
		if !ok {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authentication required", "User type not found in context")
		}

		// Check if user has required role
		hasRole := false
		for _, role := range allowedRoles {
			if strings.EqualFold(userType, role) {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Access denied", "You don't have permission to access this resource")
		}

		return c.Next()
	}
}

// JobSeekerOnly middleware allows only job seekers
func (m *AuthMiddleware) JobSeekerOnly() fiber.Handler {
	return m.RoleRequired("job_seeker")
}

// EmployerOnly middleware allows only employers
func (m *AuthMiddleware) EmployerOnly() fiber.Handler {
	return m.RoleRequired("employer")
}

// AdminOnly middleware allows only admins
func (m *AuthMiddleware) AdminOnly() fiber.Handler {
	return m.RoleRequired("admin")
}

// extractToken extracts JWT token from Authorization header
func (m *AuthMiddleware) extractToken(c *fiber.Ctx) (string, error) {
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

// GetUserID extracts user ID from context
func GetUserID(c *fiber.Ctx) int64 {
	userID, ok := c.Locals(ContextKeyUserID).(int64)
	if !ok {
		return 0
	}
	return userID
}

// GetEmail extracts email from context
func GetEmail(c *fiber.Ctx) string {
	email, ok := c.Locals(ContextKeyEmail).(string)
	if !ok {
		return ""
	}
	return email
}

// GetUserType extracts user type from context
func GetUserType(c *fiber.Ctx) string {
	userType, ok := c.Locals(ContextKeyUserType).(string)
	if !ok {
		return ""
	}
	return userType
}

// GetClaims extracts full claims from context
func GetClaims(c *fiber.Ctx) *utils.Claims {
	claims, ok := c.Locals(ContextKeyClaims).(*utils.Claims)
	if !ok {
		return nil
	}
	return claims
}

// GetAdminID extracts admin ID from admin context
func GetAdminID(c *fiber.Ctx) int64 {
	// Try to get from admin claims first
	adminClaims, ok := c.Locals("admin_claims").(*utils.AdminClaims)
	if ok && adminClaims != nil {
		return adminClaims.AdminID
	}

	// Fallback to regular user_id if admin middleware set it
	adminID, ok := c.Locals("admin_id").(int64)
	if ok {
		return adminID
	}

	return 0
}

// GetAdminClaims extracts full admin claims from context
func GetAdminClaims(c *fiber.Ctx) *utils.AdminClaims {
	claims, ok := c.Locals("admin_claims").(*utils.AdminClaims)
	if !ok {
		return nil
	}
	return claims
}
