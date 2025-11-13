package master

import (
	"context"
	"strings"

	"keerja-backend/internal/domain/master"

	"gorm.io/gorm"
)

// industryRepositoryImpl implements master.IndustryRepository
type industryRepositoryImpl struct {
	db *gorm.DB
}

// NewIndustryRepository creates a new industry repository instance
func NewIndustryRepository(db *gorm.DB) master.IndustryRepository {
	return &industryRepositoryImpl{db: db}
}

// GetAll retrieves all industries (including soft deleted if needed)
func (r *industryRepositoryImpl) GetAll(ctx context.Context) ([]master.Industry, error) {
	var industries []master.Industry
	err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Order("display_order ASC, name ASC").
		Find(&industries).Error
	return industries, err
}

// GetActive retrieves all active industries
func (r *industryRepositoryImpl) GetActive(ctx context.Context) ([]master.Industry, error) {
	var industries []master.Industry
	err := r.db.WithContext(ctx).
		Where("is_active = ? AND deleted_at IS NULL", true).
		Order("display_order ASC, name ASC").
		Find(&industries).Error
	return industries, err
}

// GetByID retrieves an industry by ID
func (r *industryRepositoryImpl) GetByID(ctx context.Context, id int64) (*master.Industry, error) {
	var industry master.Industry
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&industry).Error
	if err != nil {
		return nil, err
	}
	return &industry, nil
}

// GetBySlug retrieves an industry by slug
func (r *industryRepositoryImpl) GetBySlug(ctx context.Context, slug string) (*master.Industry, error) {
	var industry master.Industry
	err := r.db.WithContext(ctx).
		Where("slug = ? AND deleted_at IS NULL", slug).
		First(&industry).Error
	if err != nil {
		return nil, err
	}
	return &industry, nil
}

// Search searches industries by name (case-insensitive)
func (r *industryRepositoryImpl) Search(ctx context.Context, query string) ([]master.Industry, error) {
	var industries []master.Industry
	searchPattern := "%" + strings.ToLower(query) + "%"
	err := r.db.WithContext(ctx).
		Where("LOWER(name) LIKE ? AND deleted_at IS NULL", searchPattern).
		Order("display_order ASC, name ASC").
		Find(&industries).Error
	return industries, err
}

// Create creates a new industry
func (r *industryRepositoryImpl) Create(ctx context.Context, industry *master.Industry) error {
	return r.db.WithContext(ctx).Create(industry).Error
}

// Update updates an existing industry
func (r *industryRepositoryImpl) Update(ctx context.Context, industry *master.Industry) error {
	return r.db.WithContext(ctx).Save(industry).Error
}

// Delete soft deletes an industry by setting deleted_at timestamp
func (r *industryRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&master.Industry{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// ExistsByID checks if an industry exists by ID (excluding soft deleted)
func (r *industryRepositoryImpl) ExistsByID(ctx context.Context, id int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&master.Industry{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Count(&count).Error
	return count > 0, err
}

// GetByName retrieves an industry by exact name (case-insensitive)
func (r *industryRepositoryImpl) GetByName(ctx context.Context, name string) (*master.Industry, error) {
	var industry master.Industry
	err := r.db.WithContext(ctx).
		Where("LOWER(name) = ? AND deleted_at IS NULL", strings.ToLower(name)).
		First(&industry).Error
	if err != nil {
		return nil, err
	}
	return &industry, nil
}
