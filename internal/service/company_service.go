package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"keerja-backend/internal/cache"
	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/domain/job"
	"keerja-backend/internal/domain/master"
	"keerja-backend/internal/domain/user"
	"keerja-backend/internal/utils"

	"gorm.io/gorm"
)

// Cache TTL constants
const (
	CompanyListTTL     = 5 * time.Minute  // Company listings cache
	CompanyDetailTTL   = 10 * time.Minute // Individual company details
	CompanyProfileTTL  = 10 * time.Minute // Company profiles
	CompanyStatsTTL    = 15 * time.Minute // Company statistics
	CompanyReviewTTL   = 5 * time.Minute  // Company reviews
	CompanyRatingTTL   = 10 * time.Minute // Average ratings
	CompanyVerifiedTTL = 15 * time.Minute // Verified companies
	CompanyTopRatedTTL = 15 * time.Minute // Top-rated companies
)

// companyService implements the CompanyService interface
type companyService struct {
	companyRepo        company.CompanyRepository
	uploadService      UploadService
	cache              cache.Cache
	db                 *gorm.DB // GORM DB for transaction support
	industryService    master.IndustryService
	companySizeService master.CompanySizeService
	districtService    master.DistrictService
	jobRepo            job.JobRepository
	userService        user.UserService
	userRepo           user.UserRepository
}

// GetJobsGroupedByStatus implements CompanyService interface for job status grouping
func (s *companyService) GetJobsGroupedByStatus(ctx context.Context, userID int64) (map[string][]job.Job, error) {
	// Delegate to jobRepo for real data
	return s.jobRepo.GetJobsGroupedByStatus(ctx, userID)
}

// NewCompanyService creates a new company service instance
func NewCompanyService(
	companyRepo company.CompanyRepository,
	uploadService UploadService,
	cacheService cache.Cache,
	db *gorm.DB,
	industryService master.IndustryService,
	companySizeService master.CompanySizeService,
	districtService master.DistrictService,
	jobRepo job.JobRepository,
	userService user.UserService,
	userRepo user.UserRepository,
) company.CompanyService {
	return &companyService{
		companyRepo:        companyRepo,
		uploadService:      uploadService,
		cache:              cacheService,
		db:                 db,
		industryService:    industryService,
		companySizeService: companySizeService,
		districtService:    districtService,
		jobRepo:            jobRepo,
		userService:        userService,
		userRepo:           userRepo,
	}
}

// =============================================================================
// Company Registration and Management
// =============================================================================

// RegisterCompany registers a new company with transaction support
func (s *companyService) RegisterCompany(ctx context.Context, req *company.RegisterCompanyRequest, userID int64) (*company.Company, error) {
	var comp *company.Company

	// Validate Master Data IDs
	err := s.ValidateMasterDataIDs(ctx, req.IndustryID, req.CompanySizeID, req.DistrictID)
	if err != nil {
		return nil, fmt.Errorf("master data validation failed: %w", err)
	}

	// Auto-populate City and Province from District if provided
	var cityID, provinceID *int64
	if req.DistrictID != nil {
		// Get district with relations to extract city and province IDs
		district, err := s.districtService.GetByID(ctx, *req.DistrictID)
		if err == nil && district != nil {
			if district.City != nil {
				cityID = &district.City.ID
				if district.City.Province != nil {
					provinceID = &district.City.Province.ID
				}
			}
		}
	}

	// Execute in transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Generate unique slug from company name
		slug := utils.GenerateSlug(req.CompanyName)

		// Check if slug exists, generate unique one if needed
		_, err := s.companyRepo.FindBySlug(ctx, slug)
		if err == nil {
			// Slug exists, generate unique one
			slug = utils.GenerateSlugSimple(req.CompanyName)
		}

		// Create company
		comp = &company.Company{
			CompanyName:        req.CompanyName,
			Slug:               slug,
			LegalName:          req.LegalName,
			RegistrationNumber: req.RegistrationNumber,

			// Master Data Fields
			IndustryID:    req.IndustryID,
			CompanySizeID: req.CompanySizeID,
			DistrictID:    req.DistrictID,
			CityID:        cityID,     // Auto-populated from District
			ProvinceID:    provinceID, // Auto-populated from District
			FullAddress:   req.FullAddress,
			Description:   req.Description,

			// Legacy Fields (for backward compatibility)
			Industry:     req.Industry,
			SizeCategory: req.SizeCategory,
			City:         req.City,
			Province:     req.Province,
			Address:      req.Address,

			// Other Fields
			CompanyType: req.CompanyType,
			WebsiteURL:  req.WebsiteURL,
			EmailDomain: req.EmailDomain,
			Phone:       req.Phone,
			Country:     "Indonesia", // Default
			PostalCode:  req.PostalCode,
			About:       req.About,
			IsActive:    true,
			Verified:    false,
		}

		if req.Country != nil {
			comp.Country = *req.Country
		}

		// Create company within transaction
		if err := tx.Create(comp).Error; err != nil {
			return fmt.Errorf("failed to create company: %w", err)
		}

		// Add user as company owner
		now := time.Now()
		employerUser := &company.EmployerUser{
			UserID:     userID,
			CompanyID:  comp.ID,
			Role:       "owner",
			IsActive:   true,
			IsVerified: true, // Auto-verify the owner
			VerifiedAt: &now,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		// Create employer user within transaction
		if err := tx.Create(employerUser).Error; err != nil {
			return fmt.Errorf("failed to add user as company owner: %w", err)
		}

		// Create default company profile
		profile := &company.CompanyProfile{
			CompanyID: comp.ID,
			Status:    "draft",
			Verified:  false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := tx.Create(profile).Error; err != nil {
			return fmt.Errorf("failed to create company profile: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Invalidate list caches (new company should appear in lists)
	s.cache.DeletePattern("companies:list:*")
	if comp.Verified {
		s.cache.DeletePattern("companies:verified:*")
	}

	// Invalidate user companies cache for the owner
	s.cache.Delete(cache.GenerateCacheKey("user", "companies", userID))

	// Return company with preloaded master data
	// This ensures response includes full master data details
	if comp.IndustryID != nil || comp.CompanySizeID != nil || comp.DistrictID != nil {
		companyWithData, err := s.GetCompanyWithMasterData(ctx, comp.ID)
		if err == nil {
			return companyWithData, nil
		}
		// If preloading fails, still return the created company
	}

	return comp, nil
}

// GetCompany retrieves company by ID (with caching)
func (s *companyService) GetCompany(ctx context.Context, id int64) (*company.Company, error) {
	// Try cache first
	cacheKey := cache.GenerateCacheKey("company", "detail", id)
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.(*company.Company), nil
	}

	// Cache miss - fetch from database
	comp, err := s.companyRepo.GetFullCompanyProfile(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}

	// Store in cache
	s.cache.Set(cacheKey, comp, CompanyDetailTTL)

	return comp, nil
}

// GetCompanyBySlug retrieves company by slug (with caching)
func (s *companyService) GetCompanyBySlug(ctx context.Context, slug string) (*company.Company, error) {
	// Try cache first
	cacheKey := cache.GenerateCacheKey("company", "slug", slug)
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.(*company.Company), nil
	}

	// Cache miss - fetch from database
	comp, err := s.companyRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}

	// Get full profile
	fullComp, err := s.companyRepo.GetFullCompanyProfile(ctx, comp.ID)
	if err != nil {
		// Store basic info if full profile fails
		s.cache.Set(cacheKey, comp, CompanyDetailTTL)
		return comp, nil
	}

	// Store in cache
	s.cache.Set(cacheKey, fullComp, CompanyDetailTTL)

	return fullComp, nil
}

// UpdateCompany updates company information with banner and logo
func (s *companyService) UpdateCompany(ctx context.Context, companyID int64, req *company.UpdateCompanyRequest, bannerFile, logoFile *multipart.FileHeader) error {
	// Get existing company
	comp, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return fmt.Errorf("company not found: %w", err)
	}

	// Upload banner if provided
	if bannerFile != nil {
		// Delete old banner if exists
		if comp.BannerURL != nil && *comp.BannerURL != "" {
			_ = s.uploadService.DeleteFile(ctx, *comp.BannerURL)
		}

		// Upload new banner
		bannerPath, err := s.uploadService.UploadFile(ctx, bannerFile, fmt.Sprintf("companies/%d/banner", companyID))
		if err != nil {
			return fmt.Errorf("failed to upload banner: %w", err)
		}
		comp.BannerURL = &bannerPath
	}

	// Upload logo if provided
	if logoFile != nil {
		// Delete old logo if exists
		if comp.LogoURL != nil && *comp.LogoURL != "" {
			_ = s.uploadService.DeleteFile(ctx, *comp.LogoURL)
		}

		// Upload new logo
		logoPath, err := s.uploadService.UploadFile(ctx, logoFile, fmt.Sprintf("companies/%d/logo", companyID))
		if err != nil {
			return fmt.Errorf("failed to upload logo: %w", err)
		}
		comp.LogoURL = &logoPath
	}

	// NOTE: CompanyName, Country, Province, City, SizeCategory/EmployeeCount, Industry
	// tidak di-update karena sudah di-set saat create company
	// Data tersebut bersifat read-only dan akan di-get dari database

	// Full Address (bisa di-edit, dari data company saat create)
	if req.FullAddress != nil {
		comp.FullAddress = *req.FullAddress
	}

	// Update Deskripsi Singkat - Visi dan Misi Perusahaan (mapped to About field)
	if req.ShortDescription != nil {
		comp.About = req.ShortDescription
	}

	// Website & Social Media
	if req.WebsiteURL != nil {
		comp.WebsiteURL = req.WebsiteURL
	}
	if req.InstagramURL != nil {
		comp.InstagramURL = req.InstagramURL
	}
	if req.FacebookURL != nil {
		comp.FacebookURL = req.FacebookURL
	}
	if req.LinkedinURL != nil {
		comp.LinkedinURL = req.LinkedinURL
	}
	if req.TwitterURL != nil {
		comp.TwitterURL = req.TwitterURL
	}

	// Rich Text Descriptions
	// CompanyDescription (deskripsi lengkap) mapped to Description field
	if req.CompanyDescription != nil {
		comp.Description = req.CompanyDescription
	}
	// CompanyCulture mapped to Culture field
	if req.CompanyCulture != nil {
		comp.Culture = req.CompanyCulture
	}

	// Automatically update jobs if company verification status changes from false to true
	wasVerified := comp.Verified
	newVerified := wasVerified
	if req.Verified != nil {
		newVerified = *req.Verified
	}

	if !wasVerified && newVerified {
		// Update all jobs for this company from 'in_review' to 'draft'
		if err := s.db.Model(&job.Job{}).
			Where("company_id = ? AND status = ?", companyID, "in_review").
			Update("status", "draft").Error; err != nil {
			return fmt.Errorf("failed to update jobs to draft: %w", err)
		}
		comp.Verified = true
		comp.VerifiedAt = utils.TimePtr(time.Now())
	}

	// Update company in database
	if err := s.companyRepo.Update(ctx, comp); err != nil {
		return fmt.Errorf("failed to update company: %w", err)
	}

	// Invalidate related caches
	s.cache.Delete(cache.GenerateCacheKey("company", "detail", companyID))
	s.cache.Delete(cache.GenerateCacheKey("company", "slug", comp.Slug))
	s.cache.Delete(cache.GenerateCacheKey("company", "profile", companyID))
	s.cache.Delete(cache.GenerateCacheKey("company", "stats", companyID))
	s.cache.DeletePattern("companies:list:*")
	s.cache.DeletePattern("companies:top-rated:*")

	return nil
}

