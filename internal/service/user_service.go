package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"keerja-backend/internal/domain/master"
	"keerja-backend/internal/domain/user"
	"keerja-backend/internal/utils"
)

// userService implements the UserService interface
type userService struct {
	userRepo         user.UserRepository
	uploadService    UploadService
	skillsMasterRepo master.SkillsMasterRepository
}

// NewUserService creates a new user service instance
func NewUserService(
	userRepo user.UserRepository,
	uploadService UploadService,
	skillsMasterRepo master.SkillsMasterRepository,
) user.UserService {
	return &userService{
		userRepo:         userRepo,
		uploadService:    uploadService,
		skillsMasterRepo: skillsMasterRepo,
	}
}

// =============================================================================
// Registration and Verification (delegated to AuthService)
// =============================================================================

// Register - This is typically handled by AuthService, but included here for interface completeness
func (s *userService) Register(ctx context.Context, req *user.RegisterRequest) (*user.User, error) {
	return nil, fmt.Errorf("Register should be called via AuthService")
}

// VerifyEmail - This is typically handled by AuthService
func (s *userService) VerifyEmail(ctx context.Context, token string) error {
	return fmt.Errorf("VerifyEmail should be called via AuthService")
}

// ResendVerificationEmail - This is typically handled by AuthService
func (s *userService) ResendVerificationEmail(ctx context.Context, email string) error {
	return fmt.Errorf("ResendVerificationEmail should be called via AuthService")
}

// =============================================================================
// Profile Management
// =============================================================================

// GetProfile retrieves user profile by user ID
func (s *userService) GetProfile(ctx context.Context, userID int64) (*user.User, error) {
	// Get full profile with all relationships
	usr, err := s.userRepo.GetFullProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return usr, nil
}

// GetProfileBySlug retrieves user profile by slug
func (s *userService) GetProfileBySlug(ctx context.Context, slug string) (*user.User, error) {
	// Find profile by slug
	profile, err := s.userRepo.FindProfileBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("profile not found: %w", err)
	}

	// Get full profile with relationships
	usr, err := s.userRepo.GetFullProfile(ctx, profile.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	return usr, nil
}

// UpdateProfile updates user profile information
func (s *userService) UpdateProfile(ctx context.Context, userID int64, req *user.UpdateProfileRequest) error {
	// Get user data first (for full_name and phone updates)
	usr, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	if usr == nil {
		return fmt.Errorf("user not found")
	}

	// Update user table fields (full_name, phone)
	userUpdated := false
	if req.FullName != nil {
		usr.FullName = *req.FullName
		userUpdated = true
	}
	if req.Phone != nil {
		usr.Phone = req.Phone
		userUpdated = true
	}

	// Save user updates if any
	if userUpdated {
		if err := s.userRepo.Update(ctx, usr); err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	}

	// Get existing profile or create new one
	profile, err := s.userRepo.FindProfileByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find profile: %w", err)
	}

	// If profile doesn't exist, create a new one
	if profile == nil {
		profile = &user.UserProfile{
			UserID: userID,
		}
		if err := s.userRepo.CreateProfile(ctx, profile); err != nil {
			return fmt.Errorf("failed to create profile: %w", err)
		}
	}

	// Update profile fields if provided
	if req.Headline != nil {
		profile.Headline = req.Headline
	}
	if req.Bio != nil {
		profile.Bio = req.Bio
	}
	if req.Gender != nil {
		profile.Gender = req.Gender
	}
	if req.BirthDate != nil {
		birthDate, err := time.Parse("2006-01-02", *req.BirthDate)
		if err != nil {
			return fmt.Errorf("invalid birth date format: %w", err)
		}
		profile.BirthDate = &birthDate
	}
	if req.Nationality != nil {
		profile.Nationality = req.Nationality
	}
	if req.Address != nil {
		profile.Address = req.Address
	}
	if req.LocationCity != nil {
		profile.LocationCity = req.LocationCity
	}
	if req.LocationState != nil {
		profile.LocationState = req.LocationState
	}
	if req.LocationCountry != nil {
		profile.LocationCountry = req.LocationCountry
	}
	if req.PostalCode != nil {
		profile.PostalCode = req.PostalCode
	}
	if req.LinkedinURL != nil {
		profile.LinkedInURL = req.LinkedinURL
	}
	if req.PortfolioURL != nil {
		profile.PortfolioURL = req.PortfolioURL
	}
	if req.GithubURL != nil {
		profile.GithubURL = req.GithubURL
	}
	if req.DesiredPosition != nil {
		profile.DesiredPosition = req.DesiredPosition
	}
	if req.DesiredSalaryMin != nil {
		profile.DesiredSalaryMin = req.DesiredSalaryMin
	}
	if req.DesiredSalaryMax != nil {
		profile.DesiredSalaryMax = req.DesiredSalaryMax
	}
	if req.ExperienceLevel != nil {
		profile.ExperienceLevel = req.ExperienceLevel
	}
	if req.IndustryInterest != nil {
		profile.IndustryInterest = req.IndustryInterest
	}
	if req.AvailabilityStatus != nil {
		profile.AvailabilityStatus = *req.AvailabilityStatus
	}

	// Update profile
	if err := s.userRepo.UpdateProfile(ctx, profile); err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	return nil
}

