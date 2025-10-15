package master

import (
	"time"

	"gorm.io/gorm"
)

// BenefitsMaster represents a master data entry for job benefits
// Maps to: benefits_master table
type BenefitsMaster struct {
	ID              int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Code            string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"code" validate:"required,min=2,max=50"`
	Name            string         `gorm:"type:varchar(150);not null;uniqueIndex" json:"name" validate:"required,min=2,max=150"`
	Category        string         `gorm:"type:varchar(50);default:'other'" json:"category" validate:"oneof=financial health career lifestyle flexibility other"`
	Description     string         `gorm:"type:text" json:"description,omitempty"`
	Icon            string         `gorm:"type:varchar(100)" json:"icon,omitempty"`
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	PopularityScore float64        `gorm:"type:numeric(5,2);default:0.00;index:idx_benefits_master_popularity,sort:desc" json:"popularity_score" validate:"min=0,max=100"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for BenefitsMaster
func (BenefitsMaster) TableName() string {
	return "benefits_master"
}

// IsFinancial checks if the benefit is in the financial category
func (b *BenefitsMaster) IsFinancial() bool {
	return b.Category == "financial"
}

// IsHealth checks if the benefit is in the health category
func (b *BenefitsMaster) IsHealth() bool {
	return b.Category == "health"
}

// IsCareer checks if the benefit is in the career category
func (b *BenefitsMaster) IsCareer() bool {
	return b.Category == "career"
}

// IsLifestyle checks if the benefit is in the lifestyle category
func (b *BenefitsMaster) IsLifestyle() bool {
	return b.Category == "lifestyle"
}

// IsFlexibility checks if the benefit is in the flexibility category
func (b *BenefitsMaster) IsFlexibility() bool {
	return b.Category == "flexibility"
}

// IsPopular checks if the benefit has high popularity score (>= 70)
func (b *BenefitsMaster) IsPopular() bool {
	return b.PopularityScore >= 70.0
}

// IncrementPopularity increases the popularity score
func (b *BenefitsMaster) IncrementPopularity(amount float64) {
	b.PopularityScore += amount
	if b.PopularityScore > 100.0 {
		b.PopularityScore = 100.0
	}
}

// SkillsMaster represents a master data entry for skills
// Maps to: skills_master table
type SkillsMaster struct {
	ID              int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Code            string         `gorm:"type:varchar(50);uniqueIndex" json:"code,omitempty" validate:"omitempty,min=2,max=50"`
	Name            string         `gorm:"type:varchar(150);not null;uniqueIndex" json:"name" validate:"required,min=2,max=150"`
	NormalizedName  string         `gorm:"type:varchar(150)" json:"normalized_name,omitempty"`
	CategoryID      *int64         `gorm:"index" json:"category_id,omitempty"`
	Description     string         `gorm:"type:text" json:"description,omitempty"`
	SkillType       string         `gorm:"type:varchar(30);default:'technical'" json:"skill_type" validate:"oneof=technical soft language tool"`
	DifficultyLevel string         `gorm:"type:varchar(20);default:'intermediate'" json:"difficulty_level" validate:"oneof=beginner intermediate advanced"`
	PopularityScore float64        `gorm:"type:numeric(5,2);default:0.00;index:idx_skills_master_popularity,sort:desc" json:"popularity_score" validate:"min=0,max=100"`
	Aliases         []string       `gorm:"type:text[]" json:"aliases,omitempty"`
	ParentID        *int64         `gorm:"index" json:"parent_id,omitempty"`
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	// Note: Category references job_categories table (not included here to avoid circular dependency)
	// In implementation, use: Category *job.JobCategory `gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL"`
	Parent   *SkillsMaster   `gorm:"foreignKey:ParentID;constraint:OnDelete:SET NULL" json:"parent,omitempty"`
	Children []SkillsMaster  `gorm:"foreignKey:ParentID;constraint:OnDelete:SET NULL" json:"children,omitempty"`
}

// TableName specifies the table name for SkillsMaster
func (SkillsMaster) TableName() string {
	return "skills_master"
}

// IsTechnical checks if the skill is a technical skill
func (s *SkillsMaster) IsTechnical() bool {
	return s.SkillType == "technical"
}

// IsSoft checks if the skill is a soft skill
func (s *SkillsMaster) IsSoft() bool {
	return s.SkillType == "soft"
}

// IsLanguage checks if the skill is a language skill
func (s *SkillsMaster) IsLanguage() bool {
	return s.SkillType == "language"
}

// IsTool checks if the skill is a tool/software skill
func (s *SkillsMaster) IsTool() bool {
	return s.SkillType == "tool"
}

// IsBeginner checks if the skill has beginner difficulty level
func (s *SkillsMaster) IsBeginner() bool {
	return s.DifficultyLevel == "beginner"
}

// IsIntermediate checks if the skill has intermediate difficulty level
func (s *SkillsMaster) IsIntermediate() bool {
	return s.DifficultyLevel == "intermediate"
}

// IsAdvanced checks if the skill has advanced difficulty level
func (s *SkillsMaster) IsAdvanced() bool {
	return s.DifficultyLevel == "advanced"
}

// IsPopular checks if the skill has high popularity score (>= 70)
func (s *SkillsMaster) IsPopular() bool {
	return s.PopularityScore >= 70.0
}

// HasParent checks if the skill has a parent skill
func (s *SkillsMaster) HasParent() bool {
	return s.ParentID != nil
}

// HasChildren checks if the skill has child skills
func (s *SkillsMaster) HasChildren() bool {
	return len(s.Children) > 0
}

// IncrementPopularity increases the popularity score
func (s *SkillsMaster) IncrementPopularity(amount float64) {
	s.PopularityScore += amount
	if s.PopularityScore > 100.0 {
		s.PopularityScore = 100.0
	}
}

// AddAlias adds a new alias to the skill
func (s *SkillsMaster) AddAlias(alias string) {
	for _, existing := range s.Aliases {
		if existing == alias {
			return
		}
	}
	s.Aliases = append(s.Aliases, alias)
}

// RemoveAlias removes an alias from the skill
func (s *SkillsMaster) RemoveAlias(alias string) {
	for i, existing := range s.Aliases {
		if existing == alias {
			s.Aliases = append(s.Aliases[:i], s.Aliases[i+1:]...)
			return
		}
	}
}

// HasAlias checks if the skill has a specific alias
func (s *SkillsMaster) HasAlias(alias string) bool {
	for _, existing := range s.Aliases {
		if existing == alias {
			return true
		}
	}
	return false
}
