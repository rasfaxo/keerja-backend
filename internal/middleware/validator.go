package middleware

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// ValidateRequest validates request body against a struct with validation tags
func ValidateRequest(target interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse body into target struct
		if err := c.BodyParser(target); err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
		}

		// Validate struct
		if err := utils.ValidateStruct(target); err != nil {
			errors := utils.FormatValidationErrors(err)
			return utils.ValidationErrorResponse(c, "Validation failed", errors)
		}

		// Store validated data in context for handler to use
		c.Locals("validated_data", target)

		return c.Next()
	}
}

// ValidateQuery validates query parameters
func ValidateQuery(target interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse query into target struct
		if err := c.QueryParser(target); err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid query parameters", err.Error())
		}

		// Validate struct
		if err := utils.ValidateStruct(target); err != nil {
			errors := utils.FormatValidationErrors(err)
			return utils.ValidationErrorResponse(c, "Validation failed", errors)
		}

		// Store validated data in context
		c.Locals("validated_query", target)

		return c.Next()
	}
}

// ValidateParams validates path parameters
func ValidateParams(target interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse params into target struct
		if err := c.ParamsParser(target); err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid path parameters", err.Error())
		}

		// Validate struct
		if err := utils.ValidateStruct(target); err != nil {
			errors := utils.FormatValidationErrors(err)
			return utils.ValidationErrorResponse(c, "Validation failed", errors)
		}

		// Store validated data in context
		c.Locals("validated_params", target)

		return c.Next()
	}
}

// FileUploadConfig configures file upload validation
type FileUploadConfig struct {
	MaxFileSize       int64    // in bytes
	AllowedMimeTypes  []string // e.g., ["image/jpeg", "image/png", "application/pdf"]
	AllowedExtensions []string // e.g., [".jpg", ".png", ".pdf"]
	Required          bool
	FieldName         string // form field name, default is "file"
}

// ValidateFileUpload validates uploaded files
func ValidateFileUpload(config FileUploadConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Set default field name
		fieldName := config.FieldName
		if fieldName == "" {
			fieldName = "file"
		}

		// Get file from form
		fileHeader, err := c.FormFile(fieldName)
		if err != nil {
			if config.Required {
				return utils.ErrorResponse(c, fiber.StatusBadRequest, "File upload required", fmt.Sprintf("No file found with field name '%s'", fieldName))
			}
			// File is optional and not provided
			return c.Next()
		}

		// Validate file size
		if config.MaxFileSize > 0 && fileHeader.Size > config.MaxFileSize {
			maxSizeMB := float64(config.MaxFileSize) / (1024 * 1024)
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "File too large", fmt.Sprintf("File size exceeds maximum allowed size of %.2f MB", maxSizeMB))
		}

		// Validate MIME type
		if len(config.AllowedMimeTypes) > 0 {
			file, err := fileHeader.Open()
			if err != nil {
				return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to open file", err.Error())
			}
			defer file.Close()

			// Read first 512 bytes to detect content type
			buffer := make([]byte, 512)
			_, err = file.Read(buffer)
			if err != nil {
				return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to read file", err.Error())
			}

			contentType := fileHeader.Header.Get("Content-Type")
			if !isAllowedMimeType(contentType, config.AllowedMimeTypes) {
				return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid file type", fmt.Sprintf("File type '%s' is not allowed. Allowed types: %s", contentType, strings.Join(config.AllowedMimeTypes, ", ")))
			}
		}

		// Validate file extension
		if len(config.AllowedExtensions) > 0 {
			ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
			if !isAllowedExtension(ext, config.AllowedExtensions) {
				return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid file extension", fmt.Sprintf("File extension '%s' is not allowed. Allowed extensions: %s", ext, strings.Join(config.AllowedExtensions, ", ")))
			}
		}

		// Store file header in context for handler to use
		c.Locals("uploaded_file", fileHeader)

		return c.Next()
	}
}

// ValidateMultipleFiles validates multiple uploaded files
func ValidateMultipleFiles(config FileUploadConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Set default field name
		fieldName := config.FieldName
		if fieldName == "" {
			fieldName = "files"
		}

		// Get form
		form, err := c.MultipartForm()
		if err != nil {
			if config.Required {
				return utils.ErrorResponse(c, fiber.StatusBadRequest, "Files upload required", err.Error())
			}
			return c.Next()
		}

		// Get files from form
		files := form.File[fieldName]
		if len(files) == 0 {
			if config.Required {
				return utils.ErrorResponse(c, fiber.StatusBadRequest, "Files upload required", fmt.Sprintf("No files found with field name '%s'", fieldName))
			}
			return c.Next()
		}

		// Validate each file
		var validatedFiles []*multipart.FileHeader
		for _, fileHeader := range files {
			// Validate file size
			if config.MaxFileSize > 0 && fileHeader.Size > config.MaxFileSize {
				maxSizeMB := float64(config.MaxFileSize) / (1024 * 1024)
				return utils.ErrorResponse(c, fiber.StatusBadRequest, "File too large", fmt.Sprintf("File '%s' size exceeds maximum allowed size of %.2f MB", fileHeader.Filename, maxSizeMB))
			}

			// Validate MIME type
			if len(config.AllowedMimeTypes) > 0 {
				contentType := fileHeader.Header.Get("Content-Type")
				if !isAllowedMimeType(contentType, config.AllowedMimeTypes) {
					return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid file type", fmt.Sprintf("File '%s' type '%s' is not allowed. Allowed types: %s", fileHeader.Filename, contentType, strings.Join(config.AllowedMimeTypes, ", ")))
				}
			}

			// Validate file extension
			if len(config.AllowedExtensions) > 0 {
				ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
				if !isAllowedExtension(ext, config.AllowedExtensions) {
					return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid file extension", fmt.Sprintf("File '%s' extension '%s' is not allowed. Allowed extensions: %s", fileHeader.Filename, ext, strings.Join(config.AllowedExtensions, ", ")))
				}
			}

			validatedFiles = append(validatedFiles, fileHeader)
		}

		// Store validated files in context
		c.Locals("uploaded_files", validatedFiles)

		return c.Next()
	}
}

// Custom validation helpers

func isAllowedMimeType(mimeType string, allowedTypes []string) bool {
	for _, allowed := range allowedTypes {
		if strings.EqualFold(mimeType, allowed) {
			return true
		}
	}
	return false
}

func isAllowedExtension(extension string, allowedExtensions []string) bool {
	for _, allowed := range allowedExtensions {
		if strings.EqualFold(extension, allowed) {
			return true
		}
	}
	return false
}

// GetValidatedData retrieves validated data from context
func GetValidatedData(c *fiber.Ctx) interface{} {
	return c.Locals("validated_data")
}

// GetValidatedQuery retrieves validated query from context
func GetValidatedQuery(c *fiber.Ctx) interface{} {
	return c.Locals("validated_query")
}

// GetValidatedParams retrieves validated params from context
func GetValidatedParams(c *fiber.Ctx) interface{} {
	return c.Locals("validated_params")
}

// GetUploadedFile retrieves uploaded file from context
func GetUploadedFile(c *fiber.Ctx) *multipart.FileHeader {
	file, ok := c.Locals("uploaded_file").(*multipart.FileHeader)
	if !ok {
		return nil
	}
	return file
}

// GetUploadedFiles retrieves multiple uploaded files from context
func GetUploadedFiles(c *fiber.Ctx) []*multipart.FileHeader {
	files, ok := c.Locals("uploaded_files").([]*multipart.FileHeader)
	if !ok {
		return nil
	}
	return files
}