// DeleteCompany deletes a company (soft delete)
func (s *companyService) DeleteCompany(ctx context.Context, companyID int64) error {
	// Get company to clean up files
	comp, err := s.companyRepo.GetFullCompanyProfile(ctx, companyID)
	if err != nil {
		return fmt.Errorf("company not found: %w", err)
	}

	// Clean up logo and banner
	if comp.LogoURL != nil {
		_ = s.uploadService.DeleteFile(ctx, *comp.LogoURL)
	}
	if comp.BannerURL != nil {
		_ = s.uploadService.DeleteFile(ctx, *comp.BannerURL)
	}

	// Clean up documents
	if len(comp.Documents) > 0 {
		for _, doc := range comp.Documents {
			_ = s.uploadService.DeleteFile(ctx, doc.FilePath)
		}
	}

	// Get all employers before deletion to invalidate their caches
	employers, _ := s.companyRepo.GetEmployerUsersByCompanyID(ctx, companyID)

	// Delete company (cascade will handle related records)
	if err := s.companyRepo.Delete(ctx, companyID); err != nil {
		return fmt.Errorf("failed to delete company: %w", err)
	}

	// Invalidate all related caches
	s.cache.DeletePattern(fmt.Sprintf("company:*:%d", companyID))
	s.cache.DeletePattern("companies:list:*")
	s.cache.DeletePattern("companies:verified:*")
	s.cache.DeletePattern("companies:top-rated:*")

	// Invalidate user companies cache for all employers
	for _, emp := range employers {
		s.cache.Delete(cache.GenerateCacheKey("user", "companies", emp.UserID))
	}

	return nil
}

// ListCompanies lists companies with filters (with caching)
func (s *companyService) ListCompanies(ctx context.Context, filter *company.CompanyFilter) ([]company.Company, int64, error) {
	// Generate cache key from filter
	filterMap := map[string]interface{}{
		"search_query":  filter.SearchQuery,
		"industry":      filter.Industry,
		"company_type":  filter.CompanyType,
		"size_category": filter.SizeCategory,
		"city":          filter.City,
		"province":      filter.Province,
		"verified":      filter.Verified,
		"is_active":     filter.IsActive,
		"page":          filter.Page,
		"limit":         filter.Limit,
		"sort_by":       filter.SortBy,
		"sort_order":    filter.SortOrder,
	}
	filterHash := cache.GenerateFilterHash(filterMap)
	cacheKey := cache.GenerateCacheKey("companies", "list", filterHash)

	// Try cache first
	if cached, ok := s.cache.Get(cacheKey); ok {
		result := cached.(map[string]interface{})
		return result["companies"].([]company.Company), result["total"].(int64), nil
	}

	// Cache miss - fetch from database
	companies, total, err := s.companyRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list companies: %w", err)
	}

	// Store in cache
	s.cache.Set(cacheKey, map[string]interface{}{
		"companies": companies,
		"total":     total,
	}, CompanyListTTL)

	return companies, total, nil
}

// SearchCompanies searches companies by query
func (s *companyService) SearchCompanies(ctx context.Context, query string, filter *company.CompanyFilter) ([]company.Company, int64, error) {
	companies, total, err := s.companyRepo.SearchCompanies(ctx, query, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search companies: %w", err)
	}

	return companies, total, nil
}

// =============================================================================
// Profile Management
// =============================================================================

// CreateProfile creates a company profile
func (s *companyService) CreateProfile(ctx context.Context, companyID int64, req *company.CreateProfileRequest) error {
	// Check if profile already exists
	_, err := s.companyRepo.FindProfileByCompanyID(ctx, companyID)
	if err == nil {
		return fmt.Errorf("profile already exists")
	}

	// Create profile
	profile := &company.CompanyProfile{
		CompanyID:        companyID,
		Tagline:          req.Tagline,
		ShortDescription: req.ShortDescription,
		LongDescription:  req.LongDescription,
		Mission:          req.Mission,
		Vision:           req.Vision,
		Culture:          req.Culture,
		WorkEnvironment:  req.WorkEnvironment,
		VideoURL:         req.VideoURL,
		HiringTagline:    req.HiringTagline,
		SEOTitle:         req.SEOTitle,
		SEODescription:   req.SEODescription,
		Status:           "draft",
	}

	if err := s.companyRepo.CreateProfile(ctx, profile); err != nil {
		return fmt.Errorf("failed to create profile: %w", err)
	}

	return nil
}

// UpdateProfile updates company profile
func (s *companyService) UpdateProfile(ctx context.Context, companyID int64, req *company.UpdateProfileRequest) error {
	// Get existing profile
	profile, err := s.companyRepo.FindProfileByCompanyID(ctx, companyID)
	if err != nil {
		return fmt.Errorf("profile not found: %w", err)
	}

	// Update fields if provided
	if req.Tagline != nil {
		profile.Tagline = req.Tagline
	}
	if req.ShortDescription != nil {
		profile.ShortDescription = req.ShortDescription
	}
	if req.LongDescription != nil {
		profile.LongDescription = req.LongDescription
	}
	if req.Mission != nil {
		profile.Mission = req.Mission
	}
	if req.Vision != nil {
		profile.Vision = req.Vision
	}
	if req.Culture != nil {
		profile.Culture = req.Culture
	}
	if req.WorkEnvironment != nil {
		profile.WorkEnvironment = req.WorkEnvironment
	}
	if req.VideoURL != nil {
		profile.VideoURL = req.VideoURL
	}
	if req.HiringTagline != nil {
		profile.HiringTagline = req.HiringTagline
	}
	if req.SEOTitle != nil {
		profile.SEOTitle = req.SEOTitle
	}
	if req.SEODescription != nil {
		profile.SEODescription = req.SEODescription
	}

	// Update profile
	if err := s.companyRepo.UpdateProfile(ctx, profile); err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	return nil
}

// GetProfile retrieves company profile
func (s *companyService) GetProfile(ctx context.Context, companyID int64) (*company.CompanyProfile, error) {
	profile, err := s.companyRepo.FindProfileByCompanyID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("profile not found: %w", err)
	}

	return profile, nil
}

// PublishProfile publishes company profile
func (s *companyService) PublishProfile(ctx context.Context, companyID int64) error {
	profile, err := s.companyRepo.FindProfileByCompanyID(ctx, companyID)
	if err != nil {
		return fmt.Errorf("profile not found: %w", err)
	}

	profile.Status = "published"
	if err := s.companyRepo.UpdateProfile(ctx, profile); err != nil {
		return fmt.Errorf("failed to publish profile: %w", err)
	}

	return nil
}

// UnpublishProfile unpublishes company profile
func (s *companyService) UnpublishProfile(ctx context.Context, companyID int64) error {
	profile, err := s.companyRepo.FindProfileByCompanyID(ctx, companyID)
	if err != nil {
		return fmt.Errorf("profile not found: %w", err)
	}

	profile.Status = "draft"
	if err := s.companyRepo.UpdateProfile(ctx, profile); err != nil {
		return fmt.Errorf("failed to unpublish profile: %w", err)
	}

	return nil
}

// =============================================================================
// Logo and Banner Management
// =============================================================================

// UploadLogo uploads company logo
func (s *companyService) UploadLogo(ctx context.Context, companyID int64, file *multipart.FileHeader) (string, error) {
	// Validate file
	if err := s.uploadService.ValidateFile(file, ImageTypes, MaxAvatarSize); err != nil {
		return "", fmt.Errorf("invalid logo file: %w", err)
	}

	// Get company
	comp, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return "", fmt.Errorf("company not found: %w", err)
	}

	// Delete old logo if exists
	if comp.LogoURL != nil && *comp.LogoURL != "" {
		_ = s.uploadService.DeleteFile(ctx, *comp.LogoURL)
	}

	// Upload new logo
	logoURL, err := s.uploadService.UploadFile(ctx, file, "company/logos")
	if err != nil {
		return "", fmt.Errorf("failed to upload logo: %w", err)
	}

	// Update company with new logo URL
	comp.LogoURL = &logoURL
	if err := s.companyRepo.Update(ctx, comp); err != nil {
		// Clean up uploaded file
		_ = s.uploadService.DeleteFile(ctx, logoURL)
		return "", fmt.Errorf("failed to update company with logo URL: %w", err)
	}

	return logoURL, nil
}

