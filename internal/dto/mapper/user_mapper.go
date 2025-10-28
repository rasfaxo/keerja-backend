package mapper

import (
	"keerja-backend/internal/domain/user"
	"keerja-backend/internal/dto/response"
	"time"
)

// ToUserResponse converts User entity to UserResponse DTO
func ToUserResponse(u *user.User) *response.UserResponse {
	if u == nil {
		return nil
	}

	resp := &response.UserResponse{
		ID:         u.ID,
		UUID:       u.UUID.String(),
		FullName:   u.FullName,
		Email:      u.Email,
		Phone:      PtrToString(u.Phone),
		UserType:   u.UserType,
		IsVerified: u.IsVerified,
		Status:     u.Status,
		LastLogin:  u.LastLogin,
		CreatedAt:  u.CreatedAt,
	}

	// Map profile if exists
	if u.Profile != nil {
		resp.Profile = ToUserProfileResponse(u.Profile)
	}

	return resp
}

// ToUserDetailResponse converts User entity with relations to UserDetailResponse DTO
func ToUserDetailResponse(u *user.User) *response.UserDetailResponse {
	if u == nil {
		return nil
	}

	resp := &response.UserDetailResponse{
		ID:         u.ID,
		UUID:       u.UUID.String(),
		FullName:   u.FullName,
		Email:      u.Email,
		Phone:      PtrToString(u.Phone),
		UserType:   u.UserType,
		IsVerified: u.IsVerified,
		Status:     u.Status,
		LastLogin:  u.LastLogin,
		CreatedAt:  u.CreatedAt,
	}

	// Map profile
	if u.Profile != nil {
		resp.Profile = ToUserProfileResponse(u.Profile)
	}

	// Map preference
	if u.Preference != nil {
		resp.Preference = ToUserPreferenceResponse(u.Preference)
	}

	// Map educations
	if len(u.Educations) > 0 {
		resp.Educations = make([]response.UserEducationResponse, 0, len(u.Educations))
		for _, edu := range u.Educations {
			if mapped := ToUserEducationResponse(&edu); mapped != nil {
				resp.Educations = append(resp.Educations, *mapped)
			}
		}
	}

	// Map experiences
	if len(u.Experiences) > 0 {
		resp.Experiences = make([]response.UserExperienceResponse, 0, len(u.Experiences))
		for _, exp := range u.Experiences {
			if mapped := ToUserExperienceResponse(&exp); mapped != nil {
				resp.Experiences = append(resp.Experiences, *mapped)
			}
		}
	}

	// Map skills
	if len(u.Skills) > 0 {
		resp.Skills = make([]response.UserSkillResponse, 0, len(u.Skills))
		for _, skill := range u.Skills {
			if mapped := ToUserSkillResponse(&skill); mapped != nil {
				resp.Skills = append(resp.Skills, *mapped)
			}
		}
	}

	// Map certifications
	if len(u.Certifications) > 0 {
		resp.Certifications = make([]response.UserCertificationResponse, 0, len(u.Certifications))
		for _, cert := range u.Certifications {
			if mapped := ToUserCertificationResponse(&cert); mapped != nil {
				resp.Certifications = append(resp.Certifications, *mapped)
			}
		}
	}

	// Map languages
	if len(u.Languages) > 0 {
		resp.Languages = make([]response.UserLanguageResponse, 0, len(u.Languages))
		for _, lang := range u.Languages {
			if mapped := ToUserLanguageResponse(&lang); mapped != nil {
				resp.Languages = append(resp.Languages, *mapped)
			}
		}
	}

	// Map projects
	if len(u.Projects) > 0 {
		resp.Projects = make([]response.UserProjectResponse, 0, len(u.Projects))
		for _, proj := range u.Projects {
			if mapped := ToUserProjectResponse(&proj); mapped != nil {
				resp.Projects = append(resp.Projects, *mapped)
			}
		}
	}

	// Map documents
	if len(u.Documents) > 0 {
		resp.Documents = make([]response.UserDocumentResponse, 0, len(u.Documents))
		for _, doc := range u.Documents {
			if mapped := ToUserDocumentResponse(&doc); mapped != nil {
				resp.Documents = append(resp.Documents, *mapped)
			}
		}
	}

	return resp
}

