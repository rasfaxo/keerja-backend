package master

import (
	"context"
	"strings"

	"keerja-backend/internal/domain/master"

	"gorm.io/gorm"
)

// districtRepositoryImpl implements master.DistrictRepository
type districtRepositoryImpl struct {
	db *gorm.DB
}

// NewDistrictRepository creates a new district repository instance
func NewDistrictRepository(db *gorm.DB) master.DistrictRepository {
	return &districtRepositoryImpl{db: db}
}

// GetAll retrieves all districts
func (r *districtRepositoryImpl) GetAll(ctx context.Context) ([]master.District, error) {
	var districts []master.District
	err := r.db.WithContext(ctx).
		Order("name ASC").
		Find(&districts).Error
	return districts, err
}

// GetActive retrieves all active districts
func (r *districtRepositoryImpl) GetActive(ctx context.Context) ([]master.District, error) {
	var districts []master.District
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("name ASC").
		Find(&districts).Error
	return districts, err
}

// GetByID retrieves a district by ID
func (r *districtRepositoryImpl) GetByID(ctx context.Context, id int64) (*master.District, error) {
	var district master.District
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&district).Error
	if err != nil {
		return nil, err
	}
	return &district, nil
}

// GetByCode retrieves a district by code
func (r *districtRepositoryImpl) GetByCode(ctx context.Context, code string) (*master.District, error) {
	var district master.District
	err := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&district).Error
	if err != nil {
		return nil, err
	}
	return &district, nil
}

// GetByCityID retrieves all districts in a city
func (r *districtRepositoryImpl) GetByCityID(ctx context.Context, cityID int64) ([]master.District, error) {
	var districts []master.District
	err := r.db.WithContext(ctx).
		Where("city_id = ?", cityID).
		Order("name ASC").
		Find(&districts).Error
	return districts, err
}

// GetActiveByCityID retrieves all active districts in a city
func (r *districtRepositoryImpl) GetActiveByCityID(ctx context.Context, cityID int64) ([]master.District, error) {
	var districts []master.District
	err := r.db.WithContext(ctx).
		Where("city_id = ? AND is_active = ?", cityID, true).
		Order("name ASC").
		Find(&districts).Error
	return districts, err
}

// Search searches districts by name with optional city filter (case-insensitive)
func (r *districtRepositoryImpl) Search(ctx context.Context, query string, cityID *int64) ([]master.District, error) {
	var districts []master.District
	searchPattern := "%" + strings.ToLower(query) + "%"

	db := r.db.WithContext(ctx).
		Where("LOWER(name) LIKE ?", searchPattern)

	// Apply city filter if provided
	if cityID != nil && *cityID > 0 {
		db = db.Where("city_id = ?", *cityID)
	}

	err := db.Order("name ASC").Find(&districts).Error
	return districts, err
}

// GetWithFullLocation retrieves a district with city and province preloaded (3-level hierarchy)
func (r *districtRepositoryImpl) GetWithFullLocation(ctx context.Context, id int64) (*master.District, error) {
	var district master.District
	err := r.db.WithContext(ctx).
		Preload("City").          // Load city
		Preload("City.Province"). // Load province through city
		Where("id = ?", id).
		First(&district).Error
	if err != nil {
		return nil, err
	}
	return &district, nil
}

// GetByPostalCode retrieves districts by postal code
func (r *districtRepositoryImpl) GetByPostalCode(ctx context.Context, postalCode string) ([]master.District, error) {
	var districts []master.District
	err := r.db.WithContext(ctx).
		Where("postal_code = ?", postalCode).
		Order("name ASC").
		Find(&districts).Error
	return districts, err
}

// Create creates a new district
func (r *districtRepositoryImpl) Create(ctx context.Context, district *master.District) error {
	return r.db.WithContext(ctx).Create(district).Error
}

// Update updates an existing district
func (r *districtRepositoryImpl) Update(ctx context.Context, district *master.District) error {
	return r.db.WithContext(ctx).Save(district).Error
}

// Delete deletes a district by ID
func (r *districtRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Delete(&master.District{}, id).Error
}

// ExistsByID checks if a district exists by ID
func (r *districtRepositoryImpl) ExistsByID(ctx context.Context, id int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&master.District{}).
		Where("id = ?", id).
		Count(&count).Error
	return count > 0, err
}

// GetByNameAndCityID retrieves a district by exact name and city ID (case-insensitive)
func (r *districtRepositoryImpl) GetByNameAndCityID(ctx context.Context, name string, cityID int64) (*master.District, error) {
	var district master.District
	err := r.db.WithContext(ctx).
		Where("LOWER(name) = ? AND city_id = ?", strings.ToLower(name), cityID).
		First(&district).Error
	if err != nil {
		return nil, err
	}
	return &district, nil
}
