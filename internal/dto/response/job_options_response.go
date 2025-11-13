package response

import "keerja-backend/internal/domain/master"

// JobTypesOptionsResponse represents response for GET /api/v1/jobs/job-types
type JobTypesOptionsResponse struct {
	JobTypes             []master.JobType      `json:"job_types"`
	WorkPolicies         []master.WorkPolicy   `json:"work_policies"`
	WorkAddresses        []WorkAddressOption   `json:"work_addresses"`
	SalaryRanges         []SalaryRangeOption   `json:"salary_ranges"`
	SalaryDisplayOptions []SalaryDisplayOption `json:"salary_display_options"`
	SalaryInfo           string                `json:"salary_info"`
}

// WorkAddressOption represents a work address option
type WorkAddressOption struct {
	ID          int64           `json:"id"`
	CompanyID   int64           `json:"company_id"`
	CompanyName string          `json:"company_name"`
	FullAddress string          `json:"full_address"`
	Location    *LocationDetail `json:"location,omitempty"`
	Latitude    *float64        `json:"latitude,omitempty"`
	Longitude   *float64        `json:"longitude,omitempty"`
}

// LocationDetail represents location details
type LocationDetail struct {
	District string `json:"district,omitempty"`
	City     string `json:"city,omitempty"`
	Province string `json:"province,omitempty"`
}

// SalaryRangeOption represents a salary range option
type SalaryRangeOption struct {
	Label    string `json:"label"`
	MinValue int64  `json:"min_value"`
	MaxValue int64  `json:"max_value"` // 0 means unlimited
}

// SalaryDisplayOption represents how salary should be displayed
type SalaryDisplayOption struct {
	Value       string `json:"value"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

// JobRequirementsOptionsResponse represents response for GET /api/v1/jobs/job-requirements
type JobRequirementsOptionsResponse struct {
	Genders          []master.GenderPreference `json:"genders"`
	AgeRanges        []AgeRangeOption          `json:"age_ranges"`
	EducationLevels  []master.EducationLevel   `json:"education_levels"`
	ExperienceLevels []master.ExperienceLevel  `json:"experience_levels"`
	Skills           []SkillOption             `json:"skills"`
	SkillsTotal      int64                     `json:"skills_total"`
	AgeInfo          string                    `json:"age_info"`
	SkillsNote       string                    `json:"skills_note"`
}

// AgeRangeOption represents an age range option
type AgeRangeOption struct {
	Label string `json:"label"`
	Min   *int   `json:"min"`
	Max   *int   `json:"max"`
}

// SkillOption represents a skill option for selection
type SkillOption struct {
	ID              int64   `json:"id"`
	Code            string  `json:"code,omitempty"`
	Name            string  `json:"name"`
	NormalizedName  string  `json:"normalized_name,omitempty"`
	SkillType       string  `json:"skill_type"`
	DifficultyLevel string  `json:"difficulty_level"`
	PopularityScore float64 `json:"popularity_score"`
	Description     string  `json:"description,omitempty"`
}
