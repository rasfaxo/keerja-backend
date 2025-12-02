package companyhandler

import (
	"mime/multipart"
	"strconv"
	"strings"

	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/domain/master"
	"keerja-backend/internal/dto/mapper"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/dto/response"
	"keerja-backend/internal/handler/http/common"
	"keerja-backend/internal/helpers"
	"keerja-backend/internal/middleware"
	"keerja-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// CompanyBasicHandler handles basic CRUD operations for companies
type CompanyBasicHandler struct {
	companyService  company.CompanyService
	industryRepo    master.IndustryRepository
	companySizeRepo master.CompanySizeRepository
	provinceRepo    master.ProvinceRepository
	cityRepo        master.CityRepository
	districtRepo    master.DistrictRepository
}

// NewCompanyBasicHandler creates a new instance of CompanyBasicHandler
func NewCompanyBasicHandler(
	companyService company.CompanyService,
	industryRepo master.IndustryRepository,
	companySizeRepo master.CompanySizeRepository,
	provinceRepo master.ProvinceRepository,
	cityRepo master.CityRepository,
	districtRepo master.DistrictRepository,
) *CompanyBasicHandler {
	return &CompanyBasicHandler{
		companyService:  companyService,
		industryRepo:    industryRepo,
		companySizeRepo: companySizeRepo,
		provinceRepo:    provinceRepo,
		cityRepo:        cityRepo,
		districtRepo:    districtRepo,
	}
}

func (h *CompanyBasicHandler) ListCompanies(c *fiber.Ctx) error {
	ctx := c.Context()

	var q request.CompanySearchRequest
	if err := c.QueryParser(&q); err != nil {
		return utils.BadRequestResponse(c, common.ErrInvalidQueryParams)
	}
	if err := utils.ValidateStruct(&q); err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, common.ErrValidationFailed, errs)
	}
	q.Page, q.Limit = utils.ValidatePagination(q.Page, q.Limit, 100)

	q.Query = utils.SanitizeString(q.Query)
	q.Location = utils.SanitizeString(q.Location)

	filter := helpers.BuildCompanyFilter(q)

	var (
		companies []company.Company
		total     int64
		err       error
	)
	if q.Query != "" {
		companies, total, err = h.companyService.SearchCompanies(ctx, q.Query, filter)
	} else {
		companies, total, err = h.companyService.ListCompanies(ctx, filter)
	}
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, common.ErrFailedOperation, err.Error())
	}

	respList := make([]response.CompanyResponse, 0, len(companies))
	for _, comp := range companies {
		cr := mapper.ToCompanyResponse(&comp)
		if cr != nil {
			respList = append(respList, *cr)
		}
	}

	meta := utils.GetPaginationMeta(q.Page, q.Limit, total)
	payload := response.CompanyListResponse{Companies: respList}
	return utils.SuccessResponseWithMeta(c, common.MsgFetchedSuccess, payload, meta)
}