// ToUserResponseWithIncludes converts User entity with selective relations based on includes parameter
// Supports: profile, educations, experiences, skills, certifications, languages, projects, documents, preference
func ToUserResponseWithIncludes(u *user.User, includes []string) *response.UserDetailResponse {
	if u == nil {
		return nil
	}

	// Create a map for quick lookup
	includeMap := make(map[string]bool)
	for _, inc := range includes {
		includeMap[inc] = true
	}

	resp := &response.UserDetailResponse{
		ID:         u.ID,
		UUID:       u.UUID.String(),
		FullName:   u.FullName,
		Email:      u.Email,
		Phone:      PtrToString(u.Phone),
		UserType:   u.UserType,
		IsVerified: u.IsVerified,
		Status:     u.Status,
		LastLogin:  u.LastLogin,
		CreatedAt:  u.CreatedAt,
	}

	// Always include basic profile (it's lightweight)
	if u.Profile != nil {
		resp.Profile = ToUserProfileResponse(u.Profile)
	}

	// Conditionally include preference
	if includeMap["preference"] && u.Preference != nil {
		resp.Preference = ToUserPreferenceResponse(u.Preference)
	}

	// Conditionally include educations
	if includeMap["educations"] && len(u.Educations) > 0 {
		resp.Educations = make([]response.UserEducationResponse, 0, len(u.Educations))
		for _, edu := range u.Educations {
			if mapped := ToUserEducationResponse(&edu); mapped != nil {
				resp.Educations = append(resp.Educations, *mapped)
			}
		}
	}

	// Conditionally include experiences
	if includeMap["experiences"] && len(u.Experiences) > 0 {
		resp.Experiences = make([]response.UserExperienceResponse, 0, len(u.Experiences))
		for _, exp := range u.Experiences {
			if mapped := ToUserExperienceResponse(&exp); mapped != nil {
				resp.Experiences = append(resp.Experiences, *mapped)
			}
		}
	}

	// Conditionally include skills
	if includeMap["skills"] && len(u.Skills) > 0 {
		resp.Skills = make([]response.UserSkillResponse, 0, len(u.Skills))
		for _, skill := range u.Skills {
			if mapped := ToUserSkillResponse(&skill); mapped != nil {
				resp.Skills = append(resp.Skills, *mapped)
			}
		}
	}

	// Conditionally include certifications
	if includeMap["certifications"] && len(u.Certifications) > 0 {
		resp.Certifications = make([]response.UserCertificationResponse, 0, len(u.Certifications))
		for _, cert := range u.Certifications {
			if mapped := ToUserCertificationResponse(&cert); mapped != nil {
				resp.Certifications = append(resp.Certifications, *mapped)
			}
		}
	}

	// Conditionally include languages
	if includeMap["languages"] && len(u.Languages) > 0 {
		resp.Languages = make([]response.UserLanguageResponse, 0, len(u.Languages))
		for _, lang := range u.Languages {
			if mapped := ToUserLanguageResponse(&lang); mapped != nil {
				resp.Languages = append(resp.Languages, *mapped)
			}
		}
	}

	// Conditionally include projects
	if includeMap["projects"] && len(u.Projects) > 0 {
		resp.Projects = make([]response.UserProjectResponse, 0, len(u.Projects))
		for _, proj := range u.Projects {
			if mapped := ToUserProjectResponse(&proj); mapped != nil {
				resp.Projects = append(resp.Projects, *mapped)
			}
		}
	}

	// Conditionally include documents
	if includeMap["documents"] && len(u.Documents) > 0 {
		resp.Documents = make([]response.UserDocumentResponse, 0, len(u.Documents))
		for _, doc := range u.Documents {
			if mapped := ToUserDocumentResponse(&doc); mapped != nil {
				resp.Documents = append(resp.Documents, *mapped)
			}
		}
	}

	return resp
}