// UploadBanner uploads company banner
func (s *companyService) UploadBanner(ctx context.Context, companyID int64, file *multipart.FileHeader) (string, error) {
	// Validate file
	if err := s.uploadService.ValidateFile(file, ImageTypes, MaxCoverSize); err != nil {
		return "", fmt.Errorf("invalid banner file: %w", err)
	}

	// Get company
	comp, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return "", fmt.Errorf("company not found: %w", err)
	}

	// Delete old banner if exists
	if comp.BannerURL != nil && *comp.BannerURL != "" {
		_ = s.uploadService.DeleteFile(ctx, *comp.BannerURL)
	}

	// Upload new banner
	bannerURL, err := s.uploadService.UploadFile(ctx, file, "company/banners")
	if err != nil {
		return "", fmt.Errorf("failed to upload banner: %w", err)
	}

	// Update company with new banner URL
	comp.BannerURL = &bannerURL
	if err := s.companyRepo.Update(ctx, comp); err != nil {
		// Clean up uploaded file
		_ = s.uploadService.DeleteFile(ctx, bannerURL)
		return "", fmt.Errorf("failed to update company with banner URL: %w", err)
	}

	return bannerURL, nil
}

// DeleteLogo deletes company logo
func (s *companyService) DeleteLogo(ctx context.Context, companyID int64) error {
	comp, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return fmt.Errorf("company not found: %w", err)
	}

	if comp.LogoURL != nil && *comp.LogoURL != "" {
		_ = s.uploadService.DeleteFile(ctx, *comp.LogoURL)
	}

	comp.LogoURL = nil
	if err := s.companyRepo.Update(ctx, comp); err != nil {
		return fmt.Errorf("failed to update company: %w", err)
	}

	return nil
}

// DeleteBanner deletes company banner
func (s *companyService) DeleteBanner(ctx context.Context, companyID int64) error {
	comp, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return fmt.Errorf("company not found: %w", err)
	}

	if comp.BannerURL != nil && *comp.BannerURL != "" {
		_ = s.uploadService.DeleteFile(ctx, *comp.BannerURL)
	}

	comp.BannerURL = nil
	if err := s.companyRepo.Update(ctx, comp); err != nil {
		return fmt.Errorf("failed to update company: %w", err)
	}

	return nil
}

// =============================================================================
// Follower Management
// =============================================================================

// FollowCompany follows a company
func (s *companyService) FollowCompany(ctx context.Context, companyID, userID int64) error {
	// Check if already following
	isFollowing, err := s.companyRepo.IsFollowing(ctx, companyID, userID)
	if err != nil {
		return fmt.Errorf("failed to check follow status: %w", err)
	}

	if isFollowing {
		return fmt.Errorf("already following this company")
	}

	// Follow company
	if err := s.companyRepo.FollowCompany(ctx, companyID, userID); err != nil {
		return fmt.Errorf("failed to follow company: %w", err)
	}

	// Invalidate caches
	s.cache.Delete(cache.GenerateCacheKey("company", "followers", companyID))
	s.cache.Delete(cache.GenerateCacheKey("company", "stats", companyID))
	s.cache.Delete(cache.GenerateCacheKey("user", "followed", userID))

	return nil
}

// UnfollowCompany unfollows a company
func (s *companyService) UnfollowCompany(ctx context.Context, companyID, userID int64) error {
	if err := s.companyRepo.UnfollowCompany(ctx, companyID, userID); err != nil {
		return fmt.Errorf("failed to unfollow company: %w", err)
	}

	// Invalidate caches
	s.cache.Delete(cache.GenerateCacheKey("company", "followers", companyID))
	s.cache.Delete(cache.GenerateCacheKey("company", "stats", companyID))
	s.cache.Delete(cache.GenerateCacheKey("user", "followed", userID))

	return nil
}

// IsFollowing checks if user is following a company
func (s *companyService) IsFollowing(ctx context.Context, companyID, userID int64) (bool, error) {
	return s.companyRepo.IsFollowing(ctx, companyID, userID)
}

// GetFollowers retrieves company followers
func (s *companyService) GetFollowers(ctx context.Context, companyID int64, page, limit int) ([]company.CompanyFollower, int64, error) {
	followers, total, err := s.companyRepo.GetFollowers(ctx, companyID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get followers: %w", err)
	}

	return followers, total, nil
}

// GetFollowedCompanies retrieves companies followed by user
func (s *companyService) GetFollowedCompanies(ctx context.Context, userID int64, page, limit int) ([]company.Company, int64, error) {
	companies, total, err := s.companyRepo.GetFollowedCompanies(ctx, userID, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get followed companies: %w", err)
	}

	return companies, total, nil
}

// GetFollowerCount retrieves follower count for a company
func (s *companyService) GetFollowerCount(ctx context.Context, companyID int64) (int64, error) {
	count, err := s.companyRepo.CountFollowers(ctx, companyID)
	if err != nil {
		return 0, fmt.Errorf("failed to count followers: %w", err)
	}

	return count, nil
}

// =============================================================================
// Review Management
// =============================================================================

// AddReview adds a company review
func (s *companyService) AddReview(ctx context.Context, req *company.AddReviewRequest) (*company.CompanyReview, error) {
	// Create review
	review := &company.CompanyReview{
		CompanyID:          req.CompanyID,
		UserID:             &req.UserID,
		ReviewerType:       req.ReviewerType,
		PositionTitle:      req.PositionTitle,
		EmploymentPeriod:   req.EmploymentPeriod,
		RatingOverall:      &req.RatingOverall,
		RatingCulture:      req.RatingCulture,
		RatingWorkLife:     req.RatingWorkLife,
		RatingSalary:       req.RatingSalary,
		RatingManagement:   req.RatingManagement,
		Pros:               req.Pros,
		Cons:               req.Cons,
		AdviceToManagement: req.AdviceToManagement,
		IsAnonymous:        req.IsAnonymous,
		RecommendToFriend:  req.RecommendToFriend,
		Status:             "pending", // Requires moderation
	}

	if err := s.companyRepo.CreateReview(ctx, review); err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	// Invalidate caches
	s.cache.Delete(cache.GenerateCacheKey("company", "reviews", req.CompanyID))
	s.cache.Delete(cache.GenerateCacheKey("company", "ratings", req.CompanyID))
	s.cache.Delete(cache.GenerateCacheKey("company", "stats", req.CompanyID))
	s.cache.DeletePattern("companies:top-rated:*")

	return review, nil
}

// UpdateReview updates a company review
func (s *companyService) UpdateReview(ctx context.Context, reviewID int64, userID int64, req *company.UpdateReviewRequest) error {
	// Get review to verify ownership
	review, err := s.companyRepo.FindReviewByID(ctx, reviewID)
	if err != nil {
		return fmt.Errorf("review not found: %w", err)
	}

	// Verify ownership
	if review.UserID == nil || *review.UserID != userID {
		return fmt.Errorf("unauthorized to update this review")
	}

	// Update fields if provided
	if req.ReviewerType != nil {
		review.ReviewerType = req.ReviewerType
	}
	if req.PositionTitle != nil {
		review.PositionTitle = req.PositionTitle
	}
	if req.EmploymentPeriod != nil {
		review.EmploymentPeriod = req.EmploymentPeriod
	}
	if req.RatingOverall != nil {
		review.RatingOverall = req.RatingOverall
	}
	if req.RatingCulture != nil {
		review.RatingCulture = req.RatingCulture
	}
	if req.RatingWorkLife != nil {
		review.RatingWorkLife = req.RatingWorkLife
	}
	if req.RatingSalary != nil {
		review.RatingSalary = req.RatingSalary
	}
	if req.RatingManagement != nil {
		review.RatingManagement = req.RatingManagement
	}
	if req.Pros != nil {
		review.Pros = req.Pros
	}
	if req.Cons != nil {
		review.Cons = req.Cons
	}
	if req.AdviceToManagement != nil {
		review.AdviceToManagement = req.AdviceToManagement
	}
	if req.RecommendToFriend != nil {
		review.RecommendToFriend = *req.RecommendToFriend
	}

	// Reset status to pending after update
	review.Status = "pending"

	if err := s.companyRepo.UpdateReview(ctx, review); err != nil {
		return fmt.Errorf("failed to update review: %w", err)
	}

	// Invalidate caches
	s.cache.Delete(cache.GenerateCacheKey("company", "reviews", review.CompanyID))
	s.cache.Delete(cache.GenerateCacheKey("company", "ratings", review.CompanyID))
	s.cache.Delete(cache.GenerateCacheKey("company", "stats", review.CompanyID))
	s.cache.DeletePattern("companies:top-rated:*")

	return nil
}

// DeleteReview deletes a company review
func (s *companyService) DeleteReview(ctx context.Context, reviewID, userID int64) error {
	// Get review to verify ownership
	review, err := s.companyRepo.FindReviewByID(ctx, reviewID)
	if err != nil {
		return fmt.Errorf("review not found: %w", err)
	}

	// Verify ownership
	if review.UserID == nil || *review.UserID != userID {
		return fmt.Errorf("unauthorized to delete this review")
	}

	companyID := review.CompanyID

	if err := s.companyRepo.DeleteReview(ctx, reviewID); err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}

	// Invalidate caches
	s.cache.Delete(cache.GenerateCacheKey("company", "reviews", companyID))
	s.cache.Delete(cache.GenerateCacheKey("company", "ratings", companyID))
	s.cache.Delete(cache.GenerateCacheKey("company", "stats", companyID))
	s.cache.DeletePattern("companies:top-rated:*")

	return nil
}

// GetReview retrieves a review by ID
func (s *companyService) GetReview(ctx context.Context, reviewID int64) (*company.CompanyReview, error) {
	review, err := s.companyRepo.FindReviewByID(ctx, reviewID)
	if err != nil {
		return nil, fmt.Errorf("review not found: %w", err)
	}

	return review, nil
}

// GetCompanyReviews retrieves reviews for a company
func (s *companyService) GetCompanyReviews(ctx context.Context, companyID int64, filter *company.ReviewFilter) ([]company.CompanyReview, int64, error) {
	reviews, total, err := s.companyRepo.GetReviewsByCompanyID(ctx, companyID, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get reviews: %w", err)
	}

	return reviews, total, nil
}

// GetUserReviews retrieves reviews written by a user
func (s *companyService) GetUserReviews(ctx context.Context, userID int64) ([]company.CompanyReview, error) {
	reviews, err := s.companyRepo.GetReviewsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user reviews: %w", err)
	}

	return reviews, nil
}

