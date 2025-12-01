package userhandler

import (
	"keerja-backend/internal/domain/user"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/handler/http/common"
	"keerja-backend/internal/helpers"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// UserMiscHandler handles miscellaneous user data (certifications, languages, projects)
type UserMiscHandler struct {
	userService user.UserService
}

// NewUserMiscHandler creates a new instance of UserMiscHandler
func NewUserMiscHandler(userService user.UserService) *UserMiscHandler {
	return &UserMiscHandler{
		userService: userService,
	}
}

func (h *UserMiscHandler) GetCertifications(c *fiber.Ctx) error {
	usr, err := helpers.GetProfile(c, h.userService)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, common.ErrFailedOperation, err.Error())
	}

	certifications := mapper.MapEntities(usr.Certifications, mapper.ToUserCertificationResponse)
	return utils.SuccessResponse(c, common.MsgOperationSuccess, certifications)
}

func (h *UserMiscHandler) GetLanguages(c *fiber.Ctx) error {
	usr, err := helpers.GetProfile(c, h.userService)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, common.ErrFailedOperation, err.Error())
	}

	languages := mapper.MapEntities(usr.Languages, mapper.ToUserLanguageResponse)
	return utils.SuccessResponse(c, common.MsgOperationSuccess, languages)
}

func (h *UserMiscHandler) GetProjects(c *fiber.Ctx) error {
	usr, err := helpers.GetProfile(c, h.userService)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, common.ErrFailedOperation, err.Error())
	}

	projects := mapper.MapEntities(usr.Projects, mapper.ToUserProjectResponse)
	return utils.SuccessResponse(c, common.MsgOperationSuccess, projects)
}
