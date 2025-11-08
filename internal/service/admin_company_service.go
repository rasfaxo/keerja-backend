package service

import (
	"context"
	"fmt"
	"time"

	"keerja-backend/internal/cache"
	"keerja-backend/internal/domain/admin"
	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/domain/job"
)

// adminCompanyService implements admin.AdminCompanyService interface
type adminCompanyService struct {
	companyRepo  company.CompanyRepository
	jobRepo      job.JobRepository
	emailService EmailService
	cache        cache.Cache
}

// NewAdminCompanyService creates a new admin company service instance
func NewAdminCompanyService(
	companyRepo company.CompanyRepository,
	jobRepo job.JobRepository,
	emailService EmailService,
	cacheService cache.Cache,
) admin.AdminCompanyService {
	return &adminCompanyService{
		companyRepo:  companyRepo,
		jobRepo:      jobRepo,
		emailService: emailService,
		cache:        cacheService,
	}
}

// ListCompanies retrieves companies with filters, pagination, and search
// Task 2.1: List companies with advanced filtering
func (s *adminCompanyService) ListCompanies(ctx context.Context, req *admin.AdminCompanyListRequest) (*admin.AdminCompanyListResponse, error) {
	// Build company filter
	filter := &company.CompanyFilter{
		IndustryID:    req.IndustryID,
		CompanySizeID: req.CompanySizeID,
		ProvinceID:    req.ProvinceID,
		CityID:        req.CityID,
		Verified:      req.Verified,
		IsActive:      req.IsActive,
		SearchQuery:   toStringPtr(req.Search),
		Page:          req.Page,
		Limit:         req.Limit,
		SortBy:        req.SortBy,
		SortOrder:     req.SortOrder,
	}

	// Fetch companies with master data
	companies, total, err := s.companyRepo.ListWithMasterData(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch companies: %w", err)
	}

	// Map companies to response items
	items := make([]admin.AdminCompanyListItem, 0, len(companies))
	for _, c := range companies {
		// Get industry name
		industryName := ""
		if c.Industry != nil {
			industryName = *c.Industry
		}

		// Get company size name
		companySizeName := ""
		if c.SizeCategory != nil {
			companySizeName = *c.SizeCategory
		}

		// Build location string
		location := ""
		if c.Province != nil {
			location = *c.Province
			if c.City != nil {
				location = *c.City + ", " + location
			}
		}

		// Get verification status
		verificationStatus := "pending"
		if c.Verified {
			verificationStatus = "verified"
		}

		// Get legal name
		legalName := ""
		if c.LegalName != nil {
			legalName = *c.LegalName
		}

		// Get registration number
		registrationNumber := ""
		if c.RegistrationNumber != nil {
			registrationNumber = *c.RegistrationNumber
		}

		// Count total jobs for this company
		jobFilter := job.JobFilter{
			CompanyID: c.ID,
		}
		_, totalJobs, _ := s.jobRepo.List(ctx, jobFilter, 1, 1)

		item := admin.AdminCompanyListItem{
			ID:                 c.ID,
			UUID:               c.UUID.String(),
			CompanyName:        c.CompanyName,
			Slug:               c.Slug,
			LegalName:          legalName,
			RegistrationNumber: registrationNumber,
			Industry:           industryName,
			CompanySize:        companySizeName,
			Location:           location,
			Verified:           c.Verified,
			VerifiedAt:         c.VerifiedAt,
			IsActive:           c.IsActive,
			VerificationStatus: verificationStatus,
			TotalJobs:          totalJobs,
			ActiveJobs:         totalJobs, // TODO: Get active jobs count
			TotalApplications:  0,         // TODO: Get from application table
			CreatedAt:          c.CreatedAt,
			UpdatedAt:          c.UpdatedAt,
		}
		items = append(items, item)
	}

	// Calculate pagination
	totalPages := (total + int64(req.Limit) - 1) / int64(req.Limit)
	hasNext := req.Page < int(totalPages)
	hasPrev := req.Page > 1

	return &admin.AdminCompanyListResponse{
		Companies:  items,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int(totalPages),
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}, nil
}

