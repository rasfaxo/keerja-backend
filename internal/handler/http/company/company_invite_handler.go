package companyhandler

import (
	"fmt"

	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/domain/email"
	"keerja-backend/internal/domain/user"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/handler/http/common"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// CompanyInviteHandler handles company employee invitation operations
type CompanyInviteHandler struct {
	companyService company.CompanyService
	emailService   email.EmailService
	userService    user.UserService
}

// NewCompanyInviteHandler creates a new instance of CompanyInviteHandler
func NewCompanyInviteHandler(companyService company.CompanyService, emailService email.EmailService, userService user.UserService) *CompanyInviteHandler {
	return &CompanyInviteHandler{
		companyService: companyService,
		emailService:   emailService,
		userService:    userService,
	}
}

func (h *CompanyInviteHandler) InviteEmployee(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get company ID from path parameter
	companyID, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, common.ErrInvalidCompanyID, err.Error())
	}

	// Get authenticated user from context using middleware helper
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated", "userID not found in context")
	}

	// Parse and validate request body
	var req request.InviteEmployeeRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, common.ErrInvalidRequest, err.Error())
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, common.ErrValidationFailed, errors)
	}

	// Sanitize input
	req.Email = utils.SanitizeString(req.Email)
	req.FullName = utils.SanitizeString(req.FullName)
	req.Position = utils.SanitizeString(req.Position)
	req.Role = utils.SanitizeString(req.Role)

	// Get company to verify ownership and existence
	comp, err := h.companyService.GetCompany(ctx, int64(companyID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, common.ErrCompanyNotFound, err.Error())
	}

	// Check if user is authorized to invite employees (owner or admin)
	hasPermission, err := h.companyService.CheckEmployerPermission(ctx, userID, int64(companyID), "admin")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to check user permission", err.Error())
	}
	if !hasPermission {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "You don't have permission to invite employees. Only company owner or admin can perform this action.", "")
	}

	// Create invitation request
	inviteReq := &company.InviteEmployerRequest{
		CompanyID:     int64(companyID),
		Email:         req.Email,
		Role:          req.Role,
		PositionTitle: &req.Position,
	}

	// Save invitation to database
	if err := h.companyService.InviteEmployer(ctx, inviteReq); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, common.ErrFailedOperation, err.Error())
	}

	// Get invitation to get token for email
	invitations, err := h.companyService.GetPendingInvitations(ctx, int64(companyID))
	if err != nil || len(invitations) == 0 {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, common.ErrFailedOperation, "")
	}

	// Find the latest invitation for this email
	var invitation *company.CompanyInvitation
	for i := range invitations {
		if invitations[i].Email == req.Email && invitations[i].Status == "pending" {
			invitation = &invitations[i]
			break
		}
	}

	if invitation == nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, common.ErrFailedOperation, "")
	}

	// Generate invitation URL
	inviteURL := fmt.Sprintf("%s/accept-invite?token=%s", c.BaseURL(), invitation.Token)

	// Get inviter name from user service
	inviterName := "Administrator" // Default fallback
	inviterUser, err := h.userService.GetProfile(ctx, userID)
	if err == nil && inviterUser != nil {
		inviterName = inviterUser.FullName
	}

	// Send invitation email using template
	if err := h.emailService.SendCompanyInvitationEmail(
		ctx,
		req.Email,
		req.FullName,
		comp.CompanyName,
		inviterName,
		req.Position,
		req.Role,
		inviteURL,
		7, // 7 days expiry
	); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, common.ErrFailedOperation, err.Error())
	}

	return utils.SuccessResponse(c, "Invitation sent successfully", fiber.Map{
		"email":         req.Email,
		"name":          req.FullName,
		"position":      req.Position,
		"role":          req.Role,
		"company":       comp.CompanyName,
		"expires_at":    invitation.ExpiresAt,
		"invitation_id": invitation.ID,
	})
}

func (h *CompanyInviteHandler) AcceptInvitation(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get token from query parameter
	token := c.Query("token")
	if token == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invitation token is required", "")
	}

	// Get authenticated user from context
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated", "userID not found in context")
	}

	// Accept invitation
	if err := h.companyService.AcceptInvitation(ctx, token, userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, common.ErrFailedOperation, err.Error())
	}

	return utils.SuccessResponse(c, "Invitation accepted successfully", fiber.Map{
		"message": "You are now an employer of this company",
	})
}

func (h *CompanyInviteHandler) ResendInvitation(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get company ID and invitation ID from path parameters
	companyID, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, common.ErrInvalidCompanyID, err.Error())
	}

	invitationID, err := c.ParamsInt("invitationId")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid invitation ID", err.Error())
	}

	// Get authenticated user
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated", "")
	}

	// Check permission
	hasPermission, err := h.companyService.CheckEmployerPermission(ctx, userID, int64(companyID), "admin")
	if err != nil || !hasPermission {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "You don't have permission to resend invitations", "")
	}

	// Resend invitation
	if err := h.companyService.ResendInvitation(ctx, int64(invitationID), userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, common.ErrFailedOperation, err.Error())
	}

	return utils.SuccessResponse(c, "Invitation resent successfully", nil)
}

func (h *CompanyInviteHandler) CancelInvitation(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get invitation ID from path parameter (company ID validation done by permission check)
	invitationID, err := c.ParamsInt("invitationId")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid invitation ID", err.Error())
	}

	// Get authenticated user
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated", "")
	}

	// Cancel invitation (permission check inside service)
	if err := h.companyService.CancelInvitation(ctx, int64(invitationID), userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, common.ErrFailedOperation, err.Error())
	}

	return utils.SuccessResponse(c, "Invitation canceled successfully", nil)
}

func (h *CompanyInviteHandler) GetPendingInvitations(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get company ID from path parameter
	companyID, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, common.ErrInvalidCompanyID, err.Error())
	}

	// Get authenticated user
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not authenticated", "")
	}

	// Check permission
	hasPermission, err := h.companyService.CheckEmployerPermission(ctx, userID, int64(companyID), "admin")
	if err != nil || !hasPermission {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "You don't have permission to view invitations", "")
	}

	// Get invitations
	invitations, err := h.companyService.GetPendingInvitations(ctx, int64(companyID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, common.ErrFailedOperation, err.Error())
	}

	return utils.SuccessResponse(c, "Invitations retrieved successfully", fiber.Map{
		"invitations": invitations,
		"total":       len(invitations),
	})
}