// UploadAvatar uploads user avatar image
func (s *userService) UploadAvatar(ctx context.Context, userID int64, file *multipart.FileHeader) (string, error) {
	// Validate file
	if err := s.uploadService.ValidateFile(file, ImageTypes, MaxAvatarSize); err != nil {
		return "", fmt.Errorf("invalid avatar file: %w", err)
	}

	// Get existing profile
	profile, err := s.userRepo.FindProfileByUserID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("profile not found: %w", err)
	}

	// Delete old avatar if exists
	if profile.AvatarURL != nil && *profile.AvatarURL != "" {
		_ = s.uploadService.DeleteFile(ctx, *profile.AvatarURL)
	}

	// Upload new avatar
	avatarURL, err := s.uploadService.UploadFile(ctx, file, "avatars")
	if err != nil {
		return "", fmt.Errorf("failed to upload avatar: %w", err)
	}

	// Update profile with new avatar URL
	profile.AvatarURL = &avatarURL
	if err := s.userRepo.UpdateProfile(ctx, profile); err != nil {
		// Clean up uploaded file
		_ = s.uploadService.DeleteFile(ctx, avatarURL)
		return "", fmt.Errorf("failed to update profile with avatar URL: %w", err)
	}

	return avatarURL, nil
}

// UploadCover uploads user cover image
func (s *userService) UploadCover(ctx context.Context, userID int64, file *multipart.FileHeader) (string, error) {
	// Validate file
	if err := s.uploadService.ValidateFile(file, ImageTypes, MaxCoverSize); err != nil {
		return "", fmt.Errorf("invalid cover file: %w", err)
	}

	// Get existing profile
	profile, err := s.userRepo.FindProfileByUserID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("profile not found: %w", err)
	}

	// Delete old cover if exists
	if profile.CoverURL != nil && *profile.CoverURL != "" {
		_ = s.uploadService.DeleteFile(ctx, *profile.CoverURL)
	}

	// Upload new cover
	coverURL, err := s.uploadService.UploadFile(ctx, file, "covers")
	if err != nil {
		return "", fmt.Errorf("failed to upload cover: %w", err)
	}

	// Update profile with new cover URL
	profile.CoverURL = &coverURL
	if err := s.userRepo.UpdateProfile(ctx, profile); err != nil {
		// Clean up uploaded file
		_ = s.uploadService.DeleteFile(ctx, coverURL)
		return "", fmt.Errorf("failed to update profile with cover URL: %w", err)
	}

	return coverURL, nil
}

// DeleteAvatar deletes user avatar
func (s *userService) DeleteAvatar(ctx context.Context, userID int64) error {
	// Get existing profile
	profile, err := s.userRepo.FindProfileByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("profile not found: %w", err)
	}

	// Delete avatar file if exists
	if profile.AvatarURL != nil && *profile.AvatarURL != "" {
		_ = s.uploadService.DeleteFile(ctx, *profile.AvatarURL)
	}

	// Update profile to remove avatar URL
	profile.AvatarURL = nil
	if err := s.userRepo.UpdateProfile(ctx, profile); err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	return nil
}

// DeleteCover deletes user cover image
func (s *userService) DeleteCover(ctx context.Context, userID int64) error {
	// Get existing profile
	profile, err := s.userRepo.FindProfileByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("profile not found: %w", err)
	}

	// Delete cover file if exists
	if profile.CoverURL != nil && *profile.CoverURL != "" {
		_ = s.uploadService.DeleteFile(ctx, *profile.CoverURL)
	}

	// Update profile to remove cover URL
	profile.CoverURL = nil
	if err := s.userRepo.UpdateProfile(ctx, profile); err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	return nil
}

// =============================================================================
// Preference Management
// =============================================================================

// GetPreferences retrieves user preferences
func (s *userService) GetPreferences(ctx context.Context, userID int64) (*user.UserPreference, error) {
	pref, err := s.userRepo.FindPreferenceByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get preferences: %w", err)
	}

	return pref, nil
}

// UpdatePreferences updates user preferences
func (s *userService) UpdatePreferences(ctx context.Context, userID int64, req *user.UpdatePreferenceRequest) error {
	// Get existing preferences
	pref, err := s.userRepo.FindPreferenceByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("preferences not found: %w", err)
	}

	// Update fields if provided
	if req.LanguagePreference != nil {
		pref.LanguagePreference = *req.LanguagePreference
	}
	if req.ThemePreference != nil {
		pref.ThemePreference = *req.ThemePreference
	}
	if req.PreferredJobType != nil {
		pref.PreferredJobType = *req.PreferredJobType
	}
	if req.PreferredIndustry != nil {
		pref.PreferredIndustry = req.PreferredIndustry
	}
	if req.PreferredLocation != nil {
		pref.PreferredLocation = req.PreferredLocation
	}
	if req.PreferredSalaryMin != nil {
		pref.PreferredSalaryMin = req.PreferredSalaryMin
	}
	if req.PreferredSalaryMax != nil {
		pref.PreferredSalaryMax = req.PreferredSalaryMax
	}
	if req.EmailNotifications != nil {
		pref.EmailNotifications = *req.EmailNotifications
	}
	if req.SMSNotifications != nil {
		pref.SMSNotifications = *req.SMSNotifications
	}
	if req.PushNotifications != nil {
		pref.PushNotifications = *req.PushNotifications
	}
	if req.EmailMarketing != nil {
		pref.EmailMarketing = *req.EmailMarketing
	}
	if req.ProfileVisibility != nil {
		pref.ProfileVisibility = *req.ProfileVisibility
	}
	if req.ShowOnlineStatus != nil {
		pref.ShowOnlineStatus = *req.ShowOnlineStatus
	}
	if req.AllowDirectMessages != nil {
		pref.AllowDirectMessages = *req.AllowDirectMessages
	}
	if req.DataSharingConsent != nil {
		pref.DataSharingConsent = *req.DataSharingConsent
	}

	// Update preferences
	if err := s.userRepo.UpdatePreference(ctx, pref); err != nil {
		return fmt.Errorf("failed to update preferences: %w", err)
	}

	return nil
}

