package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"keerja-backend/internal/domain/company"
	"keerja-backend/internal/domain/job"
	"keerja-backend/internal/domain/master"
	"keerja-backend/internal/domain/user"
	"keerja-backend/internal/utils"
)

// jobService implements job.JobService interface
type jobService struct {
	jobRepo         job.JobRepository
	companyRepo     company.CompanyRepository
	userRepo        user.UserRepository
	industryService master.IndustryService
	districtService master.DistrictService
}

// NewJobService creates a new job service instance
func NewJobService(
	jobRepo job.JobRepository,
	companyRepo company.CompanyRepository,
	userRepo user.UserRepository,
	industryService master.IndustryService,
	districtService master.DistrictService,
) job.JobService {
	return &jobService{
		jobRepo:         jobRepo,
		companyRepo:     companyRepo,
		userRepo:        userRepo,
		industryService: industryService,
		districtService: districtService,
	}
}

// ===== Job Management (Employer) =====

// CreateJob creates a new job posting
func (s *jobService) CreateJob(ctx context.Context, req *job.CreateJobRequest) (*job.Job, error) {
	// Verify company exists
	_, err := s.companyRepo.FindByID(ctx, req.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}

	// Verify employer user has permission
	if req.EmployerUserID != 0 {
		// Get employer user
		employerUser, err := s.companyRepo.FindEmployerUserByUserAndCompany(ctx, req.EmployerUserID, req.CompanyID)
		if err != nil || employerUser == nil {
			return nil, errors.New("employer user does not have permission to create jobs")
		}

		// Check if role has permission (recruiter and above)
		if employerUser.Role != "recruiter" && employerUser.Role != "admin" && employerUser.Role != "owner" {
			return nil, errors.New("employer user does not have permission to create jobs")
		}
	}

	// Generate unique slug from title
	slug := utils.GenerateSlug(req.Title)
	// If slug exists, add random suffix
	existingJob, _ := s.jobRepo.FindBySlug(ctx, slug)
	if existingJob != nil {
		slug = utils.GenerateSlugSimple(req.Title)
	}

	// Create job entity
	newJob := &job.Job{
		CompanyID:        req.CompanyID,
		EmployerUserID:   &req.EmployerUserID,
		CategoryID:       req.CategoryID,
		Title:            req.Title,
		Slug:             slug,
		JobLevel:         req.JobLevel,
		EmploymentType:   req.EmploymentType,
		Description:      req.Description,
		RequirementsText: req.RequirementsText,
		Responsibilities: req.Responsibilities,
		RemoteOption:     req.RemoteOption,
		SalaryMin:        req.SalaryMin,
		SalaryMax:        req.SalaryMax,
		Currency:         req.Currency,
		ExperienceMin:    req.ExperienceMin,
		ExperienceMax:    req.ExperienceMax,
		EducationLevel:   req.EducationLevel,
		TotalHires:       req.TotalHires,
		Status:           "draft", // Default status
		ExpiredAt:        req.ExpiredAt,
	}

	// Set default currency if not provided
	if newJob.Currency == "" {
		newJob.Currency = "IDR"
	}

	// Set default total hires if not provided
	if newJob.TotalHires == 0 {
		newJob.TotalHires = 1
	}

	// Validate job before creation
	if err := s.ValidateJob(ctx, newJob); err != nil {
		return nil, fmt.Errorf("job validation failed: %w", err)
	}

	// Create job
	if err := s.jobRepo.Create(ctx, newJob); err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	// Add locations if provided
	if len(req.Locations) > 0 {
		for _, locReq := range req.Locations {
			location := &job.JobLocation{
				JobID:         newJob.ID,
				CompanyID:     &req.CompanyID,
				LocationType:  locReq.LocationType,
				Address:       locReq.Address,
				City:          locReq.City,
				Province:      locReq.Province,
				PostalCode:    locReq.PostalCode,
				Country:       locReq.Country,
				Latitude:      locReq.Latitude,
				Longitude:     locReq.Longitude,
				GooglePlaceID: locReq.GooglePlaceID,
				MapURL:        locReq.MapURL,
				IsPrimary:     locReq.IsPrimary,
			}
			if location.Country == "" {
				location.Country = "Indonesia"
			}
			if err := s.jobRepo.CreateLocation(ctx, location); err != nil {
				return nil, fmt.Errorf("failed to create job location: %w", err)
			}
		}
	}

	// Add benefits if provided
	if len(req.Benefits) > 0 {
		if err := s.BulkAddBenefits(ctx, newJob.ID, req.Benefits); err != nil {
			return nil, fmt.Errorf("failed to add benefits: %w", err)
		}
	}

	// Add skills if provided
	if len(req.Skills) > 0 {
		if err := s.BulkAddSkills(ctx, newJob.ID, req.Skills); err != nil {
			return nil, fmt.Errorf("failed to add skills: %w", err)
		}
	}

	// Add requirements if provided
	if len(req.JobRequirements) > 0 {
		if err := s.BulkAddRequirements(ctx, newJob.ID, req.JobRequirements); err != nil {
			return nil, fmt.Errorf("failed to add requirements: %w", err)
		}
	}

	// Reload job with all relationships
	return s.jobRepo.FindByID(ctx, newJob.ID)
}

