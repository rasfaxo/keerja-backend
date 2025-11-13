package master

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"keerja-backend/internal/cache"
	"keerja-backend/internal/domain/master"
)

// adminJobTypeServiceImpl implements AdminJobTypeService
type adminJobTypeServiceImpl struct {
	jobOptionsService master.JobOptionsService // For read operations
	repo              master.JobOptionsRepository
	db                *gorm.DB // For counting references
	cache             cache.Cache
}

// NewAdminJobTypeService creates a new AdminJobTypeService
func NewAdminJobTypeService(
	jobOptionsService master.JobOptionsService,
	repo master.JobOptionsRepository,
	db *gorm.DB,
	cache cache.Cache,
) master.AdminJobTypeService {
	return &adminJobTypeServiceImpl{
		jobOptionsService: jobOptionsService,
		repo:              repo,
		db:                db,
		cache:             cache,
	}
}

// GetJobTypes retrieves all job types
func (s *adminJobTypeServiceImpl) GetJobTypes(ctx context.Context) ([]master.JobType, error) {
	return s.jobOptionsService.GetJobTypes(ctx)
}

// GetJobTypeByID retrieves a job type by ID
func (s *adminJobTypeServiceImpl) GetJobTypeByID(ctx context.Context, id int64) (*master.JobType, error) {
	return s.repo.FindJobTypeByID(ctx, id)
}

// Create creates a new job type
func (s *adminJobTypeServiceImpl) Create(ctx context.Context, req master.CreateJobTypeRequest) (*master.JobTypeResponse, error) {
	// Check duplicate code
	exists, err := s.CheckDuplicateCode(ctx, req.Code, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check duplicate job type code: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("job type with code '%s' already exists", req.Code)
	}

	// Create job type entity
	jobType := &master.JobType{
		Name:  req.Name,
		Code:  req.Code,
		Order: req.Order,
	}

	// Create using GORM directly (since JobOptionsRepository might not have Create method)
	if err := s.db.WithContext(ctx).Table("job_types").Create(jobType).Error; err != nil {
		return nil, fmt.Errorf("failed to create job type: %w", err)
	}

	// Invalidate cache
	s.invalidateCache()

	// Map to response
	response := &master.JobTypeResponse{
		ID:    jobType.ID,
		Name:  jobType.Name,
		Code:  jobType.Code,
		Order: jobType.Order,
	}

	return response, nil
}

// Update updates an existing job type
func (s *adminJobTypeServiceImpl) Update(ctx context.Context, id int64, req master.UpdateJobTypeRequest) (*master.JobTypeResponse, error) {
	// Get existing job type
	jobType, err := s.repo.FindJobTypeByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get job type: %w", err)
	}
	if jobType == nil {
		return nil, errors.New("job type not found")
	}

	// Update fields
	if req.Name != "" {
		jobType.Name = req.Name
	}

	if req.Code != "" {
		// Check duplicate code if code is being changed
		if req.Code != jobType.Code {
			exists, err := s.CheckDuplicateCode(ctx, req.Code, &id)
			if err != nil {
				return nil, fmt.Errorf("failed to check duplicate job type code: %w", err)
			}
			if exists {
				return nil, fmt.Errorf("job type with code '%s' already exists", req.Code)
			}
		}
		jobType.Code = req.Code
	}

	if req.Order != nil {
		jobType.Order = *req.Order
	}

	// Update using GORM directly
	if err := s.db.WithContext(ctx).Table("job_types").Where("id = ?", id).Updates(jobType).Error; err != nil {
		return nil, fmt.Errorf("failed to update job type: %w", err)
	}

	// Invalidate cache
	s.invalidateCache()

	// Map to response
	response := &master.JobTypeResponse{
		ID:    jobType.ID,
		Name:  jobType.Name,
		Code:  jobType.Code,
		Order: jobType.Order,
	}

	return response, nil
}

// Delete deletes a job type if not referenced by jobs
func (s *adminJobTypeServiceImpl) Delete(ctx context.Context, id int64) error {
	// Check if job type exists
	_, err := s.repo.FindJobTypeByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get job type: %w", err)
	}

	// Check references
	jobs, err := s.CountJobReferences(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check job type references: %w", err)
	}
	if jobs > 0 {
		return fmt.Errorf("cannot delete job type: it is still referenced by %d jobs", jobs)
	}

	// Delete using GORM directly
	if err := s.db.WithContext(ctx).Table("job_types").Where("id = ?", id).Delete(&master.JobType{}).Error; err != nil {
		return fmt.Errorf("failed to delete job type: %w", err)
	}

	// Invalidate cache
	s.invalidateCache()

	return nil
}

// CheckDuplicateCode checks if a job type with the given code exists
func (s *adminJobTypeServiceImpl) CheckDuplicateCode(ctx context.Context, code string, excludeID *int64) (bool, error) {
	var count int64
	query := s.db.WithContext(ctx).Table("job_types").Where("code = ?", code)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CountJobReferences counts how many jobs reference this job type
func (s *adminJobTypeServiceImpl) CountJobReferences(ctx context.Context, id int64) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Table("jobs").
		Where("job_type_id = ?", id).
		Count(&count).Error
	return count, err
}

// invalidateCache invalidates all job type-related cache entries
func (s *adminJobTypeServiceImpl) invalidateCache() {
	// Job types are cached in job options cache
	s.cache.Delete("job_options:all")
}
