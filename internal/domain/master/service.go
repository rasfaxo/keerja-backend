package master

import (
	"context"
)

// BenefitsMasterService defines business logic for benefits master data management
type BenefitsMasterService interface {
	// Benefit Management
	CreateBenefit(ctx context.Context, req *CreateBenefitRequest) (*BenefitsMaster, error)
	UpdateBenefit(ctx context.Context, id int64, req *UpdateBenefitRequest) (*BenefitsMaster, error)
	DeleteBenefit(ctx context.Context, id int64) error
	GetBenefit(ctx context.Context, id int64) (*BenefitResponse, error)
	GetBenefitByCode(ctx context.Context, code string) (*BenefitResponse, error)
	GetBenefits(ctx context.Context, filter *BenefitsFilter) (*BenefitListResponse, error)
	SearchBenefits(ctx context.Context, query string, page, pageSize int) (*BenefitListResponse, error)

	// Category Operations
	GetBenefitsByCategory(ctx context.Context, category string) ([]BenefitResponse, error)
	GetCategories(ctx context.Context) ([]CategoryInfo, error)
	GetCategoryStats(ctx context.Context) ([]CategoryStatResponse, error)

	// Popularity Management
	UpdatePopularity(ctx context.Context, id int64, score float64) error
	IncrementPopularity(ctx context.Context, id int64) error
	GetMostPopular(ctx context.Context, limit int) ([]BenefitResponse, error)
	GetPopularByCategory(ctx context.Context, category string, limit int) ([]BenefitResponse, error)

	// Status Management
	ActivateBenefit(ctx context.Context, id int64) error
	DeactivateBenefit(ctx context.Context, id int64) error
	GetActiveBenefits(ctx context.Context) ([]BenefitResponse, error)

	// Bulk Operations
	BulkCreateBenefits(ctx context.Context, benefits []CreateBenefitRequest) ([]BenefitsMaster, error)
	BulkUpdatePopularity(ctx context.Context, updates map[int64]float64) error
	ImportBenefits(ctx context.Context, data []BenefitImportData) (*ImportResult, error)
	ExportBenefits(ctx context.Context, filter *BenefitsFilter) ([]BenefitExportData, error)

	// Statistics
	GetBenefitStats(ctx context.Context) (*BenefitStatsResponse, error)
	GetTrendingBenefits(ctx context.Context, limit int) ([]BenefitResponse, error)

	// Validation
	ValidateBenefit(ctx context.Context, req *CreateBenefitRequest) error
	CheckDuplicateBenefit(ctx context.Context, name, code string) (bool, error)
}