// ToUserProfileResponse converts UserProfile entity to UserProfileResponse DTO
func ToUserProfileResponse(p *user.UserProfile) *response.UserProfileResponse {
	if p == nil {
		return nil
	}

	return &response.UserProfileResponse{
		ID:                 p.ID,
		UserID:             p.UserID,
		Headline:           p.Headline,
		Bio:                p.Bio,
		Gender:             p.Gender,
		BirthDate:          p.BirthDate,
		LocationCity:       p.LocationCity,
		LocationCountry:    p.LocationCountry,
		DesiredPosition:    p.DesiredPosition,
		DesiredSalaryMin:   p.DesiredSalaryMin,
		DesiredSalaryMax:   p.DesiredSalaryMax,
		ExperienceLevel:    p.ExperienceLevel,
		IndustryInterest:   p.IndustryInterest,
		AvailabilityStatus: p.AvailabilityStatus,
		ProfileVisibility:  p.ProfileVisibility,
		Slug:               p.Slug,
		AvatarURL:          p.AvatarURL,
		CoverURL:           p.CoverURL,
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
	}
}

// ToUserEducationResponse converts UserEducation entity to UserEducationResponse DTO
func ToEducationResponse(e *user.UserEducation) *response.UserEducationResponse {
	if e == nil {
		return nil
	}

	return &response.UserEducationResponse{
		ID:              e.ID,
		InstitutionName: e.InstitutionName,
		Major:           e.Major,
		DegreeLevel:     e.DegreeLevel,
		StartYear:       e.StartYear,
		EndYear:         e.EndYear,
		GPA:             e.GPA,
		Activities:      e.Activities,
		Description:     e.Description,
		IsCurrent:       e.IsCurrent,
		CreatedAt:       e.CreatedAt.Format(time.RFC3339),
	}
}

// ToUserEducationResponse is an alias for ToEducationResponse for backward compatibility
func ToUserEducationResponse(e *user.UserEducation) *response.UserEducationResponse {
	return ToEducationResponse(e)
}

// ToUserExperienceResponse converts UserExperience entity to UserExperienceResponse DTO
func ToExperienceResponse(e *user.UserExperience) *response.UserExperienceResponse {
	if e == nil {
		return nil
	}

	var endDate *string
	if e.EndDate != nil {
		endDateStr := e.EndDate.Format("2006-01-02")
		endDate = &endDateStr
	}

	return &response.UserExperienceResponse{
		ID:              e.ID,
		CompanyName:     e.CompanyName,
		PositionTitle:   e.PositionTitle,
		Industry:        e.Industry,
		EmploymentType:  e.EmploymentType,
		StartDate:       e.StartDate.Format("2006-01-02"),
		EndDate:         endDate,
		IsCurrent:       e.IsCurrent,
		Description:     e.Description,
		Achievements:    e.Achievements,
		LocationCity:    e.LocationCity,
		LocationCountry: e.LocationCountry,
		CreatedAt:       e.CreatedAt.Format(time.RFC3339),
	}
}

// ToUserExperienceResponse is an alias for ToExperienceResponse for backward compatibility
func ToUserExperienceResponse(e *user.UserExperience) *response.UserExperienceResponse {
	return ToExperienceResponse(e)
}

// ToUserSkillResponse converts UserSkill entity to UserSkillResponse DTO
func ToUserSkillResponse(s *user.UserSkill) *response.UserSkillResponse {
	if s == nil {
		return nil
	}

	return &response.UserSkillResponse{
		ID:               s.ID,
		SkillName:        s.SkillName,
		ProficiencyLevel: PtrToString(s.SkillLevel),
		CreatedAt:        s.CreatedAt,
	}
}

// ToUserCertificationResponse converts UserCertification entity to UserCertificationResponse DTO
func ToUserCertificationResponse(c *user.UserCertification) *response.UserCertificationResponse {
	if c == nil {
		return nil
	}

	// Determine DoesNotExpire
	doesNotExpire := c.ExpirationDate == nil

	// Convert IssueDate from pointer to value
	var issueDate time.Time
	if c.IssueDate != nil {
		issueDate = *c.IssueDate
	}

	return &response.UserCertificationResponse{
		ID:                  c.ID,
		Name:                c.CertificationName,
		IssuingOrganization: c.IssuingOrganization,
		IssueDate:           issueDate,
		ExpiryDate:          c.ExpirationDate,
		CredentialID:        PtrToString(c.CredentialID),
		CredentialURL:       PtrToString(c.CredentialURL),
		DoesNotExpire:       doesNotExpire,
		CreatedAt:           c.CreatedAt,
	}
}