// =============================================================================
// Education Management
// =============================================================================

// AddEducation adds a new education entry
func (s *userService) AddEducation(ctx context.Context, userID int64, req *user.AddEducationRequest) (*user.UserEducation, error) {
	education := &user.UserEducation{
		UserID:          userID,
		InstitutionName: req.InstitutionName,
		Major:           req.Major,
		DegreeLevel:     req.DegreeLevel,
		StartYear:       req.StartYear,
		EndYear:         req.EndYear,
		GPA:             req.GPA,
		Activities:      req.Activities,
		Description:     req.Description,
		IsCurrent:       req.IsCurrent,
	}

	if err := s.userRepo.AddEducation(ctx, education); err != nil {
		return nil, fmt.Errorf("failed to add education: %w", err)
	}

	return education, nil
}

// UpdateEducation updates an education entry
func (s *userService) UpdateEducation(ctx context.Context, userID int64, educationID int64, req *user.UpdateEducationRequest) error {
	// Get educations to verify ownership
	educations, err := s.userRepo.GetEducationsByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get educations: %w", err)
	}

	// Find the education to update
	var education *user.UserEducation
	for i := range educations {
		if educations[i].ID == educationID {
			education = &educations[i]
			break
		}
	}

	if education == nil {
		return fmt.Errorf("education not found or unauthorized")
	}

	// Update fields if provided
	if req.InstitutionName != nil {
		education.InstitutionName = *req.InstitutionName
	}
	if req.Major != nil {
		education.Major = req.Major
	}
	if req.DegreeLevel != nil {
		education.DegreeLevel = req.DegreeLevel
	}
	if req.StartYear != nil {
		education.StartYear = req.StartYear
	}
	if req.EndYear != nil {
		education.EndYear = req.EndYear
	}
	if req.GPA != nil {
		education.GPA = req.GPA
	}
	if req.Activities != nil {
		education.Activities = req.Activities
	}
	if req.Description != nil {
		education.Description = req.Description
	}
	if req.IsCurrent != nil {
		education.IsCurrent = *req.IsCurrent
	}

	if err := s.userRepo.UpdateEducation(ctx, education); err != nil {
		return fmt.Errorf("failed to update education: %w", err)
	}

	return nil
}

// DeleteEducation deletes an education entry
func (s *userService) DeleteEducation(ctx context.Context, userID int64, educationID int64) error {
	// Verify ownership
	educations, err := s.userRepo.GetEducationsByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get educations: %w", err)
	}

	found := false
	for _, edu := range educations {
		if edu.ID == educationID {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("education not found or unauthorized")
	}

	if err := s.userRepo.DeleteEducation(ctx, educationID); err != nil {
		return fmt.Errorf("failed to delete education: %w", err)
	}

	return nil
}

// GetEducations retrieves all education entries for a user
func (s *userService) GetEducations(ctx context.Context, userID int64) ([]user.UserEducation, error) {
	educations, err := s.userRepo.GetEducationsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get educations: %w", err)
	}

	return educations, nil
}

// =============================================================================
// Experience Management
// =============================================================================

// AddExperience adds a new experience entry
func (s *userService) AddExperience(ctx context.Context, userID int64, req *user.AddExperienceRequest) (*user.UserExperience, error) {
	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}

	var endDate *time.Time
	if req.EndDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format: %w", err)
		}
		endDate = &parsed
	}

	experience := &user.UserExperience{
		UserID:          userID,
		CompanyName:     req.CompanyName,
		PositionTitle:   req.PositionTitle,
		Industry:        req.Industry,
		EmploymentType:  req.EmploymentType,
		StartDate:       startDate,
		EndDate:         endDate,
		IsCurrent:       req.IsCurrent,
		Description:     req.Description,
		Achievements:    req.Achievements,
		LocationCity:    req.LocationCity,
		LocationCountry: "Indonesia", // Default
	}

	if req.LocationCountry != nil {
		experience.LocationCountry = *req.LocationCountry
	}

	if err := s.userRepo.AddExperience(ctx, experience); err != nil {
		return nil, fmt.Errorf("failed to add experience: %w", err)
	}

	return experience, nil
}

