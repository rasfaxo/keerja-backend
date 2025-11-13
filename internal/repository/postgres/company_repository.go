package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"keerja-backend/internal/domain/company"

	"gorm.io/gorm"
)

// companyRepository implements company.CompanyRepository
type companyRepository struct {
	db *gorm.DB
}

// NewCompanyRepository creates a new instance of CompanyRepository
func NewCompanyRepository(db *gorm.DB) company.CompanyRepository {
	return &companyRepository{db: db}
}

// ===========================================
// COMPANY CRUD OPERATIONS
// ===========================================

// Create creates a new company
func (r *companyRepository) Create(ctx context.Context, c *company.Company) error {
	return r.db.WithContext(ctx).Create(c).Error
}

// FindByID finds a company by ID
func (r *companyRepository) FindByID(ctx context.Context, id int64) (*company.Company, error) {
	var c company.Company
	err := r.db.WithContext(ctx).
		Preload("Profile").
		Preload("Verification").
		First(&c, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

// FindByIDWithMasterData finds a company by ID with all master data relations preloaded
func (r *companyRepository) FindByIDWithMasterData(ctx context.Context, id int64) (*company.Company, error) {
	var c company.Company
	err := r.db.WithContext(ctx).
		Preload("Profile").
		Preload("Verification").
		Preload("IndustryRelation").
		Preload("CompanySizeRelation").
		Preload("ProvinceRelation").
		Preload("CityRelation").
		Preload("DistrictRelation").
		Preload("DistrictRelation.City").
		Preload("DistrictRelation.City.Province").
		First(&c, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

// FindByUUID finds a company by UUID
func (r *companyRepository) FindByUUID(ctx context.Context, uuid string) (*company.Company, error) {
	var c company.Company
	err := r.db.WithContext(ctx).
		Preload("Profile").
		Preload("Verification").
		Where("uuid = ?", uuid).
		First(&c).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

// FindByUUIDWithMasterData finds a company by UUID with all master data relations preloaded
func (r *companyRepository) FindByUUIDWithMasterData(ctx context.Context, uuid string) (*company.Company, error) {
	var c company.Company
	err := r.db.WithContext(ctx).
		Preload("Profile").
		Preload("Verification").
		Preload("IndustryRelation").
		Preload("CompanySizeRelation").
		Preload("ProvinceRelation").
		Preload("CityRelation").
		Preload("DistrictRelation").
		Preload("DistrictRelation.City").          // Fix: Use City, not CityRelation
		Preload("DistrictRelation.City.Province"). // Fix: Use Province, not ProvinceRelation
		Where("uuid = ?", uuid).
		First(&c).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

// FindBySlug finds a company by slug (for public profiles)
func (r *companyRepository) FindBySlug(ctx context.Context, slug string) (*company.Company, error) {
	var c company.Company
	err := r.db.WithContext(ctx).
		Preload("Profile").
		Preload("Verification").
		Where("slug = ?", slug).
		First(&c).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

// FindBySlugWithMasterData finds a company by slug with all master data relations preloaded
func (r *companyRepository) FindBySlugWithMasterData(ctx context.Context, slug string) (*company.Company, error) {
	var c company.Company
	err := r.db.WithContext(ctx).
		Preload("Profile").
		Preload("Verification").
		Preload("IndustryRelation").
		Preload("CompanySizeRelation").
		Preload("ProvinceRelation").
		Preload("CityRelation").
		Preload("DistrictRelation").
		Preload("DistrictRelation.City").          // Fix: Use City, not CityRelation
		Preload("DistrictRelation.City.Province"). // Fix: Use Province, not ProvinceRelation
		Where("slug = ?", slug).
		First(&c).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

// Update updates a company
func (r *companyRepository) Update(ctx context.Context, c *company.Company) error {
	return r.db.WithContext(ctx).Save(c).Error
}

// Delete soft deletes a company
func (r *companyRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&company.Company{}, id).Error
}

// List retrieves companies with filtering, pagination, and sorting
func (r *companyRepository) List(ctx context.Context, filter *company.CompanyFilter) ([]company.Company, int64, error) {
	var companies []company.Company
	var total int64

	query := r.db.WithContext(ctx).Model(&company.Company{})

	// Apply filters
	if filter != nil {
		if filter.Industry != nil {
			query = query.Where("industry = ?", *filter.Industry)
		}
		if filter.CompanyType != nil {
			query = query.Where("company_type = ?", *filter.CompanyType)
		}
		if filter.SizeCategory != nil {
			query = query.Where("size_category = ?", *filter.SizeCategory)
		}
		if filter.City != nil {
			query = query.Where("city = ?", *filter.City)
		}
		if filter.Province != nil {
			query = query.Where("province = ?", *filter.Province)
		}
		if filter.Verified != nil {
			query = query.Where("verified = ?", *filter.Verified)
		}
		if filter.IsActive != nil {
			query = query.Where("is_active = ?", *filter.IsActive)
		}

		// Search by company name
		if filter.SearchQuery != nil && *filter.SearchQuery != "" {
			searchPattern := "%" + strings.ToLower(*filter.SearchQuery) + "%"
			query = query.Where("LOWER(company_name) LIKE ? OR LOWER(legal_name) LIKE ?", searchPattern, searchPattern)
		}
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	page := 1
	limit := 10
	if filter != nil {
		if filter.Page > 0 {
			page = filter.Page
		}
		if filter.Limit > 0 {
			limit = filter.Limit
		}
	}
	offset := (page - 1) * limit

	// Apply sorting
	sortBy := "created_at"
	sortOrder := "DESC"
	if filter != nil {
		if filter.SortBy != "" {
			sortBy = filter.SortBy
		}
		if filter.SortOrder != "" {
			sortOrder = strings.ToUpper(filter.SortOrder)
		}
	}

	// Execute query with preloads
	err := query.
		Preload("Profile").
		Preload("Verification").
		Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
		Limit(limit).
		Offset(offset).
		Find(&companies).Error

	return companies, total, err
}

// ListWithMasterData retrieves companies with filtering, pagination, sorting, and master data preloaded
func (r *companyRepository) ListWithMasterData(ctx context.Context, filter *company.CompanyFilter) ([]company.Company, int64, error) {
	var companies []company.Company
	var total int64

	query := r.db.WithContext(ctx).Model(&company.Company{})

	// Apply filters
	if filter != nil {
		// New master data filters
		if filter.IndustryID != nil {
			query = query.Where("industry_id = ?", *filter.IndustryID)
		}
		if filter.CompanySizeID != nil {
			query = query.Where("company_size_id = ?", *filter.CompanySizeID)
		}
		if filter.ProvinceID != nil {
			query = query.Where("province_id = ?", *filter.ProvinceID)
		}
		if filter.CityID != nil {
			query = query.Where("city_id = ?", *filter.CityID)
		}
		if filter.DistrictID != nil {
			query = query.Where("district_id = ?", *filter.DistrictID)
		}

		// Legacy filters (still supported)
		if filter.Industry != nil {
			query = query.Where("industry = ?", *filter.Industry)
		}
		if filter.CompanyType != nil {
			query = query.Where("company_type = ?", *filter.CompanyType)
		}
		if filter.SizeCategory != nil {
			query = query.Where("size_category = ?", *filter.SizeCategory)
		}
		if filter.City != nil {
			query = query.Where("city = ?", *filter.City)
		}
		if filter.Province != nil {
			query = query.Where("province = ?", *filter.Province)
		}
		if filter.Verified != nil {
			query = query.Where("verified = ?", *filter.Verified)
		}
		if filter.IsActive != nil {
			query = query.Where("is_active = ?", *filter.IsActive)
		}

		// Search by company name
		if filter.SearchQuery != nil && *filter.SearchQuery != "" {
			searchPattern := "%" + strings.ToLower(*filter.SearchQuery) + "%"
			query = query.Where("LOWER(company_name) LIKE ? OR LOWER(legal_name) LIKE ?", searchPattern, searchPattern)
		}
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	page := 1
	limit := 10
	if filter != nil {
		if filter.Page > 0 {
			page = filter.Page
		}
		if filter.Limit > 0 {
			limit = filter.Limit
		}
	}
	offset := (page - 1) * limit

	// Apply sorting
	sortBy := "created_at"
	sortOrder := "DESC"
	if filter != nil {
		if filter.SortBy != "" {
			sortBy = filter.SortBy
		}
		if filter.SortOrder != "" {
			sortOrder = strings.ToUpper(filter.SortOrder)
		}
	}

	// Execute query with all preloads including master data
	err := query.
		Preload("Profile").
		Preload("Verification").
		Preload("IndustryRelation").
		Preload("CompanySizeRelation").
		Preload("ProvinceRelation").
		Preload("CityRelation").
		Preload("DistrictRelation").
		Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
		Limit(limit).
		Offset(offset).
		Find(&companies).Error

	return companies, total, err
}

// ===========================================
// PROFILE OPERATIONS
// ===========================================

// CreateProfile creates a company profile
func (r *companyRepository) CreateProfile(ctx context.Context, profile *company.CompanyProfile) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

// FindProfileByCompanyID finds a profile by company ID
func (r *companyRepository) FindProfileByCompanyID(ctx context.Context, companyID int64) (*company.CompanyProfile, error) {
	var profile company.CompanyProfile
	err := r.db.WithContext(ctx).
		Where("company_id = ?", companyID).
		First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

// UpdateProfile updates a company profile
func (r *companyRepository) UpdateProfile(ctx context.Context, profile *company.CompanyProfile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}

// ===========================================
// FOLLOWER OPERATIONS
// ===========================================

// FollowCompany adds a follower to a company
func (r *companyRepository) FollowCompany(ctx context.Context, companyID, userID int64) error {
	// Check if already following
	var existing company.CompanyFollower
	err := r.db.WithContext(ctx).
		Where("company_id = ? AND user_id = ?", companyID, userID).
		First(&existing).Error

	now := time.Now()
	if err == gorm.ErrRecordNotFound {
		// Create new follower record
		follower := &company.CompanyFollower{
			CompanyID:  companyID,
			UserID:     userID,
			FollowedAt: now,
			IsActive:   true,
		}
		return r.db.WithContext(ctx).Create(follower).Error
	}

	// Reactivate if exists
	if !existing.IsActive {
		return r.db.WithContext(ctx).
			Model(&existing).
			Updates(map[string]interface{}{
				"is_active":     true,
				"followed_at":   now,
				"unfollowed_at": nil,
			}).Error
	}

	return nil
}

// UnfollowCompany removes a follower from a company
func (r *companyRepository) UnfollowCompany(ctx context.Context, companyID, userID int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&company.CompanyFollower{}).
		Where("company_id = ? AND user_id = ?", companyID, userID).
		Updates(map[string]interface{}{
			"is_active":     false,
			"unfollowed_at": now,
		}).Error
}

// IsFollowing checks if a user is following a company
func (r *companyRepository) IsFollowing(ctx context.Context, companyID, userID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&company.CompanyFollower{}).
		Where("company_id = ? AND user_id = ? AND is_active = ?", companyID, userID, true).
		Count(&count).Error
	return count > 0, err
}

// GetFollowers retrieves followers of a company
func (r *companyRepository) GetFollowers(ctx context.Context, companyID int64, page, limit int) ([]company.CompanyFollower, int64, error) {
	var followers []company.CompanyFollower
	var total int64

	query := r.db.WithContext(ctx).
		Model(&company.CompanyFollower{}).
		Where("company_id = ? AND is_active = ?", companyID, true)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	err := query.
		Order("followed_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&followers).Error

	return followers, total, err
}

// GetFollowedCompanies retrieves companies followed by a user
func (r *companyRepository) GetFollowedCompanies(ctx context.Context, userID int64, page, limit int) ([]company.Company, int64, error) {
	var companies []company.Company
	var total int64

	query := r.db.WithContext(ctx).
		Model(&company.Company{}).
		Joins("INNER JOIN company_followers ON company_followers.company_id = companies.id").
		Where("company_followers.user_id = ? AND company_followers.is_active = ?", userID, true)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	err := query.
		Preload("Profile").
		Preload("Verification").
		Order("company_followers.followed_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&companies).Error

	return companies, total, err
}

// CountFollowers counts the number of followers for a company
func (r *companyRepository) CountFollowers(ctx context.Context, companyID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&company.CompanyFollower{}).
		Where("company_id = ? AND is_active = ?", companyID, true).
		Count(&count).Error
	return count, err
}

// ===========================================
// REVIEW OPERATIONS
// ===========================================

// CreateReview creates a company review
func (r *companyRepository) CreateReview(ctx context.Context, review *company.CompanyReview) error {
	return r.db.WithContext(ctx).Create(review).Error
}

// UpdateReview updates a company review
func (r *companyRepository) UpdateReview(ctx context.Context, review *company.CompanyReview) error {
	return r.db.WithContext(ctx).Save(review).Error
}

// DeleteReview soft deletes a review
func (r *companyRepository) DeleteReview(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&company.CompanyReview{}, id).Error
}

// FindReviewByID finds a review by ID
func (r *companyRepository) FindReviewByID(ctx context.Context, id int64) (*company.CompanyReview, error) {
	var review company.CompanyReview
	err := r.db.WithContext(ctx).
		Preload("Company").
		First(&review, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &review, nil
}

// GetReviewsByCompanyID retrieves reviews for a company with filtering
func (r *companyRepository) GetReviewsByCompanyID(ctx context.Context, companyID int64, filter *company.ReviewFilter) ([]company.CompanyReview, int64, error) {
	var reviews []company.CompanyReview
	var total int64

	query := r.db.WithContext(ctx).
		Model(&company.CompanyReview{}).
		Where("company_id = ?", companyID)

	// Apply filters
	if filter != nil {
		if filter.ReviewerType != nil {
			query = query.Where("reviewer_type = ?", *filter.ReviewerType)
		}
		if filter.Status != nil {
			query = query.Where("status = ?", *filter.Status)
		}
		if filter.MinRating != nil {
			query = query.Where("rating_overall >= ?", *filter.MinRating)
		}
		if filter.MaxRating != nil {
			query = query.Where("rating_overall <= ?", *filter.MaxRating)
		}
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	page := 1
	limit := 10
	if filter != nil {
		if filter.Page > 0 {
			page = filter.Page
		}
		if filter.Limit > 0 {
			limit = filter.Limit
		}
	}
	offset := (page - 1) * limit

	// Sorting
	sortBy := "created_at"
	sortOrder := "DESC"
	if filter != nil {
		if filter.SortBy != "" {
			sortBy = filter.SortBy
		}
		if filter.SortOrder != "" {
			sortOrder = strings.ToUpper(filter.SortOrder)
		}
	}

	err := query.
		Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
		Limit(limit).
		Offset(offset).
		Find(&reviews).Error

	return reviews, total, err
}

// GetReviewsByUserID retrieves all reviews by a user
func (r *companyRepository) GetReviewsByUserID(ctx context.Context, userID int64) ([]company.CompanyReview, error) {
	var reviews []company.CompanyReview
	err := r.db.WithContext(ctx).
		Preload("Company").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&reviews).Error
	return reviews, err
}

// ApproveReview approves a company review
func (r *companyRepository) ApproveReview(ctx context.Context, id, moderatedBy int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&company.CompanyReview{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       "approved",
			"moderated_by": moderatedBy,
			"moderated_at": now,
		}).Error
}

// RejectReview rejects a company review
func (r *companyRepository) RejectReview(ctx context.Context, id, moderatedBy int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&company.CompanyReview{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       "rejected",
			"moderated_by": moderatedBy,
			"moderated_at": now,
		}).Error
}

// CalculateAverageRatings calculates average ratings for a company
func (r *companyRepository) CalculateAverageRatings(ctx context.Context, companyID int64) (*company.AverageRatings, error) {
	var result struct {
		AvgOverall    float64
		AvgCulture    float64
		AvgWorkLife   float64
		AvgSalary     float64
		AvgManagement float64
		TotalReviews  int64
	}

	err := r.db.WithContext(ctx).
		Model(&company.CompanyReview{}).
		Where("company_id = ? AND status = ?", companyID, "approved").
		Select(`
			COALESCE(AVG(rating_overall), 0) as avg_overall,
			COALESCE(AVG(rating_culture), 0) as avg_culture,
			COALESCE(AVG(rating_worklife), 0) as avg_worklife,
			COALESCE(AVG(rating_salary), 0) as avg_salary,
			COALESCE(AVG(rating_management), 0) as avg_management,
			COUNT(*) as total_reviews
		`).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &company.AverageRatings{
		Overall:      result.AvgOverall,
		Culture:      result.AvgCulture,
		WorkLife:     result.AvgWorkLife,
		Salary:       result.AvgSalary,
		Management:   result.AvgManagement,
		TotalReviews: result.TotalReviews,
	}, nil
}

// ===========================================
// DOCUMENT OPERATIONS
// ===========================================

// CreateDocument creates a company document
func (r *companyRepository) CreateDocument(ctx context.Context, doc *company.CompanyDocument) error {
	return r.db.WithContext(ctx).Create(doc).Error
}

// UpdateDocument updates a company document
func (r *companyRepository) UpdateDocument(ctx context.Context, doc *company.CompanyDocument) error {
	return r.db.WithContext(ctx).Save(doc).Error
}

// DeleteDocument soft deletes a document
func (r *companyRepository) DeleteDocument(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&company.CompanyDocument{}, id).Error
}

// FindDocumentByID finds a document by ID
func (r *companyRepository) FindDocumentByID(ctx context.Context, id int64) (*company.CompanyDocument, error) {
	var doc company.CompanyDocument
	err := r.db.WithContext(ctx).
		Preload("Company").
		First(&doc, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &doc, nil
}

// GetDocumentsByCompanyID retrieves all documents for a company
func (r *companyRepository) GetDocumentsByCompanyID(ctx context.Context, companyID int64) ([]company.CompanyDocument, error) {
	var documents []company.CompanyDocument
	err := r.db.WithContext(ctx).
		Where("company_id = ? AND is_active = ?", companyID, true).
		Order("document_type ASC, created_at DESC").
		Find(&documents).Error
	return documents, err
}

// ApproveDocument approves a company document
func (r *companyRepository) ApproveDocument(ctx context.Context, id, verifiedBy int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&company.CompanyDocument{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      "approved",
			"verified_by": verifiedBy,
			"verified_at": now,
		}).Error
}

// RejectDocument rejects a company document
func (r *companyRepository) RejectDocument(ctx context.Context, id, verifiedBy int64, reason string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&company.CompanyDocument{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":           "rejected",
			"verified_by":      verifiedBy,
			"verified_at":      now,
			"rejection_reason": reason,
		}).Error
}

// ===========================================
// EMPLOYEE OPERATIONS
// ===========================================

// AddEmployee adds an employee record
func (r *companyRepository) AddEmployee(ctx context.Context, employee *company.CompanyEmployee) error {
	return r.db.WithContext(ctx).Create(employee).Error
}

// UpdateEmployee updates an employee record
func (r *companyRepository) UpdateEmployee(ctx context.Context, employee *company.CompanyEmployee) error {
	return r.db.WithContext(ctx).Save(employee).Error
}

// DeleteEmployee soft deletes an employee
func (r *companyRepository) DeleteEmployee(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&company.CompanyEmployee{}, id).Error
}

// GetEmployeesByCompanyID retrieves employees for a company
func (r *companyRepository) GetEmployeesByCompanyID(ctx context.Context, companyID int64, includeInactive bool) ([]company.CompanyEmployee, error) {
	var employees []company.CompanyEmployee
	query := r.db.WithContext(ctx).
		Where("company_id = ?", companyID)

	if !includeInactive {
		query = query.Where("employment_status = ?", "active")
	}

	err := query.
		Order("employment_status ASC, join_date DESC").
		Find(&employees).Error
	return employees, err
}

// CountEmployees counts employees for a company
func (r *companyRepository) CountEmployees(ctx context.Context, companyID int64, activeOnly bool) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&company.CompanyEmployee{}).
		Where("company_id = ?", companyID)

	if activeOnly {
		query = query.Where("employment_status = ?", "active")
	}

	err := query.Count(&count).Error
	return count, err
}

// ===========================================
// EMPLOYER USER OPERATIONS
// ===========================================

// CreateEmployerUser creates an employer user relationship
func (r *companyRepository) CreateEmployerUser(ctx context.Context, employerUser *company.EmployerUser) error {
	return r.db.WithContext(ctx).Create(employerUser).Error
}

// UpdateEmployerUser updates an employer user
func (r *companyRepository) UpdateEmployerUser(ctx context.Context, employerUser *company.EmployerUser) error {
	return r.db.WithContext(ctx).Save(employerUser).Error
}

// DeleteEmployerUser soft deletes an employer user
func (r *companyRepository) DeleteEmployerUser(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&company.EmployerUser{}, id).Error
}

// FindEmployerUserByID finds an employer user by ID
func (r *companyRepository) FindEmployerUserByID(ctx context.Context, id int64) (*company.EmployerUser, error) {
	var employerUser company.EmployerUser
	err := r.db.WithContext(ctx).
		Preload("Company").
		First(&employerUser, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &employerUser, nil
}

// FindEmployerUserByUserAndCompany finds employer user by user ID and company ID
func (r *companyRepository) FindEmployerUserByUserAndCompany(ctx context.Context, userID, companyID int64) (*company.EmployerUser, error) {
	var employerUser company.EmployerUser
	err := r.db.WithContext(ctx).
		Preload("Company").
		Where("user_id = ? AND company_id = ?", userID, companyID).
		First(&employerUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &employerUser, nil
}

// GetEmployerUsersByCompanyID retrieves all employer users for a company
func (r *companyRepository) GetEmployerUsersByCompanyID(ctx context.Context, companyID int64) ([]company.EmployerUser, error) {
	var employerUsers []company.EmployerUser
	err := r.db.WithContext(ctx).
		Where("company_id = ? AND is_active = ?", companyID, true).
		Order("role ASC, created_at ASC").
		Find(&employerUsers).Error
	return employerUsers, err
}

// GetCompaniesByUserID retrieves all companies managed by a user
func (r *companyRepository) GetCompaniesByUserID(ctx context.Context, userID int64) ([]company.Company, error) {
	var companies []company.Company
	err := r.db.WithContext(ctx).
		Joins("INNER JOIN employer_users ON employer_users.company_id = companies.id").
		Where("employer_users.user_id = ? AND employer_users.is_active = ?", userID, true).
		Preload("Profile").
		Preload("Verification").
		Preload("IndustryRelation").
		Preload("CompanySizeRelation").
		Preload("ProvinceRelation").
		Preload("CityRelation").
		Preload("DistrictRelation").
		Preload("Reviews").
		Order("companies.created_at DESC").
		Find(&companies).Error
	return companies, err
}

// ===========================================
// VERIFICATION OPERATIONS
// ===========================================

// CreateVerification creates a verification record
func (r *companyRepository) CreateVerification(ctx context.Context, verification *company.CompanyVerification) error {
	return r.db.WithContext(ctx).Create(verification).Error
}

// UpdateVerification updates a verification record
func (r *companyRepository) UpdateVerification(ctx context.Context, verification *company.CompanyVerification) error {
	return r.db.WithContext(ctx).
		Model(&company.CompanyVerification{}).
		Where("id = ?", verification.ID).
		Select("*").
		Updates(verification).Error
}

// FindVerificationByCompanyID finds verification by company ID
func (r *companyRepository) FindVerificationByCompanyID(ctx context.Context, companyID int64) (*company.CompanyVerification, error) {
	var verification company.CompanyVerification
	err := r.db.WithContext(ctx).
		Where("company_id = ?", companyID).
		First(&verification).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &verification, nil
}

// RequestVerification requests company verification
func (r *companyRepository) RequestVerification(ctx context.Context, companyID, requestedBy int64) error {
	verification := &company.CompanyVerification{
		CompanyID:   companyID,
		RequestedBy: &requestedBy,
		Status:      "pending",
	}
	return r.db.WithContext(ctx).Create(verification).Error
}

// ApproveVerification approves company verification
func (r *companyRepository) ApproveVerification(ctx context.Context, companyID, reviewedBy int64, notes string) error {
	now := time.Now()
	expiry := now.AddDate(1, 0, 0) // 1 year from now

	// Check if verification record exists
	var verification company.CompanyVerification
	err := r.db.WithContext(ctx).
		Where("company_id = ?", companyID).
		First(&verification).Error

	if err != nil {
		// Record not found, create new verification record
		verification = company.CompanyVerification{
			CompanyID:          companyID,
			ReviewedBy:         &reviewedBy,
			ReviewedAt:         &now,
			Status:             "verified",
			NPWPNumber:         "", // Will be empty if not provided during request
			VerificationScore:  100.0,
			VerificationNotes:  &notes,
			VerificationExpiry: &expiry,
			BadgeGranted:       true,
			CreatedAt:          now,
			UpdatedAt:          now,
		}
		if err := r.db.WithContext(ctx).Create(&verification).Error; err != nil {
			return err
		}
	} else {
		// Update existing verification
		if err := r.db.WithContext(ctx).
			Model(&verification).
			Updates(map[string]interface{}{
				"status":              "verified",
				"reviewed_by":         reviewedBy,
				"reviewed_at":         now,
				"verification_score":  100.0,
				"verification_notes":  notes,
				"verification_expiry": expiry,
				"badge_granted":       true,
				"updated_at":          now,
			}).Error; err != nil {
			return err
		}
	}

	// Update company
	return r.db.WithContext(ctx).
		Model(&company.Company{}).
		Where("id = ?", companyID).
		Updates(map[string]interface{}{
			"verified":    true,
			"verified_at": now,
			"verified_by": reviewedBy,
		}).Error
}

// RejectVerification rejects company verification
func (r *companyRepository) RejectVerification(ctx context.Context, companyID, reviewedBy int64, reason string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&company.CompanyVerification{}).
		Where("company_id = ?", companyID).
		Updates(map[string]interface{}{
			"status":           "rejected",
			"reviewed_by":      reviewedBy,
			"reviewed_at":      now,
			"rejection_reason": reason,
		}).Error
}

// GetPendingVerifications retrieves pending verifications
func (r *companyRepository) GetPendingVerifications(ctx context.Context, page, limit int) ([]company.CompanyVerification, int64, error) {
	var verifications []company.CompanyVerification
	var total int64

	query := r.db.WithContext(ctx).
		Model(&company.CompanyVerification{}).
		Where("status IN ?", []string{"pending", "under_review"})

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	err := query.
		Preload("Company").
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&verifications).Error

	return verifications, total, err
}

// ===========================================
// INDUSTRY OPERATIONS
// ===========================================

// CreateIndustry creates an industry
func (r *companyRepository) CreateIndustry(ctx context.Context, industry *company.CompanyIndustry) error {
	return r.db.WithContext(ctx).Create(industry).Error
}

// UpdateIndustry updates an industry
func (r *companyRepository) UpdateIndustry(ctx context.Context, industry *company.CompanyIndustry) error {
	return r.db.WithContext(ctx).Save(industry).Error
}

// DeleteIndustry soft deletes an industry
func (r *companyRepository) DeleteIndustry(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&company.CompanyIndustry{}, id).Error
}

