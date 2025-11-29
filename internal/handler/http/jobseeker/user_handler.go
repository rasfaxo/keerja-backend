package userhandler

import (
	"fmt"
	"strings"

	"keerja-backend/internal/domain/user"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/handler/http"
	"keerja-backend/internal/helpers"
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
// @Description Get authenticated user's profile. Use query param 'include' to fetch related data: 'all' for everything, or comma-separated list like 'educations,skills,projects'
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param include query string false "Include related data: all, educations, experiences, skills, certifications, languages, projects, documents, preference" example(educations,skills)
// @Success 200 {object} utils.Response
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/me [get]
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	// Parse include query parameter
	includeParam := c.Query("include", "")

	// Get user profile
	usr, err := h.userService.GetProfile(ctx, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get profile", err.Error())
	}

	// Determine response based on include parameter
	var response any

	switch includeParam {
	case "":
		// Default: Basic profile only (fast, minimal data)
		response = mapper.ToUserResponse(usr)
	case "all":
		// Include everything (full detail)
		response = mapper.ToUserDetailResponse(usr)
	default:
		// Include specific sections (custom includes)
		includes := strings.Split(includeParam, ",")
		// Trim spaces from each include
		for i := range includes {
			includes[i] = strings.TrimSpace(includes[i])
		}
		response = mapper.ToUserResponseWithIncludes(usr, includes)
	}

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
		return utils.ErrorResponse(c, fiber.StatusBadRequest, http.ErrInvalidRequest, err.Error())
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, http.ErrValidationFailed, errors)
	}

	// Sanitize text input fields
	req.FullName = utils.SanitizePtr(req.FullName)
	req.Headline = utils.SanitizePtr(req.Headline)
	req.Bio = utils.SanitizePtr(req.Bio)
	req.LocationCity = utils.SanitizePtr(req.LocationCity)
	req.LocationCountry = utils.SanitizePtr(req.LocationCountry)
	req.DesiredPosition = utils.SanitizePtr(req.DesiredPosition)
	req.IndustryInterest = utils.SanitizePtr(req.IndustryInterest)

	// Convert to domain request
	domainReq := &user.UpdateProfileRequest{
		FullName:           req.FullName,
		Phone:              req.Phone,
		Headline:           req.Headline,
		Bio:                req.Bio,
		Gender:             req.Gender,
		BirthDate:          req.BirthDate,
		Nationality:        req.Nationality,
		Address:            req.Address,
		LocationCity:       req.LocationCity,
		LocationState:      req.LocationState,
		LocationCountry:    req.LocationCountry,
		PostalCode:         req.PostalCode,
		LinkedinURL:        req.LinkedinURL,
		PortfolioURL:       req.PortfolioURL,
		GithubURL:          req.GithubURL,
		DesiredPosition:    req.DesiredPosition,
		DesiredSalaryMin:   req.DesiredSalaryMin,
		DesiredSalaryMax:   req.DesiredSalaryMax,
		ExperienceLevel:    req.ExperienceLevel,
		IndustryInterest:   req.IndustryInterest,
		AvailabilityStatus: req.AvailabilityStatus,
	}

	// Update profile
	if err := h.userService.UpdateProfile(ctx, userID, domainReq); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, http.ErrFailedOperation, err.Error())
	}

	return utils.SuccessResponse(c, http.MsgOperationSuccess, nil)
}

// GetPreferences godoc
// @Summary Get user preferences
// @Description Get authenticated user's preferences
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/me/preferences [get]
func (h *UserHandler) GetPreferences(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	prefs, err := h.userService.GetPreferences(ctx, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve user preferences", err.Error())
	}

	response := mapper.ToUserPreferenceResponse(prefs)
	return utils.SuccessResponse(c, "User preferences retrieved successfully", response)
}

