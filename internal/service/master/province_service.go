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

// Error definitions for Province service
var (
	ErrProvinceNotFound  = errors.New("province not found")
	ErrProvinceInactive  = errors.New("province is not active")
	ErrInvalidProvinceID = errors.New("invalid province ID")
)

// provinceServiceImpl implements the ProvinceService interface
type provinceServiceImpl struct {
	repo  master.ProvinceRepository
	cache cache.Cache
}

// NewProvinceService creates a new instance of ProvinceService
func NewProvinceService(repo master.ProvinceRepository, cache cache.Cache) master.ProvinceService {
	return &provinceServiceImpl{
		repo:  repo,
		cache: cache,
	}
}

// GetAll retrieves all provinces with optional search
func (s *provinceServiceImpl) GetAll(ctx context.Context, search string) ([]master.ProvinceResponse, error) {
	// Trim search parameter
	search = strings.TrimSpace(search)

	// Generate cache key
	cacheKey := "provinces:all"
	if search != "" {
		cacheKey = fmt.Sprintf("provinces:search:%s", strings.ToLower(search))
	}

	// Check cache
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.([]master.ProvinceResponse), nil
	}

	// Query repository
	var provinces []master.Province
	var err error

	if search != "" {
		// Validate search length
		if len(search) < 2 {
			return []master.ProvinceResponse{}, nil
		}
		provinces, err = s.repo.Search(ctx, search)
	} else {
		provinces, err = s.repo.GetAll(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get provinces: %w", err)
	}

	// Map to response DTOs
	responses := s.mapToResponses(provinces)

	// Cache results (static data, cache for 24 hours)
	ttl := 24 * time.Hour
	if search != "" {
		ttl = 1 * time.Hour // Search results cache for 1 hour
	}
	s.cache.Set(cacheKey, responses, ttl)

	return responses, nil
}

// GetActive retrieves all active provinces with optional search
func (s *provinceServiceImpl) GetActive(ctx context.Context, search string) ([]master.ProvinceResponse, error) {
	// Trim search parameter
	search = strings.TrimSpace(search)

	// Generate cache key
	cacheKey := "provinces:active"
	if search != "" {
		cacheKey = fmt.Sprintf("provinces:active:search:%s", strings.ToLower(search))
	}

	// Check cache
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.([]master.ProvinceResponse), nil
	}

	// Query repository
	var provinces []master.Province
	var err error

	if search != "" {
		// Validate search length
		if len(search) < 2 {
			return []master.ProvinceResponse{}, nil
		}
		// Get all from search and filter active in service layer
		allProvinces, err := s.repo.Search(ctx, search)
		if err != nil {
			return nil, fmt.Errorf("failed to search provinces: %w", err)
		}
		// Filter only active
		provinces = make([]master.Province, 0)
		for _, province := range allProvinces {
			if province.IsActive {
				provinces = append(provinces, province)
			}
		}
	} else {
		provinces, err = s.repo.GetActive(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get active provinces: %w", err)
		}
	}

	// Map to response DTOs
	responses := s.mapToResponses(provinces)

	// Cache results (static data, cache for 24 hours)
	ttl := 24 * time.Hour
	if search != "" {
		ttl = 1 * time.Hour // Search results cache for 1 hour
	}
	s.cache.Set(cacheKey, responses, ttl)

	return responses, nil
}

// GetByID retrieves a province by ID
func (s *provinceServiceImpl) GetByID(ctx context.Context, id int64) (*master.ProvinceResponse, error) {
	// Validate ID
	if id <= 0 {
		return nil, ErrInvalidProvinceID
	}

	// Generate cache key
	cacheKey := fmt.Sprintf("province:id:%d", id)

	// Check cache
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.(*master.ProvinceResponse), nil
	}

	// Query repository
	province, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProvinceNotFound
		}
		return nil, fmt.Errorf("failed to get province: %w", err)
	}

	// Map to response DTO
	response := s.mapToResponse(province)

	// Cache result
	s.cache.Set(cacheKey, response, 24*time.Hour)

	return response, nil
}

// ValidateProvinceID checks if a province ID exists and is active
func (s *provinceServiceImpl) ValidateProvinceID(ctx context.Context, id int64) error {
	// Validate ID format
	if id <= 0 {
		return ErrInvalidProvinceID
	}

	// Check if province exists
	province, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProvinceNotFound
		}
		return fmt.Errorf("failed to validate province ID: %w", err)
	}

	// Check if active
	if !province.IsActive {
		return ErrProvinceInactive
	}

	return nil
}

// Private helper methods

// mapToResponses converts a slice of Province entities to ProvinceResponse DTOs
func (s *provinceServiceImpl) mapToResponses(provinces []master.Province) []master.ProvinceResponse {
	responses := make([]master.ProvinceResponse, len(provinces))
	for i, province := range provinces {
		responses[i] = master.ProvinceResponse{
			ID:       province.ID,
			Name:     province.Name,
			Code:     province.Code,
			IsActive: province.IsActive,
		}
	}
	return responses
}

// mapToResponse converts a single Province entity to ProvinceResponse DTO
func (s *provinceServiceImpl) mapToResponse(province *master.Province) *master.ProvinceResponse {
	return &master.ProvinceResponse{
		ID:       province.ID,
		Name:     province.Name,
		Code:     province.Code,
		IsActive: province.IsActive,
	}
}
