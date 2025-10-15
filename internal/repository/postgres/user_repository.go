package postgres

import (
	"context"
	"fmt"
	"strings"

	"keerja-backend/internal/domain/user"

	"gorm.io/gorm"
)

// userRepository implements user.UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &userRepository{db: db}
}

// ===========================================
// USER CRUD OPERATIONS
// ===========================================

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

// FindByID finds a user by ID
func (r *userRepository) FindByID(ctx context.Context, id int64) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).
		Preload("Profile").
		Preload("Preference").
		First(&u, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// FindByUUID finds a user by UUID
func (r *userRepository) FindByUUID(ctx context.Context, uuid string) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).
		Preload("Profile").
		Preload("Preference").
		Where("uuid = ?", uuid).
		First(&u).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// FindByEmail finds a user by email (for authentication)
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).
		Preload("Profile").
		Preload("Preference").
		Where("email = ?", email).
		First(&u).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// Update updates a user
func (r *userRepository) Update(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).Save(u).Error
}

// Delete soft deletes a user
func (r *userRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&user.User{}, id).Error
}

// List retrieves users with filtering, pagination, and sorting
func (r *userRepository) List(ctx context.Context, filter *user.UserFilter) ([]user.User, int64, error) {
	var users []user.User
	var total int64

	query := r.db.WithContext(ctx).Model(&user.User{})

	// Apply filters
	if filter != nil {
		if filter.UserType != nil {
			query = query.Where("user_type = ?", *filter.UserType)
		}
		if filter.Status != nil {
			query = query.Where("status = ?", *filter.Status)
		}
		if filter.IsVerified != nil {
			query = query.Where("is_verified = ?", *filter.IsVerified)
		}

		// Join with profile for additional filters
		if filter.LocationCity != nil || filter.ExperienceLevel != nil || filter.IndustryInterest != nil {
			query = query.Joins("LEFT JOIN user_profiles ON user_profiles.user_id = users.id")

			if filter.LocationCity != nil {
				query = query.Where("user_profiles.location_city = ?", *filter.LocationCity)
			}
			if filter.ExperienceLevel != nil {
				query = query.Where("user_profiles.experience_level = ?", *filter.ExperienceLevel)
			}
			if filter.IndustryInterest != nil {
				query = query.Where("user_profiles.industry_interest = ?", *filter.IndustryInterest)
			}
		}

		// Search by name or email
		if filter.SearchQuery != nil && *filter.SearchQuery != "" {
			searchPattern := "%" + strings.ToLower(*filter.SearchQuery) + "%"
			query = query.Where("LOWER(full_name) LIKE ? OR LOWER(email) LIKE ?", searchPattern, searchPattern)
		}

		// Filter by skills
		if len(filter.SkillNames) > 0 {
			query = query.Joins("JOIN user_skills ON user_skills.user_id = users.id").
				Where("LOWER(user_skills.skill_name) IN ?", toLowerSlice(filter.SkillNames)).
				Group("users.id").
				Having("COUNT(DISTINCT user_skills.skill_name) = ?", len(filter.SkillNames))
		}
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	page := 1
	limit := 10
	if filter != nil {
		if filter.Page > 0 {
			page = filter.Page
		}
		if filter.Limit > 0 {
			limit = filter.Limit
		}
	}
	offset := (page - 1) * limit

	// Apply sorting
	sortBy := "created_at"
	sortOrder := "DESC"
	if filter != nil {
		if filter.SortBy != "" {
			sortBy = filter.SortBy
		}
		if filter.SortOrder != "" {
			sortOrder = strings.ToUpper(filter.SortOrder)
		}
	}

	// Execute query with preloads
	err := query.
		Preload("Profile").
		Preload("Preference").
		Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	return users, total, err
}

// ===========================================
// PROFILE OPERATIONS
// ===========================================

// CreateProfile creates a user profile
func (r *userRepository) CreateProfile(ctx context.Context, profile *user.UserProfile) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

// FindProfileByUserID finds a profile by user ID
func (r *userRepository) FindProfileByUserID(ctx context.Context, userID int64) (*user.UserProfile, error) {
	var profile user.UserProfile
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

// FindProfileBySlug finds a profile by slug
func (r *userRepository) FindProfileBySlug(ctx context.Context, slug string) (*user.UserProfile, error) {
	var profile user.UserProfile
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("slug = ?", slug).
		First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

// UpdateProfile updates a user profile
func (r *userRepository) UpdateProfile(ctx context.Context, profile *user.UserProfile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}

// ===========================================
// PREFERENCE OPERATIONS
// ===========================================

// CreatePreference creates user preferences
func (r *userRepository) CreatePreference(ctx context.Context, preference *user.UserPreference) error {
	return r.db.WithContext(ctx).Create(preference).Error
}

// FindPreferenceByUserID finds preferences by user ID
func (r *userRepository) FindPreferenceByUserID(ctx context.Context, userID int64) (*user.UserPreference, error) {
	var preference user.UserPreference
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&preference).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &preference, nil
}

// UpdatePreference updates user preferences
func (r *userRepository) UpdatePreference(ctx context.Context, preference *user.UserPreference) error {
	return r.db.WithContext(ctx).Save(preference).Error
}

// ===========================================
// EDUCATION OPERATIONS
// ===========================================

// AddEducation adds an education entry
func (r *userRepository) AddEducation(ctx context.Context, education *user.UserEducation) error {
	return r.db.WithContext(ctx).Create(education).Error
}

// UpdateEducation updates an education entry
func (r *userRepository) UpdateEducation(ctx context.Context, education *user.UserEducation) error {
	return r.db.WithContext(ctx).Save(education).Error
}

// DeleteEducation deletes an education entry
func (r *userRepository) DeleteEducation(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&user.UserEducation{}, id).Error
}

// GetEducationsByUserID retrieves all education entries for a user
func (r *userRepository) GetEducationsByUserID(ctx context.Context, userID int64) ([]user.UserEducation, error) {
	var educations []user.UserEducation
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_current DESC, end_year DESC NULLS FIRST, start_year DESC").
		Find(&educations).Error
	return educations, err
}

// ===========================================
// EXPERIENCE OPERATIONS
// ===========================================

// AddExperience adds a work experience entry
func (r *userRepository) AddExperience(ctx context.Context, experience *user.UserExperience) error {
	return r.db.WithContext(ctx).Create(experience).Error
}

// UpdateExperience updates a work experience entry
func (r *userRepository) UpdateExperience(ctx context.Context, experience *user.UserExperience) error {
	return r.db.WithContext(ctx).Save(experience).Error
}

// DeleteExperience deletes a work experience entry
func (r *userRepository) DeleteExperience(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&user.UserExperience{}, id).Error
}

