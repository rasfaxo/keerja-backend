package postgres

import (
	"context"

	"keerja-backend/internal/domain/master"

	"gorm.io/gorm"
)

type jobOptionsRepository struct {
	db *gorm.DB
}

// NewJobOptionsRepository creates a new JobOptions repository
func NewJobOptionsRepository(db *gorm.DB) master.JobOptionsRepository {
	return &jobOptionsRepository{db: db}
}

// GetAllJobTypes retrieves all job types
func (r *jobOptionsRepository) GetAllJobTypes(ctx context.Context) ([]master.JobType, error) {
	var jobTypes []master.JobType
	err := r.db.WithContext(ctx).
		Order("\"order\" ASC, name ASC").
		Find(&jobTypes).Error
	return jobTypes, err
}

// GetAllWorkPolicies retrieves all work policies
func (r *jobOptionsRepository) GetAllWorkPolicies(ctx context.Context) ([]master.WorkPolicy, error) {
	var workPolicies []master.WorkPolicy
	err := r.db.WithContext(ctx).
		Order("\"order\" ASC, name ASC").
		Find(&workPolicies).Error
	return workPolicies, err
}

// GetAllEducationLevels retrieves all education levels
func (r *jobOptionsRepository) GetAllEducationLevels(ctx context.Context) ([]master.EducationLevel, error) {
	var educationLevels []master.EducationLevel
	err := r.db.WithContext(ctx).
		Order("\"order\" ASC, name ASC").
		Find(&educationLevels).Error
	return educationLevels, err
}

// GetAllExperienceLevels retrieves all experience levels
func (r *jobOptionsRepository) GetAllExperienceLevels(ctx context.Context) ([]master.ExperienceLevel, error) {
	var experienceLevels []master.ExperienceLevel
	err := r.db.WithContext(ctx).
		Order("\"order\" ASC, name ASC").
		Find(&experienceLevels).Error
	return experienceLevels, err
}

// GetAllGenderPreferences retrieves all gender preferences
func (r *jobOptionsRepository) GetAllGenderPreferences(ctx context.Context) ([]master.GenderPreference, error) {
	var genderPreferences []master.GenderPreference
	err := r.db.WithContext(ctx).
		Order("\"order\" ASC, name ASC").
		Find(&genderPreferences).Error
	return genderPreferences, err
}

// GetJobOptions retrieves all job options in one query (for caching)
func (r *jobOptionsRepository) GetJobOptions(ctx context.Context) (*master.JobOptionsResponse, error) {
	response := &master.JobOptionsResponse{}

	// Get all data in parallel would be ideal, but for simplicity we'll do sequential
	jobTypes, err := r.GetAllJobTypes(ctx)
	if err != nil {
		return nil, err
	}
	response.JobTypes = jobTypes

	workPolicies, err := r.GetAllWorkPolicies(ctx)
	if err != nil {
		return nil, err
	}
	response.WorkPolicies = workPolicies

	educationLevels, err := r.GetAllEducationLevels(ctx)
	if err != nil {
		return nil, err
	}
	response.EducationLevels = educationLevels

	experienceLevels, err := r.GetAllExperienceLevels(ctx)
	if err != nil {
		return nil, err
	}
	response.ExperienceLevels = experienceLevels

	genderPreferences, err := r.GetAllGenderPreferences(ctx)
	if err != nil {
		return nil, err
	}
	response.GenderPreferences = genderPreferences

	return response, nil
}

// FindJobTypeByID finds job type by ID
func (r *jobOptionsRepository) FindJobTypeByID(ctx context.Context, id int64) (*master.JobType, error) {
	var jobType master.JobType
	err := r.db.WithContext(ctx).First(&jobType, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &jobType, nil
}

// FindWorkPolicyByID finds work policy by ID
func (r *jobOptionsRepository) FindWorkPolicyByID(ctx context.Context, id int64) (*master.WorkPolicy, error) {
	var workPolicy master.WorkPolicy
	err := r.db.WithContext(ctx).First(&workPolicy, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &workPolicy, nil
}

// FindEducationLevelByID finds education level by ID
func (r *jobOptionsRepository) FindEducationLevelByID(ctx context.Context, id int64) (*master.EducationLevel, error) {
	var educationLevel master.EducationLevel
	err := r.db.WithContext(ctx).First(&educationLevel, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &educationLevel, nil
}

// FindExperienceLevelByID finds experience level by ID
func (r *jobOptionsRepository) FindExperienceLevelByID(ctx context.Context, id int64) (*master.ExperienceLevel, error) {
	var experienceLevel master.ExperienceLevel
	err := r.db.WithContext(ctx).First(&experienceLevel, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &experienceLevel, nil
}

// FindGenderPreferenceByID finds gender preference by ID
func (r *jobOptionsRepository) FindGenderPreferenceByID(ctx context.Context, id int64) (*master.GenderPreference, error) {
	var genderPreference master.GenderPreference
	err := r.db.WithContext(ctx).First(&genderPreference, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &genderPreference, nil
}