// UpdateExperience updates an experience entry
func (s *userService) UpdateExperience(ctx context.Context, userID int64, experienceID int64, req *user.UpdateExperienceRequest) error {
	// Get experiences to verify ownership
	experiences, err := s.userRepo.GetExperiencesByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get experiences: %w", err)
	}

	// Find the experience to update
	var experience *user.UserExperience
	for i := range experiences {
		if experiences[i].ID == experienceID {
			experience = &experiences[i]
			break
		}
	}

	if experience == nil {
		return fmt.Errorf("experience not found or unauthorized")
	}

	// Update fields if provided
	if req.CompanyName != nil {
		experience.CompanyName = *req.CompanyName
	}
	if req.PositionTitle != nil {
		experience.PositionTitle = *req.PositionTitle
	}
	if req.Industry != nil {
		experience.Industry = req.Industry
	}
	if req.EmploymentType != nil {
		experience.EmploymentType = req.EmploymentType
	}
	if req.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return fmt.Errorf("invalid start date format: %w", err)
		}
		experience.StartDate = startDate
	}
	if req.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return fmt.Errorf("invalid end date format: %w", err)
		}
		experience.EndDate = &endDate
	}
	if req.IsCurrent != nil {
		experience.IsCurrent = *req.IsCurrent
	}
	if req.Description != nil {
		experience.Description = req.Description
	}
	if req.Achievements != nil {
		experience.Achievements = req.Achievements
	}
	if req.LocationCity != nil {
		experience.LocationCity = req.LocationCity
	}
	if req.LocationCountry != nil {
		experience.LocationCountry = *req.LocationCountry
	}

	if err := s.userRepo.UpdateExperience(ctx, experience); err != nil {
		return fmt.Errorf("failed to update experience: %w", err)
	}

	return nil
}

// DeleteExperience deletes an experience entry
func (s *userService) DeleteExperience(ctx context.Context, userID int64, experienceID int64) error {
	// Verify ownership
	experiences, err := s.userRepo.GetExperiencesByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get experiences: %w", err)
	}

	found := false
	for _, exp := range experiences {
		if exp.ID == experienceID {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("experience not found or unauthorized")
	}

	if err := s.userRepo.DeleteExperience(ctx, experienceID); err != nil {
		return fmt.Errorf("failed to delete experience: %w", err)
	}

	return nil
}

// GetExperiences retrieves all experience entries for a user
func (s *userService) GetExperiences(ctx context.Context, userID int64) ([]user.UserExperience, error) {
	experiences, err := s.userRepo.GetExperiencesByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get experiences: %w", err)
	}

	return experiences, nil
}

// =============================================================================
// Skill Management
// =============================================================================

// AddSkill adds a single skill for a user
// User can either select from skills_master (using skill_id) or input custom skill (using skill_name)
func (s *userService) AddSkill(ctx context.Context, userID int64, req *user.AddUserSkillRequest) error {
	skill := &user.UserSkill{
		UserID:     userID,
		SkillLevel: &req.ProficiencyLevel,
	}

	// Determine skill name: either from skills_master or custom input
	if req.SkillID != nil {
		// Query skills_master to get skill name by ID
		skillMaster, err := s.skillsMasterRepo.FindByID(ctx, *req.SkillID)
		if err != nil {
			return fmt.Errorf("skill_id %d not found in skills_master: %w", *req.SkillID, err)
		}
		skill.SkillName = skillMaster.Name
	} else {
		// Use custom skill name provided by user
		skill.SkillName = req.SkillName
	}

	if req.YearsOfExperience != nil {
		years := int(*req.YearsOfExperience)
		skill.YearsExperience = &years
	}

	if err := s.userRepo.AddSkill(ctx, skill); err != nil {
		return fmt.Errorf("failed to add skill: %w", err)
	}

	return nil
}

// AddSkills adds multiple skills for a user
// Each skill can be from skills_master or custom input
func (s *userService) AddSkills(ctx context.Context, userID int64, req *user.AddUserSkillsRequest) ([]user.UserSkill, error) {
	addedSkills := make([]user.UserSkill, 0, len(req.Skills))

	for i, skillReq := range req.Skills {
		skill := &user.UserSkill{
			UserID:     userID,
			SkillLevel: &skillReq.ProficiencyLevel,
		}

		// Determine skill name: either from skills_master or custom input
		if skillReq.SkillID != nil {
			// Query skills_master to get skill name by ID
			skillMaster, err := s.skillsMasterRepo.FindByID(ctx, *skillReq.SkillID)
			if err != nil {
				return nil, fmt.Errorf("skill #%d: skill_id %d not found in skills_master: %w", i+1, *skillReq.SkillID, err)
			}
			skill.SkillName = skillMaster.Name
		} else {
			// Use custom skill name provided by user
			skill.SkillName = skillReq.SkillName
		}

		if skillReq.YearsOfExperience != nil {
			years := int(*skillReq.YearsOfExperience)
			skill.YearsExperience = &years
		}

		if err := s.userRepo.AddSkill(ctx, skill); err != nil {
			return nil, fmt.Errorf("failed to add skill #%d: %w", i+1, err)
		}

		addedSkills = append(addedSkills, *skill)
	}

	return addedSkills, nil
}