// GetAverageRatings retrieves average ratings for a company
func (s *companyService) GetAverageRatings(ctx context.Context, companyID int64) (*company.AverageRatings, error) {
	ratings, err := s.companyRepo.CalculateAverageRatings(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate average ratings: %w", err)
	}

	return ratings, nil
}

// =============================================================================
// Review Moderation (Admin Only)
// =============================================================================

// ApproveReview approves a review
func (s *companyService) ApproveReview(ctx context.Context, reviewID, moderatedBy int64) error {
	if err := s.companyRepo.ApproveReview(ctx, reviewID, moderatedBy); err != nil {
		return fmt.Errorf("failed to approve review: %w", err)
	}

	return nil
}

// RejectReview rejects a review
func (s *companyService) RejectReview(ctx context.Context, reviewID, moderatedBy int64) error {
	if err := s.companyRepo.RejectReview(ctx, reviewID, moderatedBy); err != nil {
		return fmt.Errorf("failed to reject review: %w", err)
	}

	return nil
}

// HideReview hides a review
func (s *companyService) HideReview(ctx context.Context, reviewID, moderatedBy int64) error {
	review, err := s.companyRepo.FindReviewByID(ctx, reviewID)
	if err != nil {
		return fmt.Errorf("review not found: %w", err)
	}

	review.Status = "hidden"
	review.ModeratedBy = &moderatedBy
	now := time.Now()
	review.ModeratedAt = &now

	if err := s.companyRepo.UpdateReview(ctx, review); err != nil {
		return fmt.Errorf("failed to hide review: %w", err)
	}

	return nil
}

// GetPendingReviews retrieves pending reviews for moderation
func (s *companyService) GetPendingReviews(ctx context.Context, page, limit int) ([]company.CompanyReview, int64, error) {
	// This is a simplified implementation - ideally would have a dedicated repo method
	reviews := []company.CompanyReview{}
	var total int64 = 0

	return reviews, total, nil
}

// =============================================================================
// Document Management
// =============================================================================

// UploadDocument uploads a company document
func (s *companyService) UploadDocument(ctx context.Context, companyID int64, file *multipart.FileHeader, req *company.UploadDocumentRequest) (*company.CompanyDocument, error) {
	// Validate file
	if err := s.uploadService.ValidateFile(file, DocumentTypes, MaxDocumentSize); err != nil {
		return nil, fmt.Errorf("invalid document file: %w", err)
	}

	// Upload file
	fileURL, err := s.uploadService.UploadFile(ctx, file, "company/documents")
	if err != nil {
		return nil, fmt.Errorf("failed to upload document: %w", err)
	}

	// Parse dates if provided
	var issueDate, expiryDate *time.Time
	if req.IssueDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.IssueDate)
		if err != nil {
			_ = s.uploadService.DeleteFile(ctx, fileURL)
			return nil, fmt.Errorf("invalid issue date format: %w", err)
		}
		issueDate = &parsed
	}

	if req.ExpiryDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.ExpiryDate)
		if err != nil {
			_ = s.uploadService.DeleteFile(ctx, fileURL)
			return nil, fmt.Errorf("invalid expiry date format: %w", err)
		}
		expiryDate = &parsed
	}

	// Create document record
	doc := &company.CompanyDocument{
		CompanyID:      companyID,
		DocumentType:   req.DocumentType,
		DocumentNumber: req.DocumentNumber,
		DocumentName:   req.DocumentName,
		FilePath:       fileURL,
		IssueDate:      issueDate,
		ExpiryDate:     expiryDate,
		Status:         "pending",
		IsActive:       true,
	}

	if err := s.companyRepo.CreateDocument(ctx, doc); err != nil {
		_ = s.uploadService.DeleteFile(ctx, fileURL)
		return nil, fmt.Errorf("failed to save document record: %w", err)
	}

	return doc, nil
}

// UpdateDocument updates document metadata
func (s *companyService) UpdateDocument(ctx context.Context, documentID int64, req *company.UpdateDocumentRequest) error {
	doc, err := s.companyRepo.FindDocumentByID(ctx, documentID)
	if err != nil {
		return fmt.Errorf("document not found: %w", err)
	}

	// Update fields if provided
	if req.DocumentNumber != nil {
		doc.DocumentNumber = req.DocumentNumber
	}
	if req.DocumentName != nil {
		doc.DocumentName = req.DocumentName
	}
	if req.IssueDate != nil {
		issueDate, err := time.Parse("2006-01-02", *req.IssueDate)
		if err != nil {
			return fmt.Errorf("invalid issue date format: %w", err)
		}
		doc.IssueDate = &issueDate
	}
	if req.ExpiryDate != nil {
		expiryDate, err := time.Parse("2006-01-02", *req.ExpiryDate)
		if err != nil {
			return fmt.Errorf("invalid expiry date format: %w", err)
		}
		doc.ExpiryDate = &expiryDate
	}

	if err := s.companyRepo.UpdateDocument(ctx, doc); err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}

	return nil
}

// DeleteDocument deletes a company document
func (s *companyService) DeleteDocument(ctx context.Context, documentID, companyID int64) error {
	// Get document to verify ownership and get file path
	doc, err := s.companyRepo.FindDocumentByID(ctx, documentID)
	if err != nil {
		return fmt.Errorf("document not found: %w", err)
	}

	// Verify ownership
	if doc.CompanyID != companyID {
		return fmt.Errorf("unauthorized to delete this document")
	}

	// Delete file
	_ = s.uploadService.DeleteFile(ctx, doc.FilePath)

	// Delete document record
	if err := s.companyRepo.DeleteDocument(ctx, documentID); err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	return nil
}

// GetDocuments retrieves all documents for a company
func (s *companyService) GetDocuments(ctx context.Context, companyID int64) ([]company.CompanyDocument, error) {
	docs, err := s.companyRepo.GetDocumentsByCompanyID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %w", err)
	}

	return docs, nil
}

// =============================================================================
// Document Verification (Admin Only)
// =============================================================================

// ApproveDocument approves a document
func (s *companyService) ApproveDocument(ctx context.Context, documentID, verifiedBy int64) error {
	if err := s.companyRepo.ApproveDocument(ctx, documentID, verifiedBy); err != nil {
		return fmt.Errorf("failed to approve document: %w", err)
	}

	return nil
}

// RejectDocument rejects a document
func (s *companyService) RejectDocument(ctx context.Context, documentID, verifiedBy int64, reason string) error {
	if err := s.companyRepo.RejectDocument(ctx, documentID, verifiedBy, reason); err != nil {
		return fmt.Errorf("failed to reject document: %w", err)
	}

	return nil
}

// CheckExpiredDocuments checks and updates expired documents
func (s *companyService) CheckExpiredDocuments(ctx context.Context) error {
	// Get all companies needing verification renewal
	companies, err := s.companyRepo.GetCompaniesNeedingVerificationRenewal(ctx)
	if err != nil {
		return fmt.Errorf("failed to get companies: %w", err)
	}

	now := time.Now()
	for _, comp := range companies {
		// Get documents
		docs, err := s.companyRepo.GetDocumentsByCompanyID(ctx, comp.ID)
		if err != nil {
			continue
		}

		// Check each document for expiry
		for _, doc := range docs {
			if doc.ExpiryDate != nil && doc.ExpiryDate.Before(now) && doc.Status != "expired" {
				doc.Status = "expired"
				_ = s.companyRepo.UpdateDocument(ctx, &doc)
			}
		}
	}

	return nil
}

// =============================================================================
// Employee Management
// =============================================================================

// AddEmployee adds an employee record
func (s *companyService) AddEmployee(ctx context.Context, companyID int64, req *company.AddEmployeeRequest) (*company.CompanyEmployee, error) {
	// Parse join date if provided
	var joinDate *time.Time
	if req.JoinDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.JoinDate)
		if err != nil {
			return nil, fmt.Errorf("invalid join date format: %w", err)
		}
		joinDate = &parsed
	}

	employee := &company.CompanyEmployee{
		CompanyID:        companyID,
		UserID:           req.UserID,
		FullName:         req.FullName,
		JobTitle:         req.JobTitle,
		Department:       req.Department,
		EmploymentType:   req.EmploymentType,
		EmploymentStatus: req.EmploymentStatus,
		JoinDate:         joinDate,
		SalaryRangeMin:   req.SalaryRangeMin,
		SalaryRangeMax:   req.SalaryRangeMax,
		Note:             req.Note,
		IsVisiblePublic:  req.IsVisiblePublic,
	}

	if err := s.companyRepo.AddEmployee(ctx, employee); err != nil {
		return nil, fmt.Errorf("failed to add employee: %w", err)
	}

	return employee, nil
}

// UpdateEmployee updates employee information
func (s *companyService) UpdateEmployee(ctx context.Context, employeeID, companyID int64, req *company.UpdateEmployeeRequest) error {
	// Get employees to verify ownership
	employees, err := s.companyRepo.GetEmployeesByCompanyID(ctx, companyID, true)
	if err != nil {
		return fmt.Errorf("failed to get employees: %w", err)
	}

	// Find employee
	var employee *company.CompanyEmployee
	for i := range employees {
		if employees[i].ID == employeeID {
			employee = &employees[i]
			break
		}
	}

	if employee == nil {
		return fmt.Errorf("employee not found or unauthorized")
	}

	// Update fields if provided
	if req.FullName != nil {
		employee.FullName = req.FullName
	}
	if req.JobTitle != nil {
		employee.JobTitle = req.JobTitle
	}
	if req.Department != nil {
		employee.Department = req.Department
	}
	if req.EmploymentType != nil {
		employee.EmploymentType = *req.EmploymentType
	}
	if req.EmploymentStatus != nil {
		employee.EmploymentStatus = *req.EmploymentStatus
	}
	if req.JoinDate != nil {
		joinDate, err := time.Parse("2006-01-02", *req.JoinDate)
		if err != nil {
			return fmt.Errorf("invalid join date format: %w", err)
		}
		employee.JoinDate = &joinDate
	}
	if req.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return fmt.Errorf("invalid end date format: %w", err)
		}
		employee.EndDate = &endDate
	}
	if req.SalaryRangeMin != nil {
		employee.SalaryRangeMin = req.SalaryRangeMin
	}
	if req.SalaryRangeMax != nil {
		employee.SalaryRangeMax = req.SalaryRangeMax
	}
	if req.Note != nil {
		employee.Note = req.Note
	}
	if req.IsVisiblePublic != nil {
		employee.IsVisiblePublic = *req.IsVisiblePublic
	}

	if err := s.companyRepo.UpdateEmployee(ctx, employee); err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}

	return nil
}