// UpdateJob updates an existing job
func (s *jobService) UpdateJob(ctx context.Context, jobID int64, req *job.UpdateJobRequest) (*job.Job, error) {
	// Find existing job
	existingJob, err := s.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("job not found: %w", err)
	}

	// Update fields if provided
	if req.CategoryID != nil {
		existingJob.CategoryID = req.CategoryID
	}
	if req.Title != "" {
		existingJob.Title = req.Title
		// Regenerate slug if title changed
		newSlug := utils.GenerateSlug(req.Title)
		// Check if slug exists
		slugJob, _ := s.jobRepo.FindBySlug(ctx, newSlug)
		if slugJob != nil && slugJob.ID != existingJob.ID {
			newSlug = utils.GenerateSlugSimple(req.Title)
		}
		existingJob.Slug = newSlug
	}
	if req.JobLevel != "" {
		existingJob.JobLevel = req.JobLevel
	}
	if req.EmploymentType != "" {
		existingJob.EmploymentType = req.EmploymentType
	}
	if req.Description != "" {
		existingJob.Description = req.Description
	}
	if req.RequirementsText != "" {
		existingJob.RequirementsText = req.RequirementsText
	}
	if req.Responsibilities != "" {
		existingJob.Responsibilities = req.Responsibilities
	}
	if req.RemoteOption != nil {
		existingJob.RemoteOption = *req.RemoteOption
	}
	if req.SalaryMin != nil {
		existingJob.SalaryMin = req.SalaryMin
	}
	if req.SalaryMax != nil {
		existingJob.SalaryMax = req.SalaryMax
	}
	if req.Currency != "" {
		existingJob.Currency = req.Currency
	}
	if req.ExperienceMin != nil {
		existingJob.ExperienceMin = req.ExperienceMin
	}
	if req.ExperienceMax != nil {
		existingJob.ExperienceMax = req.ExperienceMax
	}
	if req.EducationLevel != "" {
		existingJob.EducationLevel = req.EducationLevel
	}
	if req.TotalHires != nil {
		existingJob.TotalHires = *req.TotalHires
	}
	if req.ExpiredAt != nil {
		existingJob.ExpiredAt = req.ExpiredAt
	}

	// Validate updated job
	if err := s.ValidateJob(ctx, existingJob); err != nil {
		return nil, fmt.Errorf("job validation failed: %w", err)
	}

	// Update job
	if err := s.jobRepo.Update(ctx, existingJob); err != nil {
		return nil, fmt.Errorf("failed to update job: %w", err)
	}

	// Reload job with relationships
	return s.jobRepo.FindByID(ctx, jobID)
}

// DeleteJob deletes a job (soft delete)
func (s *jobService) DeleteJob(ctx context.Context, jobID int64, employerUserID int64) error {
	// Check ownership
	if err := s.CheckJobOwnership(ctx, jobID, employerUserID); err != nil {
		return err
	}

	// Soft delete job
	return s.jobRepo.SoftDelete(ctx, jobID)
}

// GetJob retrieves a job by ID
func (s *jobService) GetJob(ctx context.Context, jobID int64) (*job.Job, error) {
	return s.jobRepo.FindByID(ctx, jobID)
}

// GetJobBySlug retrieves a job by slug
func (s *jobService) GetJobBySlug(ctx context.Context, slug string) (*job.Job, error) {
	return s.jobRepo.FindBySlug(ctx, slug)
}

// GetJobByUUID retrieves a job by UUID
func (s *jobService) GetJobByUUID(ctx context.Context, uuidStr string) (*job.Job, error) {
	return s.jobRepo.FindByUUID(ctx, uuidStr)
}

// GetMyJobs retrieves jobs created by employer user
func (s *jobService) GetMyJobs(ctx context.Context, employerUserID int64, filter job.JobFilter, page, limit int) ([]job.Job, int64, error) {
	return s.jobRepo.ListByEmployer(ctx, employerUserID, filter, page, limit)
}

// GetCompanyJobs retrieves all jobs for a company
func (s *jobService) GetCompanyJobs(ctx context.Context, companyID int64, filter job.JobFilter, page, limit int) ([]job.Job, int64, error) {
	return s.jobRepo.ListByCompany(ctx, companyID, filter, page, limit)
}

// ===== Job Status Management =====

// PublishJob publishes a job (makes it visible to job seekers)
func (s *jobService) PublishJob(ctx context.Context, jobID int64, employerUserID int64) error {
	// Check ownership
	if err := s.CheckJobOwnership(ctx, jobID, employerUserID); err != nil {
		return err
	}

	// Get job to validate
	j, err := s.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf("job not found: %w", err)
	}

	// Validate job before publishing
	if err := s.ValidateJob(ctx, j); err != nil {
		return fmt.Errorf("cannot publish job: %w", err)
	}

	// Check if job is already published
	if j.Status == "published" {
		return errors.New("job is already published")
	}

	// Publish job
	if err := s.jobRepo.PublishJob(ctx, jobID); err != nil {
		return fmt.Errorf("failed to publish job: %w", err)
	}

	return nil
}

// UnpublishJob unpublishes a job (hides from job seekers)
func (s *jobService) UnpublishJob(ctx context.Context, jobID int64, employerUserID int64) error {
	// Check ownership
	if err := s.CheckJobOwnership(ctx, jobID, employerUserID); err != nil {
		return err
	}

	// Update status to draft
	return s.jobRepo.UpdateStatus(ctx, jobID, "draft")
}

// CloseJob closes a job (no longer accepting applications)
func (s *jobService) CloseJob(ctx context.Context, jobID int64, employerUserID int64) error {
	// Check ownership
	if err := s.CheckJobOwnership(ctx, jobID, employerUserID); err != nil {
		return err
	}

	// Close job
	return s.jobRepo.CloseJob(ctx, jobID)
}

// ReopenJob reopens a closed job
func (s *jobService) ReopenJob(ctx context.Context, jobID int64, employerUserID int64) error {
	// Check ownership
	if err := s.CheckJobOwnership(ctx, jobID, employerUserID); err != nil {
		return err
	}

	// Get job
	j, err := s.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf("job not found: %w", err)
	}

	// Check if job is closed
	if j.Status != "closed" && j.Status != "expired" {
		return errors.New("only closed or expired jobs can be reopened")
	}

	// Validate job before reopening
	if err := s.ValidateJob(ctx, j); err != nil {
		return fmt.Errorf("cannot reopen job: %w", err)
	}

	// Reopen job (set to published)
	return s.jobRepo.PublishJob(ctx, jobID)
}

// SuspendJob suspends a job (admin action)
func (s *jobService) SuspendJob(ctx context.Context, jobID int64, employerUserID int64, reason string) error {
	// Note: In production, this should check for admin privileges
	// For now, we'll allow employer to suspend their own jobs

	// Check ownership
	if err := s.CheckJobOwnership(ctx, jobID, employerUserID); err != nil {
		return err
	}

	// Suspend job
	return s.jobRepo.SuspendJob(ctx, jobID)
}

// SetJobExpiry sets job expiry date
func (s *jobService) SetJobExpiry(ctx context.Context, jobID int64, expiryDate time.Time) error {
	// Get job
	j, err := s.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf("job not found: %w", err)
	}

	// Update expiry date
	j.ExpiredAt = &expiryDate

	return s.jobRepo.Update(ctx, j)
}

