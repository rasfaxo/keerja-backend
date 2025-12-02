package helpers

import (
	"keerja-backend/internal/domain/user"
	"keerja-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// GetProfile fetches authenticated user's profile using middleware to get user id
func GetProfile(c *fiber.Ctx, svc user.UserService) (*user.User, error) {
	ctx := c.Context()
	userID := middleware.GetUserID(c)
	return svc.GetProfile(ctx, userID)
}
