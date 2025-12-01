package seeders

import (
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Industry represents an industry entity
type Industry struct {
	ID           int64  `gorm:"primaryKey"`
	Name         string `gorm:"uniqueIndex;not null"`
	Slug         string `gorm:"uniqueIndex;not null"`
	Description  string
	DisplayOrder int `gorm:"default:0"`
}

func (Industry) TableName() string {
	return "industries"
}

// IndustriesSeeder seeds the industries table
func IndustriesSeeder(db *gorm.DB) error {
	log.Println("Seeding industries table...")

	industries := []Industry{
		{Name: "Technology", Slug: "technology", Description: "Information technology, software development, and IT services", DisplayOrder: 1},
		{Name: "Healthcare", Slug: "healthcare", Description: "Healthcare services, medical facilities, and pharmaceuticals", DisplayOrder: 2},
		{Name: "Education", Slug: "education", Description: "Educational institutions, training centers, and e-learning", DisplayOrder: 3},
		{Name: "Finance", Slug: "finance", Description: "Banking, insurance, investment, and financial services", DisplayOrder: 4},
		{Name: "Retail", Slug: "retail", Description: "Retail stores, e-commerce, and consumer goods", DisplayOrder: 5},
		{Name: "Manufacturing", Slug: "manufacturing", Description: "Manufacturing and industrial production", DisplayOrder: 6},
		{Name: "Construction", Slug: "construction", Description: "Construction, real estate development, and infrastructure", DisplayOrder: 7},
		{Name: "Transportation", Slug: "transportation", Description: "Logistics, shipping, and transportation services", DisplayOrder: 8},
		{Name: "Hospitality", Slug: "hospitality", Description: "Hotels, restaurants, tourism, and food & beverage", DisplayOrder: 9},
		{Name: "Agriculture", Slug: "agriculture", Description: "Agriculture, farming, and agribusiness", DisplayOrder: 10},
		{Name: "Telecommunications", Slug: "telecommunications", Description: "Telecommunications and internet service providers", DisplayOrder: 11},
		{Name: "Media & Entertainment", Slug: "media-entertainment", Description: "Media production, entertainment, and creative industries", DisplayOrder: 12},
		{Name: "Energy", Slug: "energy", Description: "Energy production, oil & gas, and renewable energy", DisplayOrder: 13},
		{Name: "Automotive", Slug: "automotive", Description: "Automotive manufacturing and services", DisplayOrder: 14},
		{Name: "Fashion", Slug: "fashion", Description: "Fashion, apparel, and textile industry", DisplayOrder: 15},
		{Name: "Consulting", Slug: "consulting", Description: "Business consulting and professional services", DisplayOrder: 16},
		{Name: "Marketing & Advertising", Slug: "marketing-advertising", Description: "Marketing, advertising, and public relations", DisplayOrder: 17},
		{Name: "Legal Services", Slug: "legal-services", Description: "Legal firms and legal services", DisplayOrder: 18},
		{Name: "Real Estate", Slug: "real-estate", Description: "Real estate agencies and property management", DisplayOrder: 19},
		{Name: "Non-Profit", Slug: "non-profit", Description: "Non-profit organizations and NGOs", DisplayOrder: 20},
		{Name: "Government", Slug: "government", Description: "Government agencies and public sector", DisplayOrder: 21},
		{Name: "Arts & Crafts", Slug: "arts-crafts", Description: "Arts, crafts, and creative industries", DisplayOrder: 22},
		{Name: "Sports & Recreation", Slug: "sports-recreation", Description: "Sports, fitness, and recreational services", DisplayOrder: 23},
		{Name: "Beauty & Wellness", Slug: "beauty-wellness", Description: "Beauty salons, spas, and wellness centers", DisplayOrder: 24},
		{Name: "Financial Technology", Slug: "fintech", Description: "Financial technology and digital payments", DisplayOrder: 25},
		{Name: "E-Commerce", Slug: "e-commerce", Description: "Online marketplace and e-commerce platforms", DisplayOrder: 26},
		{Name: "Education Technology", Slug: "edtech", Description: "Educational technology and e-learning platforms", DisplayOrder: 27},
		{Name: "Healthcare Technology", Slug: "healthtech", Description: "Healthcare technology and digital health", DisplayOrder: 28},
		{Name: "Travel & Tourism", Slug: "travel-tourism", Description: "Travel agencies, booking platforms, and tourism", DisplayOrder: 29},
		{Name: "Food & Beverage", Slug: "food-beverage", Description: "Food production, restaurants, and beverage industry", DisplayOrder: 30},
		{Name: "Pharmaceutical", Slug: "pharmaceutical", Description: "Pharmaceutical and drug manufacturing", DisplayOrder: 31},
		{Name: "Insurance", Slug: "insurance", Description: "Insurance services and products", DisplayOrder: 32},
		{Name: "Banking", Slug: "banking", Description: "Commercial and retail banking services", DisplayOrder: 33},
		{Name: "Investment", Slug: "investment", Description: "Investment management and securities", DisplayOrder: 34},
		{Name: "Human Resources", Slug: "human-resources", Description: "HR services, recruitment, and staffing", DisplayOrder: 35},
		{Name: "Other", Slug: "other", Description: "Other industries not listed above", DisplayOrder: 99},
	}

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"slug", "description", "display_order"}),
	}).Create(&industries)

	if result.Error != nil {
		log.Printf("Failed to seed industries: %v", result.Error)
		return result.Error
	}

	log.Printf("Successfully seeded %d industries", len(industries))
	return nil
}