// SkillsMasterService defines business logic for skills master data management
type SkillsMasterService interface {
	// Skill Management
	CreateSkill(ctx context.Context, req *CreateSkillRequest) (*SkillsMaster, error)
	UpdateSkill(ctx context.Context, id int64, req *UpdateSkillRequest) (*SkillsMaster, error)
	DeleteSkill(ctx context.Context, id int64) error
	GetSkill(ctx context.Context, id int64) (*SkillResponse, error)
	GetSkillByCode(ctx context.Context, code string) (*SkillResponse, error)
	GetSkills(ctx context.Context, filter *SkillsFilter) (*SkillListResponse, error)
	SearchSkills(ctx context.Context, query string, page, pageSize int) (*SkillListResponse, error)

	// Type & Difficulty Operations
	GetSkillsByType(ctx context.Context, skillType string) ([]SkillResponse, error)
	GetSkillsByDifficulty(ctx context.Context, difficultyLevel string) ([]SkillResponse, error)
	GetSkillsByCategory(ctx context.Context, categoryID int64) ([]SkillResponse, error)
	GetTypeStats(ctx context.Context) ([]TypeStatResponse, error)
	GetDifficultyStats(ctx context.Context) ([]DifficultyStatResponse, error)

	// Batch Operations
	GetSkillsByIDs(ctx context.Context, ids []int64) ([]SkillResponse, error)
	GetSkillsByNames(ctx context.Context, names []string) ([]SkillResponse, error)

	// Hierarchy Operations
	GetRootSkills(ctx context.Context) ([]SkillResponse, error)
	GetChildSkills(ctx context.Context, parentID int64) ([]SkillResponse, error)
	GetParentSkill(ctx context.Context, childID int64) (*SkillResponse, error)
	GetSkillTree(ctx context.Context, rootID int64) (*SkillTreeResponse, error)
	SetParentSkill(ctx context.Context, skillID, parentID int64) error
	RemoveParentSkill(ctx context.Context, skillID int64) error

	// Popularity Management
	UpdatePopularity(ctx context.Context, id int64, score float64) error
	IncrementPopularity(ctx context.Context, id int64) error
	GetMostPopular(ctx context.Context, limit int) ([]SkillResponse, error)
	GetPopularByType(ctx context.Context, skillType string, limit int) ([]SkillResponse, error)

	// Alias Management
	AddAlias(ctx context.Context, skillID int64, alias string) error
	RemoveAlias(ctx context.Context, skillID int64, alias string) error
	UpdateAliases(ctx context.Context, skillID int64, aliases []string) error
	SearchByAlias(ctx context.Context, alias string) ([]SkillResponse, error)

	// Status Management
	ActivateSkill(ctx context.Context, id int64) error
	DeactivateSkill(ctx context.Context, id int64) error
	GetActiveSkills(ctx context.Context) ([]SkillResponse, error)

	// Bulk Operations
	BulkCreateSkills(ctx context.Context, skills []CreateSkillRequest) ([]SkillsMaster, error)
	BulkUpdatePopularity(ctx context.Context, updates map[int64]float64) error
	ImportSkills(ctx context.Context, data []SkillImportData) (*ImportResult, error)
	ExportSkills(ctx context.Context, filter *SkillsFilter) ([]SkillExportData, error)

	// Recommendations
	GetRelatedSkills(ctx context.Context, skillID int64, limit int) ([]SkillResponse, error)
	GetComplementarySkills(ctx context.Context, skillIDs []int64, limit int) ([]SkillResponse, error)
	GetSkillSuggestions(ctx context.Context, userSkills []int64, limit int) ([]SkillResponse, error)

	// Statistics
	GetSkillStats(ctx context.Context) (*SkillStatsResponse, error)
	GetTrendingSkills(ctx context.Context, limit int) ([]SkillResponse, error)

	// Validation
	ValidateSkill(ctx context.Context, req *CreateSkillRequest) error
	CheckDuplicateSkill(ctx context.Context, name, code string) (bool, error)
	ValidateHierarchy(ctx context.Context, skillID, parentID int64) error
}

// Request DTOs

// CreateBenefitRequest represents a request to create a benefit
type CreateBenefitRequest struct {
	Code        string `json:"code" validate:"required,min=2,max=50"`
	Name        string `json:"name" validate:"required,min=2,max=150"`
	Category    string `json:"category" validate:"required,oneof=financial health career lifestyle flexibility other"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty" validate:"omitempty,max=100"`
	IsActive    bool   `json:"is_active"`
}

