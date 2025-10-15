package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/lib/pq"
	"gorm.io/gorm"

	"keerja-backend/internal/domain/master"
)

// skillsMasterRepository implements master.SkillsMasterRepository interface
type skillsMasterRepository struct {
	db *gorm.DB
}

// NewSkillsMasterRepository creates a new instance of skills master repository
func NewSkillsMasterRepository(db *gorm.DB) master.SkillsMasterRepository {
	return &skillsMasterRepository{db: db}
}

// Create creates a new skill master record
func (r *skillsMasterRepository) Create(ctx context.Context, skill *master.SkillsMaster) error {
	return r.db.WithContext(ctx).Create(skill).Error
}

// FindByID retrieves a skill master by ID
func (r *skillsMasterRepository) FindByID(ctx context.Context, id int64) (*master.SkillsMaster, error) {
	var skill master.SkillsMaster
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Where("id = ?", id).
		First(&skill).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &skill, nil
}

// FindByCode retrieves a skill master by code
func (r *skillsMasterRepository) FindByCode(ctx context.Context, code string) (*master.SkillsMaster, error) {
	var skill master.SkillsMaster
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Where("code = ?", code).
		First(&skill).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &skill, nil
}

// FindByName retrieves a skill master by name
func (r *skillsMasterRepository) FindByName(ctx context.Context, name string) (*master.SkillsMaster, error) {
	var skill master.SkillsMaster
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Where("name = ?", name).
		First(&skill).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &skill, nil
}

// FindByNormalizedName retrieves a skill by normalized name
func (r *skillsMasterRepository) FindByNormalizedName(ctx context.Context, normalizedName string) (*master.SkillsMaster, error) {
	var skill master.SkillsMaster
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Where("normalized_name = ?", normalizedName).
		First(&skill).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &skill, nil
}

// Update updates an existing skill master record
func (r *skillsMasterRepository) Update(ctx context.Context, skill *master.SkillsMaster) error {
	return r.db.WithContext(ctx).Model(skill).Updates(skill).Error
}

// Delete soft deletes a skill master record
func (r *skillsMasterRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&master.SkillsMaster{}).Error
}

