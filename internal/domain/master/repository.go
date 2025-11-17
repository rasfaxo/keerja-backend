package master

import (
	"context"
)

// BenefitsMasterRepository defines data access methods for BenefitsMaster
type BenefitsMasterRepository interface {
	// Basic CRUD
	Create(ctx context.Context, benefit *BenefitsMaster) error
	FindByID(ctx context.Context, id int64) (*BenefitsMaster, error)
	FindByCode(ctx context.Context, code string) (*BenefitsMaster, error)
	FindByName(ctx context.Context, name string) (*BenefitsMaster, error)
	Update(ctx context.Context, benefit *BenefitsMaster) error
	Delete(ctx context.Context, id int64) error

	// Listing & Search
	List(ctx context.Context, filter *BenefitsFilter) ([]BenefitsMaster, int64, error)
	ListActive(ctx context.Context) ([]BenefitsMaster, error)
	ListByCategory(ctx context.Context, category string) ([]BenefitsMaster, error)
	SearchBenefits(ctx context.Context, query string, page, pageSize int) ([]BenefitsMaster, int64, error)

	// Category Operations
	GetCategories(ctx context.Context) ([]string, error)
	CountByCategory(ctx context.Context, category string) (int64, error)

	// Popularity Operations
	UpdatePopularity(ctx context.Context, id int64, score float64) error
	IncrementPopularity(ctx context.Context, id int64, amount float64) error
	GetMostPopular(ctx context.Context, limit int) ([]BenefitsMaster, error)
	GetByPopularityRange(ctx context.Context, minScore, maxScore float64) ([]BenefitsMaster, error)

	// Status Operations
	Activate(ctx context.Context, id int64) error
	Deactivate(ctx context.Context, id int64) error

	// Bulk Operations
	BulkCreate(ctx context.Context, benefits []BenefitsMaster) error
	BulkUpdatePopularity(ctx context.Context, updates map[int64]float64) error

	// Statistics
	Count(ctx context.Context) (int64, error)
	CountActive(ctx context.Context) (int64, error)
	GetBenefitStats(ctx context.Context) (*BenefitStats, error)
}

// SkillsMasterRepository defines data access methods for SkillsMaster
type SkillsMasterRepository interface {
	// Basic CRUD
	Create(ctx context.Context, skill *SkillsMaster) error
	FindByID(ctx context.Context, id int64) (*SkillsMaster, error)
	FindByCode(ctx context.Context, code string) (*SkillsMaster, error)
	FindByName(ctx context.Context, name string) (*SkillsMaster, error)
	FindByNormalizedName(ctx context.Context, normalizedName string) (*SkillsMaster, error)
	Update(ctx context.Context, skill *SkillsMaster) error
	Delete(ctx context.Context, id int64) error

	// Listing & Search
	List(ctx context.Context, filter *SkillsFilter) ([]SkillsMaster, int64, error)
	ListActive(ctx context.Context) ([]SkillsMaster, error)
	ListByType(ctx context.Context, skillType string) ([]SkillsMaster, error)
	ListByDifficulty(ctx context.Context, difficultyLevel string) ([]SkillsMaster, error)
	ListByCategory(ctx context.Context, categoryID int64) ([]SkillsMaster, error)
	SearchSkills(ctx context.Context, query string, page, pageSize int) ([]SkillsMaster, int64, error)
	SearchByAlias(ctx context.Context, alias string) ([]SkillsMaster, error)

	// Hierarchy Operations
	GetRootSkills(ctx context.Context) ([]SkillsMaster, error)
	GetChildren(ctx context.Context, parentID int64) ([]SkillsMaster, error)
	GetParent(ctx context.Context, childID int64) (*SkillsMaster, error)
	GetSkillTree(ctx context.Context, rootID int64) (*SkillsMaster, error)
	HasChildren(ctx context.Context, id int64) (bool, error)

	// Popularity Operations
	UpdatePopularity(ctx context.Context, id int64, score float64) error
	IncrementPopularity(ctx context.Context, id int64, amount float64) error
	GetMostPopular(ctx context.Context, limit int) ([]SkillsMaster, error)
	GetByPopularityRange(ctx context.Context, minScore, maxScore float64) ([]SkillsMaster, error)
	GetPopularByType(ctx context.Context, skillType string, limit int) ([]SkillsMaster, error)

	// Alias Operations
	AddAlias(ctx context.Context, id int64, alias string) error
	RemoveAlias(ctx context.Context, id int64, alias string) error
	UpdateAliases(ctx context.Context, id int64, aliases []string) error
	FindByAliases(ctx context.Context, aliases []string) ([]SkillsMaster, error)

	// Status Operations
	Activate(ctx context.Context, id int64) error
	Deactivate(ctx context.Context, id int64) error

	// Bulk Operations
	BulkCreate(ctx context.Context, skills []SkillsMaster) error
	BulkUpdatePopularity(ctx context.Context, updates map[int64]float64) error

	// Statistics
	Count(ctx context.Context) (int64, error)
	CountActive(ctx context.Context) (int64, error)
	CountByType(ctx context.Context, skillType string) (int64, error)
	CountByDifficulty(ctx context.Context, difficultyLevel string) (int64, error)
	GetSkillStats(ctx context.Context) (*SkillStats, error)

	// Recommendations
	GetRelatedSkills(ctx context.Context, skillID int64, limit int) ([]SkillsMaster, error)
	GetComplementarySkills(ctx context.Context, skillIDs []int64, limit int) ([]SkillsMaster, error)
}

