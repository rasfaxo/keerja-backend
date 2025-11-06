package service

import (
	"context"
	"errors"
	"strings"

	"keerja-backend/internal/domain/master"
)

type jobTitleService struct {
	repo master.JobTitleRepository
}

// NewJobTitleService creates a new job title service
func NewJobTitleService(repo master.JobTitleRepository) master.JobTitleService {
	return &jobTitleService{
		repo: repo,
	}
}

// SearchJobTitles searches for job titles with query and limit
func (s *jobTitleService) SearchJobTitles(ctx context.Context, query string, limit int) ([]master.JobTitleResponse, error) {
	// Set default limit if not provided
	if limit <= 0 {
		limit = 20
	}

	// Cap limit at 100
	if limit > 100 {
		limit = 100
	}

	var jobTitles []master.JobTitle
	var err error

	if query == "" {
		// No search query - return popular titles
		jobTitles, err = s.repo.ListPopular(ctx, limit)
	} else {
		// Search with query
		jobTitles, err = s.repo.SearchJobTitles(ctx, strings.TrimSpace(query), limit)

		// Increment search count for all returned titles asynchronously
		if err == nil && len(jobTitles) > 0 {
			// Use goroutine to not block the response
			go func() {
				for _, jt := range jobTitles {
					_ = s.repo.IncrementSearchCount(context.Background(), jt.ID)
				}
			}()
		}
	}

	if err != nil {
		return nil, err
	}

	// Convert to response DTOs
	responses := make([]master.JobTitleResponse, len(jobTitles))
	for i, jt := range jobTitles {
		responses[i] = master.JobTitleResponse{
			ID:                    jt.ID,
			Name:                  jt.Name,
			RecommendedCategoryID: jt.RecommendedCategoryID,
			PopularityScore:       jt.PopularityScore,
			SearchCount:           jt.SearchCount,
		}
	}

	return responses, nil
}

// GetJobTitle retrieves a single job title by ID
func (s *jobTitleService) GetJobTitle(ctx context.Context, id int64) (*master.JobTitleResponse, error) {
	jobTitle, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if jobTitle == nil {
		return nil, errors.New("job title not found")
	}

	return &master.JobTitleResponse{
		ID:                    jobTitle.ID,
		Name:                  jobTitle.Name,
		RecommendedCategoryID: jobTitle.RecommendedCategoryID,
		PopularityScore:       jobTitle.PopularityScore,
		SearchCount:           jobTitle.SearchCount,
	}, nil
}

// ListPopularJobTitles retrieves popular job titles
func (s *jobTitleService) ListPopularJobTitles(ctx context.Context, limit int) ([]master.JobTitleResponse, error) {
	if limit <= 0 {
		limit = 10
	}

	jobTitles, err := s.repo.ListPopular(ctx, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]master.JobTitleResponse, len(jobTitles))
	for i, jt := range jobTitles {
		responses[i] = master.JobTitleResponse{
			ID:                    jt.ID,
			Name:                  jt.Name,
			RecommendedCategoryID: jt.RecommendedCategoryID,
			PopularityScore:       jt.PopularityScore,
			SearchCount:           jt.SearchCount,
		}
	}

	return responses, nil
}

// CreateJobTitle creates a new job title
func (s *jobTitleService) CreateJobTitle(ctx context.Context, req *master.CreateJobTitleRequest) (*master.JobTitle, error) {
	// Validate required fields
	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("name is required")
	}

	// Check if job title with same name already exists
	existing, err := s.repo.FindByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("job title with this name already exists")
	}

	// Create job title entity
	jobTitle := &master.JobTitle{
		Name:                  strings.TrimSpace(req.Name),
		RecommendedCategoryID: req.RecommendedCategoryID,
		PopularityScore:       req.PopularityScore,
		IsActive:              true,
		SearchCount:           0,
	}

	// Save to database
	if err := s.repo.Create(ctx, jobTitle); err != nil {
		return nil, err
	}

	return jobTitle, nil
}

// UpdateJobTitle updates an existing job title
func (s *jobTitleService) UpdateJobTitle(ctx context.Context, id int64, req *master.UpdateJobTitleRequest) (*master.JobTitle, error) {
	// Find existing job title
	jobTitle, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if jobTitle == nil {
		return nil, errors.New("job title not found")
	}

	// Update fields if provided
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, errors.New("name cannot be empty")
		}

		// Check if another job title with same name exists
		existing, err := s.repo.FindByName(ctx, name)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, errors.New("job title with this name already exists")
		}

		jobTitle.Name = name
	}

	if req.RecommendedCategoryID != nil {
		jobTitle.RecommendedCategoryID = req.RecommendedCategoryID
	}

	if req.IsActive != nil {
		jobTitle.IsActive = *req.IsActive
	}

	if req.PopularityScore != nil {
		if *req.PopularityScore < 0 || *req.PopularityScore > 100 {
			return nil, errors.New("popularity score must be between 0 and 100")
		}
		jobTitle.PopularityScore = *req.PopularityScore
	}

	// Save updates
	if err := s.repo.Update(ctx, jobTitle); err != nil {
		return nil, err
	}

	return jobTitle, nil
}

// DeleteJobTitle deletes a job title
func (s *jobTitleService) DeleteJobTitle(ctx context.Context, id int64) error {
	// Check if job title exists
	jobTitle, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if jobTitle == nil {
		return errors.New("job title not found")
	}

	// Delete the job title
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
