package postgres

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"keerja-backend/internal/domain/email"
)

// emailRepository implements email.EmailRepository
type emailRepository struct {
	db *gorm.DB
}

// NewEmailRepository creates a new email repository instance
func NewEmailRepository(db *gorm.DB) email.EmailRepository {
	return &emailRepository{db: db}
}

// Create creates a new email log
func (r *emailRepository) Create(ctx context.Context, log *email.EmailLog) error {
	if err := r.db.WithContext(ctx).Create(log).Error; err != nil {
		return fmt.Errorf("failed to create email log: %w", err)
	}
	return nil
}

// FindByID finds email log by ID
func (r *emailRepository) FindByID(ctx context.Context, id int64) (*email.EmailLog, error) {
	var log email.EmailLog
	if err := r.db.WithContext(ctx).First(&log, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("email log not found")
		}
		return nil, fmt.Errorf("failed to find email log: %w", err)
	}
	return &log, nil
}

// Update updates email log
func (r *emailRepository) Update(ctx context.Context, log *email.EmailLog) error {
	if err := r.db.WithContext(ctx).Save(log).Error; err != nil {
		return fmt.Errorf("failed to update email log: %w", err)
	}
	return nil
}

// List retrieves email logs with pagination
func (r *emailRepository) List(ctx context.Context, filter email.EmailFilter, page, limit int) ([]email.EmailLog, int64, error) {
	var logs []email.EmailLog
	var total int64

	query := r.db.WithContext(ctx).Model(&email.EmailLog{})

	// Apply filters
	if filter.Recipient != "" {
		query = query.Where("recipient LIKE ?", "%"+filter.Recipient+"%")
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Template != "" {
		query = query.Where("template = ?", filter.Template)
	}
	if filter.DateFrom != nil {
		query = query.Where("created_at >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("created_at <= ?", *filter.DateTo)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count email logs: %w", err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get paginated results
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list email logs: %w", err)
	}

	return logs, total, nil
}

// GetFailedEmails retrieves failed emails
func (r *emailRepository) GetFailedEmails(ctx context.Context, page, limit int) ([]email.EmailLog, int64, error) {
	var logs []email.EmailLog
	var total int64

	query := r.db.WithContext(ctx).Model(&email.EmailLog{}).Where("status = ?", "failed")

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count failed emails: %w", err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get paginated results
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get failed emails: %w", err)
	}

	return logs, total, nil
}

// GetPendingEmails retrieves pending emails
func (r *emailRepository) GetPendingEmails(ctx context.Context, limit int) ([]email.EmailLog, error) {
	var logs []email.EmailLog

	if err := r.db.WithContext(ctx).
		Where("status = ?", "pending").
		Order("created_at ASC").
		Limit(limit).
		Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("failed to get pending emails: %w", err)
	}

	return logs, nil
}

// Delete deletes email log
func (r *emailRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&email.EmailLog{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete email log: %w", err)
	}
	return nil
}
