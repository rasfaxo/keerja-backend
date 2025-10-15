package postgres

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"keerja-backend/internal/domain/admin"
)

// adminUserRepository implements admin.AdminUserRepository interface
type adminUserRepository struct {
	db *gorm.DB
}

// NewAdminUserRepository creates a new instance of admin user repository
func NewAdminUserRepository(db *gorm.DB) admin.AdminUserRepository {
	return &adminUserRepository{db: db}
}

// Create creates a new admin user record
func (r *adminUserRepository) Create(ctx context.Context, user *admin.AdminUser) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByID retrieves an admin user by ID
func (r *adminUserRepository) FindByID(ctx context.Context, id int64) (*admin.AdminUser, error) {
	var user admin.AdminUser
	err := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Creator").
		Where("id = ?", id).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByUUID retrieves an admin user by UUID
func (r *adminUserRepository) FindByUUID(ctx context.Context, uuid string) (*admin.AdminUser, error) {
	var user admin.AdminUser
	err := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Creator").
		Where("uuid = ?", uuid).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail retrieves an admin user by email
func (r *adminUserRepository) FindByEmail(ctx context.Context, email string) (*admin.AdminUser, error) {
	var user admin.AdminUser
	err := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Creator").
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Update updates an existing admin user record
func (r *adminUserRepository) Update(ctx context.Context, user *admin.AdminUser) error {
	return r.db.WithContext(ctx).Model(user).Updates(user).Error
}

// Delete soft deletes an admin user record
func (r *adminUserRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&admin.AdminUser{}).Error
}

// List retrieves users with filtering and pagination
func (r *adminUserRepository) List(ctx context.Context, filter *admin.AdminUserFilter) ([]admin.AdminUser, int64, error) {
	var users []admin.AdminUser
	query := r.db.WithContext(ctx).Model(&admin.AdminUser{})

	// Apply filters
	if filter != nil {
		if filter.Search != "" {
			searchPattern := "%" + strings.ToLower(filter.Search) + "%"
			query = query.Where("LOWER(full_name) LIKE ? OR LOWER(email) LIKE ?", searchPattern, searchPattern)
		}
		if filter.Status != "" {
			query = query.Where("status = ?", filter.Status)
		}
		if filter.RoleID != nil {
			query = query.Where("role_id = ?", *filter.RoleID)
		}
		if filter.MinAccessLevel != nil && filter.RoleID == nil {
			query = query.Joins("LEFT JOIN admin_roles ON admin_roles.id = admin_users.role_id").
				Where("admin_roles.access_level >= ?", *filter.MinAccessLevel)
		}
		if filter.Has2FA != nil {
			if *filter.Has2FA {
				query = query.Where("two_factor_secret IS NOT NULL AND two_factor_secret != ''")
			} else {
				query = query.Where("two_factor_secret IS NULL OR two_factor_secret = ''")
			}
		}
		if filter.CreatedBy != nil {
			query = query.Where("created_by = ?", *filter.CreatedBy)
		}
		if filter.LastLoginAfter != nil {
			query = query.Where("last_login >= ?", *filter.LastLoginAfter)
		}
		if filter.LastLoginBefore != nil {
			query = query.Where("last_login <= ?", *filter.LastLoginBefore)
		}
		if filter.CreatedAfter != nil {
			query = query.Where("created_at >= ?", *filter.CreatedAfter)
		}
		if filter.CreatedBefore != nil {
			query = query.Where("created_at <= ?", *filter.CreatedBefore)
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
		sortField := "created_at"
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
		query = query.Order("created_at DESC")
	}

	// Load relationships
	query = query.Preload("Role").Preload("Creator")

	err := query.Find(&users).Error
	return users, total, err
}

// ListByRole retrieves users by role ID with pagination
func (r *adminUserRepository) ListByRole(ctx context.Context, roleID int64, page, pageSize int) ([]admin.AdminUser, int64, error) {
	var users []admin.AdminUser
	query := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Creator").
		Where("role_id = ?", roleID)

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.
		Order("full_name ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&users).Error

	return users, total, err
}

// ListByStatus retrieves users by status with pagination
func (r *adminUserRepository) ListByStatus(ctx context.Context, status string, page, pageSize int) ([]admin.AdminUser, int64, error) {
	var users []admin.AdminUser
	query := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Creator").
		Where("status = ?", status)

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := query.
		Order("full_name ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&users).Error

	return users, total, err
}

// SearchUsers searches users by query string with pagination
func (r *adminUserRepository) SearchUsers(ctx context.Context, query string, page, pageSize int) ([]admin.AdminUser, int64, error) {
	var users []admin.AdminUser
	searchQuery := "%" + strings.ToLower(query) + "%"

	dbQuery := r.db.WithContext(ctx).Model(&admin.AdminUser{}).
		Preload("Role").
		Preload("Creator").
		Where("LOWER(full_name) LIKE ? OR LOWER(email) LIKE ? OR phone LIKE ?", searchQuery, searchQuery, searchQuery)

	// Count total
	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	err := dbQuery.
		Order("full_name ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&users).Error

	return users, total, err
}

// UpdateStatus updates the status of an admin user
func (r *adminUserRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).
		Model(&admin.AdminUser{}).
		Where("id = ?", id).
		Update("status", status).
		Error
}

// ActivateUser sets the user status to active
func (r *adminUserRepository) ActivateUser(ctx context.Context, id int64) error {
	return r.UpdateStatus(ctx, id, "active")
}

// DeactivateUser sets the user status to inactive
func (r *adminUserRepository) DeactivateUser(ctx context.Context, id int64) error {
	return r.UpdateStatus(ctx, id, "inactive")
}

// SuspendUser sets the user status to suspended
func (r *adminUserRepository) SuspendUser(ctx context.Context, id int64) error {
	return r.UpdateStatus(ctx, id, "suspended")
}

// GetActiveUsers retrieves all active users
func (r *adminUserRepository) GetActiveUsers(ctx context.Context) ([]admin.AdminUser, error) {
	var users []admin.AdminUser
	err := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Creator").
		Where("status = ?", "active").
		Order("full_name ASC").
		Find(&users).Error
	return users, err
}

// GetInactiveUsers retrieves all inactive users
func (r *adminUserRepository) GetInactiveUsers(ctx context.Context) ([]admin.AdminUser, error) {
	var users []admin.AdminUser
	err := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Creator").
		Where("status = ?", "inactive").
		Order("full_name ASC").
		Find(&users).Error
	return users, err
}

// GetSuspendedUsers retrieves all suspended users
func (r *adminUserRepository) GetSuspendedUsers(ctx context.Context) ([]admin.AdminUser, error) {
	var users []admin.AdminUser
	err := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Creator").
		Where("status = ?", "suspended").
		Order("full_name ASC").
		Find(&users).Error
	return users, err
}

// UpdateRole updates the role of an admin user
func (r *adminUserRepository) UpdateRole(ctx context.Context, userID, roleID int64) error {
	return r.db.WithContext(ctx).
		Model(&admin.AdminUser{}).
		Where("id = ?", userID).
		Update("role_id", roleID).
		Error
}

// GetUsersByRole retrieves all users with a specific role
func (r *adminUserRepository) GetUsersByRole(ctx context.Context, roleID int64) ([]admin.AdminUser, error) {
	var users []admin.AdminUser
	err := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Creator").
		Where("role_id = ?", roleID).
		Order("full_name ASC").
		Find(&users).Error
	return users, err
}

// UpdatePassword updates the password hash of an admin user
func (r *adminUserRepository) UpdatePassword(ctx context.Context, id int64, passwordHash string) error {
	return r.db.WithContext(ctx).
		Model(&admin.AdminUser{}).
		Where("id = ?", id).
		Update("password_hash", passwordHash).
		Error
}

// UpdateLastLogin updates the last login timestamp
func (r *adminUserRepository) UpdateLastLogin(ctx context.Context, id int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&admin.AdminUser{}).
		Where("id = ?", id).
		Update("last_login", now).
		Error
}

// Enable2FA enables two-factor authentication for a user
func (r *adminUserRepository) Enable2FA(ctx context.Context, id int64, secret string) error {
	return r.db.WithContext(ctx).
		Model(&admin.AdminUser{}).
		Where("id = ?", id).
		Update("two_factor_secret", secret).
		Error
}

// Disable2FA disables two-factor authentication for a user
func (r *adminUserRepository) Disable2FA(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&admin.AdminUser{}).
		Where("id = ?", id).
		Update("two_factor_secret", "").
		Error
}

// Get2FAUsers retrieves all users with 2FA enabled
func (r *adminUserRepository) Get2FAUsers(ctx context.Context) ([]admin.AdminUser, error) {
	var users []admin.AdminUser
	err := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Creator").
		Where("two_factor_secret IS NOT NULL AND two_factor_secret != ''").
		Order("full_name ASC").
		Find(&users).Error
	return users, err
}

// UpdateProfile updates profile information of an admin user
func (r *adminUserRepository) UpdateProfile(ctx context.Context, id int64, fullName, phone, profileImageURL string) error {
	updates := map[string]interface{}{
		"full_name": fullName,
		"phone":     phone,
	}
	if profileImageURL != "" {
		updates["profile_image_url"] = profileImageURL
	}

	return r.db.WithContext(ctx).
		Model(&admin.AdminUser{}).
		Where("id = ?", id).
		Updates(updates).
		Error
}

// UpdateProfileImage updates the profile image URL
func (r *adminUserRepository) UpdateProfileImage(ctx context.Context, id int64, imageURL string) error {
	return r.db.WithContext(ctx).
		Model(&admin.AdminUser{}).
		Where("id = ?", id).
		Update("profile_image_url", imageURL).
		Error
}

// Count returns total number of user records
func (r *adminUserRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&admin.AdminUser{}).Count(&count).Error
	return count, err
}

// CountByStatus returns count of users by status
func (r *adminUserRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&admin.AdminUser{}).
		Where("status = ?", status).
		Count(&count).Error
	return count, err
}

// CountByRole returns count of users by role
func (r *adminUserRepository) CountByRole(ctx context.Context, roleID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&admin.AdminUser{}).
		Where("role_id = ?", roleID).
		Count(&count).Error
	return count, err
}