// BenefitsFilter defines filter options for benefits queries
type BenefitsFilter struct {
	Search        string
	Category      string
	IsActive      *bool
	MinPopularity *float64
	MaxPopularity *float64
	Page          int
	PageSize      int
	SortBy        string
	SortOrder     string
}

// SkillsFilter defines filter options for skills queries
type SkillsFilter struct {
	Search          string
	SkillType       string
	DifficultyLevel string
	CategoryID      *int64
	ParentID        *int64
	IsActive        *bool
	MinPopularity   *float64
	MaxPopularity   *float64
	HasParent       *bool
	HasChildren     *bool
	Page            int
	PageSize        int
	SortBy          string
	SortOrder       string
}

// BenefitStats contains statistics about benefits
type BenefitStats struct {
	TotalBenefits     int64
	ActiveBenefits    int64
	InactiveBenefits  int64
	ByCategory        map[string]int64
	AveragePopularity float64
	MostPopular       *BenefitsMaster
	TopCategories     []CategoryStat
}

// SkillStats contains statistics about skills
type SkillStats struct {
	TotalSkills       int64
	ActiveSkills      int64
	InactiveSkills    int64
	ByType            map[string]int64
	ByDifficulty      map[string]int64
	ByCategory        map[int64]int64
	RootSkills        int64
	ChildSkills       int64
	AveragePopularity float64
	MostPopular       *SkillsMaster
	TopTypes          []TypeStat
}

// CategoryStat contains statistics for a benefit category
type CategoryStat struct {
	Category   string
	Count      int64
	Percentage float64
}

// TypeStat contains statistics for a skill type
type TypeStat struct {
	Type       string
	Count      int64
	Percentage float64
}

// ========================================
// Company Refactor Repository Interfaces
// ========================================

