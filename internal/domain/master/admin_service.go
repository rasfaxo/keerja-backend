package master

import "context"

// AdminIndustryService defines complete CRUD operations for industry master data
type AdminIndustryService interface {
	// Embed base IndustryService for read operations
	IndustryService

	// Create creates a new industry
	Create(ctx context.Context, req CreateIndustryRequest) (*IndustryResponse, error)

	// Update updates an existing industry
	Update(ctx context.Context, id int64, req UpdateIndustryRequest) (*IndustryResponse, error)

	// Delete deletes an industry if not referenced by companies
	Delete(ctx context.Context, id int64) error

	// CheckDuplicateName checks if an industry with the given name exists
	CheckDuplicateName(ctx context.Context, name string) (bool, error)

	// CountCompanyReferences counts how many companies reference this industry
	CountCompanyReferences(ctx context.Context, id int64) (int64, error)
}

// AdminCompanySizeService defines complete CRUD operations for company size master data
type AdminCompanySizeService interface {
	// Embed base CompanySizeService for read operations
	CompanySizeService

	// Create creates a new company size category
	Create(ctx context.Context, req CreateCompanySizeRequest) (*CompanySizeResponse, error)

	// Update updates an existing company size category
	Update(ctx context.Context, id int64, req UpdateCompanySizeRequest) (*CompanySizeResponse, error)

	// Delete deletes a company size category if not referenced by companies
	Delete(ctx context.Context, id int64) error

	// CheckDuplicateCategory checks if a size category with the given label exists
	CheckDuplicateCategory(ctx context.Context, label string) (bool, error)

	// CountCompanyReferences counts how many companies reference this size category
	CountCompanyReferences(ctx context.Context, id int64) (int64, error)
}

// AdminProvinceService defines complete CRUD operations for province master data
type AdminProvinceService interface {
	// Embed base ProvinceService for read operations
	ProvinceService

	// Create creates a new province
	Create(ctx context.Context, req CreateProvinceRequest) (*ProvinceResponse, error)

	// Update updates an existing province
	Update(ctx context.Context, id int64, req UpdateProvinceRequest) (*ProvinceResponse, error)

	// Delete deletes a province if not referenced by cities or companies
	Delete(ctx context.Context, id int64) error

	// CheckDuplicateCode checks if a province with the given code exists
	CheckDuplicateCode(ctx context.Context, code string) (bool, error)

	// CountReferences counts how many cities and companies reference this province
	CountReferences(ctx context.Context, id int64) (cities int64, companies int64, err error)
}

// AdminCityService defines complete CRUD operations for city master data
type AdminCityService interface {
	// Embed base CityService for read operations
	CityService

	// Create creates a new city
	Create(ctx context.Context, req CreateCityRequest) (*CityResponse, error)

	// Update updates an existing city
	Update(ctx context.Context, id int64, req UpdateCityRequest) (*CityResponse, error)

	// Delete deletes a city if not referenced by districts or companies
	Delete(ctx context.Context, id int64) error

	// CheckDuplicateNameInProvince checks if a city with the given name exists in the province
	CheckDuplicateNameInProvince(ctx context.Context, name string, provinceID int64) (bool, error)

	// CountReferences counts how many districts and companies reference this city
	CountReferences(ctx context.Context, id int64) (districts int64, companies int64, err error)
}

// AdminDistrictService defines complete CRUD operations for district master data
type AdminDistrictService interface {
	// Embed base DistrictService for read operations
	DistrictService

	// Create creates a new district
	Create(ctx context.Context, req CreateDistrictRequest) (*DistrictResponse, error)

	// Update updates an existing district
	Update(ctx context.Context, id int64, req UpdateDistrictRequest) (*DistrictResponse, error)

	// Delete deletes a district if not referenced by companies
	Delete(ctx context.Context, id int64) error

	// CheckDuplicateNameInCity checks if a district with the given name exists in the city
	CheckDuplicateNameInCity(ctx context.Context, name string, cityID int64) (bool, error)

	// CountCompanyReferences counts how many companies reference this district
	CountCompanyReferences(ctx context.Context, id int64) (int64, error)
}

// AdminJobTypeService defines complete CRUD operations for job type master data
type AdminJobTypeService interface {
	// Get all job types (from JobOptionsService)
	GetJobTypes(ctx context.Context) ([]JobType, error)

	// Get job type by ID
	GetJobTypeByID(ctx context.Context, id int64) (*JobType, error)

	// Create creates a new job type
	Create(ctx context.Context, req CreateJobTypeRequest) (*JobTypeResponse, error)

	// Update updates an existing job type
	Update(ctx context.Context, id int64, req UpdateJobTypeRequest) (*JobTypeResponse, error)

	// Delete deletes a job type if not referenced by jobs
	Delete(ctx context.Context, id int64) error

	// CheckDuplicateCode checks if a job type with the given code exists
	CheckDuplicateCode(ctx context.Context, code string, excludeID *int64) (bool, error)

	// CountJobReferences counts how many jobs reference this job type
	CountJobReferences(ctx context.Context, id int64) (int64, error)
}
