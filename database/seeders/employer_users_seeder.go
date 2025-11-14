package seeders

import (
	"keerja-backend/internal/domain/company"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// EmployerUsersSeeder seeds employer users for companies
func EmployerUsersSeeder(db *gorm.DB) error {
	log.Println("Seeding employer_users table...")

	// Get all companies
	var companies []company.Company
	if err := db.Find(&companies).Error; err != nil {
		log.Printf("Failed to fetch companies: %v", err)
		return err
	}

	if len(companies) == 0 {
		log.Println("No companies found, skipping employer_users seeding.")
		return nil
	}

	// Get available users from users table
	var users []struct {
		ID int64
	}
	if err := db.Table("users").Select("id").Find(&users).Error; err != nil {
		log.Printf("Failed to fetch users: %v", err)
		return err
	}

	if len(users) == 0 {
		log.Println("No users found, skipping employer_users seeding.")
		return nil
	}

	now := time.Now()
	employerUsers := make([]company.EmployerUser, 0)

	// Create employer users for each company, mapping to available users
	for i, comp := range companies {
		// Cycle through available users if companies > users
		userIndex := i % len(users)
		userID := users[userIndex].ID

		employerUser := company.EmployerUser{
			UserID:    userID,
			CompanyID: comp.ID,
			Role:      "admin",
			IsActive:  true,
			CreatedAt: now,
			UpdatedAt: now,
		}
		employerUsers = append(employerUsers, employerUser)
	}

	if len(employerUsers) == 0 {
		log.Println("No employer users to seed.")
		return nil
	}

	// Use OnConflict to handle duplicates gracefully — target unique (user_id, company_id)
	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "company_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"role", "position_title", "department", "email_company", "phone_company", "is_verified", "verified_at", "verified_by", "is_active", "last_login", "updated_at"}),
	}).Create(&employerUsers)

	if result.Error != nil {
		log.Printf("Failed to seed employer_users: %v", result.Error)
		return result.Error
	}

	log.Printf("Successfully seeded %d employer users", len(employerUsers))
	return nil
}
