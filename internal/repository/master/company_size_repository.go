package master

import (
	"context"
	"strings"

	"keerja-backend/internal/domain/master"

	"gorm.io/gorm"
)

// companySizeRepositoryImpl implements master.CompanySizeRepository
type companySizeRepositoryImpl struct {
	db *gorm.DB
}

// NewCompanySizeRepository creates a new company size repository instance
func NewCompanySizeRepository(db *gorm.DB) master.CompanySizeRepository {
	return &companySizeRepositoryImpl{db: db}
}

// GetAll retrieves all company sizes
func (r *companySizeRepositoryImpl) GetAll(ctx context.Context) ([]master.CompanySize, error) {
	var sizes []master.CompanySize
	err := r.db.WithContext(ctx).
		Order("display_order ASC, min_employees ASC").
		Find(&sizes).Error
	return sizes, err
}

// GetActive retrieves all active company sizes
func (r *companySizeRepositoryImpl) GetActive(ctx context.Context) ([]master.CompanySize, error) {
	var sizes []master.CompanySize
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("display_order ASC, min_employees ASC").
		Find(&sizes).Error
	return sizes, err
}

// GetByID retrieves a company size by ID
func (r *companySizeRepositoryImpl) GetByID(ctx context.Context, id int64) (*master.CompanySize, error) {
	var size master.CompanySize
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&size).Error
	if err != nil {
		return nil, err
	}
	return &size, nil
}

// Create creates a new company size
func (r *companySizeRepositoryImpl) Create(ctx context.Context, size *master.CompanySize) error {
	return r.db.WithContext(ctx).Create(size).Error
}

// Update updates an existing company size
func (r *companySizeRepositoryImpl) Update(ctx context.Context, size *master.CompanySize) error {
	return r.db.WithContext(ctx).Save(size).Error
}

// Delete deletes a company size by ID
func (r *companySizeRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Delete(&master.CompanySize{}, id).Error
}

// ExistsByID checks if a company size exists by ID
func (r *companySizeRepositoryImpl) ExistsByID(ctx context.Context, id int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&master.CompanySize{}).
		Where("id = ?", id).
		Count(&count).Error
	return count > 0, err
}

// GetByCategory retrieves a company size by category (exact match, case-insensitive)
func (r *companySizeRepositoryImpl) GetByCategory(ctx context.Context, category string) (*master.CompanySize, error) {
	var size master.CompanySize
	err := r.db.WithContext(ctx).
		Where("LOWER(category) = ?", strings.ToLower(category)).
		First(&size).Error
	if err != nil {
		return nil, err
	}
	return &size, nil
}
