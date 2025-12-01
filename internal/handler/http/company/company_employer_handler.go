package companyhandler

import (
	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/domain/master"
	"keerja-backend/internal/domain/user"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/handler/http/common"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// CompanyEmployerHandler handles employer profile operations
type CompanyEmployerHandler struct {
	companyService company.CompanyService
	userService    user.UserService
	provinceRepo   master.ProvinceRepository
	cityRepo       master.CityRepository
	districtRepo   master.DistrictRepository
}

// NewCompanyEmployerHandler creates a new instance of CompanyEmployerHandler
func NewCompanyEmployerHandler(
	companyService company.CompanyService,
	userService user.UserService,
	provinceRepo master.ProvinceRepository,
	cityRepo master.CityRepository,
	districtRepo master.DistrictRepository,
) *CompanyEmployerHandler {
	return &CompanyEmployerHandler{
		companyService: companyService,
		userService:    userService,
		provinceRepo:   provinceRepo,
		cityRepo:       cityRepo,
		districtRepo:   districtRepo,
	}
}

func (h *CompanyEmployerHandler) GetMyEmployerProfile(c *fiber.Ctx) error {
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
		return utils.ErrorResponse(c, fiber.StatusNotFound, common.ErrCompanyNotFound, "User is not affiliated with any company")
	}

	companyID := companies[0].ID

	employerUser, err := h.companyService.GetEmployerUser(ctx, userID, companyID)
	if err != nil || employerUser == nil {
		return utils.ErrorResponse(c, fiber.StatusForbidden, common.ErrForbidden, "You are not an employer of this company")
	}

	usr, err := h.userService.GetProfile(ctx, userID)
	if err != nil || usr == nil {
		resp := fiber.Map{
			"user_id":          employerUser.UserID,
			"employer_user_id": employerUser.ID,
			"company_id":       employerUser.CompanyID,
			"role":             employerUser.Role,
			"position_title":   employerUser.PositionTitle,
		}
		return utils.SuccessResponse(c, common.MsgFetchedSuccess, resp)
	}

	fullName := ""
	var userObj interface{}
	if usr != nil {
		fullName = usr.FullName
		if usr.Profile != nil {
			loc := fiber.Map{}
			if usr.Profile.ProvinceID != nil {
				if p, err := h.provinceRepo.GetByID(ctx, *usr.Profile.ProvinceID); err == nil && p != nil {
					loc["province"] = fiber.Map{"id": p.ID, "code": p.Code, "name": p.Name}
				}
			}
			if usr.Profile.CityID != nil {
				if cObj, err := h.cityRepo.GetByID(ctx, *usr.Profile.CityID); err == nil && cObj != nil {
					loc["city"] = fiber.Map{"id": cObj.ID, "code": cObj.Code, "name": cObj.Name, "type": cObj.Type, "province_id": cObj.ProvinceID}
				}
			}
			if usr.Profile.DistrictID != nil {
				if d, err := h.districtRepo.GetByID(ctx, *usr.Profile.DistrictID); err == nil && d != nil {
					loc["district"] = fiber.Map{"id": d.ID, "code": d.Code, "name": d.Name, "city_id": d.CityID}
				}
			}
			if len(loc) > 0 {
				userObj = loc
			}
		}
	}

	type employerResp struct {
		ID             int64   `json:"id"`
		EmployerUserID int64   `json:"employer_user_id"`
		CompanyID      int64   `json:"company_id"`
		Role           string  `json:"role"`
		PositionTitle  *string `json:"position_title,omitempty"`
		FullName       string  `json:"full_name,omitempty"`
		User           any     `json:"user,omitempty"`
	}

	resp := employerResp{
		ID:             employerUser.UserID,
		EmployerUserID: employerUser.ID,
		CompanyID:      employerUser.CompanyID,
		Role:           employerUser.Role,
		PositionTitle:  employerUser.PositionTitle,
		FullName:       fullName,
		User:           userObj,
	}

	return utils.SuccessResponse(c, common.MsgFetchedSuccess, resp)
}

