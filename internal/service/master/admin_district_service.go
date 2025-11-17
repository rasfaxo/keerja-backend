package master

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"keerja-backend/internal/cache"
	"keerja-backend/internal/domain/master"
)

// adminDistrictServiceImpl implements AdminDistrictService
type adminDistrictServiceImpl struct {
	master.DistrictService // Embed base service for read operations
	repo                   master.DistrictRepository
	db                     *gorm.DB // For counting references
	cache                  cache.Cache
}

// NewAdminDistrictService creates a new AdminDistrictService
func NewAdminDistrictService(
	baseService master.DistrictService,
	repo master.DistrictRepository,
	db *gorm.DB,
	cache cache.Cache,
) master.AdminDistrictService {
	return &adminDistrictServiceImpl{
		DistrictService: baseService,
		repo:            repo,
		db:              db,
		cache:           cache,
	}
}

// Create creates a new district
func (s *adminDistrictServiceImpl) Create(ctx context.Context, req master.CreateDistrictRequest) (*master.DistrictResponse, error) {
	// Check duplicate name in city
	existing, err := s.repo.GetByNameAndCityID(ctx, req.Name, req.CityID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check duplicate district name: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("district with name '%s' already exists in the city", req.Name)
	}

	// Create district entity
	district := &master.District{
		Name:       req.Name,
		Code:       req.Code,
		PostalCode: sql.NullString{String: req.PostalCode, Valid: req.PostalCode != ""},
		CityID:     req.CityID,
		IsActive:   req.IsActive,
	}

	if err := s.repo.Create(ctx, district); err != nil {
		return nil, fmt.Errorf("failed to create district: %w", err)
	}

	// Invalidate cache
	s.invalidateCache()

	// Get with full location for response
	districtWithLocation, err := s.repo.GetWithFullLocation(ctx, district.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created district: %w", err)
	}

	// Map to response
	response := s.mapToResponseWithFullLocation(districtWithLocation)
	return response, nil
}

// Update updates an existing district
func (s *adminDistrictServiceImpl) Update(ctx context.Context, id int64, req master.UpdateDistrictRequest) (*master.DistrictResponse, error) {
	// Get existing district
	district, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("district not found")
		}
		return nil, fmt.Errorf("failed to get district: %w", err)
	}

	// Update fields
	if req.Name != "" {
		cityID := district.CityID
		if req.CityID != nil {
			cityID = *req.CityID
		}
		// Check duplicate name if name or city is being changed
		if req.Name != district.Name || (req.CityID != nil && *req.CityID != district.CityID) {
			existing, err := s.repo.GetByNameAndCityID(ctx, req.Name, cityID)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("failed to check duplicate district name: %w", err)
			}
			if existing != nil && existing.ID != id {
				return nil, fmt.Errorf("district with name '%s' already exists in the city", req.Name)
			}
		}
		district.Name = req.Name
	}

	if req.Code != "" {
		district.Code = req.Code
	}

	if req.PostalCode != "" || req.PostalCode == "" {
		district.PostalCode = sql.NullString{String: req.PostalCode, Valid: req.PostalCode != ""}
	}

	if req.CityID != nil {
		district.CityID = *req.CityID
	}

	if req.IsActive != nil {
		district.IsActive = *req.IsActive
	}

	// Update district
	if err := s.repo.Update(ctx, district); err != nil {
		return nil, fmt.Errorf("failed to update district: %w", err)
	}

	// Invalidate cache
	s.invalidateCache()

	// Get with full location for response
	districtWithLocation, err := s.repo.GetWithFullLocation(ctx, district.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated district: %w", err)
	}

	// Map to response
	response := s.mapToResponseWithFullLocation(districtWithLocation)
	return response, nil
}

// Delete deletes a district if not referenced by companies
func (s *adminDistrictServiceImpl) Delete(ctx context.Context, id int64) error {
	// Check if district exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("district not found")
		}
		return fmt.Errorf("failed to get district: %w", err)
	}

	// Check references
	companies, err := s.CountCompanyReferences(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check district references: %w", err)
	}
	if companies > 0 {
		return fmt.Errorf("cannot delete district: it is still referenced by %d companies", companies)
	}

	// Delete district
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete district: %w", err)
	}

	// Invalidate cache
	s.invalidateCache()

	return nil
}

// CheckDuplicateNameInCity checks if a district with the given name exists in the city
func (s *adminDistrictServiceImpl) CheckDuplicateNameInCity(ctx context.Context, name string, cityID int64) (bool, error) {
	district, err := s.repo.GetByNameAndCityID(ctx, name, cityID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	return district != nil, nil
}

// CountCompanyReferences counts how many companies reference this district
func (s *adminDistrictServiceImpl) CountCompanyReferences(ctx context.Context, id int64) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Table("companies").
		Where("district_id = ?", id).
		Count(&count).Error
	return count, err
}

// invalidateCache invalidates all district-related cache entries
func (s *adminDistrictServiceImpl) invalidateCache() {
	cacheKeys := []string{
		"districts:city:",
		"district:id:",
	}
	for _, key := range cacheKeys {
		s.cache.Delete(key)
	}
}

// mapToResponseWithFullLocation converts a District entity with full hierarchy to DistrictResponse DTO
func (s *adminDistrictServiceImpl) mapToResponseWithFullLocation(district *master.District) *master.DistrictResponse {
	response := master.DistrictResponse{
		ID:         district.ID,
		Name:       district.Name,
		Code:       district.Code,
		PostalCode: district.GetPostalCode(),
		CityID:     district.CityID,
		IsActive:   district.IsActive,
	}

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
