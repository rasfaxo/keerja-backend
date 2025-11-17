package seeders

import (
	"log"

	"keerja-backend/internal/domain/job"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// JobSubcategoriesSeeder seeds the job_subcategories table and links them to job_categories
func JobSubcategoriesSeeder(db *gorm.DB) error {
	log.Println("Seeding job_subcategories table...")

	// Define subcategories grouped by category code
	subcatsByCategory := map[string][]job.JobSubcategory{
		// Technology & IT
		"TECH": {
			{Code: "SOFT_DEV", Name: "Software Development", Description: "Software engineering and development roles", IsActive: true},
			{Code: "FRONTEND", Name: "Frontend Engineer", Description: "Frontend web/mobile engineers (React, Vue, Angular)", IsActive: true},
			{Code: "BACKEND", Name: "Backend Engineer", Description: "Backend systems, APIs, and services", IsActive: true},
			{Code: "FULLSTACK", Name: "Fullstack Engineer", Description: "Fullstack web/mobile engineers", IsActive: true},
			{Code: "MOBILE", Name: "Mobile Engineer", Description: "iOS/Android native and cross-platform", IsActive: true},
		},
		"DATA_SCIENCE": {
			{Code: "DATA_ANALYTICS", Name: "Data Analyst", Description: "Data analysis and BI", IsActive: true},
			{Code: "ML_ENGINEER", Name: "Machine Learning Engineer", Description: "ML engineering and models", IsActive: true},
			{Code: "DATA_ENGINEER", Name: "Data Engineer", Description: "Data pipelines and ETL", IsActive: true},
		},
		"DEVOPS": {
			{Code: "SRE", Name: "Site Reliability Engineer", Description: "SRE and platform engineering", IsActive: true},
			{Code: "CLOUD_ENGINEER", Name: "Cloud Engineer", Description: "Cloud infrastructure and automation", IsActive: true},
		},
		"QA_TESTING": {
			{Code: "QA_MANUAL", Name: "Manual QA", Description: "Manual testing roles", IsActive: true},
			{Code: "QA_AUTOMATION", Name: "Automation QA", Description: "Test automation and frameworks", IsActive: true},
		},

		// Business & Management
		"PROJECT_MGMT": {
			{Code: "PROJECT_MANAGER", Name: "Project Manager", Description: "Project and program management", IsActive: true},
		},
		"PRODUCT_MGMT": {
			{Code: "PRODUCT_MANAGER", Name: "Product Manager", Description: "Product management and strategy", IsActive: true},
		},

		// Sales & Marketing
		"SALES": {
			{Code: "SALES_EXEC", Name: "Sales Executive", Description: "Field and inside sales", IsActive: true},
			{Code: "ACC_MANAGER", Name: "Account Manager", Description: "Customer/account management", IsActive: true},
		},
		"DIGITAL_MARKETING": {
			{Code: "DIGITAL_MARKETING_SPECIALIST", Name: "Digital Marketing Specialist", Description: "SEM, social ads, performance marketing", IsActive: true},
			{Code: "SEO_SPECIALIST", Name: "SEO Specialist", Description: "Search engine optimisation", IsActive: true},
			{Code: "CONTENT_WRITER", Name: "Content Writer", Description: "Content creation and copywriting", IsActive: true},
		},

		// Finance & Accounting
		"ACCOUNTING": {
			{Code: "ACCOUNTANT", Name: "Accountant", Description: "General accounting and bookkeeping", IsActive: true},
		},
		"FINANCIAL_ANALYSIS": {
			{Code: "FIN_ANALYST", Name: "Financial Analyst", Description: "Financial planning and analysis", IsActive: true},
		},

		// Human Resources
		"RECRUITMENT": {
			{Code: "RECRUITER", Name: "Recruiter", Description: "Talent acquisition and sourcing", IsActive: true},
		},

		// Creative & Design
		"GRAPHIC_DESIGN": {
			{Code: "GRAPHIC_DESIGNER", Name: "Graphic Designer", Description: "Visual and graphic design", IsActive: true},
		},
		"VIDEO_PRODUCTION": {
			{Code: "VIDEO_EDITOR", Name: "Video Editor", Description: "Video editing and production", IsActive: true},
		},

		// Engineering
		"MECHANICAL": {
			{Code: "MECHANICAL_ENGINEER", Name: "Mechanical Engineer", Description: "Mechanical design and development", IsActive: true},
		},
		"ELECTRICAL": {
			{Code: "ELECTRICAL_ENGINEER", Name: "Electrical Engineer", Description: "Electrical systems and design", IsActive: true},
		},

		// Healthcare
		"MEDICAL": {
			{Code: "DOCTOR", Name: "Doctor", Description: "Physicians and medical professionals", IsActive: true},
			{Code: "NURSE", Name: "Nurse", Description: "Nursing and clinical staff", IsActive: true},
		},

		// Education
		"TEACHING": {
			{Code: "SCHOOL_TEACHER", Name: "School Teacher", Description: "Primary/secondary education", IsActive: true},
			{Code: "CORPORATE_TRAINER", Name: "Corporate Trainer", Description: "Training and facilitation", IsActive: true},
		},

		// Customer Service
		"CUSTOMER_SERVICE": {
			{Code: "CUSTOMER_SUPPORT", Name: "Customer Support", Description: "Customer service and support roles", IsActive: true},
		},

		// Supply Chain & Logistics
		"LOGISTICS": {
			{Code: "LOGISTICS_COORD", Name: "Logistics Coordinator", Description: "Logistics and transportation coordination", IsActive: true},
		},

		// Manufacturing & Production
		"PRODUCTION": {
			{Code: "PRODUCTION_OPERATOR", Name: "Production Operator", Description: "Manufacturing floor roles", IsActive: true},
		},

		// Hospitality & Tourism
		"HOSPITALITY": {
			{Code: "HOTEL_MANAGER", Name: "Hotel Manager", Description: "Hotel operations and management", IsActive: true},
		},

		// Other / Administrative
		"ADMIN": {
			{Code: "ADMIN_ASSISTANT", Name: "Administrative Assistant", Description: "Admin and clerical roles", IsActive: true},
		},
		"INTERNSHIP": {
			{Code: "INTERNSHIP_POSITION", Name: "Internship", Description: "Internship and trainee positions", IsActive: true},
		},
		"FREELANCE": {
			{Code: "FREELANCE_CONTRACT", Name: "Freelance / Contract", Description: "Freelance and contract work", IsActive: true},
		},
	}

	var toCreate []job.JobSubcategory

	// For each category code, find the category and attach subcategories
	for code, subs := range subcatsByCategory {
		var cat job.JobCategory
		if err := db.WithContext(db.Statement.Context).Where("code = ?", code).First(&cat).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Printf("Category with code '%s' not found, skipping its subcategories", code)
				continue
			}
			log.Printf("Failed to lookup category '%s': %v", code, err)
			return err
		}

		for _, s := range subs {
			s.CategoryID = cat.ID
			toCreate = append(toCreate, s)
		}
	}

	if len(toCreate) == 0 {
		log.Println("No job subcategories to seed (no matching categories found)")
		return nil
	}

	// Upsert by code
	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "description", "is_active", "category_id", "updated_at"}),
	}).Create(&toCreate)

	if result.Error != nil {
		log.Printf("Failed to seed job subcategories: %v", result.Error)
		return result.Error
	}

	log.Printf("Successfully seeded %d job subcategories", len(toCreate))
	return nil
}