func (h *CompanyBasicHandler) CreateCompany(c *fiber.Ctx) error {
	ctx := c.Context()

	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, common.ErrUnauthorized, "userID not found in context")
	}

	existingCompanies, err := h.companyService.GetUserCompanies(ctx, userID)
	if err == nil && len(existingCompanies) > 0 {
		return utils.ErrorResponse(c, fiber.StatusForbidden,
			common.ErrForbidden,
			"Business rule violation: Each user can only register one company. You already own a company.")
	}

	var req request.RegisterCompanyRequest

	contentType := string(c.Request().Header.ContentType())
	isMultipart := strings.Contains(contentType, "multipart/form-data")

	if isMultipart {
		req.CompanyName = c.FormValue("company_name")

		if industryID := c.FormValue("industry_id"); industryID != "" {
			if id, err := strconv.ParseInt(industryID, 10, 64); err == nil {
				req.IndustryID = &id
			}
		}
		if companySizeID := c.FormValue("company_size_id"); companySizeID != "" {
			if id, err := strconv.ParseInt(companySizeID, 10, 64); err == nil {
				req.CompanySizeID = &id
			}
		}
		if districtID := c.FormValue("district_id"); districtID != "" {
			if id, err := strconv.ParseInt(districtID, 10, 64); err == nil {
				req.DistrictID = &id
			}
		}
		if provinceID := c.FormValue("province_id"); provinceID != "" {
			if id, err := strconv.ParseInt(provinceID, 10, 64); err == nil {
				req.ProvinceID = &id
			}
		}
		if cityID := c.FormValue("city_id"); cityID != "" {
			if id, err := strconv.ParseInt(cityID, 10, 64); err == nil {
				req.CityID = &id
			}
		}

		if industryName := c.FormValue("industry_name"); industryName != "" {
			req.IndustryName = &industryName
		}
		if companySizeName := c.FormValue("company_size_name"); companySizeName != "" {
			req.CompanySizeName = &companySizeName
		}
		if provinceName := c.FormValue("province_name"); provinceName != "" {
			req.ProvinceName = &provinceName
		}
		if cityName := c.FormValue("city_name"); cityName != "" {
			req.CityName = &cityName
		}
		if districtName := c.FormValue("district_name"); districtName != "" {
			req.DistrictName = &districtName
		}

		if fullAddress := c.FormValue("full_address"); fullAddress != "" {
			req.FullAddress = fullAddress
		}
		if latitude := c.FormValue("latitude"); latitude != "" {
			if lat, err := strconv.ParseFloat(latitude, 64); err == nil {
				req.Latitude = &lat
			}
		}
		if longitude := c.FormValue("longitude"); longitude != "" {
			if lon, err := strconv.ParseFloat(longitude, 64); err == nil {
				req.Longitude = &lon
			}
		}
		if description := c.FormValue("description"); description != "" {
			req.Description = &description
		}
		if legalName := c.FormValue("legal_name"); legalName != "" {
			req.LegalName = &legalName
		}
		if registrationNumber := c.FormValue("registration_number"); registrationNumber != "" {
			req.RegistrationNumber = &registrationNumber
		}
		if websiteURL := c.FormValue("website_url"); websiteURL != "" {
			req.WebsiteURL = &websiteURL
		}
		if phone := c.FormValue("phone"); phone != "" {
			req.Phone = &phone
		}
		if about := c.FormValue("about"); about != "" {
			req.About = &about
		}
	} else {
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, common.ErrInvalidBody)
		}
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, common.ErrValidationFailed, errors)
	}

	req.CompanyName = utils.SanitizeString(req.CompanyName)

	if req.LegalName != nil {
		sanitized := utils.SanitizeString(*req.LegalName)
		req.LegalName = &sanitized
	}
	if req.About != nil {
		sanitized := utils.SanitizeHTML(*req.About)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, common.ErrPotentialXSS)
		}
		req.About = &sanitized
	}

	if req.FullAddress != "" {
		req.FullAddress = utils.SanitizeString(req.FullAddress)
	}

	if req.Description != nil {
		sanitized := utils.SanitizeHTML(*req.Description)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, common.ErrPotentialXSS)
		}
		req.Description = &sanitized
	}

	if req.Address != nil {
		sanitized := utils.SanitizeString(*req.Address)
		req.Address = &sanitized
	}
	if req.City != nil {
		sanitized := utils.SanitizeString(*req.City)
		req.City = &sanitized
	}
	if req.Province != nil {
		sanitized := utils.SanitizeString(*req.Province)
		req.Province = &sanitized
	}
	if req.Country != nil {
		sanitized := utils.SanitizeString(*req.Country)
		req.Country = &sanitized
	}

	if req.IndustryName != nil && *req.IndustryName != "" {
		industry, err := h.industryRepo.GetByName(ctx, *req.IndustryName)
		if err != nil {
			return utils.BadRequestResponse(c, "Industry not found: "+*req.IndustryName)
		}
		req.IndustryID = &industry.ID
	}

	if req.CompanySizeName != nil && *req.CompanySizeName != "" {
		companySize, err := h.companySizeRepo.GetByCategory(ctx, *req.CompanySizeName)
		if err != nil {
			return utils.BadRequestResponse(c, "Company size not found: "+*req.CompanySizeName)
		}
		req.CompanySizeID = &companySize.ID
	}

	if req.ProvinceName != nil && *req.ProvinceName != "" {
		province, err := h.provinceRepo.GetByName(ctx, *req.ProvinceName)
		if err != nil {
			return utils.BadRequestResponse(c, "Province not found: "+*req.ProvinceName)
		}

		if req.CityName != nil && *req.CityName != "" {
			city, err := h.cityRepo.GetByNameAndProvinceID(ctx, *req.CityName, province.ID)
			if err != nil {
				return utils.BadRequestResponse(c, "City not found: "+*req.CityName+" in province: "+*req.ProvinceName)
			}

			if req.DistrictName != nil && *req.DistrictName != "" {
				district, err := h.districtRepo.GetByNameAndCityID(ctx, *req.DistrictName, city.ID)
				if err != nil {
					return utils.BadRequestResponse(c, "District not found: "+*req.DistrictName+" in city: "+*req.CityName)
				}
				req.DistrictID = &district.ID
			}
		}
	}

	domainReq := &company.RegisterCompanyRequest{
		CompanyName:        req.CompanyName,
		LegalName:          req.LegalName,
		RegistrationNumber: req.RegistrationNumber,
		IndustryID:         req.IndustryID,
		CompanySizeID:      req.CompanySizeID,
		DistrictID:         req.DistrictID,
		FullAddress:        req.FullAddress,
		Description:        req.Description,
		Industry:           req.Industry,
		CompanyType:        req.CompanyType,
		SizeCategory:       req.SizeCategory,
		Address:            req.Address,
		City:               req.City,
		Province:           req.Province,
		WebsiteURL:         req.WebsiteURL,
		EmailDomain:        req.EmailDomain,
		Phone:              req.Phone,
		Country:            req.Country,
		PostalCode:         req.PostalCode,
		About:              req.About,
	}

	createdCompany, err := h.companyService.RegisterCompany(ctx, domainReq, userID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, common.ErrFailedOperation)
	}

	if req.FullAddress != "" || req.DistrictID != nil || req.ProvinceID != nil || req.CityID != nil {
		provinceID := req.ProvinceID
		cityID := req.CityID
		districtID := req.DistrictID

		if districtID != nil {
			district, err := h.districtRepo.GetByID(ctx, *districtID)
			if err == nil && district != nil {
				city, err := h.cityRepo.GetByID(ctx, district.CityID)
				if err == nil && city != nil {
					cityID = &city.ID
					province, err := h.provinceRepo.GetByID(ctx, city.ProvinceID)
					if err == nil && province != nil {
						provinceID = &province.ID
					}
				}
			}
		}

		_, _ = h.companyService.CreateCompanyAddress(ctx, createdCompany.ID, &company.CreateCompanyAddressRequest{
			FullAddress: req.FullAddress,
			ProvinceID:  provinceID,
			CityID:      cityID,
			DistrictID:  districtID,
			Latitude:    req.Latitude,
			Longitude:   req.Longitude,
		})
	}

	if isMultipart {
		if logoFile, err := c.FormFile("logo"); err == nil && logoFile != nil {
			logoURL, uploadErr := h.companyService.UploadLogo(ctx, createdCompany.ID, logoFile)
			if uploadErr == nil && logoURL != "" {
				createdCompany.LogoURL = &logoURL
			}
		}
	}

	resp := mapper.ToCompanyDetailResponse(createdCompany)
	addrs, err := h.companyService.GetCompanyAddresses(ctx, createdCompany.ID, false)
	if err == nil && len(addrs) > 0 {
		resp.CompanyAddresses = make([]response.CompanyAddressResponse, len(addrs))
		for i, a := range addrs {
			lat := 0.0
			lon := 0.0
			if a.Latitude != nil {
				lat = *a.Latitude
			}
			if a.Longitude != nil {
				lon = *a.Longitude
			}
			resp.CompanyAddresses[i] = response.CompanyAddressResponse{
				ID:          a.ID,
				FullAddress: a.FullAddress,
				Latitude:    lat,
				Longitude:   lon,
				ProvinceID:  a.ProvinceID,
				CityID:      a.CityID,
				DistrictID:  a.DistrictID,
			}
		}
	}
	return utils.CreatedResponse(c, common.MsgCreatedSuccess, resp)
}

