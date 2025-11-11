package service

import (
	"keerja-backend/internal/cache"
	"keerja-backend/internal/domain/master"
	masterService "keerja-backend/internal/service/master"

	"gorm.io/gorm"
)

// NewIndustryService creates a new instance of IndustryService
func NewIndustryService(repo master.IndustryRepository, cache cache.Cache) master.IndustryService {
	return masterService.NewIndustryService(repo, cache)
}

// NewCompanySizeService creates a new instance of CompanySizeService
func NewCompanySizeService(repo master.CompanySizeRepository, cache cache.Cache) master.CompanySizeService {
	return masterService.NewCompanySizeService(repo, cache)
}

// NewProvinceService creates a new instance of ProvinceService
func NewProvinceService(repo master.ProvinceRepository, cache cache.Cache) master.ProvinceService {
	return masterService.NewProvinceService(repo, cache)
}

// NewCityService creates a new instance of CityService
func NewCityService(
	repo master.CityRepository,
	provinceRepo master.ProvinceRepository,
	cache cache.Cache,
) master.CityService {
	return masterService.NewCityService(repo, provinceRepo, cache)
}

// NewDistrictService creates a new instance of DistrictService
func NewDistrictService(
	repo master.DistrictRepository,
	cityRepo master.CityRepository,
	provinceRepo master.ProvinceRepository,
	cache cache.Cache,
) master.DistrictService {
	return masterService.NewDistrictService(repo, cityRepo, provinceRepo, cache)
}

// Admin Services Factory Functions

// NewAdminIndustryService creates a new AdminIndustryService
func NewAdminIndustryService(
	baseService master.IndustryService,
	repo master.IndustryRepository,
	db *gorm.DB,
	cache cache.Cache,
) master.AdminIndustryService {
	return masterService.NewAdminIndustryService(baseService, repo, db, cache)
}

// NewAdminCompanySizeService creates a new AdminCompanySizeService
func NewAdminCompanySizeService(
	baseService master.CompanySizeService,
	repo master.CompanySizeRepository,
	db *gorm.DB,
	cache cache.Cache,
) master.AdminCompanySizeService {
	return masterService.NewAdminCompanySizeService(baseService, repo, db, cache)
}

// NewAdminProvinceService creates a new AdminProvinceService
func NewAdminProvinceService(
	baseService master.ProvinceService,
	repo master.ProvinceRepository,
	db *gorm.DB,
	cache cache.Cache,
) master.AdminProvinceService {
	return masterService.NewAdminProvinceService(baseService, repo, db, cache)
}

// NewAdminCityService creates a new AdminCityService
func NewAdminCityService(
	baseService master.CityService,
	repo master.CityRepository,
	db *gorm.DB,
	cache cache.Cache,
) master.AdminCityService {
	return masterService.NewAdminCityService(baseService, repo, db, cache)
}

// NewAdminDistrictService creates a new AdminDistrictService
func NewAdminDistrictService(
	baseService master.DistrictService,
	repo master.DistrictRepository,
	db *gorm.DB,
	cache cache.Cache,
) master.AdminDistrictService {
	return masterService.NewAdminDistrictService(baseService, repo, db, cache)
}

// NewAdminJobTypeService creates a new AdminJobTypeService
func NewAdminJobTypeService(
	jobOptionsService master.JobOptionsService,
	repo master.JobOptionsRepository,
	db *gorm.DB,
	cache cache.Cache,
) master.AdminJobTypeService {
	return masterService.NewAdminJobTypeService(jobOptionsService, repo, db, cache)
}
