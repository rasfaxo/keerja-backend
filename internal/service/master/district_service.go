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

// Error definitions for District service
var (
	ErrDistrictNotFound         = errors.New("district not found")
	ErrDistrictInactive         = errors.New("district is not active")
	ErrInvalidDistrictID        = errors.New("invalid district ID")
	ErrDistrictCityIDMismatch   = errors.New("district does not belong to the specified city")
	ErrDistrictCityNotFound     = errors.New("district's city not found")
	ErrInvalidLocationHierarchy = errors.New("invalid location hierarchy: province, city, and district do not match")
)

// districtServiceImpl implements the DistrictService interface
type districtServiceImpl struct {
	repo         master.DistrictRepository
	cityRepo     master.CityRepository
	provinceRepo master.ProvinceRepository
	cache        cache.Cache
}

// NewDistrictService creates a new instance of DistrictService
func NewDistrictService(
	repo master.DistrictRepository,
	cityRepo master.CityRepository,
	provinceRepo master.ProvinceRepository,
	cache cache.Cache,
) master.DistrictService {
	return &districtServiceImpl{
		repo:         repo,
		cityRepo:     cityRepo,
		provinceRepo: provinceRepo,
		cache:        cache,
	}
}

// GetByCityID retrieves all districts in a city with optional search
func (s *districtServiceImpl) GetByCityID(ctx context.Context, cityID int64, search string) ([]master.DistrictResponse, error) {
	// Validate city ID
	if cityID <= 0 {
		return nil, ErrInvalidCityID
	}

	// Trim search parameter
	search = strings.TrimSpace(search)

	// Generate cache key
	cacheKey := fmt.Sprintf("districts:city:%d", cityID)
	if search != "" {
		cacheKey = fmt.Sprintf("districts:city:%d:search:%s", cityID, strings.ToLower(search))
	}

	// Check cache
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.([]master.DistrictResponse), nil
	}

	// Validate city exists
	cityExists, err := s.cityRepo.ExistsByID(ctx, cityID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate city: %w", err)
	}
	if !cityExists {
		return nil, ErrCityNotFound
	}

	// Query repository
	var districts []master.District
	if search != "" {
		// Validate search length
		if len(search) < 2 {
			return []master.DistrictResponse{}, nil
		}
		districts, err = s.repo.Search(ctx, search, &cityID)
	} else {
		districts, err = s.repo.GetByCityID(ctx, cityID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get districts: %w", err)
	}

	// Map to response DTOs
	responses := s.mapToResponses(districts)

	// Cache results (semi-static data, cache for 6 hours)
	ttl := 6 * time.Hour
	if search != "" {
		ttl = 1 * time.Hour // Search results cache for 1 hour
	}
	s.cache.Set(cacheKey, responses, ttl)

	return responses, nil
}

// GetActiveByCityID retrieves all active districts in a city with optional search
func (s *districtServiceImpl) GetActiveByCityID(ctx context.Context, cityID int64, search string) ([]master.DistrictResponse, error) {
	// Validate city ID
	if cityID <= 0 {
		return nil, ErrInvalidCityID
	}

	// Trim search parameter
	search = strings.TrimSpace(search)

	// Generate cache key
	cacheKey := fmt.Sprintf("districts:city:%d:active", cityID)
	if search != "" {
		cacheKey = fmt.Sprintf("districts:city:%d:active:search:%s", cityID, strings.ToLower(search))
	}

	// Check cache
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.([]master.DistrictResponse), nil
	}

	// Validate city exists
	cityExists, err := s.cityRepo.ExistsByID(ctx, cityID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate city: %w", err)
	}
	if !cityExists {
		return nil, ErrCityNotFound
	}

	// Query repository
	var districts []master.District
	if search != "" {
		// Validate search length
		if len(search) < 2 {
			return []master.DistrictResponse{}, nil
		}
		// Get all from search and filter active in service layer
		allDistricts, err := s.repo.Search(ctx, search, &cityID)
		if err != nil {
			return nil, fmt.Errorf("failed to search districts: %w", err)
		}
		// Filter only active
		districts = make([]master.District, 0)
		for _, district := range allDistricts {
			if district.IsActive {
				districts = append(districts, district)
			}
		}
	} else {
		districts, err = s.repo.GetActiveByCityID(ctx, cityID)
		if err != nil {
			return nil, fmt.Errorf("failed to get active districts: %w", err)
		}
	}

	// Map to response DTOs
	responses := s.mapToResponses(districts)

	// Cache results (semi-static data, cache for 6 hours)
	ttl := 6 * time.Hour
	if search != "" {
		ttl = 1 * time.Hour // Search results cache for 1 hour
	}
	s.cache.Set(cacheKey, responses, ttl)

	return responses, nil
}

