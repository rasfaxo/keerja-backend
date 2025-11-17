package postgres

import (
	"gorm.io/gorm"

	"keerja-backend/internal/domain/master"
	masterRepo "keerja-backend/internal/repository/master"
)

// NewIndustryRepository creates a new industry repository instance
func NewIndustryRepository(db *gorm.DB) master.IndustryRepository {
	return masterRepo.NewIndustryRepository(db)
}

// NewCompanySizeRepository creates a new company size repository instance
func NewCompanySizeRepository(db *gorm.DB) master.CompanySizeRepository {
	return masterRepo.NewCompanySizeRepository(db)
}

// NewProvinceRepository creates a new province repository instance
func NewProvinceRepository(db *gorm.DB) master.ProvinceRepository {
	return masterRepo.NewProvinceRepository(db)
}

// NewCityRepository creates a new city repository instance
func NewCityRepository(db *gorm.DB) master.CityRepository {
	return masterRepo.NewCityRepository(db)
}

// NewDistrictRepository creates a new district repository instance
func NewDistrictRepository(db *gorm.DB) master.DistrictRepository {
	return masterRepo.NewDistrictRepository(db)
}
