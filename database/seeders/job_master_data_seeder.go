package seeders

import (
	"keerja-backend/internal/domain/master"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// JobMasterDataSeeder seeds all job-related master data tables
func JobMasterDataSeeder(db *gorm.DB) error {
	if err := SeedJobTypes(db); err != nil {
		return err
	}
	if err := SeedWorkPolicies(db); err != nil {
		return err
	}
	if err := SeedEducationLevels(db); err != nil {
		return err
	}
	if err := SeedExperienceLevels(db); err != nil {
		return err
	}
	if err := SeedGenderPreferences(db); err != nil {
		return err
	}
	if err := SeedJobTitles(db); err != nil {
		return err
	}
	return nil
}

// SeedJobTypes seeds the job_types table
func SeedJobTypes(db *gorm.DB) error {
	log.Println("Seeding job_types table...")

	jobTypes := []master.JobType{
		{Name: "Full-Time", Code: "full_time", Order: 1},
		{Name: "Part-Time", Code: "part_time", Order: 2},
		{Name: "Internship", Code: "internship", Order: 3},
		{Name: "Freelance", Code: "freelance", Order: 4},
		{Name: "Contract", Code: "contract", Order: 5},
	}

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "order"}),
	}).Create(&jobTypes)

	if result.Error != nil {
		log.Printf("Error seeding job_types: %v", result.Error)
		return result.Error
	}

	log.Printf("Successfully seeded %d job types", len(jobTypes))
	return nil
}

// SeedWorkPolicies seeds the work_policies table
func SeedWorkPolicies(db *gorm.DB) error {
	log.Println("Seeding work_policies table...")

	workPolicies := []master.WorkPolicy{
		{Name: "On-site", Code: "onsite", Order: 1},
		{Name: "Remote", Code: "remote", Order: 2},
		{Name: "Hybrid", Code: "hybrid", Order: 3},
	}

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "order"}),
	}).Create(&workPolicies)

	if result.Error != nil {
		log.Printf("Error seeding work_policies: %v", result.Error)
		return result.Error
	}

	log.Printf("Successfully seeded %d work policies", len(workPolicies))
	return nil
}

// SeedEducationLevels seeds the education_levels table
func SeedEducationLevels(db *gorm.DB) error {
	log.Println("Seeding education_levels table...")

	// First, clear existing data to avoid constraint conflicts
	if err := db.Exec("DELETE FROM education_levels").Error; err != nil {
		log.Printf("Error clearing education_levels: %v", err)
		return err
	}

	educationLevels := []master.EducationLevel{
		{ID: 1, Name: "SMA/SMK", Code: "sma", Order: 1},
		{ID: 2, Name: "D3", Code: "d3", Order: 2},
		{ID: 3, Name: "D4", Code: "d4", Order: 3},
		{ID: 4, Name: "S1", Code: "s1", Order: 4},
		{ID: 5, Name: "S2", Code: "s2", Order: 5},
		{ID: 6, Name: "S3", Code: "s3", Order: 6},
		{ID: 7, Name: "Tidak Ditentukan", Code: "any", Order: 7},
	}

	result := db.Create(&educationLevels)

	if result.Error != nil {
		log.Printf("Error seeding education_levels: %v", result.Error)
		return result.Error
	}

	log.Printf("Successfully seeded %d education levels", len(educationLevels))
	return nil
}

// SeedExperienceLevels seeds the experience_levels table
func SeedExperienceLevels(db *gorm.DB) error {
	log.Println("Seeding experience_levels table...")

	experienceLevels := []master.ExperienceLevel{
		{Name: "Fresh Graduate", Code: "fresh", MinYears: 0, MaxYears: intPtr(0), Order: 1},
		{Name: "1-2 Tahun", Code: "junior", MinYears: 1, MaxYears: intPtr(2), Order: 2},
		{Name: "3-5 Tahun", Code: "mid", MinYears: 3, MaxYears: intPtr(5), Order: 3},
		{Name: "6-10 Tahun", Code: "senior", MinYears: 6, MaxYears: intPtr(10), Order: 4},
		{Name: "Lebih dari 10 Tahun", Code: "expert", MinYears: 10, MaxYears: nil, Order: 5},
		{Name: "Tidak Ditentukan", Code: "any", MinYears: 0, MaxYears: nil, Order: 6},
	}

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "min_years", "max_years", "order"}),
	}).Create(&experienceLevels)

	if result.Error != nil {
		log.Printf("Error seeding experience_levels: %v", result.Error)
		return result.Error
	}

	log.Printf("Successfully seeded %d experience levels", len(experienceLevels))
	return nil
}

// SeedGenderPreferences seeds the gender_preferences table
func SeedGenderPreferences(db *gorm.DB) error {
	log.Println("Seeding gender_preferences table...")

	genderPreferences := []master.GenderPreference{
		{Name: "Laki-laki", Code: "male", Order: 1},
		{Name: "Perempuan", Code: "female", Order: 2},
		{Name: "Semua", Code: "any", Order: 3},
	}

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "order"}),
	}).Create(&genderPreferences)

	if result.Error != nil {
		log.Printf("Error seeding gender_preferences: %v", result.Error)
		return result.Error
	}

	log.Printf("Successfully seeded %d gender preferences", len(genderPreferences))
	return nil
}

