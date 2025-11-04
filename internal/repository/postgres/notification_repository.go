package postgres

import (
	"context"
	"time"

	"keerja-backend/internal/domain/notification"

	"gorm.io/gorm"
)

// notificationRepository implements the notification.NotificationRepository interface
type notificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository creates a new notification repository instance
func NewNotificationRepository(db *gorm.DB) notification.NotificationRepository {
	return &notificationRepository{db: db}
}

// ============================================================================
// Notification CRUD Operations
// ============================================================================

// Create creates a new notification
func (r *notificationRepository) Create(ctx context.Context, notif *notification.Notification) error {
	return r.db.WithContext(ctx).Create(notif).Error
}

// FindByID finds notification by ID
func (r *notificationRepository) FindByID(ctx context.Context, id int64) (*notification.Notification, error) {
	var notif notification.Notification
	err := r.db.WithContext(ctx).First(&notif, id).Error
	if err != nil {
		return nil, err
	}
	return &notif, nil
}

// Update updates notification
func (r *notificationRepository) Update(ctx context.Context, notif *notification.Notification) error {
	return r.db.WithContext(ctx).Save(notif).Error
}

// Delete deletes notification
func (r *notificationRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&notification.Notification{}, id).Error
}

// ============================================================================
// Notification Listing and Filtering
// ============================================================================

// ListByUser lists notifications for user with filtering and pagination
func (r *notificationRepository) ListByUser(ctx context.Context, userID int64, filter notification.NotificationFilter, page, limit int) ([]notification.Notification, int64, error) {
	var notifs []notification.Notification
	var total int64

	query := r.db.WithContext(ctx).Model(&notification.Notification{}).Where("user_id = ?", userID)
	query = r.applyNotificationFilter(query, filter)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	err := query.
		Order("priority DESC, created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&notifs).Error

	if err != nil {
		return nil, 0, err
	}

	return notifs, total, nil
}

// GetUnreadByUser retrieves unread notifications for user
func (r *notificationRepository) GetUnreadByUser(ctx context.Context, userID int64, limit int) ([]notification.Notification, error) {
	var notifs []notification.Notification
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_read = ?", userID, false).
		Order("priority DESC, created_at DESC").
		Limit(limit).
		Find(&notifs).Error

	if err != nil {
		return nil, err
	}

	return notifs, nil
}

// CountUnreadByUser counts unread notifications for user
func (r *notificationRepository) CountUnreadByUser(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&notification.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error

	return count, err
}

// MarkAsRead marks notification as read
func (r *notificationRepository) MarkAsRead(ctx context.Context, id int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&notification.Notification{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

// MarkAllAsRead marks all notifications as read for user
func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&notification.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

// DeleteByUser deletes all notifications for user
func (r *notificationRepository) DeleteByUser(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&notification.Notification{}).Error
}

// GetExpiredNotifications retrieves expired notifications
func (r *notificationRepository) GetExpiredNotifications(ctx context.Context, limit int) ([]notification.Notification, error) {
	var notifs []notification.Notification
	now := time.Now()

	err := r.db.WithContext(ctx).
		Where("expires_at IS NOT NULL AND expires_at < ?", now).
		Limit(limit).
		Find(&notifs).Error

	if err != nil {
		return nil, err
	}

	return notifs, nil
}

// GetStats retrieves notification statistics
func (r *notificationRepository) GetStats(ctx context.Context, userID int64) (*notification.NotificationStats, error) {
	stats := &notification.NotificationStats{
		CategoryBreakdown: make(map[string]int64),
	}

	// Get total count
	if err := r.db.WithContext(ctx).
		Model(&notification.Notification{}).
		Where("user_id = ?", userID).
		Count(&stats.TotalCount).Error; err != nil {
		return nil, err
	}

	// Get unread count
	if err := r.db.WithContext(ctx).
		Model(&notification.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&stats.UnreadCount).Error; err != nil {
		return nil, err
	}

	// Get read count
	stats.ReadCount = stats.TotalCount - stats.UnreadCount

	// Get today's count
	today := time.Now().Truncate(24 * time.Hour)
	if err := r.db.WithContext(ctx).
		Model(&notification.Notification{}).
		Where("user_id = ? AND created_at >= ?", userID, today).
		Count(&stats.TodayCount).Error; err != nil {
		return nil, err
	}

	// Get this week's count
	weekStart := time.Now().AddDate(0, 0, -7)
	if err := r.db.WithContext(ctx).
		Model(&notification.Notification{}).
		Where("user_id = ? AND created_at >= ?", userID, weekStart).
		Count(&stats.ThisWeekCount).Error; err != nil {
		return nil, err
	}

	// Get high priority count
	if err := r.db.WithContext(ctx).
		Model(&notification.Notification{}).
		Where("user_id = ? AND (priority = ? OR priority = ?)", userID, "high", "urgent").
		Count(&stats.HighPriorityCount).Error; err != nil {
		return nil, err
	}

	// Get category breakdown
	type CategoryCount struct {
		Category string
		Count    int64
	}
	var categoryCounts []CategoryCount
	if err := r.db.WithContext(ctx).
		Model(&notification.Notification{}).
		Select("category, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("category").
		Scan(&categoryCounts).Error; err != nil {
		return nil, err
	}

	for _, cc := range categoryCounts {
		stats.CategoryBreakdown[cc.Category] = cc.Count
	}

	return stats, nil
}

// BulkCreate creates multiple notifications
func (r *notificationRepository) BulkCreate(ctx context.Context, notifications []notification.Notification) error {
	if len(notifications) == 0 {
		return nil
	}

	// Use batch insert for better performance
	return r.db.WithContext(ctx).CreateInBatches(notifications, 100).Error
}

// ============================================================================
// Notification Preferences Operations
// ============================================================================

// FindPreferenceByUser finds notification preferences for user
func (r *notificationRepository) FindPreferenceByUser(ctx context.Context, userID int64) (*notification.NotificationPreference, error) {
	var pref notification.NotificationPreference
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&pref).Error
	if err != nil {
		return nil, err
	}
	return &pref, nil
}

// CreatePreference creates notification preferences
func (r *notificationRepository) CreatePreference(ctx context.Context, preference *notification.NotificationPreference) error {
	return r.db.WithContext(ctx).Create(preference).Error
}

// UpdatePreference updates notification preferences
func (r *notificationRepository) UpdatePreference(ctx context.Context, preference *notification.NotificationPreference) error {
	return r.db.WithContext(ctx).Save(preference).Error
}

// ============================================================================
// Helper Methods
// ============================================================================

// applyNotificationFilter applies filters to notification query
func (r *notificationRepository) applyNotificationFilter(query *gorm.DB, filter notification.NotificationFilter) *gorm.DB {
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}

	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}

	if filter.IsRead != nil {
		query = query.Where("is_read = ?", *filter.IsRead)
	}

	if filter.Priority != "" {
		query = query.Where("priority = ?", filter.Priority)
	}

	if filter.DateFrom != nil {
		query = query.Where("created_at >= ?", filter.DateFrom)
	}

	if filter.DateTo != nil {
		query = query.Where("created_at <= ?", filter.DateTo)
	}

	return query
}