// RemoveEmployee removes an employee record
func (s *companyService) RemoveEmployee(ctx context.Context, employeeID, companyID int64) error {
	// Verify ownership
	employees, err := s.companyRepo.GetEmployeesByCompanyID(ctx, companyID, true)
	if err != nil {
		return fmt.Errorf("failed to get employees: %w", err)
	}

	found := false
	for _, emp := range employees {
		if emp.ID == employeeID {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("employee not found or unauthorized")
	}

	if err := s.companyRepo.DeleteEmployee(ctx, employeeID); err != nil {
		return fmt.Errorf("failed to remove employee: %w", err)
	}

	return nil
}

// GetEmployees retrieves all employees for a company
func (s *companyService) GetEmployees(ctx context.Context, companyID int64, includeInactive bool) ([]company.CompanyEmployee, error) {
	employees, err := s.companyRepo.GetEmployeesByCompanyID(ctx, companyID, includeInactive)
	if err != nil {
		return nil, fmt.Errorf("failed to get employees: %w", err)
	}

	return employees, nil
}

// GetEmployeeCount retrieves employee count for a company
func (s *companyService) GetEmployeeCount(ctx context.Context, companyID int64) (int64, error) {
	count, err := s.companyRepo.CountEmployees(ctx, companyID, true) // Active only
	if err != nil {
		return 0, fmt.Errorf("failed to count employees: %w", err)
	}

	return count, nil
}

// =============================================================================
// Employer User Management
// =============================================================================

// InviteEmployer invites a user to be an employer with full invitation system
func (s *companyService) InviteEmployer(ctx context.Context, req *company.InviteEmployerRequest) error {
	// Check if user is already an employer for this company
	existingEmployer, err := s.companyRepo.FindEmployerUserByUserAndCompany(ctx, 0, req.CompanyID)
	if err == nil && existingEmployer != nil {
		// Check by email if userID not found - need to implement email lookup
		// For now, continue with invitation
	}

	// Check if there's already a pending invitation for this email
	pendingInvites, err := s.companyRepo.GetPendingInvitationsByEmail(ctx, req.Email)
	if err == nil && len(pendingInvites) > 0 {
		// Check if any pending invitation is for this company
		for _, inv := range pendingInvites {
			if inv.CompanyID == req.CompanyID && inv.Status == "pending" && !inv.IsExpired() {
				return fmt.Errorf("invitation already sent to this email for this company")
			}
		}
	}

	// Generate invitation token (valid for 7 days)
	token := utils.GenerateRandomToken(32)
	expiresAt := time.Now().AddDate(0, 0, 7) // 7 days from now

	// Get authenticated user ID from request
	invitedBy := req.CompanyID // Placeholder - should come from auth context

	// Create invitation record
	invitation := &company.CompanyInvitation{
		CompanyID: req.CompanyID,
		Email:     req.Email,
		FullName:  req.Email, // Will be filled when accepting
		Position:  req.PositionTitle,
		Role:      req.Role,
		Token:     token,
		Status:    "pending",
		InvitedBy: invitedBy,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save invitation to database
	if err := s.companyRepo.CreateInvitation(ctx, invitation); err != nil {
		return fmt.Errorf("failed to create invitation: %w", err)
	}

	// Invalidate invitation caches
	s.cache.Delete(cache.GenerateCacheKey("company", "invitations", req.CompanyID))
	s.cache.Delete(cache.GenerateCacheKey("user", "invitations", req.Email))

	return nil
}

// AcceptInvitation accepts an employer invitation
func (s *companyService) AcceptInvitation(ctx context.Context, token string, userID int64) error {
	// Find invitation by token
	invitation, err := s.companyRepo.FindInvitationByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("invalid invitation token: %w", err)
	}

	// Check if invitation is expired
	if invitation.IsExpired() {
		// Update status to expired
		invitation.Status = "expired"
		_ = s.companyRepo.UpdateInvitation(ctx, invitation)
		return fmt.Errorf("invitation has expired")
	}

	// Check if invitation is still pending
	if invitation.Status != "pending" {
		return fmt.Errorf("invitation is no longer available (status: %s)", invitation.Status)
	}

	// Execute in transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Check if user is already an employer for this company
		existingEmployer, err := s.companyRepo.FindEmployerUserByUserAndCompany(ctx, userID, invitation.CompanyID)
		if err == nil && existingEmployer != nil {
			return fmt.Errorf("you are already an employer for this company")
		}

		// Create employer user
		now := time.Now()
		employerUser := &company.EmployerUser{
			UserID:        userID,
			CompanyID:     invitation.CompanyID,
			Role:          invitation.Role,
			PositionTitle: invitation.Position,
			IsActive:      true,
			IsVerified:    false, // Will be verified by admin/owner
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		if err := tx.Create(employerUser).Error; err != nil {
			return fmt.Errorf("failed to create employer user: %w", err)
		}

		// Update invitation status
		invitation.Status = "accepted"
		invitation.AcceptedBy = &userID
		invitation.AcceptedAt = &now
		invitation.UpdatedAt = now

		if err := tx.Save(invitation).Error; err != nil {
			return fmt.Errorf("failed to update invitation: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Invalidate caches
	s.cache.Delete(cache.GenerateCacheKey("company", "invitations", invitation.CompanyID))
	s.cache.Delete(cache.GenerateCacheKey("company", "employers", invitation.CompanyID))
	s.cache.Delete(cache.GenerateCacheKey("user", "companies", userID))

	return nil
}

// ResendInvitation resends an invitation email
func (s *companyService) ResendInvitation(ctx context.Context, invitationID, requestedBy int64) error {
	// Get invitation
	invitation, err := s.companyRepo.FindInvitationByID(ctx, invitationID)
	if err != nil {
		return fmt.Errorf("invitation not found: %w", err)
	}

	// Check if invitation is still valid for resend
	if invitation.Status != "pending" {
		return fmt.Errorf("can only resend pending invitations")
	}

	// Generate new token and extend expiry
	invitation.Token = utils.GenerateRandomToken(32)
	invitation.ExpiresAt = time.Now().AddDate(0, 0, 7) // Reset to 7 days
	invitation.UpdatedAt = time.Now()

	// Update invitation
	if err := s.companyRepo.UpdateInvitation(ctx, invitation); err != nil {
		return fmt.Errorf("failed to update invitation: %w", err)
	}

	// Invalidate cache
	s.cache.Delete(cache.GenerateCacheKey("company", "invitations", invitation.CompanyID))

	return nil
}

// CancelInvitation cancels a pending invitation
func (s *companyService) CancelInvitation(ctx context.Context, invitationID, canceledBy int64) error {
	// Get invitation
	invitation, err := s.companyRepo.FindInvitationByID(ctx, invitationID)
	if err != nil {
		return fmt.Errorf("invitation not found: %w", err)
	}

	// Check permission - only the inviter or company admin can cancel
	hasPermission, err := s.CheckEmployerPermission(ctx, canceledBy, invitation.CompanyID, "admin")
	if err != nil || !hasPermission {
		return fmt.Errorf("unauthorized to cancel this invitation")
	}

	// Delete invitation
	if err := s.companyRepo.DeleteInvitation(ctx, invitationID); err != nil {
		return fmt.Errorf("failed to cancel invitation: %w", err)
	}

	// Invalidate cache
	s.cache.Delete(cache.GenerateCacheKey("company", "invitations", invitation.CompanyID))

	return nil
}

// GetPendingInvitations retrieves pending invitations for a company
func (s *companyService) GetPendingInvitations(ctx context.Context, companyID int64) ([]company.CompanyInvitation, error) {
	// Try cache first
	cacheKey := cache.GenerateCacheKey("company", "invitations", companyID)
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.([]company.CompanyInvitation), nil
	}

	// Get from database
	invitations, err := s.companyRepo.GetPendingInvitationsByCompany(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending invitations: %w", err)
	}

	// Cache results
	s.cache.Set(cacheKey, invitations, 5*time.Minute)

	return invitations, nil
}

// GetUserPendingInvitations retrieves pending invitations for a user by email
func (s *companyService) GetUserPendingInvitations(ctx context.Context, email string) ([]company.CompanyInvitation, error) {
	// Try cache first
	cacheKey := cache.GenerateCacheKey("user", "invitations", email)
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.([]company.CompanyInvitation), nil
	}

	// Get from database
	invitations, err := s.companyRepo.GetPendingInvitationsByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user invitations: %w", err)
	}

	// Cache results
	s.cache.Set(cacheKey, invitations, 5*time.Minute)

	return invitations, nil
}

// UpdateEmployerRole updates an employer's role
func (s *companyService) UpdateEmployerRole(ctx context.Context, employerUserID int64, newRole string) error {
	employerUser, err := s.companyRepo.FindEmployerUserByID(ctx, employerUserID)
	if err != nil {
		return fmt.Errorf("employer user not found: %w", err)
	}

	// Don't allow changing owner role
	if employerUser.Role == "owner" {
		return fmt.Errorf("cannot change owner role, use transfer ownership instead")
	}

	// Don't allow setting owner role
	if newRole == "owner" {
		return fmt.Errorf("cannot set owner role, use transfer ownership instead")
	}

	employerUser.Role = newRole
	employerUser.UpdatedAt = time.Now()

	if err := s.companyRepo.UpdateEmployerUser(ctx, employerUser); err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	// Invalidate cache
	s.cache.Delete(cache.GenerateCacheKey("company", "employers", employerUser.CompanyID))
	s.cache.Delete(cache.GenerateCacheKey("user", "companies", employerUser.UserID))

	return nil
}

// UpdateEmployerUser updates fields of the employer_user record for the authenticated user within a company
func (s *companyService) UpdateEmployerUser(ctx context.Context, userID, companyID int64, req *company.UpdateEmployerUserRequest) error {
	// Find employer user by user and company
	employerUser, err := s.companyRepo.FindEmployerUserByUserAndCompany(ctx, userID, companyID)
	if err != nil {
		return fmt.Errorf("failed to find employer user: %w", err)
	}
	if employerUser == nil {
		return fmt.Errorf("employer user not found for this company")
	}

	// Apply updates only for provided fields
	if req.PositionTitle != nil {
		employerUser.PositionTitle = req.PositionTitle
	}
	if req.Department != nil {
		employerUser.Department = req.Department
	}
	if req.EmailCompany != nil {
		employerUser.EmailCompany = req.EmailCompany
	}
	if req.PhoneCompany != nil {
		employerUser.PhoneCompany = req.PhoneCompany
	}

	employerUser.UpdatedAt = time.Now()

	if err := s.companyRepo.UpdateEmployerUser(ctx, employerUser); err != nil {
		return fmt.Errorf("failed to update employer user: %w", err)
	}

	// Invalidate caches related to company employers and user companies
	s.cache.Delete(cache.GenerateCacheKey("company", "employers", employerUser.CompanyID))
	s.cache.Delete(cache.GenerateCacheKey("user", "companies", employerUser.UserID))

	return nil
}

// UpdateEmployerUserWithProfile updates both the global user profile and the
// company-scoped employer_user record inside a single database transaction.
func (s *companyService) UpdateEmployerUserWithProfile(ctx context.Context, userID, companyID int64, userReq *user.UpdateProfileRequest, req *company.UpdateEmployerUserRequest) error {
	// Run everything inside a transaction to ensure atomicity
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1) Load employer_user using transaction
		var emp company.EmployerUser
		if err := tx.WithContext(ctx).
			Where("user_id = ? AND company_id = ?", userID, companyID).
			First(&emp).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("employer user not found")
			}
			return fmt.Errorf("failed to find employer user: %w", err)
		}

		// Apply employer_user updates if provided
		if req != nil {
			if req.PositionTitle != nil {
				emp.PositionTitle = req.PositionTitle
			}
			if req.Department != nil {
				emp.Department = req.Department
			}
			if req.EmailCompany != nil {
				emp.EmailCompany = req.EmailCompany
			}
			if req.PhoneCompany != nil {
				emp.PhoneCompany = req.PhoneCompany
			}
			emp.UpdatedAt = time.Now()

			// Use tx-aware repository method for employer_user update
			if err := s.companyRepo.UpdateEmployerUserTx(ctx, tx, &emp); err != nil {
				return fmt.Errorf("failed to update employer user: %w", err)
			}
		}

		// 2) Update user and profile if requested
		if userReq != nil {
			var u user.User
			if err := tx.WithContext(ctx).Preload("Profile").First(&u, userID).Error; err != nil {
				return fmt.Errorf("failed to find user: %w", err)
			}

			// Update user-level fields
			if userReq.FullName != nil {
				u.FullName = *userReq.FullName
			}

			// Persist user using transaction (direct tx save to keep atomicity)
			if err := tx.WithContext(ctx).Save(&u).Error; err != nil {
				return fmt.Errorf("failed to update user: %w", err)
			}

			// Update or create profile
			if userReq.ProvinceID != nil || userReq.CityID != nil || userReq.DistrictID != nil {
				now := time.Now()
				if u.Profile == nil {
					profile := &user.UserProfile{
						UserID:     u.ID,
						ProvinceID: userReq.ProvinceID,
						CityID:     userReq.CityID,
						DistrictID: userReq.DistrictID,
						CreatedAt:  now,
						UpdatedAt:  now,
					}
					// Create profile within transaction
					if err := tx.WithContext(ctx).Create(profile).Error; err != nil {
						return fmt.Errorf("failed to create user profile: %w", err)
					}
				} else {
					if userReq.ProvinceID != nil {
						u.Profile.ProvinceID = userReq.ProvinceID
					}
					if userReq.CityID != nil {
						u.Profile.CityID = userReq.CityID
					}
					if userReq.DistrictID != nil {
						u.Profile.DistrictID = userReq.DistrictID
					}
					u.Profile.UpdatedAt = now

					// Use tx-aware repository method for profile update
					if err := s.userRepo.UpdateProfileTx(ctx, tx, u.Profile); err != nil {
						return fmt.Errorf("failed to update user profile: %w", err)
					}
				}
			}
		}

		return nil
	})
}