// ExtendJobExpiry extends job expiry by specified days
func (s *jobService) ExtendJobExpiry(ctx context.Context, jobID int64, days int) error {
	// Get job
	j, err := s.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf("job not found: %w", err)
	}

	// Calculate new expiry date
	var newExpiry time.Time
	if j.ExpiredAt != nil {
		newExpiry = j.ExpiredAt.AddDate(0, 0, days)
	} else {
		newExpiry = time.Now().AddDate(0, 0, days)
	}

	// Update expiry date
	j.ExpiredAt = &newExpiry

	return s.jobRepo.Update(ctx, j)
}

// AutoExpireJobs automatically expires jobs past their expiry date (cron job)
func (s *jobService) AutoExpireJobs(ctx context.Context) error {
	// Get expired jobs
	expiredJobs, err := s.jobRepo.GetExpiredJobs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get expired jobs: %w", err)
	}

	// Expire each job
	for _, j := range expiredJobs {
		if err := s.jobRepo.ExpireJob(ctx, j.ID); err != nil {
			// Log error but continue with other jobs
			fmt.Printf("failed to expire job %d: %v\n", j.ID, err)
		}
	}

	return nil
}

// ===== Job Search and Discovery (Public) =====

// ListJobs lists jobs with filters
func (s *jobService) ListJobs(ctx context.Context, filter job.JobFilter, page, limit int) ([]job.Job, int64, error) {
	// Set default filter for public listing (only show published jobs)
	if filter.Status == "" {
		filter.Status = "published"
	}

	return s.jobRepo.List(ctx, filter, page, limit)
}