// IndustryRepository defines the interface for industry data access
type IndustryRepository interface {
	// GetAll retrieves all industries
	GetAll(ctx context.Context) ([]Industry, error)

	// GetActive retrieves all active industries
	GetActive(ctx context.Context) ([]Industry, error)

	// GetByID retrieves an industry by ID
	GetByID(ctx context.Context, id int64) (*Industry, error)

	// GetBySlug retrieves an industry by slug
	GetBySlug(ctx context.Context, slug string) (*Industry, error)

	// GetByName retrieves an industry by exact name (case-insensitive)
	GetByName(ctx context.Context, name string) (*Industry, error)

	// Search searches industries by name
	Search(ctx context.Context, query string) ([]Industry, error)

	// Create creates a new industry
	Create(ctx context.Context, industry *Industry) error

	// Update updates an existing industry
	Update(ctx context.Context, industry *Industry) error

	// Delete soft deletes an industry
	Delete(ctx context.Context, id int64) error

	// ExistsByID checks if an industry exists by ID
	ExistsByID(ctx context.Context, id int64) (bool, error)
}

// CompanySizeRepository defines the interface for company size data access
type CompanySizeRepository interface {
	// GetAll retrieves all company sizes
	GetAll(ctx context.Context) ([]CompanySize, error)

	// GetActive retrieves all active company sizes
	GetActive(ctx context.Context) ([]CompanySize, error)

	// GetByID retrieves a company size by ID
	GetByID(ctx context.Context, id int64) (*CompanySize, error)

	// GetByCategory retrieves a company size by category name (exact match, case-insensitive)
	GetByCategory(ctx context.Context, category string) (*CompanySize, error)

	// Create creates a new company size
	Create(ctx context.Context, size *CompanySize) error

	// Update updates an existing company size
	Update(ctx context.Context, size *CompanySize) error

	// Delete deletes a company size
	Delete(ctx context.Context, id int64) error

	// ExistsByID checks if a company size exists by ID
	ExistsByID(ctx context.Context, id int64) (bool, error)
}

// ProvinceRepository defines the interface for province data access
type ProvinceRepository interface {
	// GetAll retrieves all provinces
	GetAll(ctx context.Context) ([]Province, error)

	// GetActive retrieves all active provinces
	GetActive(ctx context.Context) ([]Province, error)

	// GetByID retrieves a province by ID
	GetByID(ctx context.Context, id int64) (*Province, error)

	// GetByCode retrieves a province by code
	GetByCode(ctx context.Context, code string) (*Province, error)

	// GetByName retrieves a province by exact name (case-insensitive)
	GetByName(ctx context.Context, name string) (*Province, error)

	// Search searches provinces by name
	Search(ctx context.Context, query string) ([]Province, error)

	// Create creates a new province
	Create(ctx context.Context, province *Province) error

	// Update updates an existing province
	Update(ctx context.Context, province *Province) error

	// Delete deletes a province
	Delete(ctx context.Context, id int64) error

	// ExistsByID checks if a province exists by ID
	ExistsByID(ctx context.Context, id int64) (bool, error)
}

// CityRepository defines the interface for city data access
type CityRepository interface {
	// GetAll retrieves all cities
	GetAll(ctx context.Context) ([]City, error)

	// GetActive retrieves all active cities
	GetActive(ctx context.Context) ([]City, error)

	// GetByID retrieves a city by ID
	GetByID(ctx context.Context, id int64) (*City, error)

	// GetByCode retrieves a city by code
	GetByCode(ctx context.Context, code string) (*City, error)

	// GetByNameAndProvinceID retrieves a city by exact name and province ID (case-insensitive)
	GetByNameAndProvinceID(ctx context.Context, name string, provinceID int64) (*City, error)

	// GetByProvinceID retrieves all cities in a province
	GetByProvinceID(ctx context.Context, provinceID int64) ([]City, error)

	// GetActiveByProvinceID retrieves all active cities in a province
	GetActiveByProvinceID(ctx context.Context, provinceID int64) ([]City, error)

	// Search searches cities by name within a province (optional)
	Search(ctx context.Context, query string, provinceID *int64) ([]City, error)

	// GetWithProvince retrieves a city with its province preloaded
	GetWithProvince(ctx context.Context, id int64) (*City, error)

	// Create creates a new city
	Create(ctx context.Context, city *City) error

	// Update updates an existing city
	Update(ctx context.Context, city *City) error

	// Delete deletes a city
	Delete(ctx context.Context, id int64) error

	// ExistsByID checks if a city exists by ID
	ExistsByID(ctx context.Context, id int64) (bool, error)
}

