package http

import (
	"strconv"

	"keerja-backend/internal/domain/user"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService user.UserService
}

func NewUserHandler(userService user.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetProfile godoc
// @Summary Get current user profile
// @Description Get authenticated user's profile with all related data
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=response.UserProfileResponse}
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/me [get]
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	// Get user profile
	usr, err := h.userService.GetProfile(ctx, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get profile", err.Error())
	}

	// Convert to response DTO
	response := mapper.ToUserResponse(usr)

	return utils.SuccessResponse(c, "Profile retrieved successfully", response)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update authenticated user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.UpdateProfileRequest true "Update profile request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/me [put]
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	// Parse request body
	var req request.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Convert to domain request
	domainReq := &user.UpdateProfileRequest{
		Headline:           req.Headline,
		Bio:                req.Bio,
		Gender:             req.Gender,
		BirthDate:          req.BirthDate,
		LocationCity:       req.LocationCity,
		LocationCountry:    req.LocationCountry,
		DesiredPosition:    req.DesiredPosition,
		DesiredSalaryMin:   req.DesiredSalaryMin,
		DesiredSalaryMax:   req.DesiredSalaryMax,
		ExperienceLevel:    req.ExperienceLevel,
		IndustryInterest:   req.IndustryInterest,
		AvailabilityStatus: req.AvailabilityStatus,
	}

	// Update profile
	if err := h.userService.UpdateProfile(ctx, userID, domainReq); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update profile", err.Error())
	}

	return utils.SuccessResponse(c, "Profile updated successfully", nil)
}

// AddEducation godoc
// @Summary Add education
// @Description Add education to user profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.AddEducationRequest true "Add education request"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/me/education [post]
func (h *UserHandler) AddEducation(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	var req request.AddEducationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Convert to domain request
	domainReq := &user.AddEducationRequest{
		InstitutionName: req.InstitutionName,
		Major:           req.Major,
		DegreeLevel:     req.DegreeLevel,
		StartYear:       req.StartYear,
		EndYear:         req.EndYear,
		GPA:             req.GPA,
		Activities:      req.Activities,
		Description:     req.Description,
		IsCurrent:       req.IsCurrent,
	}

	education, err := h.userService.AddEducation(ctx, userID, domainReq)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to add education", err.Error())
	}

	response := mapper.ToEducationResponse(education)
	return utils.CreatedResponse(c, "Education added successfully", response)
}

// UpdateEducation godoc
// @Summary Update education
// @Description Update user's education record
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Education ID"
// @Param request body request.UpdateEducationRequest true "Update education request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/me/education/{id} [put]
func (h *UserHandler) UpdateEducation(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	educationID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid education ID", err.Error())
	}

	var req request.UpdateEducationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Convert to domain request
	domainReq := &user.UpdateEducationRequest{
		InstitutionName: req.InstitutionName,
		Major:           req.Major,
		DegreeLevel:     req.DegreeLevel,
		StartYear:       req.StartYear,
		EndYear:         req.EndYear,
		GPA:             req.GPA,
		Activities:      req.Activities,
		Description:     req.Description,
		IsCurrent:       req.IsCurrent,
	}

	if err := h.userService.UpdateEducation(ctx, userID, educationID, domainReq); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update education", err.Error())
	}

	return utils.SuccessResponse(c, "Education updated successfully", nil)
}

// DeleteEducation godoc
// @Summary Delete education
// @Description Delete user's education record
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Education ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/me/education/{id} [delete]
func (h *UserHandler) DeleteEducation(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	educationID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid education ID", err.Error())
	}

	if err := h.userService.DeleteEducation(ctx, userID, educationID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete education", err.Error())
	}

	return utils.SuccessResponse(c, "Education deleted successfully", nil)
}