// List retrieves skills with filtering and pagination
func (r *skillsMasterRepository) List(ctx context.Context, filter *master.SkillsFilter) ([]master.SkillsMaster, int64, error) {
	var skills []master.SkillsMaster
	query := r.db.WithContext(ctx).Model(&master.SkillsMaster{})

	// Apply filters
	if filter != nil {
		if filter.SkillType != "" {
			query = query.Where("skill_type = ?", filter.SkillType)
		}
		if filter.DifficultyLevel != "" {
			query = query.Where("difficulty_level = ?", filter.DifficultyLevel)
		}
		if filter.CategoryID != nil {
			query = query.Where("category_id = ?", *filter.CategoryID)
		}
		if filter.ParentID != nil {
			query = query.Where("parent_id = ?", *filter.ParentID)
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
		if filter.HasParent != nil {
			if *filter.HasParent {
				query = query.Where("parent_id IS NOT NULL")
			} else {
				query = query.Where("parent_id IS NULL")
			}
		}
		if filter.Search != "" {
			searchPattern := "%" + strings.ToLower(filter.Search) + "%"
			query = query.Where(
				"LOWER(name) LIKE ? OR LOWER(description) LIKE ? OR ? = ANY(aliases)",
				searchPattern, searchPattern, strings.ToLower(filter.Search),
			)
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

	// Load relationships
	query = query.Preload("Parent").Preload("Children")

	err := query.Find(&skills).Error
	return skills, total, err
}

// ListActive retrieves all active skills
func (r *skillsMasterRepository) ListActive(ctx context.Context) ([]master.SkillsMaster, error) {
	var skills []master.SkillsMaster
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Where("is_active = ?", true).
		Order("popularity_score DESC, name ASC").
		Find(&skills).Error
	return skills, err
}

// ListByType retrieves skills by type
func (r *skillsMasterRepository) ListByType(ctx context.Context, skillType string) ([]master.SkillsMaster, error) {
	var skills []master.SkillsMaster
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Where("skill_type = ? AND is_active = ?", skillType, true).
		Order("popularity_score DESC, name ASC").
		Find(&skills).Error
	return skills, err
}

// ListByDifficulty retrieves skills by difficulty level
func (r *skillsMasterRepository) ListByDifficulty(ctx context.Context, difficultyLevel string) ([]master.SkillsMaster, error) {
	var skills []master.SkillsMaster
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Where("difficulty_level = ? AND is_active = ?", difficultyLevel, true).
		Order("popularity_score DESC, name ASC").
		Find(&skills).Error
	return skills, err
}

// ListByCategory retrieves skills by category ID
func (r *skillsMasterRepository) ListByCategory(ctx context.Context, categoryID int64) ([]master.SkillsMaster, error) {
	var skills []master.SkillsMaster
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Where("category_id = ? AND is_active = ?", categoryID, true).
		Order("popularity_score DESC, name ASC").
		Find(&skills).Error
	return skills, err
}

// SearchSkills searches skills by query string with pagination
func (r *skillsMasterRepository) SearchSkills(ctx context.Context, query string, page, pageSize int) ([]master.SkillsMaster, int64, error) {
	var skills []master.SkillsMaster
	searchQuery := "%" + strings.ToLower(query) + "%"

	dbQuery := r.db.WithContext(ctx).Model(&master.SkillsMaster{}).
		Preload("Parent").
		Preload("Children").
		Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ? OR ? = ANY(aliases)",
			searchQuery, searchQuery, strings.ToLower(query))

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
		Find(&skills).Error

	return skills, total, err
}

// SearchByAlias searches skills by alias
func (r *skillsMasterRepository) SearchByAlias(ctx context.Context, alias string) ([]master.SkillsMaster, error) {
	var skills []master.SkillsMaster
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Where("? = ANY(aliases)", strings.ToLower(alias)).
		Order("popularity_score DESC").
		Find(&skills).Error
	return skills, err
}

// FindByAliases finds skills by one or more aliases
func (r *skillsMasterRepository) FindByAliases(ctx context.Context, aliases []string) ([]master.SkillsMaster, error) {
	if len(aliases) == 0 {
		return []master.SkillsMaster{}, nil
	}

	// Convert to lowercase
	lowerAliases := make([]string, len(aliases))
	for i, alias := range aliases {
		lowerAliases[i] = strings.ToLower(alias)
	}

	var skills []master.SkillsMaster
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Where("aliases && ?", pq.Array(lowerAliases)). // Array overlap operator
		Order("popularity_score DESC").
		Find(&skills).Error
	return skills, err
}

// GetRootSkills retrieves all root skills (without parent)
func (r *skillsMasterRepository) GetRootSkills(ctx context.Context) ([]master.SkillsMaster, error) {
	var skills []master.SkillsMaster
	err := r.db.WithContext(ctx).
		Preload("Children").
		Where("parent_id IS NULL").
		Order("popularity_score DESC, name ASC").
		Find(&skills).Error
	return skills, err
}

// GetChildren retrieves immediate children of a skill
func (r *skillsMasterRepository) GetChildren(ctx context.Context, parentID int64) ([]master.SkillsMaster, error) {
	var skills []master.SkillsMaster
	err := r.db.WithContext(ctx).
		Preload("Children").
		Where("parent_id = ?", parentID).
		Order("name ASC").
		Find(&skills).Error
	return skills, err
}

// GetParent retrieves the parent of a skill
func (r *skillsMasterRepository) GetParent(ctx context.Context, childID int64) (*master.SkillsMaster, error) {
	var child master.SkillsMaster
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Where("id = ?", childID).
		First(&child).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return child.Parent, nil
}

// GetSkillTree retrieves the complete skill tree from a root skill
func (r *skillsMasterRepository) GetSkillTree(ctx context.Context, rootID int64) (*master.SkillsMaster, error) {
	var root master.SkillsMaster
	err := r.db.WithContext(ctx).
		Preload("Children.Children.Children.Children"). // 4 levels deep
		Where("id = ?", rootID).
		First(&root).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &root, nil
}

// HasChildren checks if a skill has children
func (r *skillsMasterRepository) HasChildren(ctx context.Context, id int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&master.SkillsMaster{}).
		Where("parent_id = ?", id).
		Count(&count).Error
	return count > 0, err
}

// UpdatePopularity updates the popularity score
func (r *skillsMasterRepository) UpdatePopularity(ctx context.Context, id int64, score float64) error {
	return r.db.WithContext(ctx).
		Model(&master.SkillsMaster{}).
		Where("id = ?", id).
		Update("popularity_score", score).
		Error
}

// IncrementPopularity increases the popularity score
func (r *skillsMasterRepository) IncrementPopularity(ctx context.Context, id int64, amount float64) error {
	return r.db.WithContext(ctx).
		Model(&master.SkillsMaster{}).
		Where("id = ?", id).
		UpdateColumn("popularity_score", gorm.Expr("popularity_score + ?", amount)).
		Error
}

// GetMostPopular retrieves most popular skills
func (r *skillsMasterRepository) GetMostPopular(ctx context.Context, limit int) ([]master.SkillsMaster, error) {
	var skills []master.SkillsMaster
	query := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Where("is_active = ?", true).
		Order("popularity_score DESC, name ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&skills).Error
	return skills, err
}

// GetPopularByType retrieves popular skills filtered by type
func (r *skillsMasterRepository) GetPopularByType(ctx context.Context, skillType string, limit int) ([]master.SkillsMaster, error) {
	var skills []master.SkillsMaster
	query := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Where("is_active = ? AND skill_type = ?", true, skillType).
		Order("popularity_score DESC, name ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&skills).Error
	return skills, err
}

// GetByPopularityRange retrieves skills within popularity score range
func (r *skillsMasterRepository) GetByPopularityRange(ctx context.Context, minScore, maxScore float64) ([]master.SkillsMaster, error) {
	var skills []master.SkillsMaster
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Where("popularity_score >= ? AND popularity_score <= ?", minScore, maxScore).
		Where("is_active = ?", true).
		Order("popularity_score DESC, name ASC").
		Find(&skills).Error
	return skills, err
}

// Activate sets the skill as active
func (r *skillsMasterRepository) Activate(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&master.SkillsMaster{}).
		Where("id = ?", id).
		Update("is_active", true).
		Error
}

// Deactivate sets the skill as inactive
func (r *skillsMasterRepository) Deactivate(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&master.SkillsMaster{}).
		Where("id = ?", id).
		Update("is_active", false).
		Error
}

// AddAlias adds an alias to a skill
func (r *skillsMasterRepository) AddAlias(ctx context.Context, id int64, alias string) error {
	return r.db.WithContext(ctx).
		Model(&master.SkillsMaster{}).
		Where("id = ?", id).
		Update("aliases", gorm.Expr("array_append(aliases, ?)", strings.ToLower(alias))).
		Error
}

// RemoveAlias removes an alias from a skill
func (r *skillsMasterRepository) RemoveAlias(ctx context.Context, id int64, alias string) error {
	return r.db.WithContext(ctx).
		Model(&master.SkillsMaster{}).
		Where("id = ?", id).
		Update("aliases", gorm.Expr("array_remove(aliases, ?)", strings.ToLower(alias))).
		Error
}

// UpdateAliases replaces all aliases for a skill
func (r *skillsMasterRepository) UpdateAliases(ctx context.Context, id int64, aliases []string) error {
	// Convert to lowercase
	lowerAliases := make([]string, len(aliases))
	for i, alias := range aliases {
		lowerAliases[i] = strings.ToLower(alias)
	}

	return r.db.WithContext(ctx).
		Model(&master.SkillsMaster{}).
		Where("id = ?", id).
		Update("aliases", pq.Array(lowerAliases)).
		Error
}

// BulkAddAliases adds multiple aliases to a skill
func (r *skillsMasterRepository) BulkAddAliases(ctx context.Context, id int64, aliases []string) error {
	if len(aliases) == 0 {
		return nil
	}

	// Convert to lowercase
	lowerAliases := make([]string, len(aliases))
	for i, alias := range aliases {
		lowerAliases[i] = strings.ToLower(alias)
	}

	return r.db.WithContext(ctx).
		Model(&master.SkillsMaster{}).
		Where("id = ?", id).
		Update("aliases", gorm.Expr("aliases || ?", pq.Array(lowerAliases))).
		Error
}

// BulkCreate creates multiple skills at once
func (r *skillsMasterRepository) BulkCreate(ctx context.Context, skills []master.SkillsMaster) error {
	return r.db.WithContext(ctx).Create(&skills).Error
}

// BulkUpdatePopularity updates popularity scores for multiple skills
func (r *skillsMasterRepository) BulkUpdatePopularity(ctx context.Context, updates map[int64]float64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for id, score := range updates {
			if err := tx.Model(&master.SkillsMaster{}).
				Where("id = ?", id).
				Update("popularity_score", score).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// Count returns total number of skill records
func (r *skillsMasterRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&master.SkillsMaster{}).Count(&count).Error
	return count, err
}

// CountActive returns total number of active skill records
func (r *skillsMasterRepository) CountActive(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&master.SkillsMaster{}).
		Where("is_active = ?", true).
		Count(&count).Error
	return count, err
}

// CountByType returns count of skills by type
func (r *skillsMasterRepository) CountByType(ctx context.Context, skillType string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&master.SkillsMaster{}).
		Where("skill_type = ?", skillType).
		Count(&count).Error
	return count, err
}

