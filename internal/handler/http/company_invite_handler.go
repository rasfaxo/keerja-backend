package http

import (
	"fmt"

	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/domain/email"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// CompanyInviteHandler handles company employee invitation operations
type CompanyInviteHandler struct {
	companyService company.CompanyService
	emailService   email.EmailService
}

// NewCompanyInviteHandler creates a new instance of CompanyInviteHandler
func NewCompanyInviteHandler(companyService company.CompanyService, emailService email.EmailService) *CompanyInviteHandler {
	return &CompanyInviteHandler{
		companyService: companyService,
		emailService:   emailService,
	}
}

// InviteEmployee godoc
// @Summary Invite employee to company
// @Description Send invitation email to employee to join the company
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Param request body request.InviteEmployeeRequest true "Invite employee request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /companies/{id}/invite-employee [post]
func (h *CompanyInviteHandler) InviteEmployee(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get company ID from path parameter
	companyID, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid company ID", err.Error())
	}

	// Get authenticated user from context
	userID := c.Locals("userID").(int64)

	// Parse and validate request body
	var req request.InviteEmployeeRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Sanitize input
	req.Email = utils.SanitizeString(req.Email)
	req.FullName = utils.SanitizeString(req.FullName)
	req.Position = utils.SanitizeString(req.Position)
	req.Role = utils.SanitizeString(req.Role)

	// Get company to verify ownership and existence
	comp, err := h.companyService.GetCompany(ctx, int64(companyID))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Company not found", err.Error())
	}

	// Check if user is authorized to invite employees (owner or admin)
	hasPermission, err := h.companyService.CheckEmployerPermission(ctx, userID, int64(companyID), "admin")
	if err != nil || !hasPermission {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "You are not authorized to invite employees for this company", "")
	}

	// Generate invitation token (valid for 7 days)
	// TODO: Store token in database or Redis for verification
	inviteToken := utils.GenerateRandomToken(32)
	inviteURL := fmt.Sprintf("%s/accept-invite?token=%s", c.BaseURL(), inviteToken)

	// Send invitation email
	emailBody := h.generateInviteEmailBody(comp.CompanyName, req.FullName, req.Position, inviteURL)
	subject := fmt.Sprintf("Undangan Bergabung di %s - Keerja", comp.CompanyName)

	if err := h.emailService.SendEmail(ctx, req.Email, subject, emailBody); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to send invitation email", err.Error())
	}

	// TODO: Store invitation record in database
	// - Save invitation token
	// - Track invitation status (pending, accepted, rejected, expired)
	// - Allow resending invitation
	// - Set expiration time (7 days)

	return utils.SuccessResponse(c, "Invitation sent successfully", fiber.Map{
		"email":    req.Email,
		"name":     req.FullName,
		"position": req.Position,
		"role":     req.Role,
		"company":  comp.CompanyName,
	})
}

// generateInviteEmailBody generates HTML email body for employee invitation
func (h *CompanyInviteHandler) generateInviteEmailBody(companyName, employeeName, position, inviteURL string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Undangan Bergabung di %s</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #4CAF50;">Undangan Bergabung di %s</h2>
        <p>Halo %s,</p>
        <p>Anda telah diundang untuk bergabung dengan <strong>%s</strong> sebagai <strong>%s</strong> melalui platform Keerja.</p>
        <div style="text-align: center; margin: 30px 0;">
            <a href="%s" style="background-color: #4CAF50; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; display: inline-block;">
                Terima Undangan
            </a>
        </div>
        <p>Atau salin dan tempel link berikut di browser Anda:</p>
        <p style="word-break: break-all; color: #666; font-size: 12px;">%s</p>
        <p>Link undangan ini akan kadaluarsa dalam 7 hari.</p>
        <hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
        <p style="font-size: 12px; color: #999;">
            Jika Anda tidak mengenal perusahaan ini atau tidak mengharapkan undangan ini, abaikan email ini.<br>
            Butuh bantuan? Hubungi kami di support@keerja.com
        </p>
    </div>
</body>
</html>
	`, companyName, companyName, employeeName, companyName, position, inviteURL, inviteURL)
}