// GetUserStats returns comprehensive user statistics
func (r *adminUserRepository) GetUserStats(ctx context.Context) (*admin.AdminUserStats, error) {
	stats := &admin.AdminUserStats{
		ByRole: make(map[int64]int64),
	}

	// Total count
	if err := r.db.WithContext(ctx).Model(&admin.AdminUser{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, err
	}

	// Active users count
	if err := r.db.WithContext(ctx).
		Model(&admin.AdminUser{}).
		Where("status = ?", "active").
		Count(&stats.ActiveUsers).Error; err != nil {
		return nil, err
	}

	// Inactive users count
	if err := r.db.WithContext(ctx).
		Model(&admin.AdminUser{}).
		Where("status = ?", "inactive").
		Count(&stats.InactiveUsers).Error; err != nil {
		return nil, err
	}

	// Suspended users count
	if err := r.db.WithContext(ctx).
		Model(&admin.AdminUser{}).
		Where("status = ?", "suspended").
		Count(&stats.SuspendedUsers).Error; err != nil {
		return nil, err
	}

	// 2FA enabled users count
	if err := r.db.WithContext(ctx).
		Model(&admin.AdminUser{}).
		Where("two_factor_secret IS NOT NULL AND two_factor_secret != ''").
		Count(&stats.Users2FA).Error; err != nil {
		return nil, err
	}

	// Count by role
	type RoleCount struct {
		RoleID int64
		Count  int64
	}
	var roleCounts []RoleCount
	if err := r.db.WithContext(ctx).
		Model(&admin.AdminUser{}).
		Where("role_id IS NOT NULL").
		Select("role_id, COUNT(*) as count").
		Group("role_id").
		Scan(&roleCounts).Error; err != nil {
		return nil, err
	}

	for _, rc := range roleCounts {
		stats.ByRole[rc.RoleID] = rc.Count
	}

	// Count super admins (access level >= 9)
	if err := r.db.WithContext(ctx).
		Table("admin_users").
		Joins("LEFT JOIN admin_roles ON admin_roles.id = admin_users.role_id").
		Where("admin_users.deleted_at IS NULL AND admin_roles.access_level >= ?", 9).
		Count(&stats.SuperAdmins).Error; err != nil {
		return nil, err
	}

	// Count admins (access level >= 7 and < 9)
	if err := r.db.WithContext(ctx).
		Table("admin_users").
		Joins("LEFT JOIN admin_roles ON admin_roles.id = admin_users.role_id").
		Where("admin_users.deleted_at IS NULL AND admin_roles.access_level >= ? AND admin_roles.access_level < ?", 7, 9).
		Count(&stats.Admins).Error; err != nil {
		return nil, err
	}

	// Count moderators (access level >= 5 and < 7)
	if err := r.db.WithContext(ctx).
		Table("admin_users").
		Joins("LEFT JOIN admin_roles ON admin_roles.id = admin_users.role_id").
		Where("admin_users.deleted_at IS NULL AND admin_roles.access_level >= ? AND admin_roles.access_level < ?", 5, 7).
		Count(&stats.Moderators).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

// GetActivityStats returns activity statistics for a date range
func (r *adminUserRepository) GetActivityStats(ctx context.Context, startDate, endDate time.Time) (*admin.AdminActivityStats, error) {
	stats := &admin.AdminActivityStats{
		Period:       startDate.Format("2006-01-02") + " to " + endDate.Format("2006-01-02"),
		LoginsByDate: make(map[string]int64),
		LoginsByUser: make(map[int64]int64),
	}

	// This would require a separate login_history or activity_logs table
	// For now, we'll use last_login field with limitations

	// Count unique users who logged in during period
	if err := r.db.WithContext(ctx).
		Model(&admin.AdminUser{}).
		Where("last_login >= ? AND last_login <= ?", startDate, endDate).
		Count(&stats.UniqueUsers).Error; err != nil {
		return nil, err
	}

	// Get users with recent logins for top active users
	var topUsers []admin.AdminUser
	if err := r.db.WithContext(ctx).
		Preload("Role").
		Where("last_login >= ? AND last_login <= ?", startDate, endDate).
		Order("last_login DESC").
		Limit(10).
		Find(&topUsers).Error; err == nil {
		stats.TopActiveUsers = topUsers
	}

	// Note: For accurate login counts and history, implement a separate activity log table
	stats.TotalLogins = stats.UniqueUsers // Simplified
	if stats.UniqueUsers > 0 {
		stats.AverageLogins = float64(stats.TotalLogins) / float64(stats.UniqueUsers)
	}

	return stats, nil
}

// GetRecentLogins retrieves users with most recent logins
func (r *adminUserRepository) GetRecentLogins(ctx context.Context, limit int) ([]admin.AdminUser, error) {
	var users []admin.AdminUser
	query := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Creator").
		Where("last_login IS NOT NULL").
		Order("last_login DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&users).Error
	return users, err
}

// GetCreatedUsers retrieves all users created by a specific admin
func (r *adminUserRepository) GetCreatedUsers(ctx context.Context, creatorID int64) ([]admin.AdminUser, error) {
	var users []admin.AdminUser
	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("created_by = ?", creatorID).
		Order("created_at DESC").
		Find(&users).Error
	return users, err
}

// GetUserActivity retrieves activity information for a specific user
func (r *adminUserRepository) GetUserActivity(ctx context.Context, userID int64, startDate, endDate time.Time) (*admin.UserActivity, error) {
	activity := &admin.UserActivity{
		UserID:       userID,
		LoginHistory: []time.Time{},
	}

	// Get user
	user, err := r.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, err
	}

	// Last login
	activity.LastLogin = user.LastLogin

	// Count created users
	if err := r.db.WithContext(ctx).
		Model(&admin.AdminUser{}).
		Where("created_by = ? AND created_at >= ? AND created_at <= ?", userID, startDate, endDate).
		Count(&activity.CreatedUsers).Error; err != nil {
		return nil, err
	}

	// Count created roles
	if err := r.db.WithContext(ctx).
		Model(&admin.AdminRole{}).
		Where("created_by = ? AND created_at >= ? AND created_at <= ?", userID, startDate, endDate).
		Count(&activity.CreatedRoles).Error; err != nil {
		return nil, err
	}

	// Note: For accurate login history, implement a separate activity log table
	// For now, we just include the last_login if it's within the date range
	if user.LastLogin != nil && user.LastLogin.After(startDate) && user.LastLogin.Before(endDate) {
		activity.LoginHistory = append(activity.LoginHistory, *user.LastLogin)
		activity.TotalLogins = 1
	}

	return activity, nil
}
