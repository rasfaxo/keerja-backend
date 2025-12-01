package userhandler

import (
	"strings"

	"keerja-backend/internal/domain/user"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/handler/http/common"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// UserProfileHandler handles user profile operations
type UserProfileHandler struct {
	userService user.UserService
}

// NewUserProfileHandler creates a new instance of UserProfileHandler
func NewUserProfileHandler(userService user.UserService) *UserProfileHandler {
	return &UserProfileHandler{
		userService: userService,
	}
}

func (h *UserProfileHandler) GetProfile(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	includeParam := c.Query("include", "")

	usr, err := h.userService.GetProfile(ctx, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get profile", err.Error())
	}

	var response any

	switch includeParam {
	case "":
		response = mapper.ToUserResponse(usr)
	case "all":
		response = mapper.ToUserDetailResponse(usr)
	default:
		includes := strings.Split(includeParam, ",")
		for i := range includes {
			includes[i] = strings.TrimSpace(includes[i])
		}
		response = mapper.ToUserResponseWithIncludes(usr, includes)
	}

	return utils.SuccessResponse(c, "Profile retrieved successfully", response)
}

func (h *UserProfileHandler) UpdateProfile(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	var req request.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, common.ErrInvalidRequest, err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, common.ErrValidationFailed, errors)
	}

	req.FullName = utils.SanitizePtr(req.FullName)
	req.Headline = utils.SanitizePtr(req.Headline)
	req.Bio = utils.SanitizePtr(req.Bio)
	req.LocationCity = utils.SanitizePtr(req.LocationCity)
	req.LocationCountry = utils.SanitizePtr(req.LocationCountry)
	req.DesiredPosition = utils.SanitizePtr(req.DesiredPosition)
	req.IndustryInterest = utils.SanitizePtr(req.IndustryInterest)

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

	if err := h.userService.UpdateProfile(ctx, userID, domainReq); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, common.ErrFailedOperation, err.Error())
	}

	return utils.SuccessResponse(c, common.MsgOperationSuccess, nil)
}

func (h *UserProfileHandler) GetPreferences(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	prefs, err := h.userService.GetPreferences(ctx, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve user preferences", err.Error())
	}

	response := mapper.ToUserPreferenceResponse(prefs)
	return utils.SuccessResponse(c, "User preferences retrieved successfully", response)
}

func (h *UserProfileHandler) UpdatePreferences(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	var req request.UpdateUserPreferencesRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, common.ErrInvalidRequest, err.Error())
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, common.ErrValidationFailed, errors)
	}

	req.LanguagePreference = utils.SanitizePtr(req.LanguagePreference)
	req.ThemePreference = utils.SanitizePtr(req.ThemePreference)
	req.PreferredJobType = utils.SanitizePtr(req.PreferredJobType)
	req.PreferredIndustry = utils.SanitizePtr(req.PreferredIndustry)
	req.PreferredLocation = utils.SanitizePtr(req.PreferredLocation)
	req.ProfileVisibility = utils.SanitizePtr(req.ProfileVisibility)

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
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, common.ErrFailedOperation, err.Error())
	}

	updated, err := h.userService.GetPreferences(ctx, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve updated preferences", err.Error())
	}

	resp := mapper.ToUserPreferenceResponse(updated)
	return utils.SuccessResponse(c, common.MsgOperationSuccess, resp)
}

func (h *UserProfileHandler) UploadProfilePhoto(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := middleware.GetUserID(c)

	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, common.ErrUnauthorized, "userID not found in context")
	}

	file := middleware.GetUploadedFile(c)
	if file == nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "No file uploaded", "")
	}

	avatarURL, err := h.userService.UploadAvatar(ctx, userID, file)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to upload profile photo", err.Error())
	}

	return utils.CreatedResponse(c, "Profile photo uploaded successfully", fiber.Map{"avatar_url": avatarURL})
}