// DistrictRepository defines the interface for district data access
type DistrictRepository interface {
	// GetAll retrieves all districts
	GetAll(ctx context.Context) ([]District, error)

	// GetActive retrieves all active districts
	GetActive(ctx context.Context) ([]District, error)

	// GetByID retrieves a district by ID
	GetByID(ctx context.Context, id int64) (*District, error)

	// GetByCode retrieves a district by code
	GetByCode(ctx context.Context, code string) (*District, error)

	// GetByNameAndCityID retrieves a district by exact name and city ID (case-insensitive)
	GetByNameAndCityID(ctx context.Context, name string, cityID int64) (*District, error)

	// GetByCityID retrieves all districts in a city
	GetByCityID(ctx context.Context, cityID int64) ([]District, error)

	// GetActiveByCityID retrieves all active districts in a city
	GetActiveByCityID(ctx context.Context, cityID int64) ([]District, error)

	// Search searches districts by name within a city (optional)
	Search(ctx context.Context, query string, cityID *int64) ([]District, error)

	// GetWithFullLocation retrieves a district with city and province preloaded
	GetWithFullLocation(ctx context.Context, id int64) (*District, error)

	// GetByPostalCode retrieves districts by postal code
	GetByPostalCode(ctx context.Context, postalCode string) ([]District, error)

	// Create creates a new district
	Create(ctx context.Context, district *District) error

	// Update updates an existing district
	Update(ctx context.Context, district *District) error

	// Delete deletes a district
	Delete(ctx context.Context, id int64) error

	// ExistsByID checks if a district exists by ID
	ExistsByID(ctx context.Context, id int64) (bool, error)
}

// JobTitleRepository defines data access methods for JobTitle
type JobTitleRepository interface {
	// Basic CRUD
	Create(ctx context.Context, jobTitle *JobTitle) error
	FindByID(ctx context.Context, id int64) (*JobTitle, error)
	FindByName(ctx context.Context, name string) (*JobTitle, error)
	Update(ctx context.Context, jobTitle *JobTitle) error
	Delete(ctx context.Context, id int64) error

	// Smart Search with fuzzy matching
	SearchJobTitles(ctx context.Context, query string, limit int) ([]JobTitle, error)

	// Listing
	ListActive(ctx context.Context) ([]JobTitle, error)
	ListPopular(ctx context.Context, limit int) ([]JobTitle, error)

	// Statistics
	IncrementSearchCount(ctx context.Context, id int64) error
	UpdatePopularity(ctx context.Context, id int64, score float64) error
}

// JobOptionsRepository defines data access methods for job options (static data)
type JobOptionsRepository interface {
	// Get all options at once for caching
	GetAllJobTypes(ctx context.Context) ([]JobType, error)
	GetAllWorkPolicies(ctx context.Context) ([]WorkPolicy, error)
	GetAllEducationLevels(ctx context.Context) ([]EducationLevel, error)
	GetAllExperienceLevels(ctx context.Context) ([]ExperienceLevel, error)
	GetAllGenderPreferences(ctx context.Context) ([]GenderPreference, error)

	// Get combined options (for caching efficiency)
	GetJobOptions(ctx context.Context) (*JobOptionsResponse, error)

	// Individual lookups
	FindJobTypeByID(ctx context.Context, id int64) (*JobType, error)
	FindWorkPolicyByID(ctx context.Context, id int64) (*WorkPolicy, error)
	FindEducationLevelByID(ctx context.Context, id int64) (*EducationLevel, error)
	FindExperienceLevelByID(ctx context.Context, id int64) (*ExperienceLevel, error)
	FindGenderPreferenceByID(ctx context.Context, id int64) (*GenderPreference, error)
}
