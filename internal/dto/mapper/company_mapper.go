package mapper

import (
	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/dto/request"
	"keerja-backend/internal/dto/response"
)

// Company Entity to Response Mappers

// ToCompanyResponse maps Company entity to CompanyResponse DTO
func ToCompanyResponse(c *company.Company) *response.CompanyResponse {
	if c == nil {
		return nil
	}

	resp := &response.CompanyResponse{
		ID:           c.ID,
		UUID:         c.UUID.String(),
		CompanyName:  c.CompanyName,
		Slug:         c.Slug,
		Industry:     PtrToString(c.Industry),
		CompanyType:  PtrToString(c.CompanyType),
		SizeCategory: PtrToString(c.SizeCategory),
		WebsiteURL:   PtrToString(c.WebsiteURL),
		Phone:        PtrToString(c.Phone),
		City:         PtrToString(c.City),
		Province:     PtrToString(c.Province),
		Country:      c.Country,
		LogoURL:      PtrToString(c.LogoURL),
		BannerURL:    PtrToString(c.BannerURL),
		About:        PtrToString(c.About),
		Verified:     c.Verified,
		VerifiedAt:   c.VerifiedAt,
		IsActive:     c.IsActive,
		CreatedAt:    c.CreatedAt,
	}

	return resp
}

// ToCompanyDetailResponse maps Company entity with relations to CompanyDetailResponse DTO
func ToCompanyDetailResponse(c *company.Company) *response.CompanyDetailResponse {
	if c == nil {
		return nil
	}

	resp := &response.CompanyDetailResponse{
		ID:                 c.ID,
		UUID:               c.UUID.String(),
		CompanyName:        c.CompanyName,
		Slug:               c.Slug,
		LegalName:          PtrToString(c.LegalName),
		RegistrationNumber: PtrToString(c.RegistrationNumber),

		// Master Data Fields (NEW - Phase 9)
		FullAddress:      c.FullAddress,
		Description:      PtrToString(c.Description),
		ShortDescription: PtrToString(c.About), // short_description from about field

		// Legacy Fields (backward compatibility)
		Industry:     c.GetIndustryName(),     // fallback for industry
		SizeCategory: c.GetCompanySizeLabel(), // fallback for size category
		City:         PtrToString(c.City),
		Province:     PtrToString(c.Province),
		Address:      PtrToString(c.Address),

		// Other Fields
		CompanyType: PtrToString(c.CompanyType),
		WebsiteURL:  PtrToString(c.WebsiteURL),
		EmailDomain: PtrToString(c.EmailDomain),
		Phone:       PtrToString(c.Phone),

		// Social Media URLs
		InstagramURL: PtrToString(c.InstagramURL),
		FacebookURL:  PtrToString(c.FacebookURL),
		LinkedinURL:  PtrToString(c.LinkedinURL),
		TwitterURL:   PtrToString(c.TwitterURL),

		Country:    c.Country,
		PostalCode: PtrToString(c.PostalCode),
		Latitude:   c.Latitude,
		Longitude:  c.Longitude,
		LogoURL:    PtrToString(c.LogoURL),
		BannerURL:  PtrToString(c.BannerURL),
		Culture:    PtrToString(c.Culture),
		Benefits:   []string(c.Benefits), // Convert PostgreSQL array to Go slice
		Verified:   c.Verified,
		VerifiedAt: c.VerifiedAt,
		IsActive:   c.IsActive,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}

	// ALWAYS Map Master Data Relations (whether preloaded or not)
	// Map Industry
	if industry := c.GetIndustry(); industry != nil {
		resp.IndustryDetail = &response.MasterIndustryResponse{
			ID:          industry.ID,
			Name:        industry.Name,
			Slug:        industry.Slug,
			Description: industry.GetDescription(),
			IconURL:     industry.GetIconURL(),
		}
	}

	// Map Company Size
	if companySize := c.GetCompanySize(); companySize != nil {
		maxEmp := companySize.GetMaxEmployees()
		resp.CompanySizeDetail = &response.MasterCompanySizeResponse{
			ID:           companySize.ID,
			Label:        companySize.Label,
			Code:         companySize.Label, // Use label as code
			MinEmployees: companySize.MinEmployees,
			MaxEmployees: &maxEmp,
			Description:  companySize.GetRange(), // Use GetRange() for description
		}
	}

	// Map Location (District -> City -> Province)
	if district := c.GetDistrict(); district != nil && c.GetCity() != nil && c.GetProvince() != nil {
		city := c.GetCity()
		province := c.GetProvince()

		resp.LocationDetail = &response.CompanyLocationResponse{
			Province: response.ProvinceResponse{
				ID:   province.ID,
				Code: province.Code,
				Name: province.Name,
			},
			City: response.CityResponse{
				ID:         city.ID,
				Code:       city.Code,
				Name:       city.Name,
				Type:       city.Type,
				FullName:   city.GetFullName(),
				ProvinceID: city.ProvinceID,
			},
			District: response.DistrictResponse{
				ID:     district.ID,
				Code:   district.Code,
				Name:   district.Name,
				CityID: district.CityID,
			},
			FullLocation: c.GetFullLocation(), // e.g., "Batujajar, Kabupaten Bandung Barat, Jawa Barat"
		}
	}

	// Map profile
	if c.Profile != nil {
		resp.Profile = ToCompanyProfileResponse(c.Profile)
	}

	// Map industries (if available)
	// Note: This needs to be populated if company has Industries relation

	// Map reviews
	if len(c.Reviews) > 0 {
		resp.Reviews = make([]response.CompanyReviewResponse, len(c.Reviews))
		for i, review := range c.Reviews {
			resp.Reviews[i] = *ToCompanyReviewResponse(&review)
		}
	}

	// Note: Employees, Verification, Documents can be populated by handler if needed

	return resp
}

