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

// adminCompanySizeServiceImpl implements AdminCompanySizeService
type adminCompanySizeServiceImpl struct {
	master.CompanySizeService // Embed base service for read operations
	repo                      master.CompanySizeRepository
	db                        *gorm.DB // For counting references
	cache                     cache.Cache
}

// NewAdminCompanySizeService creates a new AdminCompanySizeService
func NewAdminCompanySizeService(
	baseService master.CompanySizeService,
	repo master.CompanySizeRepository,
	db *gorm.DB,
	cache cache.Cache,
) master.AdminCompanySizeService {
	return &adminCompanySizeServiceImpl{
		CompanySizeService: baseService,
		repo:               repo,
		db:                 db,
		cache:              cache,
	}
}

// Create creates a new company size category
func (s *adminCompanySizeServiceImpl) Create(ctx context.Context, req master.CreateCompanySizeRequest) (*master.CompanySizeResponse, error) {
	// Check duplicate label
	existing, err := s.repo.GetByCategory(ctx, req.Label)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check duplicate company size label: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("company size with label '%s' already exists", req.Label)
	}

	// Create company size entity
	companySize := &master.CompanySize{
		Label:        req.Label,
		MinEmployees: req.MinEmployees,
		MaxEmployees: s.mapMaxEmployees(req.MaxEmployees),
		IsActive:     req.IsActive,
	}

	if err := s.repo.Create(ctx, companySize); err != nil {
		return nil, fmt.Errorf("failed to create company size: %w", err)
	}

	// Invalidate cache
	s.invalidateCache()

	// Map to response
	response := &master.CompanySizeResponse{
		ID:           companySize.ID,
		Label:        companySize.Label,
		MinEmployees: companySize.MinEmployees,
		MaxEmployees: s.getMaxEmployeesPointer(companySize),
		IsActive:     companySize.IsActive,
	}

	return response, nil
}

// Update updates an existing company size category
func (s *adminCompanySizeServiceImpl) Update(ctx context.Context, id int64, req master.UpdateCompanySizeRequest) (*master.CompanySizeResponse, error) {
	// Get existing company size
	companySize, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("company size not found")
		}
		return nil, fmt.Errorf("failed to get company size: %w", err)
	}

	// Update fields
	if req.Label != "" {
		// Check duplicate label if label is being changed
		if req.Label != companySize.Label {
			existing, err := s.repo.GetByCategory(ctx, req.Label)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("failed to check duplicate company size label: %w", err)
			}
			if existing != nil && existing.ID != id {
				return nil, fmt.Errorf("company size with label '%s' already exists", req.Label)
			}
		}
		companySize.Label = req.Label
	}

	if req.MinEmployees != nil {
		companySize.MinEmployees = *req.MinEmployees
	}

	if req.MaxEmployees != nil {
		companySize.MaxEmployees = s.mapMaxEmployees(req.MaxEmployees)
	}

	if req.IsActive != nil {
		companySize.IsActive = *req.IsActive
	}

	// Update company size
	if err := s.repo.Update(ctx, companySize); err != nil {
		return nil, fmt.Errorf("failed to update company size: %w", err)
	}

	// Invalidate cache
	s.invalidateCache()

	// Map to response
	response := &master.CompanySizeResponse{
		ID:           companySize.ID,
		Label:        companySize.Label,
		MinEmployees: companySize.MinEmployees,
		MaxEmployees: s.getMaxEmployeesPointer(companySize),
		IsActive:     companySize.IsActive,
	}

	return response, nil
}

// Delete deletes a company size category if not referenced by companies
func (s *adminCompanySizeServiceImpl) Delete(ctx context.Context, id int64) error {
	// Check if company size exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("company size not found")
		}
		return fmt.Errorf("failed to get company size: %w", err)
	}

	// Check references
	companies, err := s.CountCompanyReferences(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check company size references: %w", err)
	}
	if companies > 0 {
		return fmt.Errorf("cannot delete company size: it is still referenced by %d companies", companies)
	}

	// Delete company size
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete company size: %w", err)
	}

	// Invalidate cache
	s.invalidateCache()

	return nil
}

// CheckDuplicateCategory checks if a size category with the given label exists
func (s *adminCompanySizeServiceImpl) CheckDuplicateCategory(ctx context.Context, label string) (bool, error) {
	companySize, err := s.repo.GetByCategory(ctx, label)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	return companySize != nil, nil
}

// CountCompanyReferences counts how many companies reference this size category
func (s *adminCompanySizeServiceImpl) CountCompanyReferences(ctx context.Context, id int64) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Table("companies").
		Where("company_size_id = ?", id).
		Count(&count).Error
	return count, err
}

// invalidateCache invalidates all company size-related cache entries
func (s *adminCompanySizeServiceImpl) invalidateCache() {
	cacheKeys := []string{
		"company_sizes:all",
		"company_sizes:active",
	}
	for _, key := range cacheKeys {
		s.cache.Delete(key)
	}
}

// mapMaxEmployees converts *int to sql.NullInt32
func (s *adminCompanySizeServiceImpl) mapMaxEmployees(maxEmp *int) sql.NullInt32 {
	if maxEmp == nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: int32(*maxEmp), Valid: true}
}

// getMaxEmployeesPointer returns a pointer to max employees value, or nil if unlimited
func (s *adminCompanySizeServiceImpl) getMaxEmployeesPointer(size *master.CompanySize) *int {
	maxEmp := size.GetMaxEmployees()
	if maxEmp == -1 || !size.MaxEmployees.Valid {
		return nil // Unlimited
	}
	result := int(size.MaxEmployees.Int32)
	return &result
}

