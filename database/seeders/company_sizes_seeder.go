package seeders

import (
	"log"

	"gorm.io/gorm"
)

// CompanySize represents a company size category
type CompanySize struct {
	ID           int64  `gorm:"primaryKey"`
	Label        string `gorm:"uniqueIndex;not null"`
	MinEmployees int    `gorm:"not null"`
	MaxEmployees *int   // nullable for 1000+
	DisplayOrder int    `gorm:"default:0"`
}

func (CompanySize) TableName() string {
	return "company_sizes"
}

// CompanySizesSeeder seeds the company_sizes table
func CompanySizesSeeder(db *gorm.DB) error {
	log.Println("Seeding company_sizes table...")

	// Helper for max employees pointer
	maxPtr := func(i int) *int { return &i }

	companySizes := []CompanySize{
		{Label: "1 - 10 karyawan", MinEmployees: 1, MaxEmployees: maxPtr(10), DisplayOrder: 1},
		{Label: "11 - 50 karyawan", MinEmployees: 11, MaxEmployees: maxPtr(50), DisplayOrder: 2},
		{Label: "51 - 200 karyawan", MinEmployees: 51, MaxEmployees: maxPtr(200), DisplayOrder: 3},
		{Label: "201 - 500 karyawan", MinEmployees: 201, MaxEmployees: maxPtr(500), DisplayOrder: 4},
		{Label: "501 - 1000 karyawan", MinEmployees: 501, MaxEmployees: maxPtr(1000), DisplayOrder: 5},
		{Label: "1000+ karyawan", MinEmployees: 1001, MaxEmployees: nil, DisplayOrder: 6},
	}

	// Insert one by one, skip if already exists
	insertedCount := 0
	for _, size := range companySizes {
		// Check if already exists
		var existing CompanySize
		if err := db.Where("label = ?", size.Label).First(&existing).Error; err == nil {
			// Already exists, skip
			continue
		}

		if err := db.Create(&size).Error; err != nil {
			log.Printf("Failed to seed company size %s: %v", size.Label, err)
			continue
		}
		insertedCount++
	}

	log.Printf("Successfully seeded %d company sizes", insertedCount)
	return nil
}