// FindIndustryByID finds an industry by ID
func (r *companyRepository) FindIndustryByID(ctx context.Context, id int64) (*company.CompanyIndustry, error) {
	var industry company.CompanyIndustry
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		First(&industry, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &industry, nil
}

// FindIndustryByCode finds an industry by code
func (r *companyRepository) FindIndustryByCode(ctx context.Context, code string) (*company.CompanyIndustry, error) {
	var industry company.CompanyIndustry
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Where("code = ?", code).
		First(&industry).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &industry, nil
}

// GetAllIndustries retrieves all industries
func (r *companyRepository) GetAllIndustries(ctx context.Context, activeOnly bool) ([]company.CompanyIndustry, error) {
	var industries []company.CompanyIndustry
	query := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children")

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	err := query.
		Order("code ASC").
		Find(&industries).Error
	return industries, err
}

// GetIndustryTree retrieves hierarchical industry tree
func (r *companyRepository) GetIndustryTree(ctx context.Context) ([]company.CompanyIndustry, error) {
	var industries []company.CompanyIndustry
	err := r.db.WithContext(ctx).
		Preload("Children").
		Where("parent_id IS NULL AND is_active = ?", true).
		Order("code ASC").
		Find(&industries).Error
	return industries, err
}

// ===========================================
// SEARCH AND ANALYTICS
// ===========================================

// SearchCompanies searches companies by query
func (r *companyRepository) SearchCompanies(ctx context.Context, query string, filter *company.CompanyFilter) ([]company.Company, int64, error) {
	if filter == nil {
		filter = &company.CompanyFilter{}
	}
	searchQuery := query
	filter.SearchQuery = &searchQuery
	return r.List(ctx, filter)
}

// GetVerifiedCompanies retrieves verified companies
func (r *companyRepository) GetVerifiedCompanies(ctx context.Context, page, limit int) ([]company.Company, int64, error) {
	verified := true
	active := true
	filter := &company.CompanyFilter{
		Verified:  &verified,
		IsActive:  &active,
		Page:      page,
		Limit:     limit,
		SortBy:    "verified_at",
		SortOrder: "DESC",
	}
	return r.List(ctx, filter)
}

// GetTopRatedCompanies retrieves top rated companies
func (r *companyRepository) GetTopRatedCompanies(ctx context.Context, limit int) ([]company.Company, error) {
	var companies []company.Company

	err := r.db.WithContext(ctx).
		Model(&company.Company{}).
		Select("companies.*, COALESCE(AVG(company_reviews.rating_overall), 0) as avg_rating").
		Joins("LEFT JOIN company_reviews ON company_reviews.company_id = companies.id AND company_reviews.status = ?", "approved").
		Where("companies.is_active = ?", true).
		Group("companies.id").
		Having("COUNT(company_reviews.id) >= ?", 5).
		Order("avg_rating DESC").
		Limit(limit).
		Preload("Profile").
		Preload("Verification").
		Find(&companies).Error

	return companies, err
}

// GetCompaniesNeedingVerificationRenewal retrieves companies needing renewal
func (r *companyRepository) GetCompaniesNeedingVerificationRenewal(ctx context.Context) ([]company.Company, error) {
	var companies []company.Company

	thirtyDaysFromNow := time.Now().AddDate(0, 0, 30)

	err := r.db.WithContext(ctx).
		Joins("INNER JOIN company_verifications ON company_verifications.company_id = companies.id").
		Where("company_verifications.status = ?", "verified").
		Where("company_verifications.verification_expiry <= ?", thirtyDaysFromNow).
		Where("companies.is_active = ?", true).
		Preload("Profile").
		Preload("Verification").
		Find(&companies).Error

	return companies, err
}

// GetFullCompanyProfile retrieves complete company profile with all relationships
func (r *companyRepository) GetFullCompanyProfile(ctx context.Context, companyID int64) (*company.Company, error) {
	var c company.Company
	err := r.db.WithContext(ctx).
		Preload("Profile").
		Preload("Verification").
		Preload("Followers", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = ?", true).Order("followed_at DESC").Limit(100)
		}).
		Preload("Reviews", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", "approved").Order("created_at DESC").Limit(50)
		}).
		Preload("Documents", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = ?", true).Order("document_type ASC")
		}).
		Preload("Employees", func(db *gorm.DB) *gorm.DB {
			return db.Where("employment_status = ?", "active").Order("join_date DESC")
		}).
		Preload("EmployerUsers", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = ?", true).Order("role ASC")
		}).
		First(&c, companyID).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

