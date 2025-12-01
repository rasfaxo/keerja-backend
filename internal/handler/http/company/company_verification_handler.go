package companyhandler

import (
	"fmt"
	"mime/multipart"

	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/handler/http/common"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// CompanyVerificationHandler handles company verification operations
type CompanyVerificationHandler struct {
	companyService company.CompanyService
}

// NewCompanyVerificationHandler creates a new instance of CompanyVerificationHandler
func NewCompanyVerificationHandler(companyService company.CompanyService) *CompanyVerificationHandler {
	return &CompanyVerificationHandler{
		companyService: companyService,
	}
}

func (h *CompanyVerificationHandler) GetCompanyVerificationStatus(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := utils.ParseIDParam(c, "id")
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	comp, err := h.companyService.GetCompany(ctx, companyID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, common.ErrCompanyNotFound, err.Error())
	}

	verification, err := h.companyService.GetVerificationStatus(ctx, companyID)

	resp := fiber.Map{
		"id":           comp.ID,
		"company_name": comp.CompanyName,
		"verified":     comp.Verified,
		"verified_at":  comp.VerifiedAt,
	}

	if err == nil && verification != nil {
		resp["status"] = verification.Status
		resp["verification_score"] = verification.VerificationScore
		resp["verification_notes"] = verification.VerificationNotes
		resp["npwp_number"] = verification.NPWPNumber
		resp["nib_number"] = verification.NIBNumber
		resp["reviewed_at"] = verification.ReviewedAt
		resp["verification_expiry"] = verification.VerificationExpiry
		resp["badge_granted"] = verification.BadgeGranted
		resp["rejection_reason"] = verification.RejectionReason
	} else {
		resp["status"] = "not_requested"
	}

	return utils.SuccessResponse(c, "Verification status retrieved successfully", resp)
}

func (h *CompanyVerificationHandler) GetMyCompanyVerificationStatus(c *fiber.Ctx) error {
	ctx := c.Context()

	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, common.ErrUnauthorized, "userID not found in context")
	}

	companies, err := h.companyService.GetUserCompanies(ctx, userID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, common.ErrFailedOperation)
	}

	if len(companies) == 0 {
		return utils.ErrorResponse(c, fiber.StatusNotFound, common.ErrCompanyNotFound, "You don't have any company registered")
	}

	comp := companies[0]

	verification, err := h.companyService.GetVerificationStatus(ctx, comp.ID)

	resp := fiber.Map{
		"id":           comp.ID,
		"company_name": comp.CompanyName,
		"verified":     comp.Verified,
		"verified_at":  comp.VerifiedAt,
	}

	if err == nil && verification != nil {
		resp["status"] = verification.Status
		resp["verification_score"] = verification.VerificationScore
		resp["verification_notes"] = verification.VerificationNotes
		resp["npwp_number"] = verification.NPWPNumber
		resp["nib_number"] = verification.NIBNumber
		resp["reviewed_at"] = verification.ReviewedAt
		resp["verification_expiry"] = verification.VerificationExpiry
		resp["badge_granted"] = verification.BadgeGranted
		resp["rejection_reason"] = verification.RejectionReason
	} else {
		resp["status"] = "not_requested"
	}

	return utils.SuccessResponse(c, "Verification status retrieved successfully", resp)
}

func (h *CompanyVerificationHandler) RequestVerification(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, common.ErrInvalidCompanyID, err.Error())
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, common.ErrUnauthorized, "User ID not found in context")
	}

	comp, err := h.companyService.GetCompany(ctx, companyID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, common.ErrCompanyNotFound, err.Error())
	}

	if comp.Verified {
		return utils.ErrorResponse(c, fiber.StatusConflict, "Company already verified", "This company is already verified")
	}

	verificationStatus, err := h.companyService.GetVerificationStatus(ctx, companyID)
	if err == nil && verificationStatus != nil {
		if verificationStatus.Status == "pending" || verificationStatus.Status == "under_review" {
			return utils.ErrorResponse(c, fiber.StatusConflict, "Verification request already submitted",
				"A verification request is already pending for this company")
		}
	}

	employerUserID, err := h.companyService.GetEmployerUserID(ctx, userID, companyID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Access denied",
			common.ErrNotCompanyMember)
	}

	var req request.RequestVerificationRequest

	npwpNumber := c.FormValue("npwp_number")
	if npwpNumber == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "NPWP number is required", "npwp_number field is required")
	}
	req.NPWPNumber = &npwpNumber

	nibNumber := c.FormValue("nib_number")
	if nibNumber != "" {
		req.NIBNumber = &nibNumber
	}

	*req.NPWPNumber = utils.SanitizeString(*req.NPWPNumber)
	if req.NIBNumber != nil {
		*req.NIBNumber = utils.SanitizeString(*req.NIBNumber)
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, common.ErrValidationFailed, err.Error())
	}

	npwpFile, err := c.FormFile("npwp_file")
	if err != nil || npwpFile == nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "NPWP file is required", "npwp_file must be uploaded")
	}

	if npwpFile.Size > 10*1024*1024 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "NPWP file too large", "Maximum file size is 10MB")
	}

	form, err := c.MultipartForm()
	var additionalFiles []*multipart.FileHeader
	if err == nil && form != nil {
		if files, ok := form.File["additional_documents"]; ok {
			maxFiles := 5
			if len(files) > maxFiles {
				files = files[:maxFiles]
			}

			for _, file := range files {
				if file.Size > 10*1024*1024 {
					return utils.ErrorResponse(c, fiber.StatusBadRequest, "File too large",
						fmt.Sprintf("File %s exceeds 10MB limit", file.Filename))
				}
				additionalFiles = append(additionalFiles, file)
			}
		}
	}

	if err := h.companyService.RequestVerification(
		ctx,
		companyID,
		employerUserID,
		*req.NPWPNumber,
		req.NIBNumber,
		npwpFile,
		additionalFiles,
	); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, common.ErrFailedOperation, err.Error())
	}

	return utils.SuccessResponse(c, "Verification request submitted successfully", fiber.Map{
		"company_id": companyID,
		"status":     "pending",
		"message":    "Your verification request has been submitted and will be reviewed by our admin team",
	})
}
