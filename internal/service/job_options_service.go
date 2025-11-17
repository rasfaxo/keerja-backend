package service

import (
	"context"
	"time"

	"keerja-backend/internal/cache"
	"keerja-backend/internal/domain/master"
)

const (
	jobOptionsCacheKey = "master:job_options"
	jobOptionsCacheTTL = 24 * time.Hour // Cache for 24 hours since this is static data
)

type jobOptionsService struct {
	repo  master.JobOptionsRepository
	cache cache.Cache
}

// NewJobOptionsService creates a new job options service
func NewJobOptionsService(repo master.JobOptionsRepository, cache cache.Cache) master.JobOptionsService {
	return &jobOptionsService{
		repo:  repo,
		cache: cache,
	}
}

// GetJobOptions retrieves all job options (heavily cached)
func (s *jobOptionsService) GetJobOptions(ctx context.Context) (*master.JobOptionsResponse, error) {
	// Try to get from cache first
	cachedData, ok := s.cache.Get(jobOptionsCacheKey)
	if ok {
		if response, ok := cachedData.(*master.JobOptionsResponse); ok {
			return response, nil
		}
	}

	// Cache miss or error - fetch from database
	response, err := s.repo.GetJobOptions(ctx)
	if err != nil {
		return nil, err
	}

	// Store in cache asynchronously to not block response
	go func() {
		s.cache.Set(jobOptionsCacheKey, response, jobOptionsCacheTTL)
	}()

	return response, nil
}

// GetJobTypes retrieves all job types
func (s *jobOptionsService) GetJobTypes(ctx context.Context) ([]master.JobType, error) {
	jobTypes, err := s.repo.GetAllJobTypes(ctx)
	if err != nil {
		return nil, err
	}
	return jobTypes, nil
}

// GetWorkPolicies retrieves all work policies
func (s *jobOptionsService) GetWorkPolicies(ctx context.Context) ([]master.WorkPolicy, error) {
	workPolicies, err := s.repo.GetAllWorkPolicies(ctx)
	if err != nil {
		return nil, err
	}
	return workPolicies, nil
}

// GetEducationLevels retrieves all education levels
func (s *jobOptionsService) GetEducationLevels(ctx context.Context) ([]master.EducationLevel, error) {
	educationLevels, err := s.repo.GetAllEducationLevels(ctx)
	if err != nil {
		return nil, err
	}
	return educationLevels, nil
}

// GetExperienceLevels retrieves all experience levels
func (s *jobOptionsService) GetExperienceLevels(ctx context.Context) ([]master.ExperienceLevel, error) {
	experienceLevels, err := s.repo.GetAllExperienceLevels(ctx)
	if err != nil {
		return nil, err
	}
	return experienceLevels, nil
}

// GetGenderPreferences retrieves all gender preferences
func (s *jobOptionsService) GetGenderPreferences(ctx context.Context) ([]master.GenderPreference, error) {
	genderPreferences, err := s.repo.GetAllGenderPreferences(ctx)
	if err != nil {
		return nil, err
	}
	return genderPreferences, nil
}

// InvalidateCache clears the job options cache (for admin use after updates)
func (s *jobOptionsService) InvalidateCache() error {
	s.cache.Delete(jobOptionsCacheKey)
	return nil
}