// GetExperiencesByUserID retrieves all work experiences for a user
func (r *userRepository) GetExperiencesByUserID(ctx context.Context, userID int64) ([]user.UserExperience, error) {
	var experiences []user.UserExperience
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_current DESC, start_date DESC").
		Find(&experiences).Error
	return experiences, err
}

// ===========================================
// SKILL OPERATIONS
// ===========================================

// AddSkill adds a skill entry
func (r *userRepository) AddSkill(ctx context.Context, skill *user.UserSkill) error {
	return r.db.WithContext(ctx).Create(skill).Error
}

// UpdateSkill updates a skill entry
func (r *userRepository) UpdateSkill(ctx context.Context, skill *user.UserSkill) error {
	return r.db.WithContext(ctx).Save(skill).Error
}

// DeleteSkill deletes a skill entry
func (r *userRepository) DeleteSkill(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&user.UserSkill{}, id).Error
}

// GetSkillsByUserID retrieves all skills for a user
func (r *userRepository) GetSkillsByUserID(ctx context.Context, userID int64) ([]user.UserSkill, error) {
	var skills []user.UserSkill
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("skill_level DESC, years_experience DESC NULLS LAST, skill_name ASC").
		Find(&skills).Error
	return skills, err
}

// ===========================================
// CERTIFICATION OPERATIONS
// ===========================================

// AddCertification adds a certification entry
func (r *userRepository) AddCertification(ctx context.Context, cert *user.UserCertification) error {
	return r.db.WithContext(ctx).Create(cert).Error
}

// UpdateCertification updates a certification entry
func (r *userRepository) UpdateCertification(ctx context.Context, cert *user.UserCertification) error {
	return r.db.WithContext(ctx).Save(cert).Error
}

// DeleteCertification deletes a certification entry
func (r *userRepository) DeleteCertification(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&user.UserCertification{}, id).Error
}

// GetCertificationsByUserID retrieves all certifications for a user
func (r *userRepository) GetCertificationsByUserID(ctx context.Context, userID int64) ([]user.UserCertification, error) {
	var certifications []user.UserCertification
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("is_active = ?", true).
		Order("issue_date DESC NULLS LAST, certification_name ASC").
		Find(&certifications).Error
	return certifications, err
}

