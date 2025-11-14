package seeders

import (
	"log"

	"gorm.io/gorm"
)

// RunSeeders runs all seeders
func RunSeeders(db *gorm.DB) error {
	log.Println("Starting database seeding...")

	// Run job master data seeder (Phase 1-4: job_types, work_policies, education_levels, experience_levels, gender_preferences, job_titles)
	if err := JobMasterDataSeeder(db); err != nil {
		log.Printf("Failed to run JobMasterDataSeeder: %v", err)
		return err
	}

	// Run skills master seeder
	if err := SkillsMasterSeeder(db); err != nil {
		log.Printf("Failed to run SkillsMasterSeeder: %v", err)
		return err
	}

	// Run benefits master seeder
	if err := BenefitsMasterSeeder(db); err != nil {
		log.Printf("Failed to run BenefitsMasterSeeder: %v", err)
		return err
	}

	// Run job categories seeder
	if err := JobCategoriesSeeder(db); err != nil {
		log.Printf("Failed to run JobCategoriesSeeder: %v", err)
		return err
	}

	// Run job subcategories seeder (depends on job_categories)
	if err := JobSubcategoriesSeeder(db); err != nil {
		log.Printf("Failed to run JobSubcategoriesSeeder: %v", err)
		return err
	}

	// Run admin roles seeder
	if err := AdminRolesSeeder(db); err != nil {
		log.Printf("Failed to run AdminRolesSeeder: %v", err)
		return err
	}

	// Run admin user seeder
	if err := AdminUserSeeder(db); err != nil {
		log.Printf("Failed to run AdminUserSeeder: %v", err)
		return err
	}

	// Run companies seeder
	if err := CompaniesSeeder(db); err != nil {
		log.Printf("Failed to run CompaniesSeeder: %v", err)
		return err
	}

	// Run employer users seeder
	if err := EmployerUsersSeeder(db); err != nil {
		log.Printf("Failed to run EmployerUsersSeeder: %v", err)
		return err
	}

	log.Println("Database seeding completed successfully")
	return nil
}
