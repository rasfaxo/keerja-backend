package userhandler

import (
	"strings"
	"time"

	"keerja-backend/internal/domain/user"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/handler/http/common"
	"keerja-backend/internal/helpers"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// UserDocumentHandler handles user document operations
type UserDocumentHandler struct {
	userService user.UserService
}

// NewUserDocumentHandler creates a new instance of UserDocumentHandler
func NewUserDocumentHandler(userService user.UserService) *UserDocumentHandler {
	return &UserDocumentHandler{
		userService: userService,
	}
}

func (h *UserDocumentHandler) GetDocuments(c *fiber.Ctx) error {
	usr, err := helpers.GetProfile(c, h.userService)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, common.ErrFailedOperation, err.Error())
	}

	documents := mapper.MapEntities(usr.Documents, mapper.ToUserDocumentResponse)
	return utils.SuccessResponse(c, common.MsgOperationSuccess, documents)
}

func (h *UserDocumentHandler) UploadDocument(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	file := middleware.GetUploadedFile(c)
	if file == nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "No file uploaded", "")
	}

	documentType := c.FormValue("document_type")
	if documentType == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Document type is required", "")
	}

	documentName := c.FormValue("document_name")
	if documentName == "" {
		documentName = file.Filename
	}
	description := c.FormValue("description")

	documentType = utils.SanitizeIfNonEmpty(documentType)
	documentName = utils.SanitizeIfNonEmpty(documentName)
	if description != "" {
		description = utils.SanitizeIfNonEmpty(description)
	}

	var descPtr *string
	if description != "" {
		descPtr = &description
	}

	domainReq := &user.UploadDocumentRequest{
		DocumentType: &documentType,
		DocumentName: documentName,
		Description:  descPtr,
	}

	document, err := h.userService.UploadDocument(ctx, userID, file, domainReq)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to upload document", err.Error())
	}

	return utils.CreatedResponse(c, "Document uploaded successfully", document)
}

func (h *UserDocumentHandler) DeleteDocument(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	documentID, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid document ID", "")
	}

	err = h.userService.DeleteDocument(ctx, userID, documentID)
	if err != nil {
		if strings.Contains(err.Error(), "not found or unauthorized") {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Document not found", "")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete document", err.Error())
	}

	response := map[string]interface{}{
		"id":        documentID,
		"deleted_at": time.Now().Format(time.RFC3339),
	}

	return utils.SuccessResponse(c, "Document deleted successfully", response)
}
