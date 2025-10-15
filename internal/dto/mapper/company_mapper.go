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
		Industry:           PtrToString(c.Industry),
		CompanyType:        PtrToString(c.CompanyType),
		SizeCategory:       PtrToString(c.SizeCategory),
		WebsiteURL:         PtrToString(c.WebsiteURL),
		EmailDomain:        PtrToString(c.EmailDomain),
		Phone:              PtrToString(c.Phone),
		Address:            PtrToString(c.Address),
		City:               PtrToString(c.City),
		Province:           PtrToString(c.Province),
		Country:            c.Country,
		PostalCode:         PtrToString(c.PostalCode),
		Latitude:           c.Latitude,
		Longitude:          c.Longitude,
		LogoURL:            PtrToString(c.LogoURL),
		BannerURL:          PtrToString(c.BannerURL),
		About:              PtrToString(c.About),
		Culture:            PtrToString(c.Culture),
		Verified:           c.Verified,
		VerifiedAt:         c.VerifiedAt,
		IsActive:           c.IsActive,
		CreatedAt:          c.CreatedAt,
		UpdatedAt:          c.UpdatedAt,
	}

	// Map profile
	if c.Profile != nil {
		resp.Profile = ToCompanyProfileResponse(c.Profile)
	}

	// Map reviews
	if len(c.Reviews) > 0 {
		resp.Reviews = make([]response.CompanyReviewResponse, len(c.Reviews))
		for i, review := range c.Reviews {
			resp.Reviews[i] = *ToCompanyReviewResponse(&review)
		}
	}

	// Note: Employees, Verification, Documents need to be populated by handler
	// as the response structure doesn't match entity structure directly

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

	if req.CompanyName != nil {
		c.CompanyName = *req.CompanyName
	}
	if req.LegalName != nil {
		c.LegalName = req.LegalName
	}
	if req.Industry != nil {
		c.Industry = req.Industry
	}
	if req.CompanyType != nil {
		c.CompanyType = req.CompanyType
	}
	if req.SizeCategory != nil {
		c.SizeCategory = req.SizeCategory
	}
	if req.WebsiteURL != nil {
		c.WebsiteURL = req.WebsiteURL
	}
	if req.Phone != nil {
		c.Phone = req.Phone
	}
	if req.Address != nil {
		c.Address = req.Address
	}
	if req.City != nil {
		c.City = req.City
	}
	if req.Province != nil {
		c.Province = req.Province
	}
	if req.PostalCode != nil {
		c.PostalCode = req.PostalCode
	}
	if req.About != nil {
		c.About = req.About
	}
	// Note: LogoURL and BannerURL should be updated via separate upload endpoints
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