// GetCompanyDetail retrieves detailed company information for review
func (s *adminCompanyService) GetCompanyDetail(ctx context.Context, companyID int64) (*admin.AdminCompanyDetailResponse, error) {
	// Get company from repository
	_, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}

	// TODO: Map company to AdminCompanyDetailResponse with all nested data
	return &admin.AdminCompanyDetailResponse{}, nil
}

// UpdateCompanyStatus updates company status (approve/reject/suspend)
func (s *adminCompanyService) UpdateCompanyStatus(ctx context.Context, companyID int64, req *admin.AdminCompanyStatusRequest, adminID int64) error {
	// Get company
	comp, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return fmt.Errorf("company not found: %w", err)
	}

	// Validate status transition
	validStatuses := map[string]bool{
		"pending":   true,
		"verified":  true,
		"rejected":  true,
		"suspended": true,
	}
	if !validStatuses[req.Status] {
		return fmt.Errorf("invalid status: %s", req.Status)
	}

	// Validation: rejection_reason required for rejected status
	if req.Status == "rejected" && (req.RejectionReason == nil || *req.RejectionReason == "") {
		return fmt.Errorf("rejection_reason is required when status is rejected")
	}

	// Get or create company verification record
	verification, err := s.companyRepo.FindVerificationByCompanyID(ctx, companyID)
	if err != nil {
		// Create new verification record if it doesn't exist
		verification = &company.CompanyVerification{
			CompanyID: companyID,
			Status:    req.Status,
		}
	} else {
		// Update existing verification record
		verification.Status = req.Status
	}

	// Update verification based on status
	now := time.Now()
	verification.ReviewedBy = &adminID
	verification.ReviewedAt = &now

	if req.Status == "verified" {
		// Set verified flag
		comp.Verified = true
		comp.VerifiedAt = &now
		comp.VerifiedBy = &adminID
		verification.BadgeGranted = req.GrantBadge != nil && *req.GrantBadge
	} else if req.Status == "rejected" {
		// Keep verified as false, set rejection reason
		comp.Verified = false
		verification.RejectionReason = req.RejectionReason
	} else if req.Status == "suspended" {
		// Suspend the company
		comp.Verified = false
	}

	// Add notes if provided
	if req.Notes != nil {
		verification.VerificationNotes = req.Notes
	}

	// Update company in database
	if err := s.companyRepo.Update(ctx, comp); err != nil {
		return fmt.Errorf("failed to update company: %w", err)
	}

	// Update or create verification record
	if verification.ID == 0 {
		// Create new verification record
		if err := s.companyRepo.CreateVerification(ctx, verification); err != nil {
			return fmt.Errorf("failed to create verification: %w", err)
		}
	} else {
		// Update existing verification record
		if err := s.companyRepo.UpdateVerification(ctx, verification); err != nil {
			return fmt.Errorf("failed to update verification: %w", err)
		}
	}

	// Invalidate cache for this company
	s.cache.Delete(cache.GenerateCacheKey("company", "detail", companyID))
	s.cache.Delete(cache.GenerateCacheKey("company", "slug", comp.Slug))
	s.cache.Delete(cache.GenerateCacheKey("company", "profile", companyID))
	s.cache.Delete(cache.GenerateCacheKey("company", "stats", companyID))
	s.cache.DeletePattern("companies:list:*")
	s.cache.DeletePattern("companies:verified:*")
	s.cache.DeletePattern("companies:top-rated:*")

	// TODO: Create audit log entry
	// TODO: Send email notification to company based on status

	fmt.Printf("✓ Company %d status updated to %s by admin %d\n", companyID, req.Status, adminID)
	return nil
}

