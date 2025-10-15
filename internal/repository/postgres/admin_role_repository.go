package postgres

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"keerja-backend/internal/domain/admin"
)

// adminRoleRepository implements admin.AdminRoleRepository interface
type adminRoleRepository struct {
	db *gorm.DB
}

// NewAdminRoleRepository creates a new instance of admin role repository
func NewAdminRoleRepository(db *gorm.DB) admin.AdminRoleRepository {
	return &adminRoleRepository{db: db}
}

// Create creates a new admin role record
func (r *adminRoleRepository) Create(ctx context.Context, role *admin.AdminRole) error {
	return r.db.WithContext(ctx).Create(role).Error
}

// FindByID retrieves an admin role by ID
func (r *adminRoleRepository) FindByID(ctx context.Context, id int64) (*admin.AdminRole, error) {
	var role admin.AdminRole
	err := r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Users").
		Where("id = ?", id).
		First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

// FindByName retrieves an admin role by name
func (r *adminRoleRepository) FindByName(ctx context.Context, roleName string) (*admin.AdminRole, error) {
	var role admin.AdminRole
	err := r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Users").
		Where("role_name = ?", roleName).
		First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

// Update updates an existing admin role record
func (r *adminRoleRepository) Update(ctx context.Context, role *admin.AdminRole) error {
	return r.db.WithContext(ctx).Model(role).Updates(role).Error
}

// Delete soft deletes an admin role record
func (r *adminRoleRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&admin.AdminRole{}).Error
}

// List retrieves roles with filtering and pagination
func (r *adminRoleRepository) List(ctx context.Context, filter *admin.AdminRoleFilter) ([]admin.AdminRole, int64, error) {
	var roles []admin.AdminRole
	query := r.db.WithContext(ctx).Model(&admin.AdminRole{})

	// Apply filters
	if filter != nil {
		if filter.Search != "" {
			searchPattern := "%" + strings.ToLower(filter.Search) + "%"
			query = query.Where("LOWER(role_name) LIKE ? OR LOWER(role_description) LIKE ?", searchPattern, searchPattern)
		}
		if filter.AccessLevel != nil {
			query = query.Where("access_level = ?", *filter.AccessLevel)
		}
		if filter.MinLevel != nil {
			query = query.Where("access_level >= ?", *filter.MinLevel)
		}
		if filter.MaxLevel != nil {
			query = query.Where("access_level <= ?", *filter.MaxLevel)
		}
		if filter.IsSystemRole != nil {
			query = query.Where("is_system_role = ?", *filter.IsSystemRole)
		}
		if filter.CreatedBy != nil {
			query = query.Where("created_by = ?", *filter.CreatedBy)
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
		sortField := "access_level"
		sortOrder := "DESC"
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
		query = query.Order("access_level DESC, role_name ASC")
	}

	// Load relationships
	query = query.Preload("Creator").Preload("Users")

	err := query.Find(&roles).Error
	return roles, total, err
}

// ListActive retrieves all active (non-deleted) roles
func (r *adminRoleRepository) ListActive(ctx context.Context) ([]admin.AdminRole, error) {
	var roles []admin.AdminRole
	err := r.db.WithContext(ctx).
		Preload("Creator").
		Order("access_level DESC, role_name ASC").
		Find(&roles).Error
	return roles, err
}

// GetSystemRoles retrieves all system roles
func (r *adminRoleRepository) GetSystemRoles(ctx context.Context) ([]admin.AdminRole, error) {
	var roles []admin.AdminRole
	err := r.db.WithContext(ctx).
		Preload("Creator").
		Where("is_system_role = ?", true).
		Order("access_level DESC").
		Find(&roles).Error
	return roles, err
}

// GetNonSystemRoles retrieves all non-system (custom) roles
func (r *adminRoleRepository) GetNonSystemRoles(ctx context.Context) ([]admin.AdminRole, error) {
	var roles []admin.AdminRole
	err := r.db.WithContext(ctx).
		Preload("Creator").
		Where("is_system_role = ?", false).
		Order("access_level DESC, role_name ASC").
		Find(&roles).Error
	return roles, err
}

// GetRolesByAccessLevel retrieves roles within access level range
func (r *adminRoleRepository) GetRolesByAccessLevel(ctx context.Context, minLevel, maxLevel int16) ([]admin.AdminRole, error) {
	var roles []admin.AdminRole
	err := r.db.WithContext(ctx).
		Preload("Creator").
		Where("access_level >= ? AND access_level <= ?", minLevel, maxLevel).
		Order("access_level DESC, role_name ASC").
		Find(&roles).Error
	return roles, err
}

// GetRolesByMinAccessLevel retrieves roles with minimum access level
func (r *adminRoleRepository) GetRolesByMinAccessLevel(ctx context.Context, minLevel int16) ([]admin.AdminRole, error) {
	var roles []admin.AdminRole
	err := r.db.WithContext(ctx).
		Preload("Creator").
		Where("access_level >= ?", minLevel).
		Order("access_level DESC, role_name ASC").
		Find(&roles).Error
	return roles, err
}

// Count returns total number of role records
func (r *adminRoleRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&admin.AdminRole{}).Count(&count).Error
	return count, err
}

// CountByAccessLevel returns count of roles by access level
func (r *adminRoleRepository) CountByAccessLevel(ctx context.Context, accessLevel int16) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&admin.AdminRole{}).
		Where("access_level = ?", accessLevel).
		Count(&count).Error
	return count, err
}

// GetRoleStats returns comprehensive role statistics
func (r *adminRoleRepository) GetRoleStats(ctx context.Context) (*admin.AdminRoleStats, error) {
	stats := &admin.AdminRoleStats{
		ByAccessLevel: make(map[int16]int64),
	}

	// Total count
	if err := r.db.WithContext(ctx).Model(&admin.AdminRole{}).Count(&stats.TotalRoles).Error; err != nil {
		return nil, err
	}

	// System roles count
	if err := r.db.WithContext(ctx).
		Model(&admin.AdminRole{}).
		Where("is_system_role = ?", true).
		Count(&stats.SystemRoles).Error; err != nil {
		return nil, err
	}

	stats.CustomRoles = stats.TotalRoles - stats.SystemRoles

	// Count by access level
	type AccessLevelCount struct {
		AccessLevel int16
		Count       int64
	}
	var accessLevelCounts []AccessLevelCount
	if err := r.db.WithContext(ctx).
		Model(&admin.AdminRole{}).
		Select("access_level, COUNT(*) as count").
		Group("access_level").
		Order("access_level DESC").
		Scan(&accessLevelCounts).Error; err != nil {
		return nil, err
	}

	for _, alc := range accessLevelCounts {
		stats.ByAccessLevel[alc.AccessLevel] = alc.Count
	}

	// Most used role (role with most users)
	type RoleUsage struct {
		RoleID int64
		Count  int64
	}
	var mostUsed RoleUsage
	if err := r.db.WithContext(ctx).
		Table("admin_users").
		Select("role_id, COUNT(*) as count").
		Where("role_id IS NOT NULL AND deleted_at IS NULL").
		Group("role_id").
		Order("count DESC").
		Limit(1).
		Scan(&mostUsed).Error; err == nil && mostUsed.RoleID > 0 {
		
		// Get the role details
		var role admin.AdminRole
		if err := r.db.WithContext(ctx).First(&role, mostUsed.RoleID).Error; err == nil {
			stats.MostUsedRole = &role
			stats.MostUsedRoleCount = mostUsed.Count
		}
	}

	return stats, nil
}