// UpdatePreferences godoc
// @Summary Update user preferences
// @Description Update authenticated user's preferences
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.UpdateUserPreferencesRequest true "Update preferences request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/me/preferences [put]
func (h *UserHandler) UpdatePreferences(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	var req request.UpdateUserPreferencesRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, http.ErrInvalidRequest, err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, http.ErrValidationFailed, errors)
	}

	// Sanitize pointers
	req.LanguagePreference = utils.SanitizePtr(req.LanguagePreference)
	req.ThemePreference = utils.SanitizePtr(req.ThemePreference)
	req.PreferredJobType = utils.SanitizePtr(req.PreferredJobType)
	req.PreferredIndustry = utils.SanitizePtr(req.PreferredIndustry)
	req.PreferredLocation = utils.SanitizePtr(req.PreferredLocation)
	req.ProfileVisibility = utils.SanitizePtr(req.ProfileVisibility)

	// Map to domain request
	domainReq := &user.UpdatePreferenceRequest{
		LanguagePreference:  req.LanguagePreference,
		ThemePreference:     req.ThemePreference,
		PreferredJobType:    req.PreferredJobType,
		PreferredIndustry:   req.PreferredIndustry,
		PreferredLocation:   req.PreferredLocation,
		PreferredSalaryMin:  req.PreferredSalaryMin,
		PreferredSalaryMax:  req.PreferredSalaryMax,
		EmailNotifications:  req.EmailNotifications,
		SMSNotifications:    req.SMSNotifications,
		PushNotifications:   req.PushNotifications,
		EmailMarketing:      req.EmailMarketing,
		ProfileVisibility:   req.ProfileVisibility,
		ShowOnlineStatus:    req.ShowOnlineStatus,
		AllowDirectMessages: req.AllowDirectMessages,
		DataSharingConsent:  req.DataSharingConsent,
	}

	if err := h.userService.UpdatePreferences(ctx, userID, domainReq); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, http.ErrFailedOperation, err.Error())
	}

	// Retrieve updated preferences
	updated, err := h.userService.GetPreferences(ctx, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve updated preferences", err.Error())
	}

	resp := mapper.ToUserPreferenceResponse(updated)
	return utils.SuccessResponse(c, http.MsgOperationSuccess, resp)
}

// GetEducations godoc
// @Summary Get user educations
// @Description Get list of user's educations
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/me/educations [get]
func (h *UserHandler) GetEducations(c *fiber.Ctx) error {
	usr, err := helpers.GetProfile(c, h.userService)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, http.ErrFailedOperation, err.Error())
	}

	educations := mapper.MapEntities(usr.Educations, mapper.ToUserEducationResponse)
	return utils.SuccessResponse(c, http.MsgOperationSuccess, educations)
}

// GetExperiences godoc
// @Summary Get user experiences
// @Description Get list of user's work experiences
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/me/experiences [get]
func (h *UserHandler) GetExperiences(c *fiber.Ctx) error {
	usr, err := helpers.GetProfile(c, h.userService)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, http.ErrFailedOperation, err.Error())
	}

	experiences := mapper.MapEntities(usr.Experiences, mapper.ToUserExperienceResponse)
	return utils.SuccessResponse(c, http.MsgOperationSuccess, experiences)
}

// GetSkills godoc
// @Summary Get user skills
// @Description Get list of user's skills
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]response.UserSkillResponse}
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/me/skills [get]
func (h *UserHandler) GetSkills(c *fiber.Ctx) error {
	usr, err := helpers.GetProfile(c, h.userService)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, http.ErrFailedOperation, err.Error())
	}

	skills := mapper.MapEntities(usr.Skills, mapper.ToUserSkillResponse)
	return utils.SuccessResponse(c, http.MsgOperationSuccess, skills)
}

