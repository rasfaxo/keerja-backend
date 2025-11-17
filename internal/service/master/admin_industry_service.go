package master

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"keerja-backend/internal/cache"
	"keerja-backend/internal/domain/master"
)

// adminIndustryServiceImpl implements AdminIndustryService
type adminIndustryServiceImpl struct {
	master.IndustryService // Embed base service for read operations
	repo                   master.IndustryRepository
	db                     *gorm.DB // For counting references
	cache                  cache.Cache
}

// NewAdminIndustryService creates a new AdminIndustryService
func NewAdminIndustryService(
	baseService master.IndustryService,
	repo master.IndustryRepository,
	db *gorm.DB,
	cache cache.Cache,
) master.AdminIndustryService {
	return &adminIndustryServiceImpl{
		IndustryService: baseService,
		repo:            repo,
		db:              db,
		cache:           cache,
	}
}

// Create creates a new industry
func (s *adminIndustryServiceImpl) Create(ctx context.Context, req master.CreateIndustryRequest) (*master.IndustryResponse, error) {
	// Check duplicate name
	existing, err := s.repo.GetByName(ctx, req.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check duplicate industry name: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("industry with name '%s' already exists", req.Name)
	}

	// Generate slug if not provided
	slug := req.Slug
	if slug == "" {
		slug = strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))
		// Check if slug already exists
		existingBySlug, err := s.repo.GetBySlug(ctx, slug)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to check duplicate industry slug: %w", err)
		}
		if existingBySlug != nil {
			// Append number if slug exists
			slug = fmt.Sprintf("%s-%d", slug, time.Now().Unix())
		}
	} else {
		// Check if provided slug already exists
		existingBySlug, err := s.repo.GetBySlug(ctx, slug)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to check duplicate industry slug: %w", err)
		}
		if existingBySlug != nil {
			return nil, fmt.Errorf("industry with slug '%s' already exists", slug)
		}
	}

	// Create industry entity
	industry := &master.Industry{
		Name:        req.Name,
		Slug:        slug,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		IconURL:     sql.NullString{String: req.IconURL, Valid: req.IconURL != ""},
		IsActive:    req.IsActive,
	}

	if err := s.repo.Create(ctx, industry); err != nil {
		return nil, fmt.Errorf("failed to create industry: %w", err)
	}

	// Invalidate cache
	s.invalidateCache()

	// Map to response
	response := &master.IndustryResponse{
		ID:          industry.ID,
		Name:        industry.Name,
		Slug:        industry.Slug,
		Description: industry.GetDescription(),
		IconURL:     industry.GetIconURL(),
		IsActive:    industry.IsActive,
	}

	return response, nil
}

// Update updates an existing industry
func (s *adminIndustryServiceImpl) Update(ctx context.Context, id int64, req master.UpdateIndustryRequest) (*master.IndustryResponse, error) {
	// Get existing industry
	industry, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("industry not found")
		}
		return nil, fmt.Errorf("failed to get industry: %w", err)
	}

	// Update fields
	if req.Name != "" {
		// Check duplicate name if name is being changed
		if req.Name != industry.Name {
			existing, err := s.repo.GetByName(ctx, req.Name)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("failed to check duplicate industry name: %w", err)
			}
			if existing != nil && existing.ID != id {
				return nil, fmt.Errorf("industry with name '%s' already exists", req.Name)
			}
		}
		industry.Name = req.Name
	}

	if req.Slug != "" {
		// Check duplicate slug if slug is being changed
		if req.Slug != industry.Slug {
			existingBySlug, err := s.repo.GetBySlug(ctx, req.Slug)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("failed to check duplicate industry slug: %w", err)
			}
			if existingBySlug != nil && existingBySlug.ID != id {
				return nil, fmt.Errorf("industry with slug '%s' already exists", req.Slug)
			}
		}
		industry.Slug = req.Slug
	}

	if req.Description != "" || req.Description == "" {
		industry.Description = sql.NullString{String: req.Description, Valid: req.Description != ""}
	}

	if req.IconURL != "" || req.IconURL == "" {
		industry.IconURL = sql.NullString{String: req.IconURL, Valid: req.IconURL != ""}
	}

	if req.IsActive != nil {
		industry.IsActive = *req.IsActive
	}

	// Update industry
	if err := s.repo.Update(ctx, industry); err != nil {
		return nil, fmt.Errorf("failed to update industry: %w", err)
	}

	// Invalidate cache
	s.invalidateCache()

	// Map to response
	response := &master.IndustryResponse{
		ID:          industry.ID,
		Name:        industry.Name,
		Slug:        industry.Slug,
		Description: industry.GetDescription(),
		IconURL:     industry.GetIconURL(),
		IsActive:    industry.IsActive,
	}

	return response, nil
}

// Delete deletes an industry if not referenced by companies
func (s *adminIndustryServiceImpl) Delete(ctx context.Context, id int64) error {
	// Check if industry exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("industry not found")
		}
		return fmt.Errorf("failed to get industry: %w", err)
	}

	// Check references (will be implemented with company repo)
	count, err := s.CountCompanyReferences(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check industry references: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("cannot delete industry: it is still referenced by %d companies", count)
	}

	// Delete industry (soft delete)
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete industry: %w", err)
	}

	// Invalidate cache
	s.invalidateCache()

	return nil
}

// CheckDuplicateName checks if an industry with the given name exists
func (s *adminIndustryServiceImpl) CheckDuplicateName(ctx context.Context, name string) (bool, error) {
	industry, err := s.repo.GetByName(ctx, name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	return industry != nil, nil
}

// CountCompanyReferences counts how many companies reference this industry
func (s *adminIndustryServiceImpl) CountCompanyReferences(ctx context.Context, id int64) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Table("companies").
		Where("industry_id = ?", id).
		Count(&count).Error
	return count, err
}

// invalidateCache invalidates all industry-related cache entries
func (s *adminIndustryServiceImpl) invalidateCache() {
	// Invalidate all industry cache keys
	cacheKeys := []string{
		"industries:all",
		"industries:active",
	}
	for _, key := range cacheKeys {
		s.cache.Delete(key)
	}
	// Also invalidate pattern-based keys (search, etc.)
	// This is a simplified version - in production you might want to track all cache keys
}
