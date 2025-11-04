package seeders

import (
	"keerja-backend/internal/domain/admin"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AdminRolesSeeder seeds default admin roles
func AdminRolesSeeder(db *gorm.DB) error {
	log.Println("Seeding admin_roles table...")

	roles := []admin.AdminRole{
		{
			RoleName:        "Super Admin",
			RoleDescription: "Full system access with all permissions. Can manage all users, roles, and system settings.",
			AccessLevel:     10,
			IsSystemRole:    true,
		},
		{
			RoleName:        "Admin",
			RoleDescription: "High-level administrative access. Can manage users, content, and most system features.",
			AccessLevel:     8,
			IsSystemRole:    true,
		},
		{
			RoleName:        "Content Manager",
			RoleDescription: "Manages content including jobs, companies, and user-generated content. Can approve/reject submissions.",
			AccessLevel:     7,
			IsSystemRole:    true,
		},
		{
			RoleName:        "Moderator",
			RoleDescription: "Moderate user content, reviews, and reports. Can flag inappropriate content.",
			AccessLevel:     6,
			IsSystemRole:    true,
		},
		{
			RoleName:        "Support Agent",
			RoleDescription: "Handle user support tickets and inquiries. Limited administrative access.",
			AccessLevel:     4,
			IsSystemRole:    true,
		},
		{
			RoleName:        "Analyst",
			RoleDescription: "View-only access to analytics and reports. No modification permissions.",
			AccessLevel:     3,
			IsSystemRole:    true,
		},
		{
			RoleName:        "Viewer",
			RoleDescription: "Read-only access to system data. No modification or administrative permissions.",
			AccessLevel:     1,
			IsSystemRole:    true,
		},
	}

	// Use OnConflict to update existing roles or create new ones
	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "role_name"}},
		DoUpdates: clause.AssignmentColumns([]string{"role_description", "access_level", "is_system_role", "updated_at"}),
	}).Create(&roles)

	if result.Error != nil {
		log.Printf("Failed to seed admin roles: %v", result.Error)
		return result.Error
	}

	log.Printf("Successfully seeded %d admin roles", len(roles))
	return nil
}
