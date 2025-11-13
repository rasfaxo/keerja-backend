package postgres

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"keerja-backend/internal/domain/master"
)

// benefitsMasterRepository implements master.BenefitsMasterRepository interface
type benefitsMasterRepository struct {
	db *gorm.DB
}

// NewBenefitsMasterRepository creates a new instance of benefits master repository
func NewBenefitsMasterRepository(db *gorm.DB) master.BenefitsMasterRepository {
	return &benefitsMasterRepository{db: db}
}

// Create creates a new benefit master record
func (r *benefitsMasterRepository) Create(ctx context.Context, benefit *master.BenefitsMaster) error {
	return r.db.WithContext(ctx).Create(benefit).Error
}

// FindByID retrieves a benefit master by ID
func (r *benefitsMasterRepository) FindByID(ctx context.Context, id int64) (*master.BenefitsMaster, error) {
	var benefit master.BenefitsMaster
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&benefit).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &benefit, nil
}

// FindByCode retrieves a benefit master by code
func (r *benefitsMasterRepository) FindByCode(ctx context.Context, code string) (*master.BenefitsMaster, error) {
	var benefit master.BenefitsMaster
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&benefit).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &benefit, nil
}

// FindByName retrieves a benefit master by name
func (r *benefitsMasterRepository) FindByName(ctx context.Context, name string) (*master.BenefitsMaster, error) {
	var benefit master.BenefitsMaster
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&benefit).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &benefit, nil
}

// Update updates an existing benefit master record
func (r *benefitsMasterRepository) Update(ctx context.Context, benefit *master.BenefitsMaster) error {
	return r.db.WithContext(ctx).Model(benefit).Updates(benefit).Error
}

// Delete soft deletes a benefit master record
func (r *benefitsMasterRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&master.BenefitsMaster{}).Error
}

// List retrieves benefits with filtering and pagination
func (r *benefitsMasterRepository) List(ctx context.Context, filter *master.BenefitsFilter) ([]master.BenefitsMaster, int64, error) {
	var benefits []master.BenefitsMaster
	query := r.db.WithContext(ctx).Model(&master.BenefitsMaster{})

	// Apply filters
	if filter != nil {
		if filter.Category != "" {
			query = query.Where("category = ?", filter.Category)
		}
		if filter.IsActive != nil {
			query = query.Where("is_active = ?", *filter.IsActive)
		}
		if filter.MinPopularity != nil && *filter.MinPopularity > 0.0 {
			query = query.Where("popularity_score >= ?", *filter.MinPopularity)
		}
		if filter.MaxPopularity != nil && *filter.MaxPopularity > 0.0 {
			query = query.Where("popularity_score <= ?", *filter.MaxPopularity)
		}
		if filter.Search != "" {
			searchPattern := "%" + strings.ToLower(filter.Search) + "%"
			query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchPattern, searchPattern)
		}
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and sorting
	if filter != nil {
		// Sorting
		sortField := "name"
		sortOrder := "ASC"
		if filter.SortBy != "" {
			sortField = filter.SortBy
		}
		if filter.SortOrder != "" {
			sortOrder = filter.SortOrder
		}
		query = query.Order(sortField + " " + sortOrder)

		// Pagination
		if filter.PageSize > 0 {
			offset := (filter.Page - 1) * filter.PageSize
			query = query.Limit(filter.PageSize).Offset(offset)
		}
	} else {
		query = query.Order("name ASC")
	}

	err := query.Find(&benefits).Error
	return benefits, total, err
}

// ListActive retrieves all active benefits
func (r *benefitsMasterRepository) ListActive(ctx context.Context) ([]master.BenefitsMaster, error) {
	var benefits []master.BenefitsMaster
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("popularity_score DESC, name ASC").
		Find(&benefits).Error
	return benefits, err
}

// ListByCategory retrieves benefits by category
func (r *benefitsMasterRepository) ListByCategory(ctx context.Context, category string) ([]master.BenefitsMaster, error) {
	var benefits []master.BenefitsMaster
	err := r.db.WithContext(ctx).
		Where("category = ? AND is_active = ?", category, true).
		Order("popularity_score DESC, name ASC").
		Find(&benefits).Error
	return benefits, err
}

// SearchBenefits searches benefits by query string with pagination
func (r *benefitsMasterRepository) SearchBenefits(ctx context.Context, query string, page, pageSize int) ([]master.BenefitsMaster, int64, error) {
	var benefits []master.BenefitsMaster
	searchQuery := "%" + strings.ToLower(query) + "%"

	dbQuery := r.db.WithContext(ctx).Model(&master.BenefitsMaster{}).
		Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchQuery, searchQuery)

	// Count total
	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := dbQuery.
		Order("popularity_score DESC, name ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&benefits).Error

	return benefits, total, err
}

// GetCategories retrieves all unique categories
func (r *benefitsMasterRepository) GetCategories(ctx context.Context) ([]string, error) {
	var categories []string
	err := r.db.WithContext(ctx).
		Model(&master.BenefitsMaster{}).
		Distinct("category").
		Order("category ASC").
		Pluck("category", &categories).
		Error
	return categories, err
}