// UpdateBenefitRequest represents a request to update a benefit
type UpdateBenefitRequest struct {
	Code        string `json:"code,omitempty" validate:"omitempty,min=2,max=50"`
	Name        string `json:"name,omitempty" validate:"omitempty,min=2,max=150"`
	Category    string `json:"category,omitempty" validate:"omitempty,oneof=financial health career lifestyle flexibility other"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty" validate:"omitempty,max=100"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

// CreateSkillRequest represents a request to create a skill
type CreateSkillRequest struct {
	Code            string   `json:"code,omitempty" validate:"omitempty,min=2,max=50"`
	Name            string   `json:"name" validate:"required,min=2,max=150"`
	NormalizedName  string   `json:"normalized_name,omitempty"`
	CategoryID      *int64   `json:"category_id,omitempty"`
	Description     string   `json:"description,omitempty"`
	SkillType       string   `json:"skill_type" validate:"required,oneof=technical soft language tool"`
	DifficultyLevel string   `json:"difficulty_level" validate:"required,oneof=beginner intermediate advanced"`
	Aliases         []string `json:"aliases,omitempty"`
	ParentID        *int64   `json:"parent_id,omitempty"`
	IsActive        bool     `json:"is_active"`
}

// UpdateSkillRequest represents a request to update a skill
type UpdateSkillRequest struct {
	Code            string   `json:"code,omitempty" validate:"omitempty,min=2,max=50"`
	Name            string   `json:"name,omitempty" validate:"omitempty,min=2,max=150"`
	NormalizedName  string   `json:"normalized_name,omitempty"`
	CategoryID      *int64   `json:"category_id,omitempty"`
	Description     string   `json:"description,omitempty"`
	SkillType       string   `json:"skill_type,omitempty" validate:"omitempty,oneof=technical soft language tool"`
	DifficultyLevel string   `json:"difficulty_level,omitempty" validate:"omitempty,oneof=beginner intermediate advanced"`
	Aliases         []string `json:"aliases,omitempty"`
	ParentID        *int64   `json:"parent_id,omitempty"`
	IsActive        *bool    `json:"is_active,omitempty"`
}

// BenefitImportData represents data for importing benefits
type BenefitImportData struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
	IsActive    bool   `json:"is_active"`
}

// SkillImportData represents data for importing skills
type SkillImportData struct {
	Code            string   `json:"code,omitempty"`
	Name            string   `json:"name"`
	NormalizedName  string   `json:"normalized_name,omitempty"`
	CategoryID      *int64   `json:"category_id,omitempty"`
	Description     string   `json:"description,omitempty"`
	SkillType       string   `json:"skill_type"`
	DifficultyLevel string   `json:"difficulty_level"`
	Aliases         []string `json:"aliases,omitempty"`
	ParentID        *int64   `json:"parent_id,omitempty"`
	IsActive        bool     `json:"is_active"`
}

// Response DTOs

// BenefitResponse represents a benefit response
type BenefitResponse struct {
	ID              int64   `json:"id"`
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	Category        string  `json:"category"`
	Description     string  `json:"description,omitempty"`
	Icon            string  `json:"icon,omitempty"`
	IsActive        bool    `json:"is_active"`
	PopularityScore float64 `json:"popularity_score"`
	UsageCount      int64   `json:"usage_count"`
}

// BenefitListResponse represents a paginated list of benefits
type BenefitListResponse struct {
	Benefits   []BenefitResponse `json:"benefits"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

// SkillResponse represents a skill response
type SkillResponse struct {
	ID              int64          `json:"id"`
	Code            string         `json:"code,omitempty"`
	Name            string         `json:"name"`
	NormalizedName  string         `json:"normalized_name,omitempty"`
	CategoryID      *int64         `json:"category_id,omitempty"`
	Description     string         `json:"description,omitempty"`
	SkillType       string         `json:"skill_type"`
	DifficultyLevel string         `json:"difficulty_level"`
	PopularityScore float64        `json:"popularity_score"`
	Aliases         []string       `json:"aliases,omitempty"`
	ParentID        *int64         `json:"parent_id,omitempty"`
	IsActive        bool           `json:"is_active"`
	UsageCount      int64          `json:"usage_count"`
	Parent          *SkillResponse `json:"parent,omitempty"`
	ChildrenCount   int64          `json:"children_count"`
}

// SkillListResponse represents a paginated list of skills
type SkillListResponse struct {
	Skills     []SkillResponse `json:"skills"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

// SkillTreeResponse represents a hierarchical skill tree
type SkillTreeResponse struct {
	Skill    *SkillResponse      `json:"skill"`
	Children []SkillTreeResponse `json:"children,omitempty"`
}

// BenefitStatsResponse represents statistics about benefits
type BenefitStatsResponse struct {
	TotalBenefits     int64                  `json:"total_benefits"`
	ActiveBenefits    int64                  `json:"active_benefits"`
	InactiveBenefits  int64                  `json:"inactive_benefits"`
	ByCategory        map[string]int64       `json:"by_category"`
	AveragePopularity float64                `json:"average_popularity"`
	MostPopular       *BenefitResponse       `json:"most_popular,omitempty"`
	TopCategories     []CategoryStatResponse `json:"top_categories"`
}

// SkillStatsResponse represents statistics about skills
type SkillStatsResponse struct {
	TotalSkills       int64                    `json:"total_skills"`
	ActiveSkills      int64                    `json:"active_skills"`
	InactiveSkills    int64                    `json:"inactive_skills"`
	ByType            map[string]int64         `json:"by_type"`
	ByDifficulty      map[string]int64         `json:"by_difficulty"`
	ByCategory        map[int64]int64          `json:"by_category"`
	RootSkills        int64                    `json:"root_skills"`
	ChildSkills       int64                    `json:"child_skills"`
	AveragePopularity float64                  `json:"average_popularity"`
	MostPopular       *SkillResponse           `json:"most_popular,omitempty"`
	TopTypes          []TypeStatResponse       `json:"top_types"`
	TopDifficulties   []DifficultyStatResponse `json:"top_difficulties"`
}

// CategoryInfo represents category information
type CategoryInfo struct {
	Category string `json:"category"`
	Label    string `json:"label"`
	Count    int64  `json:"count"`
}

// CategoryStatResponse represents statistics for a category
type CategoryStatResponse struct {
	Category   string  `json:"category"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

// TypeStatResponse represents statistics for a skill type
type TypeStatResponse struct {
	Type       string  `json:"type"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

// DifficultyStatResponse represents statistics for a difficulty level
type DifficultyStatResponse struct {
	Difficulty string  `json:"difficulty"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

// ImportResult represents the result of an import operation
type ImportResult struct {
	TotalRecords   int      `json:"total_records"`
	SuccessCount   int      `json:"success_count"`
	FailureCount   int      `json:"failure_count"`
	Errors         []string `json:"errors,omitempty"`
	DuplicateCount int      `json:"duplicate_count"`
}

// BenefitExportData represents data for exporting benefits
type BenefitExportData struct {
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	Category        string  `json:"category"`
	Description     string  `json:"description"`
	Icon            string  `json:"icon"`
	IsActive        bool    `json:"is_active"`
	PopularityScore float64 `json:"popularity_score"`
}

// SkillExportData represents data for exporting skills
type SkillExportData struct {
	Code            string   `json:"code"`
	Name            string   `json:"name"`
	NormalizedName  string   `json:"normalized_name"`
	CategoryID      *int64   `json:"category_id"`
	Description     string   `json:"description"`
	SkillType       string   `json:"skill_type"`
	DifficultyLevel string   `json:"difficulty_level"`
	PopularityScore float64  `json:"popularity_score"`
	Aliases         []string `json:"aliases"`
	ParentID        *int64   `json:"parent_id"`
	IsActive        bool     `json:"is_active"`
}

// ========================================
// Company Refactor Service Interfaces
// ========================================

// IndustryService defines business logic for industry master data
type IndustryService interface {
	// GetAll retrieves all industries with optional search
	GetAll(ctx context.Context, search string) ([]IndustryResponse, error)

	// GetActive retrieves all active industries with optional search
	GetActive(ctx context.Context, search string) ([]IndustryResponse, error)

	// GetByID retrieves an industry by ID
	GetByID(ctx context.Context, id int64) (*IndustryResponse, error)

	// ValidateIndustryID checks if an industry ID exists and is active
	ValidateIndustryID(ctx context.Context, id int64) error
}

// CompanySizeService defines business logic for company size master data
type CompanySizeService interface {
	// GetAll retrieves all company sizes
	GetAll(ctx context.Context) ([]CompanySizeResponse, error)

	// GetActive retrieves all active company sizes
	GetActive(ctx context.Context) ([]CompanySizeResponse, error)

	// GetByID retrieves a company size by ID
	GetByID(ctx context.Context, id int64) (*CompanySizeResponse, error)

	// ValidateCompanySizeID checks if a company size ID exists and is active
	ValidateCompanySizeID(ctx context.Context, id int64) error
}

// ProvinceService defines business logic for province master data
type ProvinceService interface {
	// GetAll retrieves all provinces with optional search
	GetAll(ctx context.Context, search string) ([]ProvinceResponse, error)

	// GetActive retrieves all active provinces with optional search
	GetActive(ctx context.Context, search string) ([]ProvinceResponse, error)

	// GetByID retrieves a province by ID
	GetByID(ctx context.Context, id int64) (*ProvinceResponse, error)

	// ValidateProvinceID checks if a province ID exists and is active
	ValidateProvinceID(ctx context.Context, id int64) error
}

// CityService defines business logic for city master data
type CityService interface {
	// GetByProvinceID retrieves all cities in a province with optional search
	GetByProvinceID(ctx context.Context, provinceID int64, search string) ([]CityResponse, error)

	// GetActiveByProvinceID retrieves all active cities in a province with optional search
	GetActiveByProvinceID(ctx context.Context, provinceID int64, search string) ([]CityResponse, error)

	// GetByID retrieves a city by ID with province info
	GetByID(ctx context.Context, id int64) (*CityResponse, error)

	// ValidateCityID checks if a city ID exists, is active, and belongs to the given province
	ValidateCityID(ctx context.Context, cityID, provinceID int64) error
}

// DistrictService defines business logic for district master data
type DistrictService interface {
	// GetByCityID retrieves all districts in a city with optional search
	GetByCityID(ctx context.Context, cityID int64, search string) ([]DistrictResponse, error)

	// GetActiveByCityID retrieves all active districts in a city with optional search
	GetActiveByCityID(ctx context.Context, cityID int64, search string) ([]DistrictResponse, error)

	// GetByID retrieves a district by ID with full location hierarchy
	GetByID(ctx context.Context, id int64) (*DistrictResponse, error)

	// ValidateDistrictID checks if a district ID exists, is active, and belongs to the given city
	ValidateDistrictID(ctx context.Context, districtID, cityID int64) error

	// ValidateLocationHierarchy validates the complete location hierarchy (province -> city -> district)
	ValidateLocationHierarchy(ctx context.Context, provinceID, cityID, districtID int64) error
}

// Response DTOs for Company Refactor

// IndustryResponse represents an industry response
type IndustryResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description,omitempty"`
	IconURL     string `json:"icon_url,omitempty"`
	IsActive    bool   `json:"is_active"`
}

// CompanySizeResponse represents a company size response
type CompanySizeResponse struct {
	ID           int64  `json:"id"`
	Label        string `json:"label"`
	MinEmployees int    `json:"min_employees"`
	MaxEmployees *int   `json:"max_employees,omitempty"` // nil = unlimited
	IsActive     bool   `json:"is_active"`
}

// ProvinceResponse represents a province response
type ProvinceResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	IsActive bool   `json:"is_active"`
}

// CityResponse represents a city response
type CityResponse struct {
	ID         int64             `json:"id"`
	Name       string            `json:"name"`
	FullName   string            `json:"full_name"` // e.g., "Kota Bandung"
	Type       string            `json:"type"`      // "Kota" or "Kabupaten"
	Code       string            `json:"code"`
	ProvinceID int64             `json:"province_id"`
	Province   *ProvinceResponse `json:"province,omitempty"`
	IsActive   bool              `json:"is_active"`
}

// DistrictResponse represents a district response
type DistrictResponse struct {
	ID               int64         `json:"id"`
	Name             string        `json:"name"`
	Code             string        `json:"code"`
	PostalCode       string        `json:"postal_code,omitempty"`
	CityID           int64         `json:"city_id"`
	City             *CityResponse `json:"city,omitempty"`
	FullLocationPath string        `json:"full_location_path,omitempty"` // e.g., "Batujajar, Kabupaten Bandung Barat, Jawa Barat"
	IsActive         bool          `json:"is_active"`
}

// JobTitleService defines business logic for job title master data
type JobTitleService interface {
	// Phase 1: GetJobTitles with smart search
	SearchJobTitles(ctx context.Context, query string, limit int) ([]JobTitleResponse, error)
	GetJobTitle(ctx context.Context, id int64) (*JobTitleResponse, error)
	ListPopularJobTitles(ctx context.Context, limit int) ([]JobTitleResponse, error)

	// Admin operations
	CreateJobTitle(ctx context.Context, req *CreateJobTitleRequest) (*JobTitle, error)
	UpdateJobTitle(ctx context.Context, id int64, req *UpdateJobTitleRequest) (*JobTitle, error)
	DeleteJobTitle(ctx context.Context, id int64) error
}

// JobOptionsService defines business logic for job options (static master data)
type JobOptionsService interface {
	// Phase 3: GetJobOptions - combined response with caching
	GetJobOptions(ctx context.Context) (*JobOptionsResponse, error)

	// Individual getters (rarely used, mostly for admin)
	GetJobTypes(ctx context.Context) ([]JobType, error)
	GetWorkPolicies(ctx context.Context) ([]WorkPolicy, error)
	GetEducationLevels(ctx context.Context) ([]EducationLevel, error)
	GetExperienceLevels(ctx context.Context) ([]ExperienceLevel, error)
	GetGenderPreferences(ctx context.Context) ([]GenderPreference, error)
}

// ===== Request DTOs =====

// CreateJobTitleRequest for creating job title
type CreateJobTitleRequest struct {
	Name                  string  `json:"name" validate:"required,min=2,max=200"`
	RecommendedCategoryID *int64  `json:"recommended_category_id,omitempty"`
	PopularityScore       float64 `json:"popularity_score" validate:"omitempty,min=0,max=100"`
}

// UpdateJobTitleRequest for updating job title
type UpdateJobTitleRequest struct {
	Name                  *string  `json:"name,omitempty" validate:"omitempty,min=2,max=200"`
	RecommendedCategoryID *int64   `json:"recommended_category_id,omitempty"`
	PopularityScore       *float64 `json:"popularity_score,omitempty" validate:"omitempty,min=0,max=100"`
	IsActive              *bool    `json:"is_active,omitempty"`
}

// ===== Response DTOs =====

// JobTitleResponse for job title with recommendations
type JobTitleResponse struct {
	ID                    int64   `json:"id"`
	Name                  string  `json:"name"`
	RecommendedCategoryID *int64  `json:"rekomendasi_kategori_id,omitempty"` // Follow naming in spec
	PopularityScore       float64 `json:"popularity_score"`
	SearchCount           int64   `json:"search_count"`
}
