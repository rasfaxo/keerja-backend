package master

import (
	"context"
	"strings"

	"keerja-backend/internal/domain/master"

	"gorm.io/gorm"
)

// provinceRepositoryImpl implements master.ProvinceRepository
type provinceRepositoryImpl struct {
	db *gorm.DB
}

// NewProvinceRepository creates a new province repository instance
func NewProvinceRepository(db *gorm.DB) master.ProvinceRepository {
	return &provinceRepositoryImpl{db: db}
}

// GetAll retrieves all provinces
func (r *provinceRepositoryImpl) GetAll(ctx context.Context) ([]master.Province, error) {
	var provinces []master.Province
	err := r.db.WithContext(ctx).
		Order("name ASC").
		Find(&provinces).Error
	return provinces, err
}

// GetActive retrieves all active provinces
func (r *provinceRepositoryImpl) GetActive(ctx context.Context) ([]master.Province, error) {
	var provinces []master.Province
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("name ASC").
		Find(&provinces).Error
	return provinces, err
}

// GetByID retrieves a province by ID
func (r *provinceRepositoryImpl) GetByID(ctx context.Context, id int64) (*master.Province, error) {
	var province master.Province
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&province).Error
	if err != nil {
		return nil, err
	}
	return &province, nil
}

// GetByCode retrieves a province by code
func (r *provinceRepositoryImpl) GetByCode(ctx context.Context, code string) (*master.Province, error) {
	var province master.Province
	err := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&province).Error
	if err != nil {
		return nil, err
	}
	return &province, nil
}

// Search searches provinces by name (case-insensitive)
func (r *provinceRepositoryImpl) Search(ctx context.Context, query string) ([]master.Province, error) {
	var provinces []master.Province
	searchPattern := "%" + strings.ToLower(query) + "%"
	err := r.db.WithContext(ctx).
		Where("LOWER(name) LIKE ?", searchPattern).
		Order("name ASC").
		Find(&provinces).Error
	return provinces, err
}

// Create creates a new province
func (r *provinceRepositoryImpl) Create(ctx context.Context, province *master.Province) error {
	return r.db.WithContext(ctx).Create(province).Error
}

// Update updates an existing province
func (r *provinceRepositoryImpl) Update(ctx context.Context, province *master.Province) error {
	return r.db.WithContext(ctx).Save(province).Error
}

// Delete deletes a province by ID
func (r *provinceRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Delete(&master.Province{}, id).Error
}

// ExistsByID checks if a province exists by ID
func (r *provinceRepositoryImpl) ExistsByID(ctx context.Context, id int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&master.Province{}).
		Where("id = ?", id).
		Count(&count).Error
	return count > 0, err
}

// GetByName retrieves a province by exact name (case-insensitive)
func (r *provinceRepositoryImpl) GetByName(ctx context.Context, name string) (*master.Province, error) {
	var province master.Province
	err := r.db.WithContext(ctx).
		Where("LOWER(name) = ?", strings.ToLower(name)).
		First(&province).Error
	if err != nil {
		return nil, err
	}
	return &province, nil
}