// UpdateCompany updates company details (admin support)
func (s *adminCompanyService) UpdateCompany(ctx context.Context, companyID int64, req *admin.AdminUpdateCompanyRequest, adminID int64) error {
	// Get company
	_, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return fmt.Errorf("company not found: %w", err)
	}

	// TODO: Implement company update logic
	// TODO: Map request fields to company entity
	// TODO: Save to database
	// TODO: Create audit log entry

	fmt.Printf("TODO: Update company %d details (admin: %d)\n", companyID, adminID)
	return nil
}

// DeleteCompany deletes a company with validation
func (s *adminCompanyService) DeleteCompany(ctx context.Context, companyID int64, req *admin.AdminDeleteCompanyRequest, adminID int64) error {
	// Get company
	_, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return fmt.Errorf("company not found: %w", err)
	}

	// TODO: Check for active jobs unless force flag is set
	// TODO: Soft delete or hard delete based on business rules
	// TODO: Create audit log entry

	fmt.Printf("TODO: Delete company %d (force: %v, admin: %d)\n", companyID, req.Force, adminID)
	return nil
}

// GetCompanyStats retrieves statistics for a specific company
func (s *adminCompanyService) GetCompanyStats(ctx context.Context, companyID int64) (*admin.AdminCompanyStatsResponse, error) {
	// Verify company exists
	_, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}

	// TODO: Implement actual statistics gathering
	return &admin.AdminCompanyStatsResponse{
		TotalJobs:           0,
		ActiveJobs:          0,
		ClosedJobs:          0,
		DraftJobs:           0,
		TotalApplications:   0,
		PendingApplications: 0,
		TotalFollowers:      0,
		TotalEmployees:      0,
		TotalReviews:        0,
		AverageRating:       0.0,
	}, nil
}

// GetDashboardStats retrieves overall admin dashboard statistics
func (s *adminCompanyService) GetDashboardStats(ctx context.Context) (*admin.AdminDashboardStatsResponse, error) {
	// TODO: Implement actual dashboard statistics
	return &admin.AdminDashboardStatsResponse{
		TotalCompanies:        0,
		VerifiedCompanies:     0,
		PendingVerification:   0,
		RejectedCompanies:     0,
		SuspendedCompanies:    0,
		NewCompaniesThisMonth: 0,
		NewCompaniesToday:     0,
		TotalJobs:             0,
		ActiveJobs:            0,
		TotalApplications:     0,
	}, nil
}

// BulkUpdateStatus updates status for multiple companies at once
func (s *adminCompanyService) BulkUpdateStatus(ctx context.Context, companyIDs []int64, status string, adminID int64) (*admin.BulkOperationResult, error) {
	result := &admin.BulkOperationResult{
		SuccessCount: 0,
		FailedCount:  0,
		Errors:       []admin.BulkError{},
	}

	// Validate status
	validStatuses := map[string]bool{
		"verified":  true,
		"rejected":  true,
		"suspended": true,
	}
	if !validStatuses[status] {
		return nil, fmt.Errorf("invalid status for bulk update: %s", status)
	}

	// Process each company
	for _, companyID := range companyIDs {
		req := &admin.AdminCompanyStatusRequest{
			Status: status,
		}

		if err := s.UpdateCompanyStatus(ctx, companyID, req, adminID); err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, admin.BulkError{
				CompanyID: companyID,
				Error:     err.Error(),
			})
		} else {
			result.SuccessCount++
		}
	}

	return result, nil
}

// GetAuditLogs retrieves audit log entries for a company
func (s *adminCompanyService) GetAuditLogs(ctx context.Context, companyID int64, page, limit int) (*admin.AuditLogListResponse, error) {
	// Verify company exists
	_, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// TODO: Implement actual audit log retrieval
	return &admin.AuditLogListResponse{
		Logs:       []admin.AuditLogEntry{},
		Total:      0,
		Page:       page,
		Limit:      limit,
		TotalPages: 0,
	}, nil
}

// Helper functions

// toStringPtr converts a string to a pointer to string
func toStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