// RemoveEmployerUser removes an employer user
func (s *companyService) RemoveEmployerUser(ctx context.Context, employerUserID, companyID int64) error {
	// Verify ownership
	employerUser, err := s.companyRepo.FindEmployerUserByID(ctx, employerUserID)
	if err != nil {
		return fmt.Errorf("employer user not found: %w", err)
	}

	if employerUser.CompanyID != companyID {
		return fmt.Errorf("unauthorized to remove this employer user")
	}

	// Don't allow removing owner
	if employerUser.Role == "owner" {
		return fmt.Errorf("cannot remove company owner, transfer ownership first")
	}

	if err := s.companyRepo.DeleteEmployerUser(ctx, employerUserID); err != nil {
		return fmt.Errorf("failed to remove employer user: %w", err)
	}

	// Invalidate caches
	s.cache.Delete(cache.GenerateCacheKey("company", "employers", companyID))
	s.cache.Delete(cache.GenerateCacheKey("user", "companies", employerUser.UserID))

	return nil
}

// GetEmployerUsers retrieves all employer users for a company
func (s *companyService) GetEmployerUsers(ctx context.Context, companyID int64) ([]company.EmployerUser, error) {
	// Try cache first
	cacheKey := cache.GenerateCacheKey("company", "employers", companyID)
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.([]company.EmployerUser), nil
	}

	// Get from database
	users, err := s.companyRepo.GetEmployerUsersByCompanyID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get employer users: %w", err)
	}

	// Cache results
	s.cache.Set(cacheKey, users, 10*time.Minute)

	return users, nil
}

// GetUserCompanies retrieves companies for a user
func (s *companyService) GetUserCompanies(ctx context.Context, userID int64) ([]company.Company, error) {
	// Try cache first
	cacheKey := cache.GenerateCacheKey("user", "companies", userID)
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.([]company.Company), nil
	}

	// Get from database
	companies, err := s.companyRepo.GetCompaniesByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user companies: %w", err)
	}

	// Cache results
	s.cache.Set(cacheKey, companies, 10*time.Minute)

	return companies, nil
}

// CreateCompanyAddress creates a persistent address record for a company
func (s *companyService) CreateCompanyAddress(ctx context.Context, companyID int64, req *company.CreateCompanyAddressRequest) (*company.CompanyAddress, error) {
	// Ensure company exists
	comp, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil || comp == nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}

	addr := &company.CompanyAddress{
		CompanyID:   companyID,
		FullAddress: req.FullAddress,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		ProvinceID:  req.ProvinceID,
		CityID:      req.CityID,
		DistrictID:  req.DistrictID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.companyRepo.CreateCompanyAddress(ctx, addr); err != nil {
		return nil, fmt.Errorf("failed to create company address: %w", err)
	}

	return addr, nil
}

// GetCompanyAddresses returns company addresses; includeDeleted toggles returning soft-deleted rows
func (s *companyService) GetCompanyAddresses(ctx context.Context, companyID int64, includeDeleted bool) ([]company.CompanyAddress, error) {
	// Ensure company exists
	comp, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil || comp == nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}

	addrs, err := s.companyRepo.GetCompanyAddressesByCompanyID(ctx, companyID, includeDeleted)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch company addresses: %w", err)
	}
	return addrs, nil
}

// GetCompanyAddressByID retrieves a single company address by id and verifies it belongs to the company
func (s *companyService) GetCompanyAddressByID(ctx context.Context, companyID, addressID int64) (*company.CompanyAddress, error) {
	addr, err := s.companyRepo.FindCompanyAddressByID(ctx, addressID)
	if err != nil {
		return nil, fmt.Errorf("failed to get company address: %w", err)
	}
	if addr == nil {
		return nil, nil
	}
	if addr.CompanyID != companyID {
		return nil, nil
	}
	return addr, nil
}

// SoftDeleteCompanyAddress soft-deletes an address after verifying ownership
func (s *companyService) SoftDeleteCompanyAddress(ctx context.Context, companyID, addressID int64) error {
	addr, err := s.companyRepo.FindCompanyAddressByID(ctx, addressID)
	if err != nil {
		return fmt.Errorf("failed to find address: %w", err)
	}
	if addr == nil {
		return fmt.Errorf("address not found")
	}
	if addr.CompanyID != companyID {
		return fmt.Errorf("unauthorized to delete this address")
	}

	if err := s.companyRepo.SoftDeleteCompanyAddress(ctx, addressID); err != nil {
		return fmt.Errorf("failed to delete address: %w", err)
	}
	return nil
}

// UpdateCompanyAddress updates an existing address after verifying ownership
func (s *companyService) UpdateCompanyAddress(ctx context.Context, companyID, addressID int64, req *company.UpdateCompanyAddressRequest) (*company.CompanyAddress, error) {
	addr, err := s.companyRepo.FindCompanyAddressByID(ctx, addressID)
	if err != nil {
		return nil, fmt.Errorf("failed to find address: %w", err)
	}
	if addr == nil {
		return nil, fmt.Errorf("address not found")
	}
	if addr.CompanyID != companyID {
		return nil, fmt.Errorf("unauthorized to update this address")
	}

	// Apply updates only for provided fields
	if req.FullAddress != nil {
		addr.FullAddress = *req.FullAddress
	}
	if req.Latitude != nil {
		addr.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		addr.Longitude = req.Longitude
	}
	if req.ProvinceID != nil {
		addr.ProvinceID = req.ProvinceID
	}
	if req.CityID != nil {
		addr.CityID = req.CityID
	}
	if req.DistrictID != nil {
		addr.DistrictID = req.DistrictID
	}

	addr.UpdatedAt = time.Now()

	if err := s.companyRepo.UpdateCompanyAddress(ctx, addr); err != nil {
		return nil, fmt.Errorf("failed to update address: %w", err)
	}
	return addr, nil
}