func (h *CompanyBasicHandler) GetCompany(c *fiber.Ctx) error {
	ctx := c.Context()

	companyID, err := utils.ParseIDParam(c, "id")
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, "Invalid company ID")
	}

	companyData, err := h.companyService.GetCompany(ctx, companyID)
	if err != nil {
		return utils.NotFoundResponse(c, common.ErrNotFound)
	}

	response := mapper.ToCompanyResponse(companyData)
	return utils.SuccessResponse(c, common.MsgFetchedSuccess, response)
}

func (h *CompanyBasicHandler) GetCompanyBySlug(c *fiber.Ctx) error {
	ctx := c.Context()
	slug := utils.SanitizeString(strings.TrimSpace(c.Params("slug")))
	if slug == "" {
		return utils.BadRequestResponse(c, common.ErrInvalidRequest)
	}

	companyData, err := h.companyService.GetCompanyBySlug(ctx, slug)
	if err != nil {
		return utils.NotFoundResponse(c, common.ErrNotFound)
	}
	responseDTO := mapper.ToCompanyResponse(companyData)
	return utils.SuccessResponse(c, common.MsgFetchedSuccess, responseDTO)
}

func (h *CompanyBasicHandler) UpdateCompany(c *fiber.Ctx) error {
	ctx := c.Context()

	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, common.ErrUnauthorized, "userID not found in context")
	}

	companyID, err := utils.ParseIDParam(c, "id")
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, "Invalid company ID")
	}

	hasPermission, err := h.companyService.CheckEmployerPermission(ctx, userID, int64(companyID), "admin")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to check user permission", err.Error())
	}
	if !hasPermission {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "You don't have permission to update this company. Only company owner or admin company can perform this action.", "")
	}

	form, err := c.MultipartForm()
	if err != nil {
		var req request.UpdateCompanyRequest
		if err := c.BodyParser(&req); err != nil {
			return utils.BadRequestResponse(c, common.ErrInvalidRequest)
		}

		if err := utils.ValidateStruct(&req); err != nil {
			errors := utils.FormatValidationErrors(err)
			return utils.ValidationErrorResponse(c, common.ErrValidationFailed, errors)
		}

		domainReq := &company.UpdateCompanyRequest{
			FullAddress:        req.FullAddress,
			ShortDescription:   req.ShortDescription,
			WebsiteURL:         req.WebsiteURL,
			InstagramURL:       req.InstagramURL,
			FacebookURL:        req.FacebookURL,
			LinkedinURL:        req.LinkedinURL,
			TwitterURL:         req.TwitterURL,
			CompanyDescription: req.CompanyDescription,
			CompanyCulture:     req.CompanyCulture,
		}

		if err := h.companyService.UpdateCompany(ctx, companyID, domainReq, nil, nil); err != nil {
			return utils.InternalServerErrorResponse(c, common.ErrFailedOperation)
		}

		return utils.SuccessResponse(c, common.MsgUpdatedSuccess, nil)
	}

	fullAddress := c.FormValue("full_address")
	shortDescription := c.FormValue("short_description")
	websiteURL := c.FormValue("website_url")
	instagramURL := c.FormValue("instagram_url")
	facebookURL := c.FormValue("facebook_url")
	linkedinURL := c.FormValue("linkedin_url")
	twitterURL := c.FormValue("twitter_url")
	companyDescription := c.FormValue("company_description")
	companyCulture := c.FormValue("company_culture")

	req := &request.UpdateCompanyRequest{}

	if fullAddress != "" {
		req.FullAddress = &fullAddress
	}
	if shortDescription != "" {
		req.ShortDescription = &shortDescription
	}
	if websiteURL != "" {
		req.WebsiteURL = &websiteURL
	}
	if instagramURL != "" {
		req.InstagramURL = &instagramURL
	}
	if facebookURL != "" {
		req.FacebookURL = &facebookURL
	}
	if linkedinURL != "" {
		req.LinkedinURL = &linkedinURL
	}
	if twitterURL != "" {
		req.TwitterURL = &twitterURL
	}
	if companyDescription != "" {
		req.CompanyDescription = &companyDescription
	}
	if companyCulture != "" {
		req.CompanyCulture = &companyCulture
	}

	if err := utils.ValidateStruct(req); err != nil {
		errors := utils.FormatValidationErrors(err)
		return utils.ValidationErrorResponse(c, common.ErrValidationFailed, errors)
	}

	if req.FullAddress != nil {
		sanitized := utils.SanitizeString(*req.FullAddress)
		req.FullAddress = &sanitized
	}
	if req.ShortDescription != nil {
		sanitized := utils.SanitizeString(*req.ShortDescription)
		req.ShortDescription = &sanitized
	}
	if req.CompanyDescription != nil {
		sanitized := utils.SanitizeHTML(*req.CompanyDescription)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, common.ErrPotentialXSS)
		}
		req.CompanyDescription = &sanitized
	}
	if req.CompanyCulture != nil {
		sanitized := utils.SanitizeHTML(*req.CompanyCulture)
		if !utils.ValidateNoXSS(sanitized) {
			return utils.BadRequestResponse(c, common.ErrPotentialXSS)
		}
		req.CompanyCulture = &sanitized
	}

	var bannerFile *multipart.FileHeader
	if files := form.File["banner"]; len(files) > 0 {
		bannerFile = files[0]
	}

	var logoFile *multipart.FileHeader
	if files := form.File["logo"]; len(files) > 0 {
		logoFile = files[0]
	}

	domainReq := &company.UpdateCompanyRequest{
		FullAddress:        req.FullAddress,
		ShortDescription:   req.ShortDescription,
		WebsiteURL:         req.WebsiteURL,
		InstagramURL:       req.InstagramURL,
		FacebookURL:        req.FacebookURL,
		LinkedinURL:        req.LinkedinURL,
		TwitterURL:         req.TwitterURL,
		CompanyDescription: req.CompanyDescription,
		CompanyCulture:     req.CompanyCulture,
	}

	if err := h.companyService.UpdateCompany(ctx, companyID, domainReq, bannerFile, logoFile); err != nil {
		return utils.InternalServerErrorResponse(c, common.ErrFailedOperation)
	}

	return utils.SuccessResponse(c, common.MsgUpdatedSuccess, nil)
}