// UpdateSkill updates a skill entry
func (s *userService) UpdateSkill(ctx context.Context, userID int64, skillID int64, req *user.UpdateSkillRequest) error {
	// Get skills to verify ownership
	skills, err := s.userRepo.GetSkillsByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get skills: %w", err)
	}

	// Find the skill to update
	var skill *user.UserSkill
	for i := range skills {
		if skills[i].ID == skillID {
			skill = &skills[i]
			break
		}
	}

	if skill == nil {
		return fmt.Errorf("skill not found or unauthorized")
	}

	// Update fields if provided
	if req.SkillLevel != nil {
		skill.SkillLevel = req.SkillLevel
	}
	if req.YearsExperience != nil {
		skill.YearsExperience = req.YearsExperience
	}
	if req.LastUsedAt != nil {
		lastUsed, err := time.Parse("2006-01-02", *req.LastUsedAt)
		if err != nil {
			return fmt.Errorf("invalid last used date format: %w", err)
		}
		skill.LastUsedAt = &lastUsed
	}

	if err := s.userRepo.UpdateSkill(ctx, skill); err != nil {
		return fmt.Errorf("failed to update skill: %w", err)
	}

	return nil
}

// DeleteSkill deletes a skill entry
func (s *userService) DeleteSkill(ctx context.Context, userID int64, skillID int64) error {
	// Verify ownership
	skills, err := s.userRepo.GetSkillsByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get skills: %w", err)
	}

	found := false
	for _, sk := range skills {
		if sk.ID == skillID {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("skill not found or unauthorized")
	}

	if err := s.userRepo.DeleteSkill(ctx, skillID); err != nil {
		return fmt.Errorf("failed to delete skill: %w", err)
	}

	return nil
}

// GetSkills retrieves all skills for a user
func (s *userService) GetSkills(ctx context.Context, userID int64) ([]user.UserSkill, error) {
	skills, err := s.userRepo.GetSkillsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get skills: %w", err)
	}

	return skills, nil
}

// =============================================================================
// Certification Management
// =============================================================================

// AddCertification adds a new certification entry
func (s *userService) AddCertification(ctx context.Context, userID int64, req *user.AddCertificationRequest) (*user.UserCertification, error) {
	cert := &user.UserCertification{
		UserID:              userID,
		CertificationName:   req.CertificationName,
		IssuingOrganization: req.IssuingOrganization,
		CredentialID:        req.CredentialID,
		CredentialURL:       req.CredentialURL,
		Description:         req.Description,
		FileURL:             req.FileURL,
	}

	// Parse dates if provided
	if req.IssueDate != nil {
		issueDate, err := time.Parse("2006-01-02", *req.IssueDate)
		if err != nil {
			return nil, fmt.Errorf("invalid issue date format: %w", err)
		}
		cert.IssueDate = &issueDate
	}

	if req.ExpirationDate != nil {
		expirationDate, err := time.Parse("2006-01-02", *req.ExpirationDate)
		if err != nil {
			return nil, fmt.Errorf("invalid expiration date format: %w", err)
		}
		cert.ExpirationDate = &expirationDate
	}

	if err := s.userRepo.AddCertification(ctx, cert); err != nil {
		return nil, fmt.Errorf("failed to add certification: %w", err)
	}

	return cert, nil
}

// UpdateCertification updates a certification entry
func (s *userService) UpdateCertification(ctx context.Context, userID int64, certID int64, req *user.UpdateCertificationRequest) error {
	// Get certifications to verify ownership
	certs, err := s.userRepo.GetCertificationsByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get certifications: %w", err)
	}

	// Find the certification to update
	var cert *user.UserCertification
	for i := range certs {
		if certs[i].ID == certID {
			cert = &certs[i]
			break
		}
	}

	if cert == nil {
		return fmt.Errorf("certification not found or unauthorized")
	}

	// Update fields if provided
	if req.CertificationName != nil {
		cert.CertificationName = *req.CertificationName
	}
	if req.IssuingOrganization != nil {
		cert.IssuingOrganization = *req.IssuingOrganization
	}
	if req.IssueDate != nil {
		issueDate, err := time.Parse("2006-01-02", *req.IssueDate)
		if err != nil {
			return fmt.Errorf("invalid issue date format: %w", err)
		}
		cert.IssueDate = &issueDate
	}
	if req.ExpirationDate != nil {
		expirationDate, err := time.Parse("2006-01-02", *req.ExpirationDate)
		if err != nil {
			return fmt.Errorf("invalid expiration date format: %w", err)
		}
		cert.ExpirationDate = &expirationDate
	}
	if req.CredentialID != nil {
		cert.CredentialID = req.CredentialID
	}
	if req.CredentialURL != nil {
		cert.CredentialURL = req.CredentialURL
	}
	if req.Description != nil {
		cert.Description = req.Description
	}
	if req.FileURL != nil {
		cert.FileURL = req.FileURL
	}

	if err := s.userRepo.UpdateCertification(ctx, cert); err != nil {
		return fmt.Errorf("failed to update certification: %w", err)
	}

	return nil
}