// GetCertifications godoc
// @Summary Get user certifications
// @Description Get list of user's certifications
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]response.UserCertificationResponse}
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/me/certifications [get]
func (h *UserHandler) GetCertifications(c *fiber.Ctx) error {
	usr, err := helpers.GetProfile(c, h.userService)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, http.ErrFailedOperation, err.Error())
	}

	certifications := mapper.MapEntities(usr.Certifications, mapper.ToUserCertificationResponse)
	return utils.SuccessResponse(c, http.MsgOperationSuccess, certifications)
}

// GetLanguages godoc
// @Summary Get user languages
// @Description Get list of user's languages
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]response.UserLanguageResponse}
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/me/languages [get]
func (h *UserHandler) GetLanguages(c *fiber.Ctx) error {
	usr, err := helpers.GetProfile(c, h.userService)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, http.ErrFailedOperation, err.Error())
	}

	languages := mapper.MapEntities(usr.Languages, mapper.ToUserLanguageResponse)
	return utils.SuccessResponse(c, http.MsgOperationSuccess, languages)
}

// GetProjects godoc
// @Summary Get user projects
// @Description Get list of user's projects/portfolio
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]response.UserProjectResponse}
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/me/projects [get]
func (h *UserHandler) GetProjects(c *fiber.Ctx) error {
	usr, err := helpers.GetProfile(c, h.userService)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, http.ErrFailedOperation, err.Error())
	}

	projects := mapper.MapEntities(usr.Projects, mapper.ToUserProjectResponse)
	return utils.SuccessResponse(c, http.MsgOperationSuccess, projects)
}

// GetDocuments godoc
// @Summary Get user documents
// @Description Get list of user's documents (CV, portfolio, etc)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]response.UserDocumentResponse}
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/me/documents [get]
func (h *UserHandler) GetDocuments(c *fiber.Ctx) error {
	usr, err := helpers.GetProfile(c, h.userService)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, http.ErrFailedOperation, err.Error())
	}

	documents := mapper.MapEntities(usr.Documents, mapper.ToUserDocumentResponse)
	return utils.SuccessResponse(c, http.MsgOperationSuccess, documents)
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
		return utils.ErrorResponse(c, fiber.StatusBadRequest, http.ErrInvalidRequest, err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, http.ErrValidationFailed, errors)
	}

	// Sanitize text input fields
	req.InstitutionName = utils.SanitizeIfNonEmpty(req.InstitutionName)
	req.Major = utils.SanitizePtr(req.Major)
	req.DegreeLevel = utils.SanitizePtr(req.DegreeLevel)
	req.Activities = utils.SanitizePtr(req.Activities)
	req.Description = utils.SanitizePtr(req.Description)

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
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, http.ErrFailedOperation, err.Error())
	}

	response := mapper.ToEducationResponse(education)
	return utils.CreatedResponse(c, http.MsgOperationSuccess, response)
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

	educationID, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, http.ErrInvalidRequest, err.Error())
	}

	var req request.UpdateEducationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, http.ErrInvalidRequest, err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, http.ErrValidationFailed, errors)
	}

	// Sanitize text input fields
	req.InstitutionName = utils.SanitizePtr(req.InstitutionName)
	req.Major = utils.SanitizePtr(req.Major)
	req.DegreeLevel = utils.SanitizePtr(req.DegreeLevel)
	req.Activities = utils.SanitizePtr(req.Activities)
	req.Description = utils.SanitizePtr(req.Description)

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

	educationID, err := utils.ParseIDParam(c, "id")
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

	// Sanitize text input fields
	req.CompanyName = utils.SanitizeIfNonEmpty(req.CompanyName)
	req.PositionTitle = utils.SanitizeIfNonEmpty(req.PositionTitle)
	req.Industry = utils.SanitizePtr(req.Industry)
	req.EmploymentType = utils.SanitizePtr(req.EmploymentType)
	req.Description = utils.SanitizePtr(req.Description)
	req.Achievements = utils.SanitizePtr(req.Achievements)
	req.LocationCity = utils.SanitizePtr(req.LocationCity)
	req.LocationCountry = utils.SanitizePtr(req.LocationCountry)

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

	experienceID, err := utils.ParseIDParam(c, "id")
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

	// Sanitize text input fields
	req.CompanyName = utils.SanitizePtr(req.CompanyName)
	req.PositionTitle = utils.SanitizePtr(req.PositionTitle)
	req.Industry = utils.SanitizePtr(req.Industry)
	req.EmploymentType = utils.SanitizePtr(req.EmploymentType)
	req.Description = utils.SanitizePtr(req.Description)
	req.Achievements = utils.SanitizePtr(req.Achievements)
	req.LocationCity = utils.SanitizePtr(req.LocationCity)
	req.LocationCountry = utils.SanitizePtr(req.LocationCountry)

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

	experienceID, err := utils.ParseIDParam(c, "id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid experience ID", err.Error())
	}

	if err := h.userService.DeleteExperience(ctx, userID, experienceID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete experience", err.Error())
	}

	return utils.SuccessResponse(c, "Experience deleted successfully", nil)
}