// CountByCategory returns count of benefits by category
func (r *benefitsMasterRepository) CountByCategory(ctx context.Context, category string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&master.BenefitsMaster{}).
		Where("category = ?", category).
		Count(&count).Error
	return count, err
}

// UpdatePopularity updates the popularity score
func (r *benefitsMasterRepository) UpdatePopularity(ctx context.Context, id int64, score float64) error {
	return r.db.WithContext(ctx).
		Model(&master.BenefitsMaster{}).
		Where("id = ?", id).
		Update("popularity_score", score).
		Error
}

// IncrementPopularity increases the popularity score
func (r *benefitsMasterRepository) IncrementPopularity(ctx context.Context, id int64, amount float64) error {
	return r.db.WithContext(ctx).
		Model(&master.BenefitsMaster{}).
		Where("id = ?", id).
		UpdateColumn("popularity_score", gorm.Expr("popularity_score + ?", amount)).
		Error
}

// GetMostPopular retrieves most popular benefits
func (r *benefitsMasterRepository) GetMostPopular(ctx context.Context, limit int) ([]master.BenefitsMaster, error) {
	var benefits []master.BenefitsMaster
	query := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("popularity_score DESC, name ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&benefits).Error
	return benefits, err
}

// GetByPopularityRange retrieves benefits within popularity score range
func (r *benefitsMasterRepository) GetByPopularityRange(ctx context.Context, minScore, maxScore float64) ([]master.BenefitsMaster, error) {
	var benefits []master.BenefitsMaster
	err := r.db.WithContext(ctx).
		Where("popularity_score >= ? AND popularity_score <= ?", minScore, maxScore).
		Where("is_active = ?", true).
		Order("popularity_score DESC, name ASC").
		Find(&benefits).Error
	return benefits, err
}

// Activate sets the benefit as active
func (r *benefitsMasterRepository) Activate(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&master.BenefitsMaster{}).
		Where("id = ?", id).
		Update("is_active", true).
		Error
}

// Deactivate sets the benefit as inactive
func (r *benefitsMasterRepository) Deactivate(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&master.BenefitsMaster{}).
		Where("id = ?", id).
		Update("is_active", false).
		Error
}

// BulkCreate creates multiple benefits at once
func (r *benefitsMasterRepository) BulkCreate(ctx context.Context, benefits []master.BenefitsMaster) error {
	return r.db.WithContext(ctx).Create(&benefits).Error
}

// BulkUpdatePopularity updates popularity scores for multiple benefits
func (r *benefitsMasterRepository) BulkUpdatePopularity(ctx context.Context, updates map[int64]float64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for id, score := range updates {
			if err := tx.Model(&master.BenefitsMaster{}).
				Where("id = ?", id).
				Update("popularity_score", score).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// Count returns total number of benefit records
func (r *benefitsMasterRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&master.BenefitsMaster{}).Count(&count).Error
	return count, err
}

// CountActive returns total number of active benefit records
func (r *benefitsMasterRepository) CountActive(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&master.BenefitsMaster{}).
		Where("is_active = ?", true).
		Count(&count).Error
	return count, err
}

// GetBenefitStats returns comprehensive benefit statistics
func (r *benefitsMasterRepository) GetBenefitStats(ctx context.Context) (*master.BenefitStats, error) {
	stats := &master.BenefitStats{
		ByCategory: make(map[string]int64),
	}

	// Total count
	if err := r.db.WithContext(ctx).Model(&master.BenefitsMaster{}).Count(&stats.TotalBenefits).Error; err != nil {
		return nil, err
	}

	// Active count
	if err := r.db.WithContext(ctx).
		Model(&master.BenefitsMaster{}).
		Where("is_active = ?", true).
		Count(&stats.ActiveBenefits).Error; err != nil {
		return nil, err
	}

	stats.InactiveBenefits = stats.TotalBenefits - stats.ActiveBenefits

	// Count by category
	type CategoryCount struct {
		Category string
		Count    int64
	}
	var categoryCounts []CategoryCount
	if err := r.db.WithContext(ctx).
		Model(&master.BenefitsMaster{}).
		Select("category, COUNT(*) as count").
		Group("category").
		Order("count DESC").
		Scan(&categoryCounts).Error; err != nil {
		return nil, err
	}

	// Build category map and top categories
	for _, cc := range categoryCounts {
		stats.ByCategory[cc.Category] = cc.Count
		percentage := float64(cc.Count) / float64(stats.TotalBenefits) * 100
		stats.TopCategories = append(stats.TopCategories, master.CategoryStat{
			Category:   cc.Category,
			Count:      cc.Count,
			Percentage: percentage,
		})
	}

	// Average popularity score
	if err := r.db.WithContext(ctx).
		Model(&master.BenefitsMaster{}).
		Select("COALESCE(AVG(popularity_score), 0)").
		Scan(&stats.AveragePopularity).Error; err != nil {
		return nil, err
	}

	// Most popular benefit
	var mostPopular master.BenefitsMaster
	if err := r.db.WithContext(ctx).
		Order("popularity_score DESC, name ASC").
		First(&mostPopular).Error; err == nil {
		stats.MostPopular = &mostPopular
	}

	return stats, nil
}