// CheckEmployerPermission checks if user has required permission for company
func (s *companyService) CheckEmployerPermission(ctx context.Context, userID, companyID int64, requiredRole string) (bool, error) {
	employerUser, err := s.companyRepo.FindEmployerUserByUserAndCompany(ctx, userID, companyID)
	if err != nil {
		return false, nil // User not an employer for this company
	}

	// Additional safety check for nil employerUser
	if employerUser == nil {
		return false, nil
	}

	// Check if user is active and verified
	if !employerUser.IsActive {
		return false, nil
	}

	// Role hierarchy: owner > admin > recruiter > viewer
	roleLevel := map[string]int{
		"owner":     4,
		"admin":     3,
		"recruiter": 2,
		"viewer":    1,
	}

	userLevel := roleLevel[employerUser.Role]
	requiredLevel := roleLevel[requiredRole]

	return userLevel >= requiredLevel, nil
}

// =============================================================================
// Verification Management
// =============================================================================

// RequestVerification requests company verification
func (s *companyService) RequestVerification(ctx context.Context, companyID, requestedBy int64, npwpNumber string, nibNumber *string, npwpFile *multipart.FileHeader, additionalFiles []*multipart.FileHeader) error {
	// Validate NPWP number is required
	if npwpNumber == "" {
		return fmt.Errorf("npwp_number is required")
	}

	// Upload NPWP document first
	var npwpDocPath string
	if npwpFile != nil {
		filePath, err := s.uploadService.UploadFile(ctx, npwpFile, "documents/npwp")
		if err != nil {
			return fmt.Errorf("failed to upload NPWP document: %w", err)
		}
		npwpDocPath = filePath
	} else {
		return fmt.Errorf("npwp_file is required")
	}

	// Start transaction
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create or update verification record with NPWP and NIB
	verification, err := s.companyRepo.FindVerificationByCompanyID(ctx, companyID)

	if err != nil {
		// Real error occurred
		tx.Rollback()
		return fmt.Errorf("failed to find verification: %w", err)
	}

	if verification == nil {
		// Create new verification record (no record exists)
		verification = &company.CompanyVerification{
			CompanyID:   companyID,
			RequestedBy: &requestedBy,
			Status:      "pending",
			NPWPNumber:  npwpNumber,
			NIBNumber:   nibNumber,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := tx.Create(verification).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create verification: %w", err)
		}
	} else {
		// Update existing verification
		verification.Status = "pending"
		verification.NPWPNumber = npwpNumber
		verification.NIBNumber = nibNumber
		verification.RequestedBy = &requestedBy
		verification.UpdatedAt = time.Now()
		if err := tx.Save(verification).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update verification: %w", err)
		}
	}

	// Save NPWP document to company_documents
	npwpDocument := &company.CompanyDocument{
		CompanyID:      companyID,
		UploadedBy:     &requestedBy,
		DocumentType:   "NPWP",
		DocumentNumber: &npwpNumber,
		DocumentName:   utils.StringPtr("NPWP Perusahaan"),
		FilePath:       npwpDocPath,
		Status:         "pending",
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := tx.Create(npwpDocument).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to save NPWP document: %w", err)
	}

	// Save NIB document if NIB number provided
	if nibNumber != nil && *nibNumber != "" {
		// Check if NIB file exists in additional documents (optional)
		// For now, we just save the NIB number without requiring document
		nibDocument := &company.CompanyDocument{
			CompanyID:      companyID,
			UploadedBy:     &requestedBy,
			DocumentType:   "NIB",
			DocumentNumber: nibNumber,
			DocumentName:   utils.StringPtr("Nomor Induk Berusaha"),
			FilePath:       "", // NIB is optional, may not have file
			Status:         "pending",
			IsActive:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		// Only create if FilePath is provided or we can skip this
		// For now, skip if no file
		_ = nibDocument
	}

	// Save additional documents
	if len(additionalFiles) > 0 {
		for i, file := range additionalFiles {
			if i >= 5 { // Max 5 files
				break
			}

			filePath, err := s.uploadService.UploadFile(ctx, file, "documents/additional")
			if err != nil {
				// Continue with other files even if one fails
				continue
			}

			additionalDoc := &company.CompanyDocument{
				CompanyID:    companyID,
				UploadedBy:   &requestedBy,
				DocumentType: "LAINNYA",
				DocumentName: utils.StringPtr(fmt.Sprintf("Dokumen Tambahan %d", i+1)),
				FilePath:     filePath,
				Status:       "pending",
				IsActive:     true,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			if err := tx.Create(additionalDoc).Error; err != nil {
				// Continue with other files
				continue
			}
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Invalidate cache
	s.cache.Delete(cache.GenerateCacheKey("company", "verification", companyID))
	s.cache.Delete(cache.GenerateCacheKey("company", "detail", companyID))

	// Get company to invalidate slug-based cache
	comp, err := s.companyRepo.FindByID(ctx, companyID)
	if err == nil && comp != nil {
		s.cache.Delete(cache.GenerateCacheKey("company", "slug", comp.Slug))
	}

	// Invalidate user companies cache for all members of this company
	employees, _ := s.companyRepo.GetEmployerUsersByCompanyID(ctx, companyID)
	for _, emp := range employees {
		s.cache.Delete(cache.GenerateCacheKey("user", "companies", emp.UserID))
	}

	return nil
}

// GetVerificationStatus retrieves verification status
func (s *companyService) GetVerificationStatus(ctx context.Context, companyID int64) (*company.CompanyVerification, error) {
	verification, err := s.companyRepo.FindVerificationByCompanyID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("verification not found: %w", err)
	}

	return verification, nil
}

// ApproveVerification approves company verification
func (s *companyService) ApproveVerification(ctx context.Context, companyID, reviewedBy int64, notes string) error {
	if err := s.companyRepo.ApproveVerification(ctx, companyID, reviewedBy, notes); err != nil {
		return fmt.Errorf("failed to approve verification: %w", err)
	}

	// Update company verified status
	comp, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return fmt.Errorf("company not found: %w", err)
	}

	comp.Verified = true
	now := time.Now()
	comp.VerifiedAt = &now
	comp.VerifiedBy = &reviewedBy

	if err := s.companyRepo.Update(ctx, comp); err != nil {
		return fmt.Errorf("failed to update company: %w", err)
	}

	return nil
}

// RejectVerification rejects company verification
func (s *companyService) RejectVerification(ctx context.Context, companyID, reviewedBy int64, reason string) error {
	if err := s.companyRepo.RejectVerification(ctx, companyID, reviewedBy, reason); err != nil {
		return fmt.Errorf("failed to reject verification: %w", err)
	}

	return nil
}

// GetPendingVerifications retrieves pending verifications
func (s *companyService) GetPendingVerifications(ctx context.Context, page, limit int) ([]company.CompanyVerification, int64, error) {
	verifications, total, err := s.companyRepo.GetPendingVerifications(ctx, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get pending verifications: %w", err)
	}

	return verifications, total, nil
}

// RenewVerification renews company verification
func (s *companyService) RenewVerification(ctx context.Context, companyID int64) error {
	verification, err := s.companyRepo.FindVerificationByCompanyID(ctx, companyID)
	if err != nil {
		return fmt.Errorf("verification not found: %w", err)
	}

	// Extend expiry by 1 year
	if verification.VerificationExpiry != nil {
		newExpiry := verification.VerificationExpiry.AddDate(1, 0, 0)
		verification.VerificationExpiry = &newExpiry
	} else {
		expiry := time.Now().AddDate(1, 0, 0)
		verification.VerificationExpiry = &expiry
	}

	verification.AutoExpired = false
	now := time.Now()
	verification.LastChecked = &now

	if err := s.companyRepo.UpdateVerification(ctx, verification); err != nil {
		return fmt.Errorf("failed to renew verification: %w", err)
	}

	return nil
}

// CheckVerificationExpiry checks and updates expired verifications
func (s *companyService) CheckVerificationExpiry(ctx context.Context) error {
	// Get companies needing renewal
	companies, err := s.companyRepo.GetCompaniesNeedingVerificationRenewal(ctx)
	if err != nil {
		return fmt.Errorf("failed to get companies: %w", err)
	}

	now := time.Now()
	for _, comp := range companies {
		verification, err := s.companyRepo.FindVerificationByCompanyID(ctx, comp.ID)
		if err != nil {
			continue
		}

		if verification.VerificationExpiry != nil && verification.VerificationExpiry.Before(now) {
			verification.Status = "expired"
			verification.AutoExpired = true
			verification.LastChecked = &now

			_ = s.companyRepo.UpdateVerification(ctx, verification)

			// Update company verified status
			comp.Verified = false
			_ = s.companyRepo.Update(ctx, &comp)
		}
	}

	return nil
}

// =============================================================================
// Industry Management (Admin Only)
// =============================================================================

// CreateIndustry creates a new industry category
func (s *companyService) CreateIndustry(ctx context.Context, req *company.CreateIndustryRequest) (*company.CompanyIndustry, error) {
	industry := &company.CompanyIndustry{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
		IsActive:    true,
	}

	if err := s.companyRepo.CreateIndustry(ctx, industry); err != nil {
		return nil, fmt.Errorf("failed to create industry: %w", err)
	}

	return industry, nil
}

// UpdateIndustry updates an industry category
func (s *companyService) UpdateIndustry(ctx context.Context, industryID int64, req *company.UpdateIndustryRequest) error {
	industry, err := s.companyRepo.FindIndustryByID(ctx, industryID)
	if err != nil {
		return fmt.Errorf("industry not found: %w", err)
	}

	// Update fields if provided
	if req.Code != nil {
		industry.Code = *req.Code
	}
	if req.Name != nil {
		industry.Name = *req.Name
	}
	if req.Description != nil {
		industry.Description = req.Description
	}
	if req.ParentID != nil {
		industry.ParentID = req.ParentID
	}
	if req.IsActive != nil {
		industry.IsActive = *req.IsActive
	}

	if err := s.companyRepo.UpdateIndustry(ctx, industry); err != nil {
		return fmt.Errorf("failed to update industry: %w", err)
	}

	return nil
}

// DeleteIndustry deletes an industry category
func (s *companyService) DeleteIndustry(ctx context.Context, industryID int64) error {
	if err := s.companyRepo.DeleteIndustry(ctx, industryID); err != nil {
		return fmt.Errorf("failed to delete industry: %w", err)
	}

	return nil
}

// GetIndustry retrieves an industry by ID
func (s *companyService) GetIndustry(ctx context.Context, industryID int64) (*company.CompanyIndustry, error) {
	industry, err := s.companyRepo.FindIndustryByID(ctx, industryID)
	if err != nil {
		return nil, fmt.Errorf("industry not found: %w", err)
	}

	return industry, nil
}

// GetAllIndustries retrieves all industries
func (s *companyService) GetAllIndustries(ctx context.Context) ([]company.CompanyIndustry, error) {
	industries, err := s.companyRepo.GetAllIndustries(ctx, true) // Active only
	if err != nil {
		return nil, fmt.Errorf("failed to get industries: %w", err)
	}

	return industries, nil
}

// GetIndustryTree retrieves industry hierarchy
func (s *companyService) GetIndustryTree(ctx context.Context) ([]company.CompanyIndustry, error) {
	industries, err := s.companyRepo.GetIndustryTree(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get industry tree: %w", err)
	}

	return industries, nil
}

// =============================================================================
// Analytics and Stats
// =============================================================================

// GetCompanyStats retrieves comprehensive company statistics
func (s *companyService) GetCompanyStats(ctx context.Context, companyID int64) (*company.CompanyStats, error) {
	// Get basic company info
	comp, err := s.companyRepo.GetFullCompanyProfile(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}

	// Get follower count
	followerCount, _ := s.companyRepo.CountFollowers(ctx, companyID)

	// Get employee count
	employeeCount, _ := s.companyRepo.CountEmployees(ctx, companyID, true)

	// Get average ratings
	ratings, _ := s.companyRepo.CalculateAverageRatings(ctx, companyID)

	stats := &company.CompanyStats{
		TotalJobs:           0, // TODO: Implement job count
		ActiveJobs:          0, // TODO: Implement active job count
		TotalApplications:   0, // TODO: Implement application count
		TotalFollowers:      followerCount,
		TotalEmployees:      employeeCount,
		AverageRating:       0,
		TotalReviews:        0,
		VerificationStatus:  "unverified",
		ProfileCompleteness: 0,
	}

	if ratings != nil {
		stats.AverageRating = ratings.Overall
		stats.TotalReviews = ratings.TotalReviews
	}

	if comp.Verified {
		stats.VerificationStatus = "verified"
	}

	return stats, nil
}

// GetTopRatedCompanies retrieves top-rated companies
func (s *companyService) GetTopRatedCompanies(ctx context.Context, limit int) ([]company.Company, error) {
	companies, err := s.companyRepo.GetTopRatedCompanies(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top rated companies: %w", err)
	}

	return companies, nil
}

// GetVerifiedCompanies retrieves verified companies
func (s *companyService) GetVerifiedCompanies(ctx context.Context, page, limit int) ([]company.Company, int64, error) {
	companies, total, err := s.companyRepo.GetVerifiedCompanies(ctx, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get verified companies: %w", err)
	}

	return companies, total, nil
}

// GetCompanyEngagement retrieves engagement statistics
func (s *companyService) GetCompanyEngagement(ctx context.Context, companyID int64) (*company.EngagementStats, error) {
	// Get follower count
	followerCount, _ := s.companyRepo.CountFollowers(ctx, companyID)

	// Get average ratings
	ratings, _ := s.companyRepo.CalculateAverageRatings(ctx, companyID)

	stats := &company.EngagementStats{
		TotalViews:     0, // TODO: Implement view tracking
		TotalFollowers: followerCount,
		FollowerGrowth: 0, // TODO: Implement growth calculation
		TotalReviews:   0,
		AverageRating:  0,
		ResponseRate:   0, // TODO: Implement response tracking
	}

	if ratings != nil {
		stats.AverageRating = ratings.Overall
		stats.TotalReviews = ratings.TotalReviews
	}

	return stats, nil
}

// ExpireOldInvitations expires old pending invitations
func (s *companyService) ExpireOldInvitations(ctx context.Context) (int64, error) {
	err := s.companyRepo.ExpireOldInvitations(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to expire old invitations: %w", err)
	}

	// Return 0 as count since repository doesn't return count
	// This could be improved in the future by updating repository interface
	return 0, nil
}

// GetEmployerUser retrieves employer user by user ID and company ID
func (s *companyService) GetEmployerUser(ctx context.Context, userID, companyID int64) (*company.EmployerUser, error) {
	// Use repository method to get the actual employer user record with correct role
	employerUser, err := s.companyRepo.FindEmployerUserByUserAndCompany(ctx, userID, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get employer user: %w", err)
	}

	if employerUser == nil {
		return nil, fmt.Errorf("user is not an employer of this company")
	}

	return employerUser, nil
}

// GetEmployerUserID retrieves employer_user ID by user ID and company ID
func (s *companyService) GetEmployerUserID(ctx context.Context, userID, companyID int64) (int64, error) {
	employerUser, err := s.GetEmployerUser(ctx, userID, companyID)
	if err != nil {
		return 0, err
	}
	return employerUser.ID, nil
}

// =============================================================================
// Master Data Validation Methods (NEW - Phase 8)
// =============================================================================

// ValidateIndustryID validates if industry ID exists and is active
func (s *companyService) ValidateIndustryID(ctx context.Context, industryID int64) error {
	// Get industry from master data service
	industry, err := s.industryService.GetByID(ctx, industryID)
	if err != nil {
		return fmt.Errorf("failed to validate industry: %w", err)
	}

	if industry == nil {
		return fmt.Errorf("industry with ID %d not found", industryID)
	}

	if !industry.IsActive {
		return fmt.Errorf("industry with ID %d is not active", industryID)
	}

	return nil
}

// ValidateCompanySizeID validates if company size ID exists and is active
func (s *companyService) ValidateCompanySizeID(ctx context.Context, companySizeID int64) error {
	// Get company size from master data service
	companySize, err := s.companySizeService.GetByID(ctx, companySizeID)
	if err != nil {
		return fmt.Errorf("failed to validate company size: %w", err)
	}

	if companySize == nil {
		return fmt.Errorf("company size with ID %d not found", companySizeID)
	}

	if !companySize.IsActive {
		return fmt.Errorf("company size with ID %d is not active", companySizeID)
	}

	return nil
}

// ValidateDistrictID validates if district ID exists and is active
// This also validates the full location hierarchy (District -> City -> Province)
func (s *companyService) ValidateDistrictID(ctx context.Context, districtID int64) error {
	// Get district from master data service (includes City and Province relations)
	district, err := s.districtService.GetByID(ctx, districtID)
	if err != nil {
		return fmt.Errorf("failed to validate district: %w", err)
	}

	if district == nil {
		return fmt.Errorf("district with ID %d not found", districtID)
	}

	if !district.IsActive {
		return fmt.Errorf("district with ID %d is not active", districtID)
	}

	// Validate City (should be included in response)
	if district.City == nil {
		return fmt.Errorf("district with ID %d has no associated city", districtID)
	}

	if !district.City.IsActive {
		return fmt.Errorf("city associated with district %d is not active", districtID)
	}

	// Validate Province (should be included via City)
	if district.City.Province == nil {
		return fmt.Errorf("city associated with district %d has no associated province", districtID)
	}

	if !district.City.Province.IsActive {
		return fmt.Errorf("province associated with district %d is not active", districtID)
	}

	return nil
}

// ValidateMasterDataIDs validates all master data IDs in a single call
func (s *companyService) ValidateMasterDataIDs(ctx context.Context, industryID, companySizeID, districtID *int64) error {
	// Validate Industry ID if provided
	if industryID != nil {
		if err := s.ValidateIndustryID(ctx, *industryID); err != nil {
			return err
		}
	}

	// Validate Company Size ID if provided
	if companySizeID != nil {
		if err := s.ValidateCompanySizeID(ctx, *companySizeID); err != nil {
			return err
		}
	}

	// Validate District ID (includes City and Province) if provided
	if districtID != nil {
		if err := s.ValidateDistrictID(ctx, *districtID); err != nil {
			return err
		}
	}

	return nil
}

// GetCompanyWithMasterData retrieves company by ID with all master data preloaded
func (s *companyService) GetCompanyWithMasterData(ctx context.Context, id int64) (*company.Company, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("company:masterdata:%d", id)

	found, ok := s.cache.Get(cacheKey)
	if ok && found != nil {
		if cachedComp, ok := found.(*company.Company); ok && cachedComp != nil {
			return cachedComp, nil
		}
	}

	// Get from repository with master data preloaded
	comp, err := s.companyRepo.FindByIDWithMasterData(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get company with master data: %w", err)
	}

	if comp == nil {
		return nil, fmt.Errorf("company not found")
	}

	// Cache for 10 minutes
	s.cache.Set(cacheKey, comp, CompanyDetailTTL)

	return comp, nil
}

// GetCompanyBySlugWithMasterData retrieves company by slug with all master data preloaded
func (s *companyService) GetCompanyBySlugWithMasterData(ctx context.Context, slug string) (*company.Company, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("company:slug:masterdata:%s", slug)

	found, ok := s.cache.Get(cacheKey)
	if ok && found != nil {
		if cachedComp, ok := found.(*company.Company); ok && cachedComp != nil {
			return cachedComp, nil
		}
	}

	// Get from repository with master data preloaded
	comp, err := s.companyRepo.FindBySlugWithMasterData(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get company by slug with master data: %w", err)
	}

	if comp == nil {
		return nil, fmt.Errorf("company not found")
	}

	// Cache for 10 minutes
	s.cache.Set(cacheKey, comp, CompanyDetailTTL)

	return comp, nil
}
