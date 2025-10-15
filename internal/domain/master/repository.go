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
	Search          string
	Category        string
	IsActive        *bool
	MinPopularity   *float64
	MaxPopularity   *float64
	Page            int
	PageSize        int
	SortBy          string
	SortOrder       string
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
	TotalBenefits    int64
	ActiveBenefits   int64
	InactiveBenefits int64
	ByCategory       map[string]int64
	AveragePopularity float64
	MostPopular      *BenefitsMaster
	TopCategories    []CategoryStat
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
	Category string
	Count    int64
	Percentage float64
}

// TypeStat contains statistics for a skill type
type TypeStat struct {
	Type       string
	Count      int64
	Percentage float64
}