// CountByDifficulty returns count of skills by difficulty
func (r *skillsMasterRepository) CountByDifficulty(ctx context.Context, difficultyLevel string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&master.SkillsMaster{}).
		Where("difficulty_level = ?", difficultyLevel).
		Count(&count).Error
	return count, err
}

// GetSkillStats returns comprehensive skill statistics
func (r *skillsMasterRepository) GetSkillStats(ctx context.Context) (*master.SkillStats, error) {
	stats := &master.SkillStats{
		ByType:       make(map[string]int64),
		ByDifficulty: make(map[string]int64),
		ByCategory:   make(map[int64]int64),
	}

	// Total count
	if err := r.db.WithContext(ctx).Model(&master.SkillsMaster{}).Count(&stats.TotalSkills).Error; err != nil {
		return nil, err
	}

	// Active count
	if err := r.db.WithContext(ctx).
		Model(&master.SkillsMaster{}).
		Where("is_active = ?", true).
		Count(&stats.ActiveSkills).Error; err != nil {
		return nil, err
	}

	stats.InactiveSkills = stats.TotalSkills - stats.ActiveSkills

	// Count by type
	type TypeCount struct {
		Type  string
		Count int64
	}
	var typeCounts []TypeCount
	if err := r.db.WithContext(ctx).
		Model(&master.SkillsMaster{}).
		Select("skill_type as type, COUNT(*) as count").
		Group("skill_type").
		Order("count DESC").
		Scan(&typeCounts).Error; err != nil {
		return nil, err
	}

	for _, tc := range typeCounts {
		stats.ByType[tc.Type] = tc.Count
		percentage := float64(tc.Count) / float64(stats.TotalSkills) * 100
		stats.TopTypes = append(stats.TopTypes, master.TypeStat{
			Type:       tc.Type,
			Count:      tc.Count,
			Percentage: percentage,
		})
	}

	// Count by difficulty
	type DifficultyCount struct {
		Difficulty string
		Count      int64
	}
	var difficultyCounts []DifficultyCount
	if err := r.db.WithContext(ctx).
		Model(&master.SkillsMaster{}).
		Select("difficulty_level as difficulty, COUNT(*) as count").
		Group("difficulty_level").
		Order("count DESC").
		Scan(&difficultyCounts).Error; err != nil {
		return nil, err
	}

	for _, dc := range difficultyCounts {
		stats.ByDifficulty[dc.Difficulty] = dc.Count
	}

	// Count by category
	type CategoryCount struct {
		CategoryID int64
		Count      int64
	}
	var categoryCounts []CategoryCount
	if err := r.db.WithContext(ctx).
		Model(&master.SkillsMaster{}).
		Where("category_id IS NOT NULL").
		Select("category_id, COUNT(*) as count").
		Group("category_id").
		Order("count DESC").
		Scan(&categoryCounts).Error; err != nil {
		return nil, err
	}

	for _, cc := range categoryCounts {
		stats.ByCategory[cc.CategoryID] = cc.Count
	}

	// Root skills count
	if err := r.db.WithContext(ctx).
		Model(&master.SkillsMaster{}).
		Where("parent_id IS NULL").
		Count(&stats.RootSkills).Error; err != nil {
		return nil, err
	}

	// Child skills count
	if err := r.db.WithContext(ctx).
		Model(&master.SkillsMaster{}).
		Where("parent_id IS NOT NULL").
		Count(&stats.ChildSkills).Error; err != nil {
		return nil, err
	}

	// Average popularity score
	if err := r.db.WithContext(ctx).
		Model(&master.SkillsMaster{}).
		Select("COALESCE(AVG(popularity_score), 0)").
		Scan(&stats.AveragePopularity).Error; err != nil {
		return nil, err
	}

	// Most popular skill
	var mostPopular master.SkillsMaster
	if err := r.db.WithContext(ctx).
		Order("popularity_score DESC, name ASC").
		First(&mostPopular).Error; err == nil {
		stats.MostPopular = &mostPopular
	}

	return stats, nil
}