// SearchJobs performs advanced job search
func (s *jobService) SearchJobs(ctx context.Context, filter job.JobSearchFilter, page, limit int) (*job.JobSearchResponse, error) {
	// Perform search
	jobs, total, err := s.jobRepo.SearchJobs(ctx, filter, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search jobs: %w", err)
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// Build response
	response := &job.JobSearchResponse{
		Jobs:       jobs,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	// TODO: Add facets and suggestions (requires additional implementation)

	return response, nil
}

// SearchJobsByLocation searches jobs by geographic location
func (s *jobService) SearchJobsByLocation(ctx context.Context, latitude, longitude, radius float64, filter job.JobFilter, page, limit int) ([]job.Job, int64, error) {
	return s.jobRepo.SearchByLocation(ctx, latitude, longitude, radius, filter, page, limit)
}

// GetFeaturedJobs retrieves featured jobs
func (s *jobService) GetFeaturedJobs(ctx context.Context, limit int) ([]job.Job, error) {
	// Featured jobs logic: published jobs from verified companies, sorted by views
	filter := job.JobFilter{
		Status: "published",
		SortBy: "views",
	}

	jobs, _, err := s.jobRepo.List(ctx, filter, 1, limit)
	return jobs, err
}

// GetLatestJobs retrieves latest published jobs
func (s *jobService) GetLatestJobs(ctx context.Context, limit int) ([]job.Job, error) {
	filter := job.JobFilter{
		Status: "published",
		SortBy: "latest",
	}

	jobs, _, err := s.jobRepo.List(ctx, filter, 1, limit)
	return jobs, err
}

// GetTrendingJobs retrieves trending jobs (most viewed recently)
func (s *jobService) GetTrendingJobs(ctx context.Context, limit int) ([]job.Job, error) {
	return s.jobRepo.GetTrendingJobs(ctx, limit)
}

// GetRecommendedJobs retrieves recommended jobs for a user
func (s *jobService) GetRecommendedJobs(ctx context.Context, userID int64, limit int) ([]job.Job, error) {
	return s.jobRepo.GetRecommendedJobs(ctx, userID, limit)
}

// GetSimilarJobs retrieves jobs similar to a given job
func (s *jobService) GetSimilarJobs(ctx context.Context, jobID int64, limit int) ([]job.Job, error) {
	return s.jobRepo.GetSimilarJobs(ctx, jobID, limit)
}

// ===== Job Matching =====

// CalculateMatchScore calculates match score between a job and user
func (s *jobService) CalculateMatchScore(ctx context.Context, jobID, userID int64) (*job.MatchScore, error) {
	// Get job with skills and requirements
	j, err := s.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("job not found: %w", err)
	}

	// Get user profile with skills, education, and experience
	userProfile, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Initialize match score
	matchScore := &job.MatchScore{
		JobID:  jobID,
		UserID: userID,
	}

	// Calculate skill score
	skillScore, matchedSkills, missingSkills := s.calculateSkillScore(j, userProfile)
	matchScore.SkillScore = skillScore
	matchScore.MatchedSkills = matchedSkills
	matchScore.MissingSkills = missingSkills

	// Calculate experience score
	matchScore.ExperienceScore = s.calculateExperienceScore(j, userProfile)

	// Calculate education score
	matchScore.EducationScore = s.calculateEducationScore(j, userProfile)

	// Calculate location score
	matchScore.LocationScore = s.calculateLocationScore(j, userProfile)

	// Calculate overall score (weighted average)
	matchScore.OverallScore = (skillScore*0.4 + matchScore.ExperienceScore*0.3 +
		matchScore.EducationScore*0.2 + matchScore.LocationScore*0.1)

	// Generate recommendation
	matchScore.Recommendation = s.generateRecommendation(matchScore)

	return matchScore, nil
}

// calculateSkillScore calculates skill match score
func (s *jobService) calculateSkillScore(j *job.Job, userProfile *user.User) (float64, []string, []string) {
	if len(j.Skills) == 0 {
		return 1.0, []string{}, []string{} // No skills required, perfect match
	}

	// Create map of user skills (by skill name since UserSkill doesn't have SkillID)
	userSkills := make(map[string]bool)
	for _, skill := range userProfile.Skills {
		userSkills[strings.ToLower(skill.SkillName)] = true
	}

	// Calculate matched and missing skills
	var matchedSkills, missingSkills []string
	var totalWeight, matchedWeight float64

	for _, jobSkill := range j.Skills {
		totalWeight += jobSkill.Weight

		// For now, we match by skill ID
		// TODO: In production, implement proper skill matching logic
		skillName := fmt.Sprintf("Skill-%d", jobSkill.SkillID)

		// Check if user has this skill (simplified matching)
		if len(userSkills) > 0 {
			matchedWeight += jobSkill.Weight
			matchedSkills = append(matchedSkills, skillName)
		} else if jobSkill.ImportanceLevel == "required" {
			missingSkills = append(missingSkills, skillName)
		}
	}

	// Calculate score
	score := 0.0
	if totalWeight > 0 {
		score = matchedWeight / totalWeight
	}

	return score, matchedSkills, missingSkills
}

// calculateExperienceScore calculates experience match score
func (s *jobService) calculateExperienceScore(j *job.Job, userProfile *user.User) float64 {
	// Calculate total user experience in years
	totalExperience := 0.0
	for _, exp := range userProfile.Experiences {
		endDate := exp.EndDate
		if endDate == nil {
			now := time.Now()
			endDate = &now
		}
		years := endDate.Sub(exp.StartDate).Hours() / (24 * 365)
		totalExperience += years
	}

	// If no experience requirement, perfect match
	if j.ExperienceMin == nil && j.ExperienceMax == nil {
		return 1.0
	}

	// Check if user meets minimum experience
	minExp := 0
	if j.ExperienceMin != nil {
		minExp = int(*j.ExperienceMin)
	}

	maxExp := 99999
	if j.ExperienceMax != nil {
		maxExp = int(*j.ExperienceMax)
	}

	// Calculate score based on experience range
	if totalExperience < float64(minExp) {
		// Below minimum: score based on how close to minimum
		return math.Max(0, totalExperience/float64(minExp))
	} else if totalExperience > float64(maxExp) {
		// Above maximum: slightly penalized but still high score
		return 0.9
	} else {
		// Within range: perfect score
		return 1.0
	}
}

// calculateEducationScore calculates education match score
func (s *jobService) calculateEducationScore(j *job.Job, userProfile *user.User) float64 {
	// If no education requirement, perfect match
	if j.EducationLevel == "" {
		return 1.0
	}

	// Education level hierarchy
	educationLevels := map[string]int{
		"High School": 1,
		"Associate":   2,
		"Bachelor":    3,
		"Master":      4,
		"Doctorate":   5,
	}

	requiredLevel := educationLevels[j.EducationLevel]

	// Find user's highest education level
	highestLevel := 0
	for _, edu := range userProfile.Educations {
		// Map Indonesian degree levels to international levels
		degreeLevel := ""
		if edu.DegreeLevel != nil {
			switch *edu.DegreeLevel {
			case "SMA":
				degreeLevel = "High School"
			case "D1", "D2", "D3":
				degreeLevel = "Associate"
			case "S1":
				degreeLevel = "Bachelor"
			case "S2":
				degreeLevel = "Master"
			case "S3":
				degreeLevel = "Doctorate"
			}
		}

		if level, ok := educationLevels[degreeLevel]; ok {
			if level > highestLevel {
				highestLevel = level
			}
		}
	}

	// Calculate score
	if highestLevel == 0 {
		return 0.0 // No education info
	} else if highestLevel >= requiredLevel {
		return 1.0 // Meets or exceeds requirement
	} else {
		// Below requirement: score based on how close
		return float64(highestLevel) / float64(requiredLevel)
	}
}

// calculateLocationScore calculates location match score
func (s *jobService) calculateLocationScore(j *job.Job, userProfile *user.User) float64 {
	// If job is remote, perfect match
	if j.RemoteOption {
		return 1.0
	}

	// If no user location info, neutral score
	if userProfile.Profile == nil || userProfile.Profile.LocationCity == nil {
		return 0.5
	}

	// Simple location matching (can be enhanced with geocoding)
	userLocation := strings.ToLower(*userProfile.Profile.LocationCity)
	jobLocation := strings.ToLower(j.Location)

	// Check if locations match
	if strings.Contains(userLocation, jobLocation) || strings.Contains(jobLocation, userLocation) {
		return 1.0
	}

	// Check city match
	if j.City != "" && strings.Contains(userLocation, strings.ToLower(j.City)) {
		return 0.8
	}

	// Check province match
	if j.Province != "" && strings.Contains(userLocation, strings.ToLower(j.Province)) {
		return 0.6
	}

	// No match
	return 0.3
}

// generateRecommendation generates recommendation text based on match score
func (s *jobService) generateRecommendation(score *job.MatchScore) string {
	if score.OverallScore >= 0.8 {
		return "Excellent match! You meet most of the requirements for this position."
	} else if score.OverallScore >= 0.6 {
		return "Good match. You have many of the skills and qualifications needed."
	} else if score.OverallScore >= 0.4 {
		return "Fair match. Consider highlighting your relevant experience in your application."
	} else {
		return "This position may be challenging, but don't let that stop you from applying if you're interested."
	}
}

// GetMatchingJobs retrieves jobs matching user profile
func (s *jobService) GetMatchingJobs(ctx context.Context, userID int64, filter job.JobFilter, page, limit int) (*job.MatchResponse, error) {
	// Get matching jobs from repository
	jobs, total, err := s.jobRepo.GetMatchingJobs(ctx, userID, filter, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get matching jobs: %w", err)
	}

	// Calculate match scores for each job
	jobsWithScores := make([]job.JobWithScore, 0, len(jobs))
	for _, j := range jobs {
		matchScore, err := s.CalculateMatchScore(ctx, j.ID, userID)
		if err != nil {
			// Log error but continue with other jobs
			fmt.Printf("failed to calculate match score for job %d: %v\n", j.ID, err)
			continue
		}

		jobsWithScores = append(jobsWithScores, job.JobWithScore{
			Job:        j,
			MatchScore: matchScore.OverallScore,
		})
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &job.MatchResponse{
		Jobs:       jobsWithScores,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// ===== Job Views and Interactions =====

// IncrementView increments job view count
func (s *jobService) IncrementView(ctx context.Context, jobID int64, userID *int64) error {
	// TODO: Implement view tracking with user ID to prevent duplicate counts
	// For now, just increment the counter
	return s.jobRepo.IncrementViews(ctx, jobID)
}

// GetJobStats retrieves job statistics
func (s *jobService) GetJobStats(ctx context.Context, jobID int64) (*job.JobStats, error) {
	return s.jobRepo.GetJobStats(ctx, jobID)
}

// GetCompanyJobStats retrieves company job statistics
func (s *jobService) GetCompanyJobStats(ctx context.Context, companyID int64) (*job.CompanyJobStats, error) {
	return s.jobRepo.GetCompanyJobStats(ctx, companyID)
}

// ===== Job Details Management =====

// AddLocation adds a location to a job
func (s *jobService) AddLocation(ctx context.Context, jobID int64, req *job.AddLocationRequest) (*job.JobLocation, error) {
	// Verify job exists
	_, err := s.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("job not found: %w", err)
	}

	// Create location
	location := &job.JobLocation{
		JobID:         jobID,
		LocationType:  req.LocationType,
		Address:       req.Address,
		City:          req.City,
		Province:      req.Province,
		PostalCode:    req.PostalCode,
		Country:       req.Country,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		GooglePlaceID: req.GooglePlaceID,
		MapURL:        req.MapURL,
		IsPrimary:     req.IsPrimary,
	}

	// Set default country
	if location.Country == "" {
		location.Country = "Indonesia"
	}

	// Set default location type
	if location.LocationType == "" {
		location.LocationType = "onsite"
	}

	if err := s.jobRepo.CreateLocation(ctx, location); err != nil {
		return nil, fmt.Errorf("failed to create location: %w", err)
	}

	return location, nil
}

// UpdateLocation updates a job location
func (s *jobService) UpdateLocation(ctx context.Context, locationID int64, req *job.UpdateLocationRequest) (*job.JobLocation, error) {
	// Find existing location
	location, err := s.jobRepo.FindLocationByID(ctx, locationID)
	if err != nil {
		return nil, fmt.Errorf("location not found: %w", err)
	}

	// Update fields if provided
	if req.LocationType != "" {
		location.LocationType = req.LocationType
	}
	if req.Address != "" {
		location.Address = req.Address
	}
	if req.City != "" {
		location.City = req.City
	}
	if req.Province != "" {
		location.Province = req.Province
	}
	if req.PostalCode != "" {
		location.PostalCode = req.PostalCode
	}
	if req.Country != "" {
		location.Country = req.Country
	}
	if req.Latitude != nil {
		location.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		location.Longitude = req.Longitude
	}
	if req.GooglePlaceID != "" {
		location.GooglePlaceID = req.GooglePlaceID
	}
	if req.MapURL != "" {
		location.MapURL = req.MapURL
	}
	if req.IsPrimary != nil {
		location.IsPrimary = *req.IsPrimary
	}

	if err := s.jobRepo.UpdateLocation(ctx, location); err != nil {
		return nil, fmt.Errorf("failed to update location: %w", err)
	}

	return location, nil
}

// DeleteLocation deletes a job location
func (s *jobService) DeleteLocation(ctx context.Context, locationID int64) error {
	return s.jobRepo.DeleteLocation(ctx, locationID)
}

// SetPrimaryLocation sets a location as primary for a job
func (s *jobService) SetPrimaryLocation(ctx context.Context, jobID, locationID int64) error {
	return s.jobRepo.SetPrimaryLocation(ctx, jobID, locationID)
}

// AddBenefit adds a benefit to a job
func (s *jobService) AddBenefit(ctx context.Context, jobID int64, req *job.AddBenefitRequest) (*job.JobBenefit, error) {
	// Verify job exists
	_, err := s.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("job not found: %w", err)
	}

	// Create benefit
	benefit := &job.JobBenefit{
		JobID:       jobID,
		BenefitID:   req.BenefitID,
		BenefitName: req.BenefitName,
		Description: req.Description,
		IsHighlight: req.IsHighlight,
	}

	if err := s.jobRepo.CreateBenefit(ctx, benefit); err != nil {
		return nil, fmt.Errorf("failed to create benefit: %w", err)
	}

	return benefit, nil
}

// UpdateBenefit updates a job benefit
func (s *jobService) UpdateBenefit(ctx context.Context, benefitID int64, req *job.UpdateBenefitRequest) (*job.JobBenefit, error) {
	// Find existing benefit
	benefit, err := s.jobRepo.FindBenefitByID(ctx, benefitID)
	if err != nil {
		return nil, fmt.Errorf("benefit not found: %w", err)
	}

	// Update fields if provided
	if req.BenefitName != "" {
		benefit.BenefitName = req.BenefitName
	}
	if req.Description != "" {
		benefit.Description = req.Description
	}
	if req.IsHighlight != nil {
		benefit.IsHighlight = *req.IsHighlight
	}

	if err := s.jobRepo.UpdateBenefit(ctx, benefit); err != nil {
		return nil, fmt.Errorf("failed to update benefit: %w", err)
	}

	return benefit, nil
}

// DeleteBenefit deletes a job benefit
func (s *jobService) DeleteBenefit(ctx context.Context, benefitID int64) error {
	return s.jobRepo.DeleteBenefit(ctx, benefitID)
}

// BulkAddBenefits adds multiple benefits to a job
func (s *jobService) BulkAddBenefits(ctx context.Context, jobID int64, benefits []job.AddBenefitRequest) error {
	// Verify job exists
	_, err := s.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf("job not found: %w", err)
	}

	// Create benefits
	jobBenefits := make([]job.JobBenefit, 0, len(benefits))
	for _, req := range benefits {
		jobBenefits = append(jobBenefits, job.JobBenefit{
			JobID:       jobID,
			BenefitID:   req.BenefitID,
			BenefitName: req.BenefitName,
			Description: req.Description,
			IsHighlight: req.IsHighlight,
		})
	}

	return s.jobRepo.BulkCreateBenefits(ctx, jobBenefits)
}

// AddSkill adds a skill requirement to a job
func (s *jobService) AddSkill(ctx context.Context, jobID int64, req *job.AddSkillRequest) (*job.JobSkill, error) {
	// Verify job exists
	_, err := s.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("job not found: %w", err)
	}

	// Create job skill
	jobSkill := &job.JobSkill{
		JobID:           jobID,
		SkillID:         req.SkillID,
		ImportanceLevel: req.ImportanceLevel,
		Weight:          req.Weight,
	}

	// Set defaults
	if jobSkill.ImportanceLevel == "" {
		jobSkill.ImportanceLevel = "required"
	}
	if jobSkill.Weight == 0 {
		jobSkill.Weight = 1.0
	}

	if err := s.jobRepo.CreateSkill(ctx, jobSkill); err != nil {
		return nil, fmt.Errorf("failed to create job skill: %w", err)
	}

	return jobSkill, nil
}

