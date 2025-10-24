package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"keerja-backend/internal/cache"
	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/utils"
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
	companyRepo   company.CompanyRepository
	uploadService UploadService
	cache         cache.Cache
}

// NewCompanyService creates a new company service instance
func NewCompanyService(companyRepo company.CompanyRepository, uploadService UploadService, cacheService cache.Cache) company.CompanyService {
	return &companyService{
		companyRepo:   companyRepo,
		uploadService: uploadService,
		cache:         cacheService,
	}
}

// =============================================================================
// Company Registration and Management
// =============================================================================

// RegisterCompany registers a new company
func (s *companyService) RegisterCompany(ctx context.Context, req *company.RegisterCompanyRequest) (*company.Company, error) {
	// Generate unique slug from company name
	slug := utils.GenerateSlug(req.CompanyName)

	// Check if slug exists, generate unique one if needed
	_, err := s.companyRepo.FindBySlug(ctx, slug)
	if err == nil {
		// Slug exists, generate unique one
		slug = utils.GenerateSlugSimple(req.CompanyName)
	}

	// Create company
	comp := &company.Company{
		CompanyName:        req.CompanyName,
		Slug:               slug,
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
		Country:            "Indonesia", // Default
		PostalCode:         req.PostalCode,
		About:              req.About,
		IsActive:           true,
		Verified:           false,
	}

	if req.Country != nil {
		comp.Country = *req.Country
	}

	if err := s.companyRepo.Create(ctx, comp); err != nil {
		return nil, fmt.Errorf("failed to create company: %w", err)
	}

	// Invalidate list caches (new company should appear in lists)
	s.cache.DeletePattern("companies:list:*")
	if comp.Verified {
		s.cache.DeletePattern("companies:verified:*")
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

// UpdateCompany updates company information
func (s *companyService) UpdateCompany(ctx context.Context, companyID int64, req *company.UpdateCompanyRequest) error {
	// Get existing company
	comp, err := s.companyRepo.FindByID(ctx, companyID)
	if err != nil {
		return fmt.Errorf("company not found: %w", err)
	}

	// Update fields if provided
	if req.CompanyName != nil {
		comp.CompanyName = *req.CompanyName
		// Regenerate slug if name changes
		newSlug := utils.GenerateSlug(*req.CompanyName)
		if newSlug != comp.Slug {
			// Check uniqueness
			_, err := s.companyRepo.FindBySlug(ctx, newSlug)
			if err == nil {
				newSlug = utils.GenerateSlugSimple(*req.CompanyName)
			}
			comp.Slug = newSlug
		}
	}
	if req.LegalName != nil {
		comp.LegalName = req.LegalName
	}
	if req.RegistrationNumber != nil {
		comp.RegistrationNumber = req.RegistrationNumber
	}
	if req.Industry != nil {
		comp.Industry = req.Industry
	}
	if req.CompanyType != nil {
		comp.CompanyType = req.CompanyType
	}
	if req.SizeCategory != nil {
		comp.SizeCategory = req.SizeCategory
	}
	if req.WebsiteURL != nil {
		comp.WebsiteURL = req.WebsiteURL
	}
	if req.EmailDomain != nil {
		comp.EmailDomain = req.EmailDomain
	}
	if req.Phone != nil {
		comp.Phone = req.Phone
	}
	if req.Address != nil {
		comp.Address = req.Address
	}
	if req.City != nil {
		comp.City = req.City
	}
	if req.Province != nil {
		comp.Province = req.Province
	}
	if req.PostalCode != nil {
		comp.PostalCode = req.PostalCode
	}
	if req.Latitude != nil {
		comp.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		comp.Longitude = req.Longitude
	}
	if req.About != nil {
		comp.About = req.About
	}
	if req.Culture != nil {
		comp.Culture = req.Culture
	}
	if req.Benefits != nil {
		comp.Benefits = req.Benefits
	}

	// Update company
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

	// Delete company (cascade will handle related records)
	if err := s.companyRepo.Delete(ctx, companyID); err != nil {
		return fmt.Errorf("failed to delete company: %w", err)
	}

	// Invalidate all related caches
	s.cache.DeletePattern(fmt.Sprintf("company:*:%d", companyID))
	s.cache.DeletePattern("companies:list:*")
	s.cache.DeletePattern("companies:verified:*")
	s.cache.DeletePattern("companies:top-rated:*")

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

	return nil
}

// UnfollowCompany unfollows a company
func (s *companyService) UnfollowCompany(ctx context.Context, companyID, userID int64) error {
	if err := s.companyRepo.UnfollowCompany(ctx, companyID, userID); err != nil {
		return fmt.Errorf("failed to unfollow company: %w", err)
	}

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

	if err := s.companyRepo.DeleteReview(ctx, reviewID); err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}

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

// InviteEmployer invites a user to be an employer
func (s *companyService) InviteEmployer(ctx context.Context, req *company.InviteEmployerRequest) error {
	// TODO: Implement email invitation system
	// For now, just create a placeholder
	return fmt.Errorf("employer invitation system not yet implemented")
}

// AcceptInvitation accepts an employer invitation
func (s *companyService) AcceptInvitation(ctx context.Context, userID, companyID int64) error {
	// TODO: Implement invitation acceptance
	return fmt.Errorf("invitation acceptance not yet implemented")
}

// UpdateEmployerRole updates an employer's role
func (s *companyService) UpdateEmployerRole(ctx context.Context, employerUserID int64, newRole string) error {
	employerUser, err := s.companyRepo.FindEmployerUserByID(ctx, employerUserID)
	if err != nil {
		return fmt.Errorf("employer user not found: %w", err)
	}

	employerUser.Role = newRole
	if err := s.companyRepo.UpdateEmployerUser(ctx, employerUser); err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	return nil
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

	if err := s.companyRepo.DeleteEmployerUser(ctx, employerUserID); err != nil {
		return fmt.Errorf("failed to remove employer user: %w", err)
	}

	return nil
}

// GetEmployerUsers retrieves all employer users for a company
func (s *companyService) GetEmployerUsers(ctx context.Context, companyID int64) ([]company.EmployerUser, error) {
	users, err := s.companyRepo.GetEmployerUsersByCompanyID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get employer users: %w", err)
	}

	return users, nil
}

// GetUserCompanies retrieves companies for a user
func (s *companyService) GetUserCompanies(ctx context.Context, userID int64) ([]company.Company, error) {
	companies, err := s.companyRepo.GetCompaniesByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user companies: %w", err)
	}

	return companies, nil
}

// CheckEmployerPermission checks if user has required permission for company
func (s *companyService) CheckEmployerPermission(ctx context.Context, userID, companyID int64, requiredRole string) (bool, error) {
	employerUser, err := s.companyRepo.FindEmployerUserByUserAndCompany(ctx, userID, companyID)
	if err != nil {
		return false, nil // User not an employer for this company
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
func (s *companyService) RequestVerification(ctx context.Context, companyID, requestedBy int64) error {
	if err := s.companyRepo.RequestVerification(ctx, companyID, requestedBy); err != nil {
		return fmt.Errorf("failed to request verification: %w", err)
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