// ToUserLanguageResponse converts UserLanguage entity to UserLanguageResponse DTO
func ToUserLanguageResponse(l *user.UserLanguage) *response.UserLanguageResponse {
	if l == nil {
		return nil
	}

	return &response.UserLanguageResponse{
		ID:               l.ID,
		LanguageName:     l.LanguageName,
		ProficiencyLevel: PtrToString(l.ProficiencyLevel),
		CreatedAt:        l.CreatedAt,
	}
}

// ToUserProjectResponse converts UserProject entity to UserProjectResponse DTO
func ToUserProjectResponse(p *user.UserProject) *response.UserProjectResponse {
	if p == nil {
		return nil
	}

	// Convert StartDate from pointer to value
	var startDate time.Time
	if p.StartDate != nil {
		startDate = *p.StartDate
	}

	// Get technologies from SkillsUsed array
	var technologies string
	if p.SkillsUsed != nil {
		technologies = *p.SkillsUsed
	}

	return &response.UserProjectResponse{
		ID:           p.ID,
		ProjectName:  p.ProjectTitle,
		Description:  PtrToString(p.Description),
		Role:         PtrToString(p.RoleInProject),
		StartDate:    startDate,
		EndDate:      p.EndDate,
		ProjectURL:   PtrToString(p.ProjectURL),
		IsOngoing:    p.IsCurrent,
		Technologies: technologies,
		CreatedAt:    p.CreatedAt,
	}
}

// ToUserDocumentResponse converts UserDocument entity to UserDocumentResponse DTO
func ToUserDocumentResponse(d *user.UserDocument) *response.UserDocumentResponse {
	if d == nil {
		return nil
	}

	// Convert DocumentType from pointer to value
	var docType string
	if d.DocumentType != nil {
		docType = *d.DocumentType
	}

	// Convert FileSize from pointer to value
	var fileSize int64
	if d.FileSize != nil {
		fileSize = *d.FileSize
	}

	return &response.UserDocumentResponse{
		ID:           d.ID,
		DocumentType: docType,
		Title:        d.DocumentName,
		FileURL:      d.FileURL,
		FileName:     d.DocumentName,
		FileSize:     fileSize,
		IsVerified:   d.Verified,
		UploadedAt:   d.UploadedAt,
	}
}

// ToUserPreferenceResponse converts UserPreference entity to UserPreferenceResponse DTO
func ToUserPreferenceResponse(p *user.UserPreference) *response.UserPreferenceResponse {
	if p == nil {
		return nil
	}

	// Convert salary from float to int64
	var salaryMin, salaryMax *int64
	if p.PreferredSalaryMin != nil {
		val := int64(*p.PreferredSalaryMin)
		salaryMin = &val
	}
	if p.PreferredSalaryMax != nil {
		val := int64(*p.PreferredSalaryMax)
		salaryMax = &val
	}

	// Determine WillingToRelocate (if PreferredLocation is not empty)
	willingToRelocate := p.PreferredLocation != nil && *p.PreferredLocation != ""

	// Determine AvailableForRemote (based on PreferredJobType)
	availableForRemote := p.PreferredJobType == "remote" || p.PreferredJobType == "hybrid"

	// Determine IsOpenToWork (based on ShowOnlineStatus and AllowDirectMessages)
	isOpenToWork := p.ShowOnlineStatus && p.AllowDirectMessages

	return &response.UserPreferenceResponse{
		ID:                 p.ID,
		ExpectedSalaryMin:  salaryMin,
		ExpectedSalaryMax:  salaryMax,
		WillingToRelocate:  willingToRelocate,
		AvailableForRemote: availableForRemote,
		IsOpenToWork:       isOpenToWork,
		ProfileVisibility:  p.ProfileVisibility,
		UpdatedAt:          p.UpdatedAt,
	}
}