// GetByID retrieves a district by ID with full location hierarchy
func (s *districtServiceImpl) GetByID(ctx context.Context, id int64) (*master.DistrictResponse, error) {
	// Validate ID
	if id <= 0 {
		return nil, ErrInvalidDistrictID
	}

	// Generate cache key
	cacheKey := fmt.Sprintf("district:id:%d", id)

	// Check cache
	if cached, ok := s.cache.Get(cacheKey); ok {
		return cached.(*master.DistrictResponse), nil
	}

	// Query repository with full location hierarchy preloaded
	district, err := s.repo.GetWithFullLocation(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDistrictNotFound
		}
		return nil, fmt.Errorf("failed to get district: %w", err)
	}

	// Map to response DTO with full hierarchy
	response := s.mapToResponseWithFullLocation(district)

	// Cache result
	s.cache.Set(cacheKey, response, 6*time.Hour)

	return response, nil
}

// ValidateDistrictID checks if a district ID exists, is active, and belongs to the given city
func (s *districtServiceImpl) ValidateDistrictID(ctx context.Context, districtID, cityID int64) error {
	// Validate ID formats
	if districtID <= 0 {
		return ErrInvalidDistrictID
	}
	if cityID <= 0 {
		return ErrInvalidCityID
	}

	// Check if district exists
	district, err := s.repo.GetByID(ctx, districtID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDistrictNotFound
		}
		return fmt.Errorf("failed to validate district ID: %w", err)
	}

	// Check if active
	if !district.IsActive {
		return ErrDistrictInactive
	}

	// Check if district belongs to the specified city
	if district.CityID != cityID {
		return ErrDistrictCityIDMismatch
	}

	// Validate city exists and is active
	city, err := s.cityRepo.GetByID(ctx, cityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDistrictCityNotFound
		}
		return fmt.Errorf("failed to validate district's city: %w", err)
	}

	if !city.IsActive {
		return ErrCityInactive
	}

	return nil
}

// ValidateLocationHierarchy validates the complete location hierarchy (province -> city -> district)
func (s *districtServiceImpl) ValidateLocationHierarchy(ctx context.Context, provinceID, cityID, districtID int64) error {
	// Validate all IDs
	if provinceID <= 0 {
		return ErrInvalidProvinceID
	}
	if cityID <= 0 {
		return ErrInvalidCityID
	}
	if districtID <= 0 {
		return ErrInvalidDistrictID
	}

	// Step 1: Validate district exists and is active
	district, err := s.repo.GetByID(ctx, districtID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDistrictNotFound
		}
		return fmt.Errorf("failed to validate district: %w", err)
	}

	if !district.IsActive {
		return ErrDistrictInactive
	}

	// Step 2: Check if district belongs to city
	if district.CityID != cityID {
		return ErrInvalidLocationHierarchy
	}

	// Step 3: Validate city exists and is active
	city, err := s.cityRepo.GetByID(ctx, cityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCityNotFound
		}
		return fmt.Errorf("failed to validate city: %w", err)
	}

	if !city.IsActive {
		return ErrInvalidLocationHierarchy
	}

	// Step 4: Check if city belongs to province
	if city.ProvinceID != provinceID {
		return ErrInvalidLocationHierarchy
	}

	// Step 5: Validate province exists and is active
	province, err := s.provinceRepo.GetByID(ctx, provinceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProvinceNotFound
		}
		return fmt.Errorf("failed to validate province: %w", err)
	}

	if !province.IsActive {
		return ErrInvalidLocationHierarchy
	}

	// All validations passed
	return nil
}

// Private helper methods

// mapToResponses converts a slice of District entities to DistrictResponse DTOs
func (s *districtServiceImpl) mapToResponses(districts []master.District) []master.DistrictResponse {
	responses := make([]master.DistrictResponse, len(districts))
	for i, district := range districts {
		responses[i] = s.buildDistrictResponse(&district)
	}
	return responses
}

// mapToResponseWithFullLocation converts a District entity with full hierarchy to DistrictResponse DTO
func (s *districtServiceImpl) mapToResponseWithFullLocation(district *master.District) *master.DistrictResponse {
	response := s.buildDistrictResponse(district)

	// Add city info if preloaded
	if district.City != nil {
		response.City = &master.CityResponse{
			ID:         district.City.ID,
			Name:       district.City.Name,
			FullName:   district.City.GetFullName(),
			Type:       district.City.Type,
			Code:       district.City.Code,
			ProvinceID: district.City.ProvinceID,
			IsActive:   district.City.IsActive,
		}

		// Add province info if preloaded through city
		if district.City.Province != nil {
			response.City.Province = &master.ProvinceResponse{
				ID:       district.City.Province.ID,
				Name:     district.City.Province.Name,
				Code:     district.City.Province.Code,
				IsActive: district.City.Province.IsActive,
			}

			// Build full location path
			response.FullLocationPath = district.GetFullLocationPath()
		}
	}

	return &response
}

// buildDistrictResponse builds a DistrictResponse from a District entity
func (s *districtServiceImpl) buildDistrictResponse(district *master.District) master.DistrictResponse {
	return master.DistrictResponse{
		ID:               district.ID,
		Name:             district.Name,
		Code:             district.Code,
		PostalCode:       district.GetPostalCode(),
		CityID:           district.CityID,
		City:             nil, // Will be set separately if needed
		FullLocationPath: "",  // Will be set separately if needed
		IsActive:         district.IsActive,
	}
}