// UpdateSkill updates a job skill requirement
func (s *jobService) UpdateSkill(ctx context.Context, jobSkillID int64, req *job.UpdateSkillRequest) (*job.JobSkill, error) {
	// Find existing job skill
	jobSkill, err := s.jobRepo.FindSkillByID(ctx, jobSkillID)
	if err != nil {
		return nil, fmt.Errorf("job skill not found: %w", err)
	}

	// Update fields if provided
	if req.ImportanceLevel != "" {
		jobSkill.ImportanceLevel = req.ImportanceLevel
	}
	if req.Weight != nil {
		jobSkill.Weight = *req.Weight
	}

	if err := s.jobRepo.UpdateSkill(ctx, jobSkill); err != nil {
		return nil, fmt.Errorf("failed to update job skill: %w", err)
	}

	return jobSkill, nil
}

// DeleteSkill deletes a job skill requirement
func (s *jobService) DeleteSkill(ctx context.Context, jobSkillID int64) error {
	return s.jobRepo.DeleteSkill(ctx, jobSkillID)
}

// BulkAddSkills adds multiple skills to a job
func (s *jobService) BulkAddSkills(ctx context.Context, jobID int64, skills []job.AddSkillRequest) error {
	// Verify job exists
	_, err := s.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf("job not found: %w", err)
	}

	// Create job skills
	jobSkills := make([]job.JobSkill, 0, len(skills))
	for _, req := range skills {
		importanceLevel := req.ImportanceLevel
		if importanceLevel == "" {
			importanceLevel = "required"
		}

		weight := req.Weight
		if weight == 0 {
			weight = 1.0
		}

		jobSkills = append(jobSkills, job.JobSkill{
			JobID:           jobID,
			SkillID:         req.SkillID,
			ImportanceLevel: importanceLevel,
			Weight:          weight,
		})
	}

	return s.jobRepo.BulkCreateSkills(ctx, jobSkills)
}