// ===========================================
// COMPANY INVITATION OPERATIONS
// ===========================================

// CreateInvitation creates a new company invitation
func (r *companyRepository) CreateInvitation(ctx context.Context, invitation *company.CompanyInvitation) error {
	return r.db.WithContext(ctx).Create(invitation).Error
}

// FindInvitationByToken finds an invitation by token
func (r *companyRepository) FindInvitationByToken(ctx context.Context, token string) (*company.CompanyInvitation, error) {
	var invitation company.CompanyInvitation
	err := r.db.WithContext(ctx).
		Preload("Company").
		Where("token = ?", token).
		First(&invitation).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invitation not found")
		}
		return nil, err
	}
	return &invitation, nil
}

// FindInvitationByID finds an invitation by ID
func (r *companyRepository) FindInvitationByID(ctx context.Context, id int64) (*company.CompanyInvitation, error) {
	var invitation company.CompanyInvitation
	err := r.db.WithContext(ctx).
		Preload("Company").
		First(&invitation, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invitation not found")
		}
		return nil, err
	}
	return &invitation, nil
}

// UpdateInvitation updates an invitation record
func (r *companyRepository) UpdateInvitation(ctx context.Context, invitation *company.CompanyInvitation) error {
	return r.db.WithContext(ctx).Save(invitation).Error
}