// ===========================================
// LANGUAGE OPERATIONS
// ===========================================

// AddLanguage adds a language entry
func (r *userRepository) AddLanguage(ctx context.Context, lang *user.UserLanguage) error {
	return r.db.WithContext(ctx).Create(lang).Error
}

// UpdateLanguage updates a language entry
func (r *userRepository) UpdateLanguage(ctx context.Context, lang *user.UserLanguage) error {
	return r.db.WithContext(ctx).Save(lang).Error
}

// DeleteLanguage deletes a language entry
func (r *userRepository) DeleteLanguage(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&user.UserLanguage{}, id).Error
}

// GetLanguagesByUserID retrieves all languages for a user
func (r *userRepository) GetLanguagesByUserID(ctx context.Context, userID int64) ([]user.UserLanguage, error) {
	var languages []user.UserLanguage
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("is_active = ?", true).
		Order("proficiency_level DESC, language_name ASC").
		Find(&languages).Error
	return languages, err
}

// ===========================================
// PROJECT OPERATIONS
// ===========================================

// AddProject adds a project entry
func (r *userRepository) AddProject(ctx context.Context, project *user.UserProject) error {
	return r.db.WithContext(ctx).Create(project).Error
}

// UpdateProject updates a project entry
func (r *userRepository) UpdateProject(ctx context.Context, project *user.UserProject) error {
	return r.db.WithContext(ctx).Save(project).Error
}

// DeleteProject deletes a project entry
func (r *userRepository) DeleteProject(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&user.UserProject{}, id).Error
}

// GetProjectsByUserID retrieves all projects for a user
func (r *userRepository) GetProjectsByUserID(ctx context.Context, userID int64) ([]user.UserProject, error) {
	var projects []user.UserProject
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("featured DESC, is_current DESC, start_date DESC NULLS LAST").
		Find(&projects).Error
	return projects, err
}

// ===========================================
// DOCUMENT OPERATIONS
// ===========================================

// AddDocument adds a document entry
func (r *userRepository) AddDocument(ctx context.Context, doc *user.UserDocument) error {
	return r.db.WithContext(ctx).Create(doc).Error
}

// UpdateDocument updates a document entry
func (r *userRepository) UpdateDocument(ctx context.Context, doc *user.UserDocument) error {
	return r.db.WithContext(ctx).Save(doc).Error
}

// DeleteDocument deletes a document entry
func (r *userRepository) DeleteDocument(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&user.UserDocument{}, id).Error
}

// GetDocumentsByUserID retrieves all documents for a user
func (r *userRepository) GetDocumentsByUserID(ctx context.Context, userID int64) ([]user.UserDocument, error) {
	var documents []user.UserDocument
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("is_active = ?", true).
		Order("verified DESC, uploaded_at DESC").
		Find(&documents).Error
	return documents, err
}

// ===========================================
// FULL PROFILE WITH ALL RELATIONSHIPS
// ===========================================

// GetFullProfile retrieves a complete user profile with all relationships
func (r *userRepository) GetFullProfile(ctx context.Context, userID int64) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).
		Preload("Profile").
		Preload("Preference").
		Preload("Educations", func(db *gorm.DB) *gorm.DB {
			return db.Order("is_current DESC, end_year DESC NULLS FIRST, start_year DESC")
		}).
		Preload("Experiences", func(db *gorm.DB) *gorm.DB {
			return db.Order("is_current DESC, start_date DESC")
		}).
		Preload("Skills", func(db *gorm.DB) *gorm.DB {
			return db.Order("skill_level DESC, years_experience DESC NULLS LAST")
		}).
		Preload("Certifications", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = ?", true).Order("issue_date DESC NULLS LAST")
		}).
		Preload("Languages", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = ?", true).Order("proficiency_level DESC")
		}).
		Preload("Projects", func(db *gorm.DB) *gorm.DB {
			return db.Order("featured DESC, is_current DESC, start_date DESC NULLS LAST")
		}).
		Preload("Documents", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = ?", true).Order("verified DESC, uploaded_at DESC")
		}).
		First(&u, userID).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// ===========================================
// HELPER FUNCTIONS
// ===========================================

// toLowerSlice converts a string slice to lowercase
func toLowerSlice(strs []string) []string {
	result := make([]string, len(strs))
	for i, str := range strs {
		result[i] = strings.ToLower(str)
	}
	return result
}