// AddRequirement adds a requirement to a job
func (s *jobService) AddRequirement(ctx context.Context, jobID int64, req *job.AddRequirementRequest) (*job.JobRequirement, error) {
	// Verify job exists
	_, err := s.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("job not found: %w", err)
	}

	// Create requirement
	requirement := &job.JobRequirement{
		JobID:           jobID,
		RequirementType: req.RequirementType,
		RequirementText: req.RequirementText,
		SkillID:         req.SkillID,
		MinExperience:   req.MinExperience,
		MaxExperience:   req.MaxExperience,
		EducationLevel:  req.EducationLevel,
		Language:        req.Language,
		IsMandatory:     req.IsMandatory,
		Priority:        req.Priority,
	}

	// Set defaults
	if requirement.RequirementType == "" {
		requirement.RequirementType = "other"
	}
	if requirement.Priority == 0 {
		requirement.Priority = 1
	}

	if err := s.jobRepo.CreateRequirement(ctx, requirement); err != nil {
		return nil, fmt.Errorf("failed to create requirement: %w", err)
	}

	return requirement, nil
}

// UpdateRequirement updates a job requirement
func (s *jobService) UpdateRequirement(ctx context.Context, requirementID int64, req *job.UpdateRequirementRequest) (*job.JobRequirement, error) {
	// Find existing requirement
	requirement, err := s.jobRepo.FindRequirementByID(ctx, requirementID)
	if err != nil {
		return nil, fmt.Errorf("requirement not found: %w", err)
	}

	// Update fields if provided
	if req.RequirementType != "" {
		requirement.RequirementType = req.RequirementType
	}
	if req.RequirementText != "" {
		requirement.RequirementText = req.RequirementText
	}
	if req.SkillID != nil {
		requirement.SkillID = req.SkillID
	}
	if req.MinExperience != nil {
		requirement.MinExperience = req.MinExperience
	}
	if req.MaxExperience != nil {
		requirement.MaxExperience = req.MaxExperience
	}
	if req.EducationLevel != "" {
		requirement.EducationLevel = req.EducationLevel
	}
	if req.Language != "" {
		requirement.Language = req.Language
	}
	if req.IsMandatory != nil {
		requirement.IsMandatory = *req.IsMandatory
	}
	if req.Priority != nil {
		requirement.Priority = *req.Priority
	}

	if err := s.jobRepo.UpdateRequirement(ctx, requirement); err != nil {
		return nil, fmt.Errorf("failed to update requirement: %w", err)
	}

	return requirement, nil
}

// DeleteRequirement deletes a job requirement
func (s *jobService) DeleteRequirement(ctx context.Context, requirementID int64) error {
	return s.jobRepo.DeleteRequirement(ctx, requirementID)
}

// BulkAddRequirements adds multiple requirements to a job
func (s *jobService) BulkAddRequirements(ctx context.Context, jobID int64, requirements []job.AddRequirementRequest) error {
	// Verify job exists
	_, err := s.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf("job not found: %w", err)
	}

	// Create requirements
	jobRequirements := make([]job.JobRequirement, 0, len(requirements))
	for _, req := range requirements {
		requirementType := req.RequirementType
		if requirementType == "" {
			requirementType = "other"
		}

		priority := req.Priority
		if priority == 0 {
			priority = 1
		}

		jobRequirements = append(jobRequirements, job.JobRequirement{
			JobID:           jobID,
			RequirementType: requirementType,
			RequirementText: req.RequirementText,
			SkillID:         req.SkillID,
			MinExperience:   req.MinExperience,
			MaxExperience:   req.MaxExperience,
			EducationLevel:  req.EducationLevel,
			Language:        req.Language,
			IsMandatory:     req.IsMandatory,
			Priority:        priority,
		})
	}

	return s.jobRepo.BulkCreateRequirements(ctx, jobRequirements)
}

// ===== Category Management (Admin) =====

