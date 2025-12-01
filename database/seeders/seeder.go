package seeders

import (
	"log"

	"gorm.io/gorm"
)

// RunSeeders runs all seeders in the correct order
func RunSeeders(db *gorm.DB) error {
	log.Println("Starting database seeding...")

	// Phase 1: Location data (provinces, cities, districts)
	if err := LocationSeeder(db); err != nil {
		log.Printf("Failed to run LocationSeeder: %v", err)
		return err
	}

	// Phase 2: Industries and Company Sizes
	if err := IndustriesSeeder(db); err != nil {
		log.Printf("Failed to run IndustriesSeeder: %v", err)
		return err
	}

	if err := CompanySizesSeeder(db); err != nil {
		log.Printf("Failed to run CompanySizesSeeder: %v", err)
		return err
	}

	// Phase 3: Job master data (job_types, work_policies, education_levels, experience_levels, gender_preferences, job_titles)
	if err := JobMasterDataSeeder(db); err != nil {
		log.Printf("Failed to run JobMasterDataSeeder: %v", err)
		return err
	}

	// Phase 4: Skills and Benefits
	if err := SkillsMasterSeeder(db); err != nil {
		log.Printf("Failed to run SkillsMasterSeeder: %v", err)
		return err
	}

	if err := BenefitsMasterSeeder(db); err != nil {
		log.Printf("Failed to run BenefitsMasterSeeder: %v", err)
		return err
	}

	// Phase 5: Job categories and subcategories
	if err := JobCategoriesSeeder(db); err != nil {
		log.Printf("Failed to run JobCategoriesSeeder: %v", err)
		return err
	}

	if err := JobSubcategoriesSeeder(db); err != nil {
		log.Printf("Failed to run JobSubcategoriesSeeder: %v", err)
		return err
	}

	// Phase 6: Admin roles and users
	if err := AdminRolesSeeder(db); err != nil {
		log.Printf("Failed to run AdminRolesSeeder: %v", err)
		return err
	}

	if err := AdminUserSeeder(db); err != nil {
		log.Printf("Failed to run AdminUserSeeder: %v", err)
		return err
	}

	// Phase 7: Companies and employer users (depends on location & industries)
	if err := CompaniesSeeder(db); err != nil {
		log.Printf("Failed to run CompaniesSeeder: %v", err)
		return err
	}

	if err := EmployerUsersSeeder(db); err != nil {
		log.Printf("Failed to run EmployerUsersSeeder: %v", err)
		return err
	}

	log.Println("Database seeding completed successfully")
	return nil
}