// DeleteCertification deletes a certification entry
func (s *userService) DeleteCertification(ctx context.Context, userID int64, certID int64) error {
	// Verify ownership
	certs, err := s.userRepo.GetCertificationsByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get certifications: %w", err)
	}

	found := false
	for _, c := range certs {
		if c.ID == certID {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("certification not found or unauthorized")
	}

	if err := s.userRepo.DeleteCertification(ctx, certID); err != nil {
		return fmt.Errorf("failed to delete certification: %w", err)
	}

	return nil
}

// GetCertifications retrieves all certifications for a user
func (s *userService) GetCertifications(ctx context.Context, userID int64) ([]user.UserCertification, error) {
	certs, err := s.userRepo.GetCertificationsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get certifications: %w", err)
	}

	return certs, nil
}

// =============================================================================
// Language Management
// =============================================================================

// AddLanguage adds a new language entry
func (s *userService) AddLanguage(ctx context.Context, userID int64, req *user.AddLanguageRequest) (*user.UserLanguage, error) {
	lang := &user.UserLanguage{
		UserID:             userID,
		LanguageName:       req.LanguageName,
		ProficiencyLevel:   req.ProficiencyLevel,
		CertificationName:  req.CertificationName,
		CertificationScore: req.CertificationScore,
		Notes:              req.Notes,
	}

	// Parse certification date if provided
	if req.CertificationDate != nil {
		certDate, err := time.Parse("2006-01-02", *req.CertificationDate)
		if err != nil {
			return nil, fmt.Errorf("invalid certification date format: %w", err)
		}
		lang.CertificationDate = &certDate
	}

	if err := s.userRepo.AddLanguage(ctx, lang); err != nil {
		return nil, fmt.Errorf("failed to add language: %w", err)
	}

	return lang, nil
}

// UpdateLanguage updates a language entry
func (s *userService) UpdateLanguage(ctx context.Context, userID int64, langID int64, req *user.UpdateLanguageRequest) error {
	// Get languages to verify ownership
	languages, err := s.userRepo.GetLanguagesByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get languages: %w", err)
	}

	// Find the language to update
	var lang *user.UserLanguage
	for i := range languages {
		if languages[i].ID == langID {
			lang = &languages[i]
			break
		}
	}

	if lang == nil {
		return fmt.Errorf("language not found or unauthorized")
	}

	// Update fields if provided
	if req.LanguageName != nil {
		lang.LanguageName = *req.LanguageName
	}
	if req.ProficiencyLevel != nil {
		lang.ProficiencyLevel = req.ProficiencyLevel
	}
	if req.CertificationName != nil {
		lang.CertificationName = req.CertificationName
	}
	if req.CertificationScore != nil {
		lang.CertificationScore = req.CertificationScore
	}
	if req.CertificationDate != nil {
		certDate, err := time.Parse("2006-01-02", *req.CertificationDate)
		if err != nil {
			return fmt.Errorf("invalid certification date format: %w", err)
		}
		lang.CertificationDate = &certDate
	}
	if req.Notes != nil {
		lang.Notes = req.Notes
	}

	if err := s.userRepo.UpdateLanguage(ctx, lang); err != nil {
		return fmt.Errorf("failed to update language: %w", err)
	}

	return nil
}

// DeleteLanguage deletes a language entry
func (s *userService) DeleteLanguage(ctx context.Context, userID int64, langID int64) error {
	// Verify ownership
	languages, err := s.userRepo.GetLanguagesByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get languages: %w", err)
	}

	found := false
	for _, l := range languages {
		if l.ID == langID {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("language not found or unauthorized")
	}

	if err := s.userRepo.DeleteLanguage(ctx, langID); err != nil {
		return fmt.Errorf("failed to delete language: %w", err)
	}

	return nil
}

// GetLanguages retrieves all languages for a user
func (s *userService) GetLanguages(ctx context.Context, userID int64) ([]user.UserLanguage, error) {
	languages, err := s.userRepo.GetLanguagesByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get languages: %w", err)
	}

	return languages, nil
}

// =============================================================================
// Project Management
// =============================================================================

// AddProject adds a new project entry
func (s *userService) AddProject(ctx context.Context, userID int64, req *user.AddProjectRequest) (*user.UserProject, error) {
	project := &user.UserProject{
		UserID:        userID,
		ProjectTitle:  req.ProjectTitle,
		RoleInProject: req.RoleInProject,
		ProjectType:   req.ProjectType,
		Description:   req.Description,
		Industry:      req.Industry,
		IsCurrent:     req.IsCurrent,
		ProjectURL:    req.ProjectURL,
		RepoURL:       req.RepoURL,
		Visibility:    "public", // Default
	}

	// Parse dates if provided
	if req.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start date format: %w", err)
		}
		project.StartDate = &startDate
	}

	if req.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format: %w", err)
		}
		project.EndDate = &endDate
	}

	if req.Visibility != nil {
		project.Visibility = *req.Visibility
	}

	// Handle arrays (converted to JSON strings for PostgreSQL array type)
	// Note: GORM handles PostgreSQL arrays, so we'll store as pointers to strings
	// The actual conversion is handled by GORM

	if err := s.userRepo.AddProject(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to add project: %w", err)
	}

	return project, nil
}

