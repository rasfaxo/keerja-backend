package master

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"keerja-backend/internal/cache"
	"keerja-backend/internal/domain/master"
)

// adminCityServiceImpl implements AdminCityService
type adminCityServiceImpl struct {
	master.CityService // Embed base service for read operations
	repo               master.CityRepository
	db                 *gorm.DB // For counting references
	cache              cache.Cache
}

// NewAdminCityService creates a new AdminCityService
func NewAdminCityService(
	baseService master.CityService,
	repo master.CityRepository,
	db *gorm.DB,
	cache cache.Cache,
) master.AdminCityService {
	return &adminCityServiceImpl{
		CityService: baseService,
		repo:        repo,
		db:          db,
		cache:       cache,
	}
}

// Create creates a new city
func (s *adminCityServiceImpl) Create(ctx context.Context, req master.CreateCityRequest) (*master.CityResponse, error) {
	// Check duplicate name in province
	existing, err := s.repo.GetByNameAndProvinceID(ctx, req.Name, req.ProvinceID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check duplicate city name: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("city with name '%s' already exists in the province", req.Name)
	}

	// Create city entity
	city := &master.City{
		Name:       req.Name,
		Type:       req.Type,
		Code:       req.Code,
		ProvinceID: req.ProvinceID,
		IsActive:   req.IsActive,
	}

	if err := s.repo.Create(ctx, city); err != nil {
		return nil, fmt.Errorf("failed to create city: %w", err)
	}

	// Invalidate cache
	s.invalidateCache()

	// Get with province for response
	cityWithProvince, err := s.repo.GetWithProvince(ctx, city.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created city: %w", err)
	}

	// Map to response
	response := s.mapToResponseWithProvince(cityWithProvince)
	return response, nil
}

// Update updates an existing city
func (s *adminCityServiceImpl) Update(ctx context.Context, id int64, req master.UpdateCityRequest) (*master.CityResponse, error) {
	// Get existing city
	city, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("city not found")
		}
		return nil, fmt.Errorf("failed to get city: %w", err)
	}

	// Update fields
	if req.Name != "" {
		provinceID := city.ProvinceID
		if req.ProvinceID != nil {
			provinceID = *req.ProvinceID
		}
		// Check duplicate name if name or province is being changed
		if req.Name != city.Name || (req.ProvinceID != nil && *req.ProvinceID != city.ProvinceID) {
			existing, err := s.repo.GetByNameAndProvinceID(ctx, req.Name, provinceID)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("failed to check duplicate city name: %w", err)
			}
			if existing != nil && existing.ID != id {
				return nil, fmt.Errorf("city with name '%s' already exists in the province", req.Name)
			}
		}
		city.Name = req.Name
	}

	if req.Type != "" {
		city.Type = req.Type
	}

	if req.Code != "" {
		city.Code = req.Code
	}

	if req.ProvinceID != nil {
		city.ProvinceID = *req.ProvinceID
	}

	if req.IsActive != nil {
		city.IsActive = *req.IsActive
	}

	// Update city
	if err := s.repo.Update(ctx, city); err != nil {
		return nil, fmt.Errorf("failed to update city: %w", err)
	}

	// Invalidate cache
	s.invalidateCache()

	// Get with province for response
	cityWithProvince, err := s.repo.GetWithProvince(ctx, city.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated city: %w", err)
	}

	// Map to response
	response := s.mapToResponseWithProvince(cityWithProvince)
	return response, nil
}

// Delete deletes a city if not referenced by districts or companies
func (s *adminCityServiceImpl) Delete(ctx context.Context, id int64) error {
	// Check if city exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("city not found")
		}
		return fmt.Errorf("failed to get city: %w", err)
	}

	// Check references
	districts, companies, err := s.CountReferences(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check city references: %w", err)
	}
	if districts > 0 || companies > 0 {
		return fmt.Errorf("cannot delete city: it is still referenced by %d districts and %d companies", districts, companies)
	}

	// Delete city
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete city: %w", err)
	}

	// Invalidate cache
	s.invalidateCache()

	return nil
}

// CheckDuplicateNameInProvince checks if a city with the given name exists in the province
func (s *adminCityServiceImpl) CheckDuplicateNameInProvince(ctx context.Context, name string, provinceID int64) (bool, error) {
	city, err := s.repo.GetByNameAndProvinceID(ctx, name, provinceID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	return city != nil, nil
}

// CountReferences counts how many districts and companies reference this city
func (s *adminCityServiceImpl) CountReferences(ctx context.Context, id int64) (districts int64, companies int64, err error) {
	// Count districts
	err = s.db.WithContext(ctx).Table("districts").
		Where("city_id = ?", id).
		Count(&districts).Error
	if err != nil {
		return 0, 0, err
	}

	// Count companies
	err = s.db.WithContext(ctx).Table("companies").
		Where("city_id = ?", id).
		Count(&companies).Error
	if err != nil {
		return 0, 0, err
	}

	return districts, companies, nil
}

// invalidateCache invalidates all city-related cache entries
func (s *adminCityServiceImpl) invalidateCache() {
	cacheKeys := []string{
		"cities:province:",
		"city:id:",
	}
	for _, key := range cacheKeys {
		// Note: In production, you might want to track all cache keys with patterns
		s.cache.Delete(key)
	}
}

// mapToResponseWithProvince converts a City entity with province to CityResponse DTO
func (s *adminCityServiceImpl) mapToResponseWithProvince(city *master.City) *master.CityResponse {
	response := master.CityResponse{
		ID:         city.ID,
		Name:       city.Name,
		FullName:   city.GetFullName(),
		Type:       city.Type,
		Code:       city.Code,
		ProvinceID: city.ProvinceID,
		IsActive:   city.IsActive,
	}

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

