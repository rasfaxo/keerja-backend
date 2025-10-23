package seeders

import (
	"keerja-backend/internal/domain/master"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// BenefitsMasterSeeder seeds the benefits_master table
func BenefitsMasterSeeder(db *gorm.DB) error {
	log.Println("Seeding benefits_master table...")

	benefits := []master.BenefitsMaster{
		// Financial Benefits
		{Code: "COMPETITIVE_SALARY", Name: "Competitive Salary", Category: "financial", Description: "Above-market salary package", Icon: "💰", PopularityScore: 95, IsActive: true},
		{Code: "PERFORMANCE_BONUS", Name: "Performance Bonus", Category: "financial", Description: "Annual or quarterly performance-based bonuses", Icon: "🎯", PopularityScore: 90, IsActive: true},
		{Code: "STOCK_OPTIONS", Name: "Stock Options", Category: "financial", Description: "Employee stock ownership plan", Icon: "📈", PopularityScore: 85, IsActive: true},
		{Code: "ANNUAL_BONUS", Name: "Annual Bonus", Category: "financial", Description: "Yearly bonus package", Icon: "💵", PopularityScore: 88, IsActive: true},
		{Code: "PROFIT_SHARING", Name: "Profit Sharing", Category: "financial", Description: "Share in company profits", Icon: "🤝", PopularityScore: 82, IsActive: true},
		{Code: "SIGNING_BONUS", Name: "Signing Bonus", Category: "financial", Description: "One-time bonus upon joining", Icon: "✍️", PopularityScore: 75, IsActive: true},
		{Code: "RELOCATION_ASSISTANCE", Name: "Relocation Assistance", Category: "financial", Description: "Support for relocation expenses", Icon: "🏠", PopularityScore: 70, IsActive: true},
		{Code: "MEAL_ALLOWANCE", Name: "Meal Allowance", Category: "financial", Description: "Daily or monthly meal subsidy", Icon: "🍽️", PopularityScore: 80, IsActive: true},
		{Code: "TRANSPORT_ALLOWANCE", Name: "Transport Allowance", Category: "financial", Description: "Transportation or commute subsidy", Icon: "🚗", PopularityScore: 82, IsActive: true},
		{Code: "PHONE_ALLOWANCE", Name: "Phone Allowance", Category: "financial", Description: "Mobile phone and data plan allowance", Icon: "📱", PopularityScore: 75, IsActive: true},

		// Health Benefits
		{Code: "HEALTH_INSURANCE", Name: "Health Insurance", Category: "health", Description: "Comprehensive medical coverage", Icon: "🏥", PopularityScore: 98, IsActive: true},
		{Code: "DENTAL_INSURANCE", Name: "Dental Insurance", Category: "health", Description: "Dental care coverage", Icon: "🦷", PopularityScore: 85, IsActive: true},
		{Code: "VISION_INSURANCE", Name: "Vision Insurance", Category: "health", Description: "Eye care and glasses coverage", Icon: "👓", PopularityScore: 80, IsActive: true},
		{Code: "LIFE_INSURANCE", Name: "Life Insurance", Category: "health", Description: "Life insurance policy", Icon: "🛡️", PopularityScore: 88, IsActive: true},
		{Code: "MENTAL_HEALTH", Name: "Mental Health Support", Category: "health", Description: "Counseling and therapy services", Icon: "🧠", PopularityScore: 90, IsActive: true},
		{Code: "GYM_MEMBERSHIP", Name: "Gym Membership", Category: "health", Description: "Fitness center membership", Icon: "💪", PopularityScore: 75, IsActive: true},
		{Code: "WELLNESS_PROGRAM", Name: "Wellness Program", Category: "health", Description: "Health and wellness initiatives", Icon: "🌟", PopularityScore: 78, IsActive: true},
		{Code: "ANNUAL_CHECKUP", Name: "Annual Health Checkup", Category: "health", Description: "Comprehensive yearly medical examination", Icon: "🩺", PopularityScore: 82, IsActive: true},

		// Career Development
		{Code: "TRAINING_BUDGET", Name: "Training Budget", Category: "career", Description: "Professional development and training allowance", Icon: "📚", PopularityScore: 92, IsActive: true},
		{Code: "CONFERENCE_ATTENDANCE", Name: "Conference Attendance", Category: "career", Description: "Support for attending industry conferences", Icon: "🎤", PopularityScore: 85, IsActive: true},
		{Code: "CERTIFICATION_SUPPORT", Name: "Certification Support", Category: "career", Description: "Funding for professional certifications", Icon: "🎓", PopularityScore: 88, IsActive: true},
		{Code: "MENTORSHIP_PROGRAM", Name: "Mentorship Program", Category: "career", Description: "Access to mentors and career guidance", Icon: "🤝", PopularityScore: 80, IsActive: true},
		{Code: "CAREER_COACHING", Name: "Career Coaching", Category: "career", Description: "Professional career development coaching", Icon: "👨‍🏫", PopularityScore: 75, IsActive: true},
		{Code: "TUITION_REIMBURSEMENT", Name: "Tuition Reimbursement", Category: "career", Description: "Support for continuing education", Icon: "🎓", PopularityScore: 82, IsActive: true},
		{Code: "ONLINE_COURSES", Name: "Online Learning Platforms", Category: "career", Description: "Access to platforms like Udemy, Coursera, etc.", Icon: "💻", PopularityScore: 85, IsActive: true},

		// Lifestyle Benefits
		{Code: "PAID_TIME_OFF", Name: "Paid Time Off", Category: "lifestyle", Description: "Generous vacation days", Icon: "🏖️", PopularityScore: 95, IsActive: true},
		{Code: "PARENTAL_LEAVE", Name: "Parental Leave", Category: "lifestyle", Description: "Paid maternity/paternity leave", Icon: "👶", PopularityScore: 92, IsActive: true},
		{Code: "SABBATICAL_LEAVE", Name: "Sabbatical Leave", Category: "lifestyle", Description: "Extended leave for personal projects", Icon: "🌍", PopularityScore: 70, IsActive: true},
		{Code: "BIRTHDAY_LEAVE", Name: "Birthday Leave", Category: "lifestyle", Description: "Day off on your birthday", Icon: "🎂", PopularityScore: 75, IsActive: true},
		{Code: "COMPANY_EVENTS", Name: "Company Events", Category: "lifestyle", Description: "Team building and social activities", Icon: "🎉", PopularityScore: 80, IsActive: true},
		{Code: "FREE_SNACKS", Name: "Free Snacks & Drinks", Category: "lifestyle", Description: "Complimentary food and beverages", Icon: "☕", PopularityScore: 78, IsActive: true},
		{Code: "GAME_ROOM", Name: "Game Room", Category: "lifestyle", Description: "Recreation area with games", Icon: "🎮", PopularityScore: 65, IsActive: true},
		{Code: "PET_FRIENDLY", Name: "Pet-Friendly Office", Category: "lifestyle", Description: "Bring your pet to work", Icon: "🐕", PopularityScore: 68, IsActive: true},

		// Flexibility Benefits
		{Code: "REMOTE_WORK", Name: "Remote Work", Category: "flexibility", Description: "Work from home or anywhere", Icon: "🏡", PopularityScore: 96, IsActive: true},
		{Code: "HYBRID_WORK", Name: "Hybrid Work", Category: "flexibility", Description: "Mix of office and remote work", Icon: "🔄", PopularityScore: 94, IsActive: true},
		{Code: "FLEXIBLE_HOURS", Name: "Flexible Hours", Category: "flexibility", Description: "Choose your working hours", Icon: "⏰", PopularityScore: 93, IsActive: true},
		{Code: "FOUR_DAY_WEEK", Name: "4-Day Work Week", Category: "flexibility", Description: "Work 4 days instead of 5", Icon: "📅", PopularityScore: 85, IsActive: true},
		{Code: "UNLIMITED_PTO", Name: "Unlimited PTO", Category: "flexibility", Description: "Unlimited paid time off policy", Icon: "🌴", PopularityScore: 88, IsActive: true},
		{Code: "NO_OVERTIME", Name: "No Overtime Policy", Category: "flexibility", Description: "Strict work-life balance", Icon: "⛔", PopularityScore: 82, IsActive: true},
		{Code: "COMPRESSED_WORKWEEK", Name: "Compressed Workweek", Category: "flexibility", Description: "Longer days, shorter week", Icon: "📊", PopularityScore: 75, IsActive: true},

		// Other Benefits
		{Code: "LAPTOP_PROVIDED", Name: "Laptop Provided", Category: "other", Description: "Company-provided work laptop", Icon: "💻", PopularityScore: 90, IsActive: true},
		{Code: "EQUIPMENT_BUDGET", Name: "Equipment Budget", Category: "other", Description: "Budget for home office setup", Icon: "🖥️", PopularityScore: 82, IsActive: true},
		{Code: "PARKING_SPACE", Name: "Parking Space", Category: "other", Description: "Free parking at office", Icon: "🅿️", PopularityScore: 70, IsActive: true},
		{Code: "COMPANY_CAR", Name: "Company Car", Category: "other", Description: "Company-provided vehicle", Icon: "🚙", PopularityScore: 75, IsActive: true},
		{Code: "EMPLOYEE_DISCOUNT", Name: "Employee Discount", Category: "other", Description: "Discounts on company products/services", Icon: "🏷️", PopularityScore: 72, IsActive: true},
		{Code: "CHILDCARE_SUPPORT", Name: "Childcare Support", Category: "other", Description: "Childcare facilities or subsidies", Icon: "👨‍👩‍👧", PopularityScore: 80, IsActive: true},
		{Code: "RETIREMENT_PLAN", Name: "Retirement Plan", Category: "other", Description: "401k or pension contributions", Icon: "👴", PopularityScore: 88, IsActive: true},
		{Code: "EQUITY_GRANT", Name: "Equity Grant", Category: "other", Description: "Company equity allocation", Icon: "📊", PopularityScore: 85, IsActive: true},
		{Code: "REFERRAL_BONUS", Name: "Referral Bonus", Category: "other", Description: "Bonus for successful referrals", Icon: "👥", PopularityScore: 70, IsActive: true},
	}

	// Use OnConflict to update existing benefits or create new ones
	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "category", "description", "icon", "popularity_score", "is_active", "updated_at"}),
	}).Create(&benefits)

	if result.Error != nil {
		log.Printf("Failed to seed benefits: %v", result.Error)
		return result.Error
	}

	log.Printf("Successfully seeded %d benefits", len(benefits))
	return nil
}
