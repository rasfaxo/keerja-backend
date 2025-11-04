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

// Error definitions for City service
var (
	ErrCityNotFound           = errors.New("city not found")
	ErrCityInactive           = errors.New("city is not active")
	ErrInvalidCityID          = errors.New("invalid city ID")
	ErrCityProvinceIDMismatch = errors.New("city does not belong to the specified province")
	ErrCityProvinceNotFound   = errors.New("city's province not found")
)

// cityServiceImpl implements the CityService interface
type cityServiceImpl struct {
	repo         master.CityRepository
	provinceRepo master.ProvinceRepository
	cache        cache.Cache
}

// NewCityService creates a new instance of CityService
func NewCityService(
	repo master.CityRepository,
	provinceRepo master.ProvinceRepository,
	cache cache.Cache,
) master.CityService {
	return &cityServiceImpl{
		repo:         repo,
		provinceRepo: provinceRepo,
		cache:        cache,
	}
}

// GetByProvinceID retrieves all cities in a province with optional search
func (s *cityServiceImpl) GetByProvinceID(ctx context.Context, provinceID int64, search string) ([]master.CityResponse, error) {
	// Validate province ID
	if provinceID <= 0 {
		return nil, ErrInvalidProvinceID
	}

	// Trim search parameter
	search = strings.TrimSpace(search)

	// Generate cache key
	cacheKey := fmt.Sprintf("cities:province:%d", provinceID)
	if search != "" {
		cacheKey = fmt.Sprintf("cities:province:%d:search:%s", provinceID, strings.ToLower(search))
	}

	// Check cache
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.([]master.CityResponse), nil
	}

	// Validate province exists
	provinceExists, err := s.provinceRepo.ExistsByID(ctx, provinceID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate province: %w", err)
	}
	if !provinceExists {
		return nil, ErrProvinceNotFound
	}

	// Query repository
	var cities []master.City
	if search != "" {
		// Validate search length
		if len(search) < 2 {
			return []master.CityResponse{}, nil
		}
		cities, err = s.repo.Search(ctx, search, &provinceID)
	} else {
		cities, err = s.repo.GetByProvinceID(ctx, provinceID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get cities: %w", err)
	}

	// Map to response DTOs
	responses := s.mapToResponses(cities)

	// Cache results (semi-static data, cache for 12 hours)
	ttl := 12 * time.Hour
	if search != "" {
		ttl = 1 * time.Hour // Search results cache for 1 hour
	}
	s.cache.Set(cacheKey, responses, ttl)

	return responses, nil
}

// GetActiveByProvinceID retrieves all active cities in a province with optional search
func (s *cityServiceImpl) GetActiveByProvinceID(ctx context.Context, provinceID int64, search string) ([]master.CityResponse, error) {
	// Validate province ID
	if provinceID <= 0 {
		return nil, ErrInvalidProvinceID
	}

	// Trim search parameter
	search = strings.TrimSpace(search)

	// Generate cache key
	cacheKey := fmt.Sprintf("cities:province:%d:active", provinceID)
	if search != "" {
		cacheKey = fmt.Sprintf("cities:province:%d:active:search:%s", provinceID, strings.ToLower(search))
	}

	// Check cache
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.([]master.CityResponse), nil
	}

	// Validate province exists
	provinceExists, err := s.provinceRepo.ExistsByID(ctx, provinceID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate province: %w", err)
	}
	if !provinceExists {
		return nil, ErrProvinceNotFound
	}

	// Query repository
	var cities []master.City
	if search != "" {
		// Validate search length
		if len(search) < 2 {
			return []master.CityResponse{}, nil
		}
		// Get all from search and filter active in service layer
		allCities, err := s.repo.Search(ctx, search, &provinceID)
		if err != nil {
			return nil, fmt.Errorf("failed to search cities: %w", err)
		}
		// Filter only active
		cities = make([]master.City, 0)
		for _, city := range allCities {
			if city.IsActive {
				cities = append(cities, city)
			}
		}
	} else {
		cities, err = s.repo.GetActiveByProvinceID(ctx, provinceID)
		if err != nil {
			return nil, fmt.Errorf("failed to get active cities: %w", err)
		}
	}

	// Map to response DTOs
	responses := s.mapToResponses(cities)

	// Cache results (semi-static data, cache for 12 hours)
	ttl := 12 * time.Hour
	if search != "" {
		ttl = 1 * time.Hour // Search results cache for 1 hour
	}
	s.cache.Set(cacheKey, responses, ttl)

	return responses, nil
}

// GetByID retrieves a city by ID with province info
func (s *cityServiceImpl) GetByID(ctx context.Context, id int64) (*master.CityResponse, error) {
	// Validate ID
	if id <= 0 {
		return nil, ErrInvalidCityID
	}

	// Generate cache key
	cacheKey := fmt.Sprintf("city:id:%d", id)

	// Check cache
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.(*master.CityResponse), nil
	}

	// Query repository with province preloaded
	city, err := s.repo.GetWithProvince(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCityNotFound
		}
		return nil, fmt.Errorf("failed to get city: %w", err)
	}

	// Map to response DTO
	response := s.mapToResponseWithProvince(city)

	// Cache result
	s.cache.Set(cacheKey, response, 12*time.Hour)

	return response, nil
}

// ValidateCityID checks if a city ID exists, is active, and belongs to the given province
func (s *cityServiceImpl) ValidateCityID(ctx context.Context, cityID, provinceID int64) error {
	// Validate ID formats
	if cityID <= 0 {
		return ErrInvalidCityID
	}
	if provinceID <= 0 {
		return ErrInvalidProvinceID
	}

	// Check if city exists
	city, err := s.repo.GetByID(ctx, cityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCityNotFound
		}
		return fmt.Errorf("failed to validate city ID: %w", err)
	}

	// Check if active
	if !city.IsActive {
		return ErrCityInactive
	}

	// Check if city belongs to the specified province
	if city.ProvinceID != provinceID {
		return ErrCityProvinceIDMismatch
	}

	// Validate province exists and is active
	province, err := s.provinceRepo.GetByID(ctx, provinceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCityProvinceNotFound
		}
		return fmt.Errorf("failed to validate city's province: %w", err)
	}

	if !province.IsActive {
		return ErrProvinceInactive
	}

	return nil
}

// Private helper methods

// mapToResponses converts a slice of City entities to CityResponse DTOs
func (s *cityServiceImpl) mapToResponses(cities []master.City) []master.CityResponse {
	responses := make([]master.CityResponse, len(cities))
	for i, city := range cities {
		responses[i] = s.buildCityResponse(&city)
	}
	return responses
}

// mapToResponseWithProvince converts a City entity with province to CityResponse DTO
func (s *cityServiceImpl) mapToResponseWithProvince(city *master.City) *master.CityResponse {
	response := s.buildCityResponse(city)

	// Add province info if preloaded
	if city.Province != nil {
		response.Province = &master.ProvinceResponse{
			ID:       city.Province.ID,
			Name:     city.Province.Name,
			Code:     city.Province.Code,
			IsActive: city.Province.IsActive,
		}
	}

	return &response
}

// buildCityResponse builds a CityResponse from a City entity
func (s *cityServiceImpl) buildCityResponse(city *master.City) master.CityResponse {
	return master.CityResponse{
		ID:         city.ID,
		Name:       city.Name,
		FullName:   city.GetFullName(),
		Type:       city.Type,
		Code:       city.Code,
		ProvinceID: city.ProvinceID,
		Province:   nil, // Will be set separately if needed
		IsActive:   city.IsActive,
	}
}
