package companyhandler

import (
	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/handler/http/common"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// CompanyImageHandler handles company image operations (logo and banner)
type CompanyImageHandler struct {
	companyService company.CompanyService
}

// NewCompanyImageHandler creates a new instance of CompanyImageHandler
func NewCompanyImageHandler(companyService company.CompanyService) *CompanyImageHandler {
	return &CompanyImageHandler{
		companyService: companyService,
	}
}

func (h *CompanyImageHandler) UploadLogo(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := utils.ParseIDParam(c, "id")
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidCompanyID)
	}

	file := middleware.GetUploadedFile(c)
	if file == nil {
		return utils.BadRequestResponse(c, common.ErrNoFileUploaded)
	}

	url, err := h.companyService.UploadLogo(ctx, companyID, file)
	if err != nil {
		return utils.InternalServerErrorResponse(c, common.ErrFileUploadFailed)
	}
	return utils.CreatedResponse(c, common.MsgUploadSuccess, fiber.Map{"logo_url": url})
}

func (h *CompanyImageHandler) UploadBanner(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := utils.ParseIDParam(c, "id")
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidCompanyID)
	}

	file := middleware.GetUploadedFile(c)
	if file == nil {
		return utils.BadRequestResponse(c, common.ErrNoFileUploaded)
	}

	url, err := h.companyService.UploadBanner(ctx, companyID, file)
	if err != nil {
		return utils.InternalServerErrorResponse(c, common.ErrFileUploadFailed)
	}
	return utils.CreatedResponse(c, common.MsgUploadSuccess, fiber.Map{"banner_url": url})
}

func (h *CompanyImageHandler) DeleteLogo(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := utils.ParseIDParam(c, "id")
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidCompanyID)
	}

	if err := h.companyService.DeleteLogo(ctx, int64(companyID)); err != nil {
		return utils.InternalServerErrorResponse(c, common.ErrFailedOperation)
	}
	return utils.SuccessResponse(c, common.MsgDeletedSuccess, nil)
}

func (h *CompanyImageHandler) DeleteBanner(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := utils.ParseIDParam(c, "id")
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidCompanyID)
	}

	if err := h.companyService.DeleteBanner(ctx, int64(companyID)); err != nil {
		return utils.InternalServerErrorResponse(c, common.ErrFailedOperation)
	}
	return utils.SuccessResponse(c, common.MsgDeletedSuccess, nil)
}