// GetRelatedSkills retrieves related skills based on category and type
func (r *skillsMasterRepository) GetRelatedSkills(ctx context.Context, skillID int64, limit int) ([]master.SkillsMaster, error) {
	// First get the skill to find its category and type
	skill, err := r.FindByID(ctx, skillID)
	if err != nil || skill == nil {
		return nil, err
	}

	var skills []master.SkillsMaster
	query := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Where("id != ?", skillID).
		Where("is_active = ?", true)

	// Match by category or skill type
	if skill.CategoryID != nil {
		query = query.Where("(category_id = ? OR skill_type = ?)", *skill.CategoryID, skill.SkillType)
	} else {
		query = query.Where("skill_type = ?", skill.SkillType)
	}

	query = query.Order("popularity_score DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err = query.Find(&skills).Error
	return skills, err
}

// GetComplementarySkills retrieves skills that complement the given skills
func (r *skillsMasterRepository) GetComplementarySkills(ctx context.Context, skillIDs []int64, limit int) ([]master.SkillsMaster, error) {
	if len(skillIDs) == 0 {
		return []master.SkillsMaster{}, nil
	}

	// Get the primary skill type from the first skill
	var primarySkill master.SkillsMaster
	if err := r.db.WithContext(ctx).First(&primarySkill, skillIDs[0]).Error; err != nil {
		return nil, err
	}

	var skills []master.SkillsMaster
	query := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Where("id NOT IN ?", skillIDs).
		Where("is_active = ?", true)

	// Get skills from different types (complementary)
	if primarySkill.CategoryID != nil {
		query = query.Where("category_id = ? AND skill_type != ?", *primarySkill.CategoryID, primarySkill.SkillType)
	} else {
		query = query.Where("skill_type != ?", primarySkill.SkillType)
	}

	query = query.Order("popularity_score DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&skills).Error
	return skills, err
}