// SeedJobTitles seeds the job_titles table with common job titles
func SeedJobTitles(db *gorm.DB) error {
	log.Println("Seeding job_titles table...")

	jobTitles := []master.JobTitle{
		// Technology & Software Development
		{Name: "Software Engineer", NormalizedName: "software engineer", PopularityScore: 100, SearchCount: 1200, IsActive: true},
		{Name: "Frontend Developer", NormalizedName: "frontend developer", PopularityScore: 95, SearchCount: 950, IsActive: true},
		{Name: "Backend Developer", NormalizedName: "backend developer", PopularityScore: 95, SearchCount: 1000, IsActive: true},
		{Name: "Full Stack Developer", NormalizedName: "full stack developer", PopularityScore: 98, SearchCount: 1100, IsActive: true},
		{Name: "Mobile Developer", NormalizedName: "mobile developer", PopularityScore: 90, SearchCount: 800, IsActive: true},
		{Name: "DevOps Engineer", NormalizedName: "devops engineer", PopularityScore: 85, SearchCount: 700, IsActive: true},
		{Name: "Data Scientist", NormalizedName: "data scientist", PopularityScore: 92, SearchCount: 850, IsActive: true},
		{Name: "Data Analyst", NormalizedName: "data analyst", PopularityScore: 88, SearchCount: 780, IsActive: true},
		{Name: "UI/UX Designer", NormalizedName: "ui ux designer", PopularityScore: 87, SearchCount: 820, IsActive: true},
		{Name: "QA Engineer", NormalizedName: "qa engineer", PopularityScore: 80, SearchCount: 650, IsActive: true},

		// Business & Management
		{Name: "Product Manager", NormalizedName: "product manager", PopularityScore: 93, SearchCount: 920, IsActive: true},
		{Name: "Project Manager", NormalizedName: "project manager", PopularityScore: 90, SearchCount: 900, IsActive: true},
		{Name: "Business Analyst", NormalizedName: "business analyst", PopularityScore: 85, SearchCount: 810, IsActive: true},
		{Name: "Operations Manager", NormalizedName: "operations manager", PopularityScore: 82, SearchCount: 750, IsActive: true},

		// Sales & Marketing
		{Name: "Sales Executive", NormalizedName: "sales executive", PopularityScore: 88, SearchCount: 880, IsActive: true},
		{Name: "Marketing Manager", NormalizedName: "marketing manager", PopularityScore: 86, SearchCount: 860, IsActive: true},
		{Name: "Digital Marketing Specialist", NormalizedName: "digital marketing specialist", PopularityScore: 90, SearchCount: 1050, IsActive: true},
		{Name: "Content Writer", NormalizedName: "content writer", PopularityScore: 83, SearchCount: 750, IsActive: true},
		{Name: "Social Media Manager", NormalizedName: "social media manager", PopularityScore: 85, SearchCount: 820, IsActive: true},

		// Finance & Accounting
		{Name: "Accountant", NormalizedName: "accountant", PopularityScore: 87, SearchCount: 870, IsActive: true},
		{Name: "Financial Analyst", NormalizedName: "financial analyst", PopularityScore: 84, SearchCount: 780, IsActive: true},
		{Name: "Tax Specialist", NormalizedName: "tax specialist", PopularityScore: 78, SearchCount: 600, IsActive: true},

		// Human Resources
		{Name: "HR Manager", NormalizedName: "hr manager", PopularityScore: 85, SearchCount: 850, IsActive: true},
		{Name: "Recruiter", NormalizedName: "recruiter", PopularityScore: 82, SearchCount: 760, IsActive: true},
		{Name: "HR Generalist", NormalizedName: "hr generalist", PopularityScore: 80, SearchCount: 720, IsActive: true},

		// Customer Service
		{Name: "Customer Service Representative", NormalizedName: "customer service representative", PopularityScore: 85, SearchCount: 810, IsActive: true},
		{Name: "Customer Success Manager", NormalizedName: "customer success manager", PopularityScore: 83, SearchCount: 780, IsActive: true},

		// Administrative
		{Name: "Administrative Assistant", NormalizedName: "administrative assistant", PopularityScore: 80, SearchCount: 720, IsActive: true},
		{Name: "Executive Assistant", NormalizedName: "executive assistant", PopularityScore: 82, SearchCount: 750, IsActive: true},
		{Name: "Office Manager", NormalizedName: "office manager", PopularityScore: 78, SearchCount: 680, IsActive: true},
	}

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"normalized_name", "popularity_score", "search_count"}),
	}).Create(&jobTitles)

	if result.Error != nil {
		log.Printf("Error seeding job_titles: %v", result.Error)
		return result.Error
	}

	log.Printf("Successfully seeded %d job titles", len(jobTitles))
	return nil
}

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
}