func (h *CompanyEmployerHandler) UpdateMyEmployerProfile(c *fiber.Ctx) error {
	ctx := c.Context()

	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, common.ErrUnauthorized, "userID not found in context")
	}

	var req request.UpdateEmployerUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidBody)
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, common.ErrValidationFailed, errs)
	}

	companies, err := h.companyService.GetUserCompanies(ctx, userID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, common.ErrFailedOperation)
	}
	if len(companies) == 0 {
		return utils.ErrorResponse(c, fiber.StatusNotFound, common.ErrCompanyNotFound, "User is not affiliated with any company")
	}

	companyID := companies[0].ID

	_, err = h.companyService.GetEmployerUser(ctx, userID, companyID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusForbidden, common.ErrForbidden, "You are not an employer of this company")
	}

	domainReq := &company.UpdateEmployerUserRequest{}
	if req.PositionTitle != nil {
		domainReq.PositionTitle = req.PositionTitle
	}

	var userReq *user.UpdateProfileRequest
	if req.Name != nil || req.ProvinceID != nil || req.CityID != nil || req.DistrictID != nil {
		userReq = &user.UpdateProfileRequest{}
		if req.Name != nil {
			userReq.FullName = req.Name
		}
		if req.ProvinceID != nil {
			userReq.ProvinceID = req.ProvinceID
		}
		if req.CityID != nil {
			userReq.CityID = req.CityID
		}
		if req.DistrictID != nil {
			userReq.DistrictID = req.DistrictID
		}
	}

	if err := h.companyService.UpdateEmployerUserWithProfile(ctx, userID, companyID, userReq, domainReq); err != nil {
		return utils.InternalServerErrorResponse(c, common.ErrFailedOperation)
	}

	updated, err := h.companyService.GetEmployerUser(ctx, userID, companyID)
	if err != nil || updated == nil {
		return utils.SuccessResponse(c, common.MsgUpdatedSuccess, nil)
	}

	usr, _ := h.userService.GetProfile(ctx, userID)

	var userObj interface{}
	if usr != nil && usr.Profile != nil {
		loc := fiber.Map{}
		if usr.Profile.ProvinceID != nil {
			if p, err := h.provinceRepo.GetByID(ctx, *usr.Profile.ProvinceID); err == nil && p != nil {
				loc["province"] = fiber.Map{"id": p.ID, "code": p.Code, "name": p.Name}
			}
		}
		if usr.Profile.CityID != nil {
			if cObj, err := h.cityRepo.GetByID(ctx, *usr.Profile.CityID); err == nil && cObj != nil {
				loc["city"] = fiber.Map{"id": cObj.ID, "code": cObj.Code, "name": cObj.Name, "type": cObj.Type, "province_id": cObj.ProvinceID}
			}
		}
		if usr.Profile.DistrictID != nil {
			if d, err := h.districtRepo.GetByID(ctx, *usr.Profile.DistrictID); err == nil && d != nil {
				loc["district"] = fiber.Map{"id": d.ID, "code": d.Code, "name": d.Name, "city_id": d.CityID}
			}
		}
		if len(loc) > 0 {
			userObj = loc
		}
	}

	type employerResp struct {
		ID             int64   `json:"id"`
		EmployerUserID int64   `json:"employer_user_id"`
		CompanyID      int64   `json:"company_id"`
		Role           string  `json:"role"`
		PositionTitle  *string `json:"position_title,omitempty"`
		FullName       string  `json:"full_name,omitempty"`
		User           any     `json:"user,omitempty"`
	}

	fullName := ""
	if usr != nil {
		fullName = usr.FullName
	}

	resp := employerResp{
		ID:             updated.UserID,
		EmployerUserID: updated.ID,
		CompanyID:      updated.CompanyID,
		Role:           updated.Role,
		PositionTitle:  updated.PositionTitle,
		FullName:       fullName,
		User:           userObj,
	}

	return utils.SuccessResponse(c, common.MsgUpdatedSuccess, resp)
}
