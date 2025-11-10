package http

import (
	"keerja-backend/internal/domain/job"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// AdminHandler handles admin operations
type AdminHandler struct {
	jobService job.AdminJobService
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(jobService job.AdminJobService) *AdminHandler {
	return &AdminHandler{
		jobService: jobService,
	}
}

// ApproveJob handles PATCH /api/v1/admin/jobs/:id/approve
func (h *AdminHandler) ApproveJob(c *fiber.Ctx) error {
	jobID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid job ID",
		})
	}

	if err := h.jobService.ApproveJob(c.Context(), jobID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to approve job",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Job approved successfully",
	})
}

// RejectJob handles PATCH /api/v1/admin/jobs/:id/reject
func (h *AdminHandler) RejectJob(c *fiber.Ctx) error {
	jobID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid job ID",
		})
	}

	var req struct {
		Reason string `json:"reason" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := h.jobService.RejectJob(c.Context(), jobID, req.Reason); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to reject job",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Job rejected successfully",
	})
}