package companyhandler

import (
	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/domain/master"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/dto/response"
	"keerja-backend/internal/handler/http/common"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// CompanyAddressHandler handles company address CRUD operations
type CompanyAddressHandler struct {
	companyService company.CompanyService
	provinceRepo   master.ProvinceRepository
	cityRepo       master.CityRepository
	districtRepo   master.DistrictRepository
}

// NewCompanyAddressHandler creates a new instance of CompanyAddressHandler
func NewCompanyAddressHandler(
	companyService company.CompanyService,
	provinceRepo master.ProvinceRepository,
	cityRepo master.CityRepository,
	districtRepo master.DistrictRepository,
) *CompanyAddressHandler {
	return &CompanyAddressHandler{
		companyService: companyService,
		provinceRepo:   provinceRepo,
		cityRepo:       cityRepo,
		districtRepo:   districtRepo,
	}
}

func (h *CompanyAddressHandler) GetMyAddresses(c *fiber.Ctx) error {
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

	companyData := companies[0]

	includeDeleted := false
	if val := c.Query("include_deleted"); val == "true" {
		includeDeleted = true
	}

	addrs, err := h.companyService.GetCompanyAddresses(ctx, companyData.ID, includeDeleted)
	if err != nil {
		return utils.InternalServerErrorResponse(c, common.ErrFailedOperation)
	}

	responses := make([]interface{}, 0, len(addrs))
	for _, a := range addrs {
		addrResp := response.CompanyAddressResponse{
			ID:          a.ID,
			FullAddress: a.FullAddress,
		}
		if a.Latitude != nil {
			addrResp.Latitude = *a.Latitude
		}
		if a.Longitude != nil {
			addrResp.Longitude = *a.Longitude
		}
		if a.ProvinceID != nil {
			addrResp.ProvinceID = a.ProvinceID
		}
		if a.CityID != nil {
			addrResp.CityID = a.CityID
		}
		if a.DistrictID != nil {
			addrResp.DistrictID = a.DistrictID
		}

		var provResp *response.ProvinceResponse
		var cityResp *response.CityResponse
		var distResp *response.DistrictResponse

		if a.ProvinceID != nil {
			if p, err := h.provinceRepo.GetByID(ctx, *a.ProvinceID); err == nil && p != nil {
				provResp = &response.ProvinceResponse{ID: p.ID, Code: p.Code, Name: p.Name}
			}
		}
		if a.CityID != nil {
			if cobj, err := h.cityRepo.GetByID(ctx, *a.CityID); err == nil && cobj != nil {
				cityResp = &response.CityResponse{ID: cobj.ID, Code: cobj.Code, Name: cobj.Name, Type: cobj.Type, ProvinceID: cobj.ProvinceID}
			}
		}
		if a.DistrictID != nil {
			if d, err := h.districtRepo.GetWithFullLocation(ctx, *a.DistrictID); err == nil && d != nil {
				distResp = &response.DistrictResponse{ID: d.ID, Code: d.Code, Name: d.Name, CityID: d.CityID}
				if d.City != nil {
					cityResp = &response.CityResponse{ID: d.City.ID, Code: d.City.Code, Name: d.City.Name, Type: d.City.Type, ProvinceID: d.City.ProvinceID}
					if d.City.Province != nil {
						provResp = &response.ProvinceResponse{ID: d.City.Province.ID, Code: d.City.Province.Code, Name: d.City.Province.Name}
					}
				}
			}
		}

		responses = append(responses, addrResp.WithLocations(provResp, cityResp, distResp))
	}

	return utils.SuccessResponse(c, common.MsgFetchedSuccess, responses)
}

func (h *CompanyAddressHandler) CreateMyAddress(c *fiber.Ctx) error {
	ctx := c.Context()

	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, common.ErrUnauthorized, "userID not found in context")
	}

	var req request.CreateCompanyAddressRequest
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

	hasPermission, err := h.companyService.CheckEmployerPermission(ctx, userID, companyID, "admin")
	if err != nil || !hasPermission {
		return utils.ErrorResponse(c, fiber.StatusForbidden, common.ErrForbidden, "You don't have permission to create company addresses")
	}

	domainReq := &company.CreateCompanyAddressRequest{
		FullAddress: req.FullAddress,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		ProvinceID:  req.ProvinceID,
		CityID:      req.CityID,
		DistrictID:  req.DistrictID,
	}

	addr, err := h.companyService.CreateCompanyAddress(ctx, companyID, domainReq)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, common.ErrFailedOperation, err.Error())
	}

	resp := response.CompanyAddressResponse{
		ID:          addr.ID,
		FullAddress: addr.FullAddress,
	}
	if addr.Latitude != nil {
		resp.Latitude = *addr.Latitude
	}
	if addr.Longitude != nil {
		resp.Longitude = *addr.Longitude
	}

	return utils.CreatedResponse(c, common.MsgCreatedSuccess, resp)
}

func (h *CompanyAddressHandler) UpdateMyAddress(c *fiber.Ctx) error {
	ctx := c.Context()

	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, common.ErrUnauthorized, "userID not found in context")
	}

	addrID, err := utils.ParseIDParam(c, "id")
	if err != nil || addrID <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	var req request.UpdateCompanyAddressRequest
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

	hasPermission, err := h.companyService.CheckEmployerPermission(ctx, userID, companyID, "admin")
	if err != nil || !hasPermission {
		return utils.ErrorResponse(c, fiber.StatusForbidden, common.ErrForbidden, "You don't have permission to update company addresses")
	}

	domainReq := &company.UpdateCompanyAddressRequest{}
	if req.FullAddress != nil {
		domainReq.FullAddress = req.FullAddress
	}
	if req.Latitude != nil {
		domainReq.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		domainReq.Longitude = req.Longitude
	}
	if req.ProvinceID != nil {
		domainReq.ProvinceID = req.ProvinceID
	}
	if req.CityID != nil {
		domainReq.CityID = req.CityID
	}
	if req.DistrictID != nil {
		domainReq.DistrictID = req.DistrictID
	}

	updated, err := h.companyService.UpdateCompanyAddress(ctx, companyID, addrID, domainReq)
	if err != nil {
		return utils.InternalServerErrorResponse(c, common.ErrFailedOperation)
	}

	resp := response.CompanyAddressResponse{
		ID:          updated.ID,
		FullAddress: updated.FullAddress,
	}
	if updated.Latitude != nil {
		resp.Latitude = *updated.Latitude
	}
	if updated.Longitude != nil {
		resp.Longitude = *updated.Longitude
	}
	if updated.ProvinceID != nil {
		resp.ProvinceID = updated.ProvinceID
	}
	if updated.CityID != nil {
		resp.CityID = updated.CityID
	}
	if updated.DistrictID != nil {
		resp.DistrictID = updated.DistrictID
	}

	return utils.SuccessResponse(c, common.MsgUpdatedSuccess, resp)
}

func (h *CompanyAddressHandler) DeleteMyAddress(c *fiber.Ctx) error {
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

	hasPermission, err := h.companyService.CheckEmployerPermission(ctx, userID, companyID, "admin")
	if err != nil || !hasPermission {
		return utils.ErrorResponse(c, fiber.StatusForbidden, common.ErrForbidden, "You don't have permission to delete company addresses")
	}

	addrID, err := utils.ParseIDParam(c, "id")
	if err != nil || addrID <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidID)
	}

	if err := h.companyService.SoftDeleteCompanyAddress(ctx, companyID, addrID); err != nil {
		return utils.InternalServerErrorResponse(c, common.ErrFailedOperation)
	}

	return utils.SuccessResponse(c, common.MsgDeletedSuccess, nil)
}