// UpdateProject updates a project entry
func (s *userService) UpdateProject(ctx context.Context, userID int64, projectID int64, req *user.UpdateProjectRequest) error {
	// Get projects to verify ownership
	projects, err := s.userRepo.GetProjectsByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get projects: %w", err)
	}

	// Find the project to update
	var project *user.UserProject
	for i := range projects {
		if projects[i].ID == projectID {
			project = &projects[i]
			break
		}
	}

	if project == nil {
		return fmt.Errorf("project not found or unauthorized")
	}

	// Update fields if provided
	if req.ProjectTitle != nil {
		project.ProjectTitle = *req.ProjectTitle
	}
	if req.RoleInProject != nil {
		project.RoleInProject = req.RoleInProject
	}
	if req.ProjectType != nil {
		project.ProjectType = req.ProjectType
	}
	if req.Description != nil {
		project.Description = req.Description
	}
	if req.Industry != nil {
		project.Industry = req.Industry
	}
	if req.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return fmt.Errorf("invalid start date format: %w", err)
		}
		project.StartDate = &startDate
	}
	if req.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return fmt.Errorf("invalid end date format: %w", err)
		}
		project.EndDate = &endDate
	}
	if req.IsCurrent != nil {
		project.IsCurrent = *req.IsCurrent
	}
	if req.ProjectURL != nil {
		project.ProjectURL = req.ProjectURL
	}
	if req.RepoURL != nil {
		project.RepoURL = req.RepoURL
	}
	if req.Visibility != nil {
		project.Visibility = *req.Visibility
	}

	if err := s.userRepo.UpdateProject(ctx, project); err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	return nil
}

// DeleteProject deletes a project entry
func (s *userService) DeleteProject(ctx context.Context, userID int64, projectID int64) error {
	// Verify ownership
	projects, err := s.userRepo.GetProjectsByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get projects: %w", err)
	}

	found := false
	for _, p := range projects {
		if p.ID == projectID {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("project not found or unauthorized")
	}

	if err := s.userRepo.DeleteProject(ctx, projectID); err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	return nil
}

// GetProjects retrieves all projects for a user
func (s *userService) GetProjects(ctx context.Context, userID int64) ([]user.UserProject, error) {
	projects, err := s.userRepo.GetProjectsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}

	return projects, nil
}

// =============================================================================
// Document Management
// =============================================================================

// UploadDocument uploads a document for the user
func (s *userService) UploadDocument(ctx context.Context, userID int64, file *multipart.FileHeader, req *user.UploadDocumentRequest) (*user.UserDocument, error) {
	// Validate file
	if err := s.uploadService.ValidateFile(file, DocumentTypes, MaxDocumentSize); err != nil {
		return nil, fmt.Errorf("invalid document file: %w", err)
	}

	// Upload file
	fileURL, err := s.uploadService.UploadFile(ctx, file, "documents")
	if err != nil {
		return nil, fmt.Errorf("failed to upload document: %w", err)
	}

	// Calculate checksum
	checksum, err := s.uploadService.CalculateChecksum(file)
	if err != nil {
		// Not critical, continue without checksum
		checksum = ""
	}

	// Create document record
	doc := &user.UserDocument{
		UserID:       userID,
		DocumentType: req.DocumentType,
		DocumentName: req.DocumentName,
		FileURL:      fileURL,
		FileSize:     &file.Size,
		MimeType:     utils.StringPtr(file.Header.Get("Content-Type")),
		Description:  req.Description,
		Checksum:     utils.StringPtr(checksum),
	}

	if err := s.userRepo.AddDocument(ctx, doc); err != nil {
		// Clean up uploaded file
		_ = s.uploadService.DeleteFile(ctx, fileURL)
		return nil, fmt.Errorf("failed to save document record: %w", err)
	}

	return doc, nil
}

// DeleteDocument deletes a document
func (s *userService) DeleteDocument(ctx context.Context, userID int64, documentID int64) error {
	// Verify ownership
	documents, err := s.userRepo.GetDocumentsByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get documents: %w", err)
	}

	var docToDelete *user.UserDocument
	for i := range documents {
		if documents[i].ID == documentID {
			docToDelete = &documents[i]
			break
		}
	}

	if docToDelete == nil {
		return fmt.Errorf("document not found or unauthorized")
	}

	// Delete file from storage
	_ = s.uploadService.DeleteFile(ctx, docToDelete.FileURL)

	// Delete document record
	if err := s.userRepo.DeleteDocument(ctx, documentID); err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	return nil
}

// GetDocuments retrieves all documents for a user
func (s *userService) GetDocuments(ctx context.Context, userID int64) ([]user.UserDocument, error) {
	documents, err := s.userRepo.GetDocumentsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %w", err)
	}

	return documents, nil
}

// =============================================================================
// Search and Discovery
// =============================================================================

// SearchUsers searches for users based on filters
func (s *userService) SearchUsers(ctx context.Context, filter *user.UserFilter) ([]user.User, int64, error) {
	users, total, err := s.userRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search users: %w", err)
	}

	return users, total, nil
}

// GetUsersBySkills finds users with specific skills
func (s *userService) GetUsersBySkills(ctx context.Context, skillNames []string) ([]user.User, error) {
	filter := &user.UserFilter{
		SkillNames: skillNames,
		Status:     utils.StringPtr("active"),
		IsVerified: utils.BoolPtr(true),
	}

	users, _, err := s.userRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by skills: %w", err)
	}

	return users, nil
}