// CreateCategory creates a new job category
func (s *jobService) CreateCategory(ctx context.Context, req *job.CreateCategoryRequest) (*job.JobCategory, error) {
	// Create category
	category := &job.JobCategory{
		ParentID:    req.ParentID,
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	if err := s.jobRepo.CreateCategory(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

// UpdateCategory updates a job category
func (s *jobService) UpdateCategory(ctx context.Context, categoryID int64, req *job.UpdateCategoryRequest) (*job.JobCategory, error) {
	// Find existing category
	category, err := s.jobRepo.FindCategoryByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	// Update fields if provided
	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}
	if req.Code != "" {
		category.Code = req.Code
	}
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := s.jobRepo.UpdateCategory(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return category, nil
}

// DeleteCategory deletes a job category
func (s *jobService) DeleteCategory(ctx context.Context, categoryID int64) error {
	return s.jobRepo.DeleteCategory(ctx, categoryID)
}

// GetCategory retrieves a job category by ID
func (s *jobService) GetCategory(ctx context.Context, categoryID int64) (*job.JobCategory, error) {
	return s.jobRepo.FindCategoryByID(ctx, categoryID)
}

// GetCategoryByCode retrieves a job category by code
func (s *jobService) GetCategoryByCode(ctx context.Context, code string) (*job.JobCategory, error) {
	return s.jobRepo.FindCategoryByCode(ctx, code)
}

// ListCategories lists job categories with filters
func (s *jobService) ListCategories(ctx context.Context, filter job.CategoryFilter, page, limit int) ([]job.JobCategory, int64, error) {
	return s.jobRepo.ListCategories(ctx, filter, page, limit)
}

// GetCategoryTree retrieves hierarchical category tree
func (s *jobService) GetCategoryTree(ctx context.Context) ([]job.JobCategory, error) {
	return s.jobRepo.GetCategoryTree(ctx)
}

// GetActiveCategories retrieves all active categories
func (s *jobService) GetActiveCategories(ctx context.Context) ([]job.JobCategory, error) {
	return s.jobRepo.GetActiveCategories(ctx)
}

// ===== Subcategory Management (Admin) =====

// CreateSubcategory creates a new job subcategory
func (s *jobService) CreateSubcategory(ctx context.Context, req *job.CreateSubcategoryRequest) (*job.JobSubcategory, error) {
	// Verify parent category exists
	_, err := s.jobRepo.FindCategoryByID(ctx, req.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("parent category not found: %w", err)
	}

	// Create subcategory
	subcategory := &job.JobSubcategory{
		CategoryID:  req.CategoryID,
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	if err := s.jobRepo.CreateSubcategory(ctx, subcategory); err != nil {
		return nil, fmt.Errorf("failed to create subcategory: %w", err)
	}

	return subcategory, nil
}

// UpdateSubcategory updates a job subcategory
func (s *jobService) UpdateSubcategory(ctx context.Context, subcategoryID int64, req *job.UpdateSubcategoryRequest) (*job.JobSubcategory, error) {
	// Find existing subcategory
	subcategory, err := s.jobRepo.FindSubcategoryByID(ctx, subcategoryID)
	if err != nil {
		return nil, fmt.Errorf("subcategory not found: %w", err)
	}

	// Update fields if provided
	if req.CategoryID != nil {
		// Verify new parent category exists
		_, err := s.jobRepo.FindCategoryByID(ctx, *req.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("parent category not found: %w", err)
		}
		subcategory.CategoryID = *req.CategoryID
	}
	if req.Code != "" {
		subcategory.Code = req.Code
	}
	if req.Name != "" {
		subcategory.Name = req.Name
	}
	if req.Description != "" {
		subcategory.Description = req.Description
	}
	if req.IsActive != nil {
		subcategory.IsActive = *req.IsActive
	}

	if err := s.jobRepo.UpdateSubcategory(ctx, subcategory); err != nil {
		return nil, fmt.Errorf("failed to update subcategory: %w", err)
	}

	return subcategory, nil
}

// DeleteSubcategory deletes a job subcategory
func (s *jobService) DeleteSubcategory(ctx context.Context, subcategoryID int64) error {
	return s.jobRepo.DeleteSubcategory(ctx, subcategoryID)
}

// GetSubcategory retrieves a job subcategory by ID
func (s *jobService) GetSubcategory(ctx context.Context, subcategoryID int64) (*job.JobSubcategory, error) {
	return s.jobRepo.FindSubcategoryByID(ctx, subcategoryID)
}

// ListSubcategories lists subcategories for a category
func (s *jobService) ListSubcategories(ctx context.Context, categoryID int64) ([]job.JobSubcategory, error) {
	return s.jobRepo.ListSubcategories(ctx, categoryID)
}

// GetActiveSubcategories retrieves active subcategories for a category
func (s *jobService) GetActiveSubcategories(ctx context.Context, categoryID int64) ([]job.JobSubcategory, error) {
	return s.jobRepo.GetActiveSubcategories(ctx, categoryID)
}

// ===== Analytics and Reporting =====

// GetJobAnalytics retrieves job analytics data
func (s *jobService) GetJobAnalytics(ctx context.Context, jobID int64, startDate, endDate time.Time) (*job.JobAnalytics, error) {
	// Get job stats
	stats, err := s.jobRepo.GetJobStats(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("failed to get job stats: %w", err)
	}

	// Build analytics response
	analytics := &job.JobAnalytics{
		JobID:             jobID,
		Period:            fmt.Sprintf("%s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")),
		TotalViews:        stats.ViewsCount,
		TotalApplications: stats.ApplicationsCount,
		ConversionRate:    stats.ConversionRate,
		// TODO: Implement time series data, unique viewers, and top sources
		ViewsData:        []job.TimeSeriesData{},
		ApplicationsData: []job.TimeSeriesData{},
		TopSources:       []job.SourceStats{},
	}

	return analytics, nil
}

// GetCompanyAnalytics retrieves company job analytics
func (s *jobService) GetCompanyAnalytics(ctx context.Context, companyID int64, startDate, endDate time.Time) (*job.CompanyAnalytics, error) {
	// Get company job stats
	stats, err := s.jobRepo.GetCompanyJobStats(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get company job stats: %w", err)
	}

	// Build analytics response
	analytics := &job.CompanyAnalytics{
		CompanyID:         companyID,
		Period:            fmt.Sprintf("%s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")),
		TotalJobs:         stats.TotalJobs,
		ActiveJobs:        stats.ActiveJobs,
		TotalViews:        stats.TotalViews,
		TotalApplications: stats.TotalApplications,
		// TODO: Implement top jobs and category breakdown
		TopJobs:           []job.JobPerformance{},
		CategoryBreakdown: []job.CategoryStats{},
	}

	return analytics, nil
}

// GetCategoryAnalytics retrieves category analytics
func (s *jobService) GetCategoryAnalytics(ctx context.Context, categoryID int64, startDate, endDate time.Time) (*job.CategoryAnalytics, error) {
	// Get category
	category, err := s.jobRepo.FindCategoryByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	// Get jobs in date range
	filter := job.JobFilter{
		CategoryID: categoryID,
	}
	jobs, err := s.jobRepo.GetJobsByDateRange(ctx, startDate, endDate, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs: %w", err)
	}

	// Calculate stats
	totalJobs := int64(len(jobs))
	activeJobs := int64(0)
	totalViews := int64(0)
	totalApplications := int64(0)

	for _, j := range jobs {
		if j.IsActive() {
			activeJobs++
		}
		totalViews += j.ViewsCount
		totalApplications += j.ApplicationsCount
	}

	// Build analytics response
	analytics := &job.CategoryAnalytics{
		CategoryID:        categoryID,
		CategoryName:      category.Name,
		Period:            fmt.Sprintf("%s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")),
		TotalJobs:         totalJobs,
		ActiveJobs:        activeJobs,
		TotalViews:        totalViews,
		TotalApplications: totalApplications,
		// TODO: Implement top companies
		TopCompanies: []job.CompanyStats{},
	}

	return analytics, nil
}

// GetPopularCategories retrieves popular categories
func (s *jobService) GetPopularCategories(ctx context.Context, limit int) ([]job.CategoryStats, error) {
	return s.jobRepo.GetPopularCategories(ctx, limit)
}

// GetTopCompanies retrieves top companies by job performance
func (s *jobService) GetTopCompanies(ctx context.Context, limit int) ([]job.CompanyStats, error) {
	// Get all companies
	companies, _, err := s.companyRepo.List(ctx, &company.CompanyFilter{})
	if err != nil {
		return nil, fmt.Errorf("failed to get companies: %w", err)
	}

	// Limit the number of companies to process
	if len(companies) > limit {
		companies = companies[:limit]
	}

	// Build stats for each company
	companyStats := make([]job.CompanyStats, 0, len(companies))
	for _, comp := range companies {
		stats, err := s.jobRepo.GetCompanyJobStats(ctx, comp.ID)
		if err != nil {
			continue
		}

		companyStats = append(companyStats, job.CompanyStats{
			CompanyID:         comp.ID,
			CompanyName:       comp.CompanyName,
			TotalJobs:         stats.TotalJobs,
			ActiveJobs:        stats.ActiveJobs,
			TotalViews:        stats.TotalViews,
			TotalApplications: stats.TotalApplications,
		})
	}

	return companyStats, nil
}

// ===== Bulk Operations =====

// BulkPublishJobs publishes multiple jobs
func (s *jobService) BulkPublishJobs(ctx context.Context, jobIDs []int64) error {
	for _, jobID := range jobIDs {
		if err := s.jobRepo.PublishJob(ctx, jobID); err != nil {
			return fmt.Errorf("failed to publish job %d: %w", jobID, err)
		}
	}
	return nil
}

// BulkCloseJobs closes multiple jobs
func (s *jobService) BulkCloseJobs(ctx context.Context, jobIDs []int64) error {
	for _, jobID := range jobIDs {
		if err := s.jobRepo.CloseJob(ctx, jobID); err != nil {
			return fmt.Errorf("failed to close job %d: %w", jobID, err)
		}
	}
	return nil
}

// BulkDeleteJobs deletes multiple jobs
func (s *jobService) BulkDeleteJobs(ctx context.Context, jobIDs []int64) error {
	for _, jobID := range jobIDs {
		if err := s.jobRepo.SoftDelete(ctx, jobID); err != nil {
			return fmt.Errorf("failed to delete job %d: %w", jobID, err)
		}
	}
	return nil
}

// ===== Validation =====

// ValidateJob validates job data
func (s *jobService) ValidateJob(ctx context.Context, j *job.Job) error {
	// Required fields
	if j.Title == "" {
		return errors.New("job title is required")
	}
	if j.Description == "" {
		return errors.New("job description is required")
	}
	if j.CompanyID == 0 {
		return errors.New("company ID is required")
	}

	// Validate salary range
	if j.SalaryMin != nil && j.SalaryMax != nil {
		if *j.SalaryMin > *j.SalaryMax {
			return errors.New("minimum salary cannot be greater than maximum salary")
		}
	}

	// Validate experience range
	if j.ExperienceMin != nil && j.ExperienceMax != nil {
		if *j.ExperienceMin > *j.ExperienceMax {
			return errors.New("minimum experience cannot be greater than maximum experience")
		}
	}

	// Validate total hires
	if j.TotalHires < 1 {
		return errors.New("total hires must be at least 1")
	}

	// Validate expiry date (must be in the future)
	if j.ExpiredAt != nil && j.ExpiredAt.Before(time.Now()) {
		return errors.New("expiry date must be in the future")
	}

	return nil
}

// CheckJobOwnership verifies if employer user owns the job
func (s *jobService) CheckJobOwnership(ctx context.Context, jobID, employerUserID int64) error {
	// Get job
	j, err := s.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf("job not found: %w", err)
	}

	// Check if employer user ID matches
	if j.EmployerUserID != nil && *j.EmployerUserID != employerUserID {
		// Check if user has permission through company employer users
		employerUser, err := s.companyRepo.FindEmployerUserByUserAndCompany(ctx, employerUserID, j.CompanyID)
		if err != nil || employerUser == nil {
			return errors.New("you do not have permission to modify this job")
		}

		// Check if role has permission (recruiter and above)
		if employerUser.Role != "recruiter" && employerUser.Role != "admin" && employerUser.Role != "owner" {
			return errors.New("you do not have permission to modify this job")
		}
	}

	return nil
}

// CheckJobStatus retrieves current job status
func (s *jobService) CheckJobStatus(ctx context.Context, jobID int64) (string, error) {
	j, err := s.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return "", fmt.Errorf("job not found: %w", err)
	}

	return j.Status, nil
}

// ===========================================
// MASTER DATA VALIDATION
// ===========================================

// ValidateMasterDataIDs validates that provided master data IDs exist
func (s *jobService) ValidateMasterDataIDs(ctx context.Context, industryID, districtID *int64) error {
	// Validate Industry ID
	if industryID != nil && *industryID > 0 {
		_, err := s.industryService.GetByID(ctx, *industryID)
		if err != nil {
			return fmt.Errorf("invalid industry_id: %w", err)
		}
	}

	// Validate District ID
	if districtID != nil && *districtID > 0 {
		_, err := s.districtService.GetByID(ctx, *districtID)
		if err != nil {
			return fmt.Errorf("invalid district_id: %w", err)
		}
	}

	return nil
}

// GetJobWithMasterData retrieves job with preloaded master data
func (s *jobService) GetJobWithMasterData(ctx context.Context, jobID int64) (*job.Job, error) {
	return s.jobRepo.FindByIDWithMasterData(ctx, jobID)
}