// GetPendingInvitationsByCompany retrieves pending invitations for a company
func (r *companyRepository) GetPendingInvitationsByCompany(ctx context.Context, companyID int64) ([]company.CompanyInvitation, error) {
	var invitations []company.CompanyInvitation
	err := r.db.WithContext(ctx).
		Where("company_id = ? AND status = ?", companyID, "pending").
		Where("expires_at > ?", time.Now()).
		Order("created_at DESC").
		Find(&invitations).Error
	return invitations, err
}

// GetPendingInvitationsByEmail retrieves pending invitations for an email
func (r *companyRepository) GetPendingInvitationsByEmail(ctx context.Context, email string) ([]company.CompanyInvitation, error) {
	var invitations []company.CompanyInvitation
	err := r.db.WithContext(ctx).
		Preload("Company").
		Where("email = ? AND status = ?", email, "pending").
		Where("expires_at > ?", time.Now()).
		Order("created_at DESC").
		Find(&invitations).Error
	return invitations, err
}

// ExpireOldInvitations marks old invitations as expired
func (r *companyRepository) ExpireOldInvitations(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Model(&company.CompanyInvitation{}).
		Where("status = ? AND expires_at < ?", "pending", time.Now()).
		Update("status", "expired").Error
}

// DeleteInvitation deletes an invitation
func (r *companyRepository) DeleteInvitation(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&company.CompanyInvitation{}, id).Error
}