func (h *CompanyBasicHandler) DeleteCompany(c *fiber.Ctx) error {
	ctx := c.Context()

	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, common.ErrUnauthorized, "userID not found in context")
	}

	companyID, err := utils.ParseIDParam(c, "id")
	if err != nil || companyID <= 0 {
		return utils.BadRequestResponse(c, common.ErrInvalidCompanyID)
	}

	isOwner, err := h.companyService.CheckEmployerPermission(ctx, userID, int64(companyID), "owner")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to check user permission", err.Error())
	}
	if !isOwner {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "You don't have permission to delete this company. Only company owner can perform this action.", "")
	}

	if err := h.companyService.DeleteCompany(ctx, int64(companyID)); err != nil {
		return utils.InternalServerErrorResponse(c, common.ErrFailedOperation)
	}

	return utils.SuccessResponse(c, common.MsgDeletedSuccess, nil)
}

func (h *CompanyBasicHandler) GetMyCompanies(c *fiber.Ctx) error {
	ctx := c.Context()

	userID := middleware.GetUserID(c)
	if userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, common.ErrUnauthorized, "userID not found in context")
	}

	companies, err := h.companyService.GetUserCompanies(ctx, userID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, common.ErrFailedOperation)
	}

	responses := mapper.MapEntities(companies, func(comp *company.Company) *response.CompanyDetailResponse {
		detail := mapper.ToCompanyDetailResponse(comp)
		if detail == nil {
			return nil
		}

		addrs, err := h.companyService.GetCompanyAddresses(ctx, comp.ID, false)
		if err == nil && len(addrs) > 0 {
			detail.CompanyAddresses = make([]response.CompanyAddressResponse, len(addrs))
			for i, a := range addrs {
				lat := 0.0
				lon := 0.0
				if a.Latitude != nil {
					lat = *a.Latitude
				}
				if a.Longitude != nil {
					lon = *a.Longitude
				}
				detail.CompanyAddresses[i] = response.CompanyAddressResponse{
					ID:          a.ID,
					FullAddress: a.FullAddress,
					Latitude:    lat,
					Longitude:   lon,
					ProvinceID:  a.ProvinceID,
					CityID:      a.CityID,
					DistrictID:  a.DistrictID,
				}
			}
		}
		return detail
	})

	return utils.SuccessResponse(c, common.MsgFetchedSuccess, responses)
}
