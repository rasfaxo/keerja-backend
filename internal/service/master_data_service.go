package service

import (
	"keerja-backend/internal/cache"
	"keerja-backend/internal/domain/master"
	masterService "keerja-backend/internal/service/master"
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