// AddSkills godoc
// @Summary Add multiple skills
// @Description Add multiple skills to user profile in batch
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.AddUserSkillsRequest true "Add Skills Request"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/me/skills/batch [post]
func (h *UserHandler) AddSkills(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	var req request.AddUserSkillsRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, "Validation failed", errors)
	}

	// Convert DTO request to domain request
	domainSkills := make([]user.AddUserSkillRequest, len(req.Skills))
	for i, skill := range req.Skills {
		// Validate: each skill must have either skill_id or skill_name
		if skill.SkillID == nil && skill.SkillName == "" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed",
				fmt.Sprintf("Skill #%d: Either skill_id or skill_name must be provided", i+1))
		}

		domainSkills[i] = user.AddUserSkillRequest{
			SkillID:           skill.SkillID,
			SkillName:         skill.SkillName,
			ProficiencyLevel:  skill.ProficiencyLevel,
			YearsOfExperience: skill.YearsOfExperience,
		}
	}

	domainReq := &user.AddUserSkillsRequest{
		Skills: domainSkills,
	}

	addedSkills, err := h.userService.AddSkills(ctx, userID, domainReq)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to add skills", err.Error())
	}

	// Format response with IDs
	skillsResponse := make([]map[string]any, len(addedSkills))
	for i, skill := range addedSkills {
		skillsResponse[i] = map[string]any{
			"id":               skill.ID,
			"skill_name":       skill.SkillName,
			"skill_level":      skill.SkillLevel,
			"years_experience": skill.YearsExperience,
		}
	}

	return utils.CreatedResponse(c,
		fmt.Sprintf("Successfully added %d skills", len(addedSkills)),
		map[string]any{
			"skills": skillsResponse,
			"total":  len(addedSkills),
		},
	)
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

	skillID, err := utils.ParseIDParam(c, "id")
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

	// Sanitize text input fields
	documentType = utils.SanitizeIfNonEmpty(documentType)
	documentName = utils.SanitizeIfNonEmpty(documentName)
	if description != "" {
		description = utils.SanitizeIfNonEmpty(description)
	}

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

// UploadProfilePhoto godoc
// @Summary Upload user profile photo
// @Description Upload a profile photo for the authenticated user
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Profile photo file (jpg, png, webp)"
// @Success 201 {object} utils.Response{data=fiber.Map} "Photo uploaded successfully with avatar_url"
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/profile-photo [post]
func (h *UserHandler) UploadProfilePhoto(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, http.ErrUnauthorized, "userID not found in context")
	}

	// Get file from form
	file := middleware.GetUploadedFile(c)
	if file == nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "No file uploaded", "")
	}

	// Upload profile photo using existing UploadAvatar service
	avatarURL, err := h.userService.UploadAvatar(ctx, userID, file)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to upload profile photo", err.Error())
	}

	return utils.CreatedResponse(c, "Profile photo uploaded successfully", fiber.Map{"avatar_url": avatarURL})
}
