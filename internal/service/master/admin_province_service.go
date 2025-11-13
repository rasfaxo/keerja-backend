package master

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"keerja-backend/internal/cache"
	"keerja-backend/internal/domain/master"
)

// adminProvinceServiceImpl implements AdminProvinceService
type adminProvinceServiceImpl struct {
	master.ProvinceService // Embed base service for read operations
	repo                   master.ProvinceRepository
	db                     *gorm.DB // For counting references
	cache                  cache.Cache
}

// NewAdminProvinceService creates a new AdminProvinceService
func NewAdminProvinceService(
	baseService master.ProvinceService,
	repo master.ProvinceRepository,
	db *gorm.DB,
	cache cache.Cache,
) master.AdminProvinceService {
	return &adminProvinceServiceImpl{
		ProvinceService: baseService,
		repo:            repo,
		db:              db,
		cache:           cache,
	}
}

// Create creates a new province
func (s *adminProvinceServiceImpl) Create(ctx context.Context, req master.CreateProvinceRequest) (*master.ProvinceResponse, error) {
	// Check duplicate code
	existing, err := s.repo.GetByCode(ctx, req.Code)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check duplicate province code: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("province with code '%s' already exists", req.Code)
	}

	// Create province entity
	province := &master.Province{
		Name:     req.Name,
		Code:     req.Code,
		IsActive: req.IsActive,
	}

	if err := s.repo.Create(ctx, province); err != nil {
		return nil, fmt.Errorf("failed to create province: %w", err)
	}

	// Invalidate cache
	s.invalidateCache()

	// Map to response
	response := &master.ProvinceResponse{
		ID:       province.ID,
		Name:     province.Name,
		Code:     province.Code,
		IsActive: province.IsActive,
	}

	return response, nil
}

// Update updates an existing province
func (s *adminProvinceServiceImpl) Update(ctx context.Context, id int64, req master.UpdateProvinceRequest) (*master.ProvinceResponse, error) {
	// Get existing province
	province, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("province not found")
		}
		return nil, fmt.Errorf("failed to get province: %w", err)
	}

	// Update fields
	if req.Name != "" {
		province.Name = req.Name
	}

	if req.Code != "" {
		// Check duplicate code if code is being changed
		if req.Code != province.Code {
			existing, err := s.repo.GetByCode(ctx, req.Code)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("failed to check duplicate province code: %w", err)
			}
			if existing != nil && existing.ID != id {
				return nil, fmt.Errorf("province with code '%s' already exists", req.Code)
			}
		}
		province.Code = req.Code
	}

	if req.IsActive != nil {
		province.IsActive = *req.IsActive
	}

	// Update province
	if err := s.repo.Update(ctx, province); err != nil {
		return nil, fmt.Errorf("failed to update province: %w", err)
	}

	// Invalidate cache
	s.invalidateCache()

	// Map to response
	response := &master.ProvinceResponse{
		ID:       province.ID,
		Name:     province.Name,
		Code:     province.Code,
		IsActive: province.IsActive,
	}

	return response, nil
}

// Delete deletes a province if not referenced by cities or companies
func (s *adminProvinceServiceImpl) Delete(ctx context.Context, id int64) error {
	// Check if province exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("province not found")
		}
		return fmt.Errorf("failed to get province: %w", err)
	}

	// Check references
	cities, companies, err := s.CountReferences(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check province references: %w", err)
	}
	if cities > 0 || companies > 0 {
		return fmt.Errorf("cannot delete province: it is still referenced by %d cities and %d companies", cities, companies)
	}

	// Delete province
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete province: %w", err)
	}

	// Invalidate cache
	s.invalidateCache()

	return nil
}

// CheckDuplicateCode checks if a province with the given code exists
func (s *adminProvinceServiceImpl) CheckDuplicateCode(ctx context.Context, code string) (bool, error) {
	province, err := s.repo.GetByCode(ctx, code)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	return province != nil, nil
}

// CountReferences counts how many cities and companies reference this province
func (s *adminProvinceServiceImpl) CountReferences(ctx context.Context, id int64) (cities int64, companies int64, err error) {
	// Count cities
	err = s.db.WithContext(ctx).Table("cities").
		Where("province_id = ?", id).
		Count(&cities).Error
	if err != nil {
		return 0, 0, err
	}

	// Count companies
	err = s.db.WithContext(ctx).Table("companies").
		Where("province_id = ?", id).
		Count(&companies).Error
	if err != nil {
		return 0, 0, err
	}

	return cities, companies, nil
}

// invalidateCache invalidates all province-related cache entries
func (s *adminProvinceServiceImpl) invalidateCache() {
	cacheKeys := []string{
		"provinces:all",
		"provinces:active",
	}
	for _, key := range cacheKeys {
		s.cache.Delete(key)
	}
}