// ToCompanyProfileResponse maps CompanyProfile entity to CompanyProfileResponse DTO
func ToCompanyProfileResponse(p *company.CompanyProfile) *response.CompanyProfileResponse {
	if p == nil {
		return nil
	}

	return &response.CompanyProfileResponse{
		ID:          p.ID,
		Description: PtrToString(p.LongDescription),
		Mission:     PtrToString(p.Mission),
		Vision:      PtrToString(p.Vision),
		UpdatedAt:   p.UpdatedAt,
		// Note: Other fields need to be mapped manually or use entity fields differently
	}
}

// ToCompanyReviewResponse maps CompanyReview entity to CompanyReviewResponse DTO
// Note: Fields may need manual mapping due to entity/DTO structure differences
func ToCompanyReviewResponse(r *company.CompanyReview) *response.CompanyReviewResponse {
	if r == nil {
		return nil
	}

	// Basic mapping - handlers should fill in remaining fields from related data
	return &response.CompanyReviewResponse{
		ID:          r.ID,
		IsAnonymous: r.IsAnonymous,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// ToCompanyEmployeeResponse maps CompanyEmployee entity to CompanyEmployeeResponse DTO
// Note: Fields may need manual mapping due to entity/DTO structure differences
func ToCompanyEmployeeResponse(e *company.CompanyEmployee) *response.CompanyEmployeeResponse {
	if e == nil {
		return nil
	}

	// Basic mapping - handlers should fill in remaining fields
	return &response.CompanyEmployeeResponse{
		ID:        e.ID,
		CreatedAt: e.CreatedAt,
	}
}

// ToCompanyVerificationResponse maps CompanyVerification entity to CompanyVerificationResponse DTO
func ToCompanyVerificationResponse(v *company.CompanyVerification) *response.CompanyVerificationResponse {
	if v == nil {
		return nil
	}

	// Basic mapping - handler should fill remaining fields
	return &response.CompanyVerificationResponse{
		ID:              v.ID,
		Status:          v.Status,
		RejectionReason: PtrToString(v.RejectionReason),
	}
}

// ToCompanyDocumentResponse maps CompanyDocument entity to CompanyDocumentResponse DTO
func ToCompanyDocumentResponse(d *company.CompanyDocument) *response.CompanyDocumentResponse {
	if d == nil {
		return nil
	}

	// Basic mapping - handler should fill remaining fields
	return &response.CompanyDocumentResponse{
		ID:           d.ID,
		DocumentType: d.DocumentType,
	}
}

// ToCompanyFollowerResponse maps CompanyFollower entity to CompanyFollowerResponse DTO
func ToCompanyFollowerResponse(f *company.CompanyFollower) *response.CompanyFollowerResponse {
	if f == nil {
		return nil
	}

	// Basic mapping - handler should fill remaining fields
	return &response.CompanyFollowerResponse{
		ID:         f.ID,
		UserID:     f.UserID,
		FollowedAt: f.FollowedAt,
	}
}

// Request DTO to Entity Mappers

// RegisterCompanyRequestToEntity converts RegisterCompanyRequest to Company entity
func RegisterCompanyRequestToEntity(req *request.RegisterCompanyRequest) *company.Company {
	if req == nil {
		return nil
	}

	c := &company.Company{
		CompanyName:        req.CompanyName,
		LegalName:          req.LegalName,
		RegistrationNumber: req.RegistrationNumber,
		Industry:           req.Industry,
		CompanyType:        req.CompanyType,
		SizeCategory:       req.SizeCategory,
		WebsiteURL:         req.WebsiteURL,
		EmailDomain:        req.EmailDomain,
		Phone:              req.Phone,
		Address:            req.Address,
		City:               req.City,
		Province:           req.Province,
		PostalCode:         req.PostalCode,
		Country:            "Indonesia",
		About:              req.About,
	}

	return c
}

// UpdateCompanyRequestToEntity updates Company entity from UpdateCompanyRequest
func UpdateCompanyRequestToEntity(req *request.UpdateCompanyRequest, c *company.Company) {
	if req == nil || c == nil {
		return
	}

	// NOTE: CompanyName, Country, Province, City, SizeCategory/EmployeeCount, Industry
	// tidak di-update karena sudah di-set saat create company (read-only)

	// Full Address (bisa di-edit, dari data company saat create)
	if req.FullAddress != nil {
		c.FullAddress = *req.FullAddress
	}

	// Deskripsi Singkat - Visi dan Misi Perusahaan
	if req.ShortDescription != nil {
		c.ShortDescription = req.ShortDescription
	}

	// Website & Social Media
	if req.WebsiteURL != nil {
		c.WebsiteURL = req.WebsiteURL
	}
	if req.InstagramURL != nil {
		c.InstagramURL = req.InstagramURL
	}
	if req.FacebookURL != nil {
		c.FacebookURL = req.FacebookURL
	}
	if req.LinkedinURL != nil {
		c.LinkedinURL = req.LinkedinURL
	}
	if req.TwitterURL != nil {
		c.TwitterURL = req.TwitterURL
	}

	// Rich Text Descriptions
	if req.CompanyDescription != nil {
		c.About = req.CompanyDescription
	}
	if req.CompanyCulture != nil {
		c.Culture = req.CompanyCulture
	}

	// Note: LogoURL and BannerURL should be updated via UpdateCompany service
}

// UpdateCompanyProfileRequestToEntity updates CompanyProfile entity from request
func UpdateCompanyProfileRequestToEntity(req *request.UpdateCompanyProfileRequest, profile *company.CompanyProfile) {
	if req == nil || profile == nil {
		return
	}

	if req.Mission != nil {
		profile.Mission = req.Mission
	}
	if req.Vision != nil {
		profile.Vision = req.Vision
	}
	// Note: Additional fields should be mapped based on actual request structure
}
