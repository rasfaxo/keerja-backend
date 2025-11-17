package seeders

import (
	"keerja-backend/internal/domain/job"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// JobCategoriesSeeder seeds the job_categories table
func JobCategoriesSeeder(db *gorm.DB) error {
	log.Println("Seeding job_categories table...")

	categories := []job.JobCategory{
		// Technology & IT
		{Code: "TECH", Name: "Technology & IT", Description: "Information Technology, Software Development, and related fields", IsActive: true},
		{Code: "SOFTWARE_DEV", Name: "Software Development", Description: "Software engineering and development roles", IsActive: true},
		{Code: "DATA_SCIENCE", Name: "Data Science & Analytics", Description: "Data analysis, machine learning, and AI", IsActive: true},
		{Code: "DEVOPS", Name: "DevOps & Infrastructure", Description: "System administration, cloud, and infrastructure", IsActive: true},
		{Code: "CYBERSECURITY", Name: "Cybersecurity", Description: "Information security and cybersecurity roles", IsActive: true},
		{Code: "QA_TESTING", Name: "QA & Testing", Description: "Quality assurance and software testing", IsActive: true},
		{Code: "IT_SUPPORT", Name: "IT Support", Description: "Technical support and helpdesk", IsActive: true},
		{Code: "UI_UX", Name: "UI/UX Design", Description: "User interface and user experience design", IsActive: true},

		// Business & Management
		{Code: "BUSINESS", Name: "Business & Management", Description: "Business operations and management roles", IsActive: true},
		{Code: "PROJECT_MGMT", Name: "Project Management", Description: "Project and program management", IsActive: true},
		{Code: "PRODUCT_MGMT", Name: "Product Management", Description: "Product management and strategy", IsActive: true},
		{Code: "BUSINESS_ANALYSIS", Name: "Business Analysis", Description: "Business analysis and consulting", IsActive: true},
		{Code: "OPERATIONS", Name: "Operations", Description: "Business operations and process management", IsActive: true},
		{Code: "STRATEGY", Name: "Strategy & Consulting", Description: "Business strategy and consulting", IsActive: true},

		// Sales & Marketing
		{Code: "SALES_MARKETING", Name: "Sales & Marketing", Description: "Sales, marketing, and business development", IsActive: true},
		{Code: "SALES", Name: "Sales", Description: "Sales and business development roles", IsActive: true},
		{Code: "MARKETING", Name: "Marketing", Description: "Marketing and brand management", IsActive: true},
		{Code: "DIGITAL_MARKETING", Name: "Digital Marketing", Description: "Online marketing and social media", IsActive: true},
		{Code: "CONTENT_CREATION", Name: "Content Creation", Description: "Content writing and creation", IsActive: true},
		{Code: "SEO_SEM", Name: "SEO & SEM", Description: "Search engine optimization and marketing", IsActive: true},

		// Finance & Accounting
		{Code: "FINANCE", Name: "Finance & Accounting", Description: "Financial services and accounting", IsActive: true},
		{Code: "ACCOUNTING", Name: "Accounting", Description: "Accounting and bookkeeping", IsActive: true},
		{Code: "FINANCIAL_ANALYSIS", Name: "Financial Analysis", Description: "Financial planning and analysis", IsActive: true},
		{Code: "AUDIT", Name: "Audit & Compliance", Description: "Internal audit and compliance", IsActive: true},
		{Code: "TAX", Name: "Tax", Description: "Tax planning and compliance", IsActive: true},

		// Human Resources
		{Code: "HR", Name: "Human Resources", Description: "HR management and talent acquisition", IsActive: true},
		{Code: "RECRUITMENT", Name: "Recruitment", Description: "Talent acquisition and recruitment", IsActive: true},
		{Code: "HR_OPERATIONS", Name: "HR Operations", Description: "HR administration and operations", IsActive: true},
		{Code: "LEARNING_DEV", Name: "Learning & Development", Description: "Training and development", IsActive: true},
		{Code: "COMPENSATION", Name: "Compensation & Benefits", Description: "Compensation and benefits management", IsActive: true},

		// Creative & Design
		{Code: "CREATIVE", Name: "Creative & Design", Description: "Creative and design roles", IsActive: true},
		{Code: "GRAPHIC_DESIGN", Name: "Graphic Design", Description: "Visual and graphic design", IsActive: true},
		{Code: "VIDEO_PRODUCTION", Name: "Video Production", Description: "Video editing and production", IsActive: true},
		{Code: "ANIMATION", Name: "Animation", Description: "2D/3D animation and motion graphics", IsActive: true},
		{Code: "PHOTOGRAPHY", Name: "Photography", Description: "Professional photography", IsActive: true},

		// Engineering
		{Code: "ENGINEERING", Name: "Engineering", Description: "Engineering roles across disciplines", IsActive: true},
		{Code: "MECHANICAL", Name: "Mechanical Engineering", Description: "Mechanical engineering roles", IsActive: true},
		{Code: "ELECTRICAL", Name: "Electrical Engineering", Description: "Electrical engineering roles", IsActive: true},
		{Code: "CIVIL", Name: "Civil Engineering", Description: "Civil engineering and construction", IsActive: true},
		{Code: "INDUSTRIAL", Name: "Industrial Engineering", Description: "Industrial engineering and process optimization", IsActive: true},

		// Healthcare
		{Code: "HEALTHCARE", Name: "Healthcare", Description: "Medical and healthcare services", IsActive: true},
		{Code: "MEDICAL", Name: "Medical Professionals", Description: "Doctors, nurses, and medical staff", IsActive: true},
		{Code: "PHARMACY", Name: "Pharmacy", Description: "Pharmaceutical services", IsActive: true},
		{Code: "HEALTH_ADMIN", Name: "Healthcare Administration", Description: "Healthcare management and administration", IsActive: true},

		// Education
		{Code: "EDUCATION", Name: "Education", Description: "Teaching and educational services", IsActive: true},
		{Code: "TEACHING", Name: "Teaching", Description: "Teaching and instruction", IsActive: true},
		{Code: "TRAINING", Name: "Corporate Training", Description: "Corporate training and facilitation", IsActive: true},
		{Code: "EDUCATION_ADMIN", Name: "Education Administration", Description: "Educational management and administration", IsActive: true},

		// Customer Service
		{Code: "CUSTOMER_SERVICE", Name: "Customer Service", Description: "Customer support and service", IsActive: true},
		{Code: "CALL_CENTER", Name: "Call Center", Description: "Call center and phone support", IsActive: true},
		{Code: "CLIENT_SUCCESS", Name: "Client Success", Description: "Customer success and account management", IsActive: true},

		// Legal
		{Code: "LEGAL", Name: "Legal", Description: "Legal services and compliance", IsActive: true},
		{Code: "CORPORATE_LAW", Name: "Corporate Law", Description: "Corporate legal services", IsActive: true},
		{Code: "LEGAL_COMPLIANCE", Name: "Legal Compliance", Description: "Legal compliance and regulatory", IsActive: true},

		// Supply Chain & Logistics
		{Code: "SUPPLY_CHAIN", Name: "Supply Chain & Logistics", Description: "Supply chain management and logistics", IsActive: true},
		{Code: "PROCUREMENT", Name: "Procurement", Description: "Purchasing and procurement", IsActive: true},
		{Code: "WAREHOUSE", Name: "Warehouse & Distribution", Description: "Warehouse management and distribution", IsActive: true},
		{Code: "LOGISTICS", Name: "Logistics", Description: "Logistics and transportation", IsActive: true},

		// Manufacturing & Production
		{Code: "MANUFACTURING", Name: "Manufacturing & Production", Description: "Manufacturing and production operations", IsActive: true},
		{Code: "PRODUCTION", Name: "Production", Description: "Production management", IsActive: true},
		{Code: "QUALITY_CONTROL", Name: "Quality Control", Description: "Quality assurance and control", IsActive: true},

		// Hospitality & Tourism
		{Code: "HOSPITALITY", Name: "Hospitality & Tourism", Description: "Hotel, restaurant, and tourism services", IsActive: true},
		{Code: "HOTEL_MGMT", Name: "Hotel Management", Description: "Hotel operations and management", IsActive: true},
		{Code: "FOOD_BEVERAGE", Name: "Food & Beverage", Description: "F&B services and management", IsActive: true},

		// Other
		{Code: "ADMIN", Name: "Administrative", Description: "Administrative and clerical roles", IsActive: true},
		{Code: "GENERAL_LABOR", Name: "General Labor", Description: "General labor and operational roles", IsActive: true},
		{Code: "INTERNSHIP", Name: "Internship", Description: "Internship and trainee positions", IsActive: true},
		{Code: "FREELANCE", Name: "Freelance & Contract", Description: "Freelance and contract work", IsActive: true},
	}

	// Use OnConflict to update existing categories or create new ones
	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "description", "is_active", "updated_at"}),
	}).Create(&categories)

	if result.Error != nil {
		log.Printf("Failed to seed job categories: %v", result.Error)
		return result.Error
	}

	log.Printf("Successfully seeded %d job categories", len(categories))
	return nil
}
