package master

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"keerja-backend/internal/cache"
	"keerja-backend/internal/domain/master"
)

// Error definitions for CompanySize service
var (
	ErrCompanySizeNotFound  = errors.New("company size not found")
	ErrCompanySizeInactive  = errors.New("company size is not active")
	ErrInvalidCompanySizeID = errors.New("invalid company size ID")
)

// companySizeServiceImpl implements the CompanySizeService interface
type companySizeServiceImpl struct {
	repo  master.CompanySizeRepository
	cache cache.Cache
}

// NewCompanySizeService creates a new instance of CompanySizeService
func NewCompanySizeService(repo master.CompanySizeRepository, cache cache.Cache) master.CompanySizeService {
	return &companySizeServiceImpl{
		repo:  repo,
		cache: cache,
	}
}

// GetAll retrieves all company sizes
func (s *companySizeServiceImpl) GetAll(ctx context.Context) ([]master.CompanySizeResponse, error) {
	// Generate cache key
	cacheKey := "company_sizes:all"

	// Check cache
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.([]master.CompanySizeResponse), nil
	}

	// Query repository
	companySizes, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get company sizes: %w", err)
	}

	// Map to response DTOs
	responses := s.mapToResponses(companySizes)

	// Cache results (static data, cache for 24 hours)
	s.cache.Set(cacheKey, responses, 24*time.Hour)

	return responses, nil
}

// GetActive retrieves all active company sizes
func (s *companySizeServiceImpl) GetActive(ctx context.Context) ([]master.CompanySizeResponse, error) {
	// Generate cache key
	cacheKey := "company_sizes:active"

	// Check cache
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.([]master.CompanySizeResponse), nil
	}

	// Query repository
	companySizes, err := s.repo.GetActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active company sizes: %w", err)
	}

	// Map to response DTOs
	responses := s.mapToResponses(companySizes)

	// Cache results (static data, cache for 24 hours)
	s.cache.Set(cacheKey, responses, 24*time.Hour)

	return responses, nil
}

// GetByID retrieves a company size by ID
func (s *companySizeServiceImpl) GetByID(ctx context.Context, id int64) (*master.CompanySizeResponse, error) {
	// Validate ID
	if id <= 0 {
		return nil, ErrInvalidCompanySizeID
	}

	// Generate cache key
	cacheKey := fmt.Sprintf("company_size:id:%d", id)

	// Check cache
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.(*master.CompanySizeResponse), nil
	}

	// Query repository
	companySize, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCompanySizeNotFound
		}
		return nil, fmt.Errorf("failed to get company size: %w", err)
	}

	// Map to response DTO
	response := s.mapToResponse(companySize)

	// Cache result
	s.cache.Set(cacheKey, response, 24*time.Hour)

	return response, nil
}

// ValidateCompanySizeID checks if a company size ID exists and is active
func (s *companySizeServiceImpl) ValidateCompanySizeID(ctx context.Context, id int64) error {
	// Validate ID format
	if id <= 0 {
		return ErrInvalidCompanySizeID
	}

	// Check if company size exists
	companySize, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCompanySizeNotFound
		}
		return fmt.Errorf("failed to validate company size ID: %w", err)
	}

	// Check if active
	if !companySize.IsActive {
		return ErrCompanySizeInactive
	}

	return nil
}

// Private helper methods

// mapToResponses converts a slice of CompanySize entities to CompanySizeResponse DTOs
func (s *companySizeServiceImpl) mapToResponses(companySizes []master.CompanySize) []master.CompanySizeResponse {
	responses := make([]master.CompanySizeResponse, len(companySizes))
	for i, size := range companySizes {
		responses[i] = master.CompanySizeResponse{
			ID:           size.ID,
			Label:        size.Label,
			MinEmployees: size.MinEmployees,
			MaxEmployees: s.getMaxEmployeesPointer(&size),
			IsActive:     size.IsActive,
		}
	}
	return responses
}

// mapToResponse converts a single CompanySize entity to CompanySizeResponse DTO
func (s *companySizeServiceImpl) mapToResponse(size *master.CompanySize) *master.CompanySizeResponse {
	return &master.CompanySizeResponse{
		ID:           size.ID,
		Label:        size.Label,
		MinEmployees: size.MinEmployees,
		MaxEmployees: s.getMaxEmployeesPointer(size),
		IsActive:     size.IsActive,
	}
}

// getMaxEmployeesPointer returns a pointer to max employees value, or nil if unlimited
func (s *companySizeServiceImpl) getMaxEmployeesPointer(size *master.CompanySize) *int {
	maxEmp := size.GetMaxEmployees()
	if maxEmp == -1 {
		return nil // Unlimited
	}
	return &maxEmp
}
