package master

import (
	"context"
	"strings"

	"keerja-backend/internal/domain/master"

	"gorm.io/gorm"
)

// cityRepositoryImpl implements master.CityRepository
type cityRepositoryImpl struct {
	db *gorm.DB
}

// NewCityRepository creates a new city repository instance
func NewCityRepository(db *gorm.DB) master.CityRepository {
	return &cityRepositoryImpl{db: db}
}

// GetAll retrieves all cities
func (r *cityRepositoryImpl) GetAll(ctx context.Context) ([]master.City, error) {
	var cities []master.City
	err := r.db.WithContext(ctx).
		Order("name ASC").
		Find(&cities).Error
	return cities, err
}

// GetActive retrieves all active cities
func (r *cityRepositoryImpl) GetActive(ctx context.Context) ([]master.City, error) {
	var cities []master.City
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("name ASC").
		Find(&cities).Error
	return cities, err
}

// GetByID retrieves a city by ID
func (r *cityRepositoryImpl) GetByID(ctx context.Context, id int64) (*master.City, error) {
	var city master.City
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&city).Error
	if err != nil {
		return nil, err
	}
	return &city, nil
}

// GetByCode retrieves a city by code
func (r *cityRepositoryImpl) GetByCode(ctx context.Context, code string) (*master.City, error) {
	var city master.City
	err := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&city).Error
	if err != nil {
		return nil, err
	}
	return &city, nil
}

// GetByProvinceID retrieves all cities in a province
func (r *cityRepositoryImpl) GetByProvinceID(ctx context.Context, provinceID int64) ([]master.City, error) {
	var cities []master.City
	err := r.db.WithContext(ctx).
		Where("province_id = ?", provinceID).
		Order("type DESC, name ASC"). // Kota first, then Kabupaten
		Find(&cities).Error
	return cities, err
}

// GetActiveByProvinceID retrieves all active cities in a province
func (r *cityRepositoryImpl) GetActiveByProvinceID(ctx context.Context, provinceID int64) ([]master.City, error) {
	var cities []master.City
	err := r.db.WithContext(ctx).
		Where("province_id = ? AND is_active = ?", provinceID, true).
		Order("type DESC, name ASC"). // Kota first, then Kabupaten
		Find(&cities).Error
	return cities, err
}

// Search searches cities by name with optional province filter (case-insensitive)
func (r *cityRepositoryImpl) Search(ctx context.Context, query string, provinceID *int64) ([]master.City, error) {
	var cities []master.City
	searchPattern := "%" + strings.ToLower(query) + "%"

	db := r.db.WithContext(ctx).
		Where("LOWER(name) LIKE ?", searchPattern)

	// Apply province filter if provided
	if provinceID != nil && *provinceID > 0 {
		db = db.Where("province_id = ?", *provinceID)
	}

	err := db.Order("type DESC, name ASC").Find(&cities).Error
	return cities, err
}

// GetWithProvince retrieves a city with its province preloaded
func (r *cityRepositoryImpl) GetWithProvince(ctx context.Context, id int64) (*master.City, error) {
	var city master.City
	err := r.db.WithContext(ctx).
		Preload("Province").
		Where("id = ?", id).
		First(&city).Error
	if err != nil {
		return nil, err
	}
	return &city, nil
}

// Create creates a new city
func (r *cityRepositoryImpl) Create(ctx context.Context, city *master.City) error {
	return r.db.WithContext(ctx).Create(city).Error
}

// Update updates an existing city
func (r *cityRepositoryImpl) Update(ctx context.Context, city *master.City) error {
	return r.db.WithContext(ctx).Save(city).Error
}

// Delete deletes a city by ID
func (r *cityRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Delete(&master.City{}, id).Error
}

// ExistsByID checks if a city exists by ID
func (r *cityRepositoryImpl) ExistsByID(ctx context.Context, id int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&master.City{}).
		Where("id = ?", id).
		Count(&count).Error
	return count > 0, err
}

// GetByNameAndProvinceID retrieves a city by exact name and province ID (case-insensitive)
func (r *cityRepositoryImpl) GetByNameAndProvinceID(ctx context.Context, name string, provinceID int64) (*master.City, error) {
	var city master.City
	err := r.db.WithContext(ctx).
		Where("LOWER(name) = ? AND province_id = ?", strings.ToLower(name), provinceID).
		First(&city).Error
	if err != nil {
		return nil, err
	}
	return &city, nil
}
