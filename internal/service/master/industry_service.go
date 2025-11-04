package master

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"keerja-backend/internal/cache"
	"keerja-backend/internal/domain/master"
)

// Error definitions for Industry service
var (
	ErrIndustryNotFound  = errors.New("industry not found")
	ErrIndustryInactive  = errors.New("industry is not active")
	ErrInvalidIndustryID = errors.New("invalid industry ID")
)

// industryServiceImpl implements the IndustryService interface
type industryServiceImpl struct {
	repo  master.IndustryRepository
	cache cache.Cache
}

// NewIndustryService creates a new instance of IndustryService
func NewIndustryService(repo master.IndustryRepository, cache cache.Cache) master.IndustryService {
	return &industryServiceImpl{
		repo:  repo,
		cache: cache,
	}
}

// GetAll retrieves all industries with optional search
func (s *industryServiceImpl) GetAll(ctx context.Context, search string) ([]master.IndustryResponse, error) {
	// Trim search parameter
	search = strings.TrimSpace(search)

	// Generate cache key
	cacheKey := "industries:all"
	if search != "" {
		cacheKey = fmt.Sprintf("industries:search:%s", strings.ToLower(search))
	}

	// Check cache
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.([]master.IndustryResponse), nil
	}

	// Query repository
	var industries []master.Industry
	var err error

	if search != "" {
		// Validate search length
		if len(search) < 2 {
			return []master.IndustryResponse{}, nil
		}
		industries, err = s.repo.Search(ctx, search)
	} else {
		industries, err = s.repo.GetAll(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get industries: %w", err)
	}

	// Map to response DTOs
	responses := s.mapToResponses(industries)

	// Cache results
	ttl := 24 * time.Hour // Static data, cache for 24 hours
	if search != "" {
		ttl = 1 * time.Hour // Search results cache for 1 hour
	}
	s.cache.Set(cacheKey, responses, ttl)

	return responses, nil
}

// GetActive retrieves all active industries with optional search
func (s *industryServiceImpl) GetActive(ctx context.Context, search string) ([]master.IndustryResponse, error) {
	// Trim search parameter
	search = strings.TrimSpace(search)

	// Generate cache key
	cacheKey := "industries:active"
	if search != "" {
		cacheKey = fmt.Sprintf("industries:active:search:%s", strings.ToLower(search))
	}

	// Check cache
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.([]master.IndustryResponse), nil
	}

	// Query repository
	var industries []master.Industry
	var err error

	if search != "" {
		// Validate search length
		if len(search) < 2 {
			return []master.IndustryResponse{}, nil
		}
		// Get all from search and filter active in service layer
		allIndustries, err := s.repo.Search(ctx, search)
		if err != nil {
			return nil, fmt.Errorf("failed to search industries: %w", err)
		}
		// Filter only active
		industries = make([]master.Industry, 0)
		for _, industry := range allIndustries {
			if industry.IsActive {
				industries = append(industries, industry)
			}
		}
	} else {
		industries, err = s.repo.GetActive(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get active industries: %w", err)
		}
	}

	// Map to response DTOs
	responses := s.mapToResponses(industries)

	// Cache results
	ttl := 24 * time.Hour // Static data, cache for 24 hours
	if search != "" {
		ttl = 1 * time.Hour // Search results cache for 1 hour
	}
	s.cache.Set(cacheKey, responses, ttl)

	return responses, nil
}

// GetByID retrieves an industry by ID
func (s *industryServiceImpl) GetByID(ctx context.Context, id int64) (*master.IndustryResponse, error) {
	// Validate ID
	if id <= 0 {
		return nil, ErrInvalidIndustryID
	}

	// Generate cache key
	cacheKey := fmt.Sprintf("industry:id:%d", id)

	// Check cache
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.(*master.IndustryResponse), nil
	}

	// Query repository
	industry, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrIndustryNotFound
		}
		return nil, fmt.Errorf("failed to get industry: %w", err)
	}

	// Map to response DTO
	response := s.mapToResponse(industry)

	// Cache result
	s.cache.Set(cacheKey, response, 24*time.Hour)

	return response, nil
}

// ValidateIndustryID checks if an industry ID exists and is active
func (s *industryServiceImpl) ValidateIndustryID(ctx context.Context, id int64) error {
	// Validate ID format
	if id <= 0 {
		return ErrInvalidIndustryID
	}

	// Check if industry exists
	industry, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrIndustryNotFound
		}
		return fmt.Errorf("failed to validate industry ID: %w", err)
	}

	// Check if active
	if !industry.IsActive {
		return ErrIndustryInactive
	}

	// Check if soft deleted
	if industry.IsDeleted() {
		return ErrIndustryNotFound
	}

	return nil
}

// Private helper methods

// mapToResponses converts a slice of Industry entities to IndustryResponse DTOs
func (s *industryServiceImpl) mapToResponses(industries []master.Industry) []master.IndustryResponse {
	responses := make([]master.IndustryResponse, len(industries))
	for i, industry := range industries {
		responses[i] = master.IndustryResponse{
			ID:          industry.ID,
			Name:        industry.Name,
			Slug:        industry.Slug,
			Description: industry.GetDescription(),
			IconURL:     industry.GetIconURL(),
			IsActive:    industry.IsActive,
		}
	}
	return responses
}

// mapToResponse converts a single Industry entity to IndustryResponse DTO
func (s *industryServiceImpl) mapToResponse(industry *master.Industry) *master.IndustryResponse {
	return &master.IndustryResponse{
		ID:          industry.ID,
		Name:        industry.Name,
		Slug:        industry.Slug,
		Description: industry.GetDescription(),
		IconURL:     industry.GetIconURL(),
		IsActive:    industry.IsActive,
	}
}
