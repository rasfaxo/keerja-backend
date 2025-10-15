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

// ToUserProfileResponse converts UserProfile entity to UserProfileResponse DTO
func ToUserProfileResponse(p *user.UserProfile) *response.UserProfileResponse {
	if p == nil {
		return nil
	}

	return &response.UserProfileResponse{
		ID:                  p.ID,
		Headline:            PtrToString(p.Headline),
		Summary:             PtrToString(p.Bio),
		DateOfBirth:         p.BirthDate,
		Gender:              PtrToString(p.Gender),
		City:                PtrToString(p.LocationCity),
		Country:             PtrToString(p.LocationCountry),
		ProfilePictureURL:   PtrToString(p.AvatarURL),
		CoverImageURL:       PtrToString(p.CoverURL),
		ProfileCompleteness: 0, // Calculate in handler
		UpdatedAt:           p.UpdatedAt,
	}
}

// ToUserEducationResponse converts UserEducation entity to UserEducationResponse DTO
func ToUserEducationResponse(e *user.UserEducation) *response.UserEducationResponse {
	if e == nil {
		return nil
	}

	// Convert int years to dates
	var startDate time.Time
	if e.StartYear != nil {
		startDate = time.Date(*e.StartYear, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	var endDatePtr *time.Time
	if e.EndYear != nil {
		endDateValue := time.Date(*e.EndYear, 1, 1, 0, 0, 0, 0, time.UTC)
		endDatePtr = &endDateValue
	}

	// Convert GPA to string
	var grade string
	if e.GPA != nil {
		grade = Float64ToString(*e.GPA)
	}

	return &response.UserEducationResponse{
		ID:                  e.ID,
		InstitutionName:     e.InstitutionName,
		Degree:              PtrToString(e.DegreeLevel),
		FieldOfStudy:        PtrToString(e.Major),
		StartDate:           startDate,
		EndDate:             endDatePtr,
		Grade:               grade,
		Description:         PtrToString(e.Description),
		IsCurrentlyStudying: e.IsCurrent,
		CreatedAt:           e.CreatedAt,
	}
}

// ToUserExperienceResponse converts UserExperience entity to UserExperienceResponse DTO
func ToUserExperienceResponse(e *user.UserExperience) *response.UserExperienceResponse {
	if e == nil {
		return nil
	}

	return &response.UserExperienceResponse{
		ID:                 e.ID,
		CompanyName:        e.CompanyName,
		JobTitle:           e.PositionTitle,
		EmploymentType:     PtrToString(e.EmploymentType),
		Location:           PtrToString(e.LocationCity),
		StartDate:          e.StartDate,
		EndDate:            e.EndDate,
		IsCurrentlyWorking: e.IsCurrent,
		Description:        PtrToString(e.Description),
		CreatedAt:          e.CreatedAt,
	}
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
