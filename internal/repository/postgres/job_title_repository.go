package postgres

import (
	"context"
	"strings"

	"keerja-backend/internal/domain/master"

	"gorm.io/gorm"
)

type jobTitleRepository struct {
	db *gorm.DB
}

// NewJobTitleRepository creates a new JobTitle repository
func NewJobTitleRepository(db *gorm.DB) master.JobTitleRepository {
	return &jobTitleRepository{db: db}
}

// Create creates a new job title
func (r *jobTitleRepository) Create(ctx context.Context, jobTitle *master.JobTitle) error {
	// Normalize name for search
	jobTitle.NormalizedName = strings.ToLower(strings.TrimSpace(jobTitle.Name))
	return r.db.WithContext(ctx).Create(jobTitle).Error
}

// FindByID finds job title by ID
func (r *jobTitleRepository) FindByID(ctx context.Context, id int64) (*master.JobTitle, error) {
	var jobTitle master.JobTitle
	err := r.db.WithContext(ctx).First(&jobTitle, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &jobTitle, nil
}

// FindByName finds job title by exact name
func (r *jobTitleRepository) FindByName(ctx context.Context, name string) (*master.JobTitle, error) {
	var jobTitle master.JobTitle
	err := r.db.WithContext(ctx).
		Where("LOWER(name) = ?", strings.ToLower(name)).
		First(&jobTitle).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &jobTitle, nil
}

// Update updates job title
func (r *jobTitleRepository) Update(ctx context.Context, jobTitle *master.JobTitle) error {
	// Update normalized name if name changed
	jobTitle.NormalizedName = strings.ToLower(strings.TrimSpace(jobTitle.Name))
	return r.db.WithContext(ctx).Save(jobTitle).Error
}

// Delete deletes job title
func (r *jobTitleRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&master.JobTitle{}, id).Error
}

// SearchJobTitles performs smart fuzzy search on job titles
func (r *jobTitleRepository) SearchJobTitles(ctx context.Context, query string, limit int) ([]master.JobTitle, error) {
	var jobTitles []master.JobTitle

	if query == "" {
		// If no query, return most popular
		err := r.db.WithContext(ctx).
			Where("is_active = ?", true).
			Order("popularity_score DESC, search_count DESC").
			Limit(limit).
			Find(&jobTitles).Error
		return jobTitles, err
	}

	// Normalize query
	normalizedQuery := strings.ToLower(strings.TrimSpace(query))

	// Smart search with multiple strategies:
	// 1. Exact match (highest priority)
	// 2. Starts with (high priority)
	// 3. Contains (medium priority)
	// 4. Fuzzy match using ILIKE (lower priority)

	// Use Raw SQL for complex CASE ordering
	err := r.db.WithContext(ctx).
		Raw(`
			SELECT * FROM job_titles
			WHERE is_active = true
				AND (
					LOWER(name) = ?
					OR LOWER(name) LIKE ?
					OR LOWER(name) LIKE ?
					OR normalized_name ILIKE ?
				)
			ORDER BY
				CASE 
					WHEN LOWER(name) = ? THEN 1
					WHEN LOWER(name) LIKE ? THEN 2
					WHEN LOWER(name) LIKE ? THEN 3
					ELSE 4
				END,
				popularity_score DESC,
				search_count DESC
			LIMIT ?
		`,
			normalizedQuery, normalizedQuery+"%", "%"+normalizedQuery+"%", "%"+normalizedQuery+"%",
			normalizedQuery, normalizedQuery+"%", "%"+normalizedQuery+"%",
			limit,
		).
		Scan(&jobTitles).Error

	return jobTitles, err
}

// ListActive lists all active job titles
func (r *jobTitleRepository) ListActive(ctx context.Context) ([]master.JobTitle, error) {
	var jobTitles []master.JobTitle
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("name ASC").
		Find(&jobTitles).Error
	return jobTitles, err
}

// ListPopular lists popular job titles
func (r *jobTitleRepository) ListPopular(ctx context.Context, limit int) ([]master.JobTitle, error) {
	var jobTitles []master.JobTitle
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("popularity_score DESC, search_count DESC").
		Limit(limit).
		Find(&jobTitles).Error
	return jobTitles, err
}

// IncrementSearchCount increments search count for analytics
func (r *jobTitleRepository) IncrementSearchCount(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&master.JobTitle{}).
		Where("id = ?", id).
		UpdateColumn("search_count", gorm.Expr("search_count + ?", 1)).Error
}

// UpdatePopularity updates popularity score
func (r *jobTitleRepository) UpdatePopularity(ctx context.Context, id int64, score float64) error {
	// Ensure score is within bounds
	if score < 0 {
		score = 0
	} else if score > 100 {
		score = 100
	}

	return r.db.WithContext(ctx).
		Model(&master.JobTitle{}).
		Where("id = ?", id).
		Update("popularity_score", score).Error
}