// AddExperience godoc
// @Summary Add experience
// @Description Add work experience to user profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.AddExperienceRequest true "Add experience request"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/me/experience [post]
func (h *UserHandler) AddExperience(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	var req request.AddExperienceRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Convert to domain request
	domainReq := &user.AddExperienceRequest{
		CompanyName:     req.CompanyName,
		PositionTitle:   req.PositionTitle,
		Industry:        req.Industry,
		EmploymentType:  req.EmploymentType,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		IsCurrent:       req.IsCurrent,
		Description:     req.Description,
		Achievements:    req.Achievements,
		LocationCity:    req.LocationCity,
		LocationCountry: req.LocationCountry,
	}

	experience, err := h.userService.AddExperience(ctx, userID, domainReq)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to add experience", err.Error())
	}

	response := mapper.ToExperienceResponse(experience)
	return utils.CreatedResponse(c, "Experience added successfully", response)
}

// UpdateExperience godoc
// @Summary Update experience
// @Description Update user's work experience
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Experience ID"
// @Param request body request.UpdateExperienceRequest true "Update experience request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/me/experience/{id} [put]
func (h *UserHandler) UpdateExperience(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	experienceID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid experience ID", err.Error())
	}

	var req request.UpdateExperienceRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Convert to domain request
	domainReq := &user.UpdateExperienceRequest{
		CompanyName:     req.CompanyName,
		PositionTitle:   req.PositionTitle,
		Industry:        req.Industry,
		EmploymentType:  req.EmploymentType,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		IsCurrent:       req.IsCurrent,
		Description:     req.Description,
		Achievements:    req.Achievements,
		LocationCity:    req.LocationCity,
		LocationCountry: req.LocationCountry,
	}

	if err := h.userService.UpdateExperience(ctx, userID, experienceID, domainReq); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update experience", err.Error())
	}

	return utils.SuccessResponse(c, "Experience updated successfully", nil)
}

// DeleteExperience godoc
// @Summary Delete experience
// @Description Delete user's work experience
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Experience ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/me/experience/{id} [delete]
func (h *UserHandler) DeleteExperience(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	experienceID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid experience ID", err.Error())
	}

	if err := h.userService.DeleteExperience(ctx, userID, experienceID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete experience", err.Error())
	}

	return utils.SuccessResponse(c, "Experience deleted successfully", nil)
}

// AddSkill godoc
// @Summary Add skill
// @Description Add skills to user profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.AddUserSkillRequest true "Add skill request"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/me/skills [post]
func (h *UserHandler) AddSkill(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	var req struct {
		SkillNames []string `json:"skill_names" validate:"required,min=1,dive,min=1"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Add skills
	if err := h.userService.AddSkills(ctx, userID, req.SkillNames); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to add skills", err.Error())
	}

	return utils.CreatedResponse(c, "Skills added successfully", nil)
}

// DeleteSkill godoc
// @Summary Delete skill
// @Description Delete user's skill
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Skill ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/me/skills/{id} [delete]
func (h *UserHandler) DeleteSkill(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	skillID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid skill ID", err.Error())
	}

	if err := h.userService.DeleteSkill(ctx, userID, skillID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete skill", err.Error())
	}

	return utils.SuccessResponse(c, "Skill deleted successfully", nil)
}

// UploadDocument godoc
// @Summary Upload document
// @Description Upload document to user profile (resume, certificate, etc.)
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Document file"
// @Param document_type formData string true "Document type (resume, certificate, portfolio, etc.)"
// @Param title formData string false "Document title"
// @Param description formData string false "Document description"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/me/documents [post]
func (h *UserHandler) UploadDocument(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	// Get file from form
	file := middleware.GetUploadedFile(c)
	if file == nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "No file uploaded", "")
	}

	// Get form fields
	documentType := c.FormValue("document_type")
	if documentType == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Document type is required", "")
	}

	documentName := c.FormValue("document_name")
	if documentName == "" {
		documentName = file.Filename
	}
	description := c.FormValue("description")

	// Convert to domain request
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