// =============================================================================
// Profile Completion and Analytics
// =============================================================================

// GetProfileCompletionPercentage calculates profile completion percentage
func (s *userService) GetProfileCompletionPercentage(ctx context.Context, userID int64) (int, error) {
	// Get full profile
	usr, err := s.userRepo.GetFullProfile(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user profile: %w", err)
	}

	// Calculate completion percentage
	totalFields := 0
	completedFields := 0

	// Basic profile fields (weight: 30%)
	basicFields := 10
	basicCompleted := 0

	totalFields += basicFields

	if usr.FullName != "" {
		basicCompleted++
	}
	if usr.Email != "" {
		basicCompleted++
	}
	if usr.Phone != nil {
		basicCompleted++
	}
	if usr.Profile != nil {
		if usr.Profile.Headline != nil && *usr.Profile.Headline != "" {
			basicCompleted++
		}
		if usr.Profile.Bio != nil && *usr.Profile.Bio != "" {
			basicCompleted++
		}
		if usr.Profile.LocationCity != nil {
			basicCompleted++
		}
		if usr.Profile.AvatarURL != nil {
			basicCompleted++
		}
		if usr.Profile.DesiredPosition != nil {
			basicCompleted++
		}
		if usr.Profile.ExperienceLevel != nil {
			basicCompleted++
		}
		if usr.Profile.IndustryInterest != nil {
			basicCompleted++
		}
	}

	completedFields += basicCompleted

	// Education (weight: 15%)
	educationFields := 2
	totalFields += educationFields
	if len(usr.Educations) > 0 {
		completedFields += educationFields
	}

	// Experience (weight: 20%)
	experienceFields := 3
	totalFields += experienceFields
	if len(usr.Experiences) > 0 {
		completedFields += experienceFields
	}

	// Skills (weight: 20%)
	skillFields := 3
	totalFields += skillFields
	if len(usr.Skills) >= 3 {
		completedFields += skillFields
	} else if len(usr.Skills) > 0 {
		completedFields += len(usr.Skills)
	}

	// Certifications (weight: 10%)
	certFields := 1
	totalFields += certFields
	if len(usr.Certifications) > 0 {
		completedFields += certFields
	}

	// Languages (weight: 5%)
	langFields := 1
	totalFields += langFields
	if len(usr.Languages) > 0 {
		completedFields += langFields
	}

	// Calculate percentage
	percentage := (completedFields * 100) / totalFields

	return percentage, nil
}

// UpdateLastLogin updates the user's last login timestamp
func (s *userService) UpdateLastLogin(ctx context.Context, userID int64) error {
	// Get user
	usr, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Update last login
	now := time.Now()
	usr.LastLogin = &now

	if err := s.userRepo.Update(ctx, usr); err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

// =============================================================================
// Account Management
// =============================================================================

// UpdateStatus updates user status (active, inactive, suspended)
func (s *userService) UpdateStatus(ctx context.Context, userID int64, status string) error {
	// Validate status
	validStatuses := map[string]bool{
		"active":    true,
		"inactive":  true,
		"suspended": true,
	}

	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s", status)
	}

	// Get user
	usr, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Update status
	usr.Status = status

	if err := s.userRepo.Update(ctx, usr); err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	return nil
}

// SuspendAccount suspends a user account
func (s *userService) SuspendAccount(ctx context.Context, userID int64, reason string) error {
	// Get user
	usr, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Update status to suspended
	usr.Status = "suspended"

	if err := s.userRepo.Update(ctx, usr); err != nil {
		return fmt.Errorf("failed to suspend account: %w", err)
	}

	// TODO: Send suspension notification email
	// emailService.SendAccountSuspensionEmail(usr.Email, reason)

	return nil
}

// DeactivateAccount deactivates a user account (soft delete)
func (s *userService) DeactivateAccount(ctx context.Context, userID int64) error {
	// Get user
	usr, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Update status to inactive
	usr.Status = "inactive"
	usr.IsVerified = false

	if err := s.userRepo.Update(ctx, usr); err != nil {
		return fmt.Errorf("failed to deactivate account: %w", err)
	}

	// TODO: Send deactivation confirmation email
	// emailService.SendAccountDeactivationEmail(usr.Email)

	return nil
}

// DeleteAccount permanently deletes a user account
func (s *userService) DeleteAccount(ctx context.Context, userID int64) error {
	// Get full profile to clean up files
	usr, err := s.userRepo.GetFullProfile(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Clean up uploaded files
	if usr.Profile != nil {
		if usr.Profile.AvatarURL != nil {
			_ = s.uploadService.DeleteFile(ctx, *usr.Profile.AvatarURL)
		}
		if usr.Profile.CoverURL != nil {
			_ = s.uploadService.DeleteFile(ctx, *usr.Profile.CoverURL)
		}
	}

	// Clean up documents
	if len(usr.Documents) > 0 {
		for _, doc := range usr.Documents {
			_ = s.uploadService.DeleteFile(ctx, doc.FileURL)
		}
	}

	// Delete user (cascade will handle related records)
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	// TODO: Send account deletion confirmation email
	// emailService.SendAccountDeletionEmail(usr.Email)

	return nil
}
