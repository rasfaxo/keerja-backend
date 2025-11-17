package master

// CRUD Request/Response DTOs for Master Data Entities

// CRUD Request DTOs

// CreateProvinceRequest for creating province
type CreateProvinceRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=255"`
	Code     string `json:"code" validate:"required,min=2,max=50"`
	IsActive bool   `json:"is_active"`
}

// UpdateProvinceRequest for updating province
type UpdateProvinceRequest struct {
	Name     string `json:"name" validate:"omitempty,min=2,max=255"`
	Code     string `json:"code" validate:"omitempty,min=2,max=50"`
	IsActive *bool  `json:"is_active,omitempty"`
}

// CreateCityRequest for creating city
type CreateCityRequest struct {
	Name       string `json:"name" validate:"required,min=2,max=255"`
	Type       string `json:"type" validate:"required,oneof=Kota Kabupaten"`
	Code       string `json:"code" validate:"required,min=2,max=50"`
	ProvinceID int64  `json:"province_id" validate:"required,min=1"`
	IsActive   bool   `json:"is_active"`
}

// UpdateCityRequest for updating city
type UpdateCityRequest struct {
	Name       string `json:"name" validate:"omitempty,min=2,max=255"`
	Type       string `json:"type" validate:"omitempty,oneof=Kota Kabupaten"`
	Code       string `json:"code" validate:"omitempty,min=2,max=50"`
	ProvinceID *int64 `json:"province_id" validate:"omitempty,min=1"`
	IsActive   *bool  `json:"is_active,omitempty"`
}

// CreateDistrictRequest for creating district
type CreateDistrictRequest struct {
	Name       string `json:"name" validate:"required,min=2,max=255"`
	Code       string `json:"code" validate:"required,min=2,max=50"`
	PostalCode string `json:"postal_code" validate:"omitempty,len=5"`
	CityID     int64  `json:"city_id" validate:"required,min=1"`
	IsActive   bool   `json:"is_active"`
}

// UpdateDistrictRequest for updating district
type UpdateDistrictRequest struct {
	Name       string `json:"name" validate:"omitempty,min=2,max=255"`
	Code       string `json:"code" validate:"omitempty,min=2,max=50"`
	PostalCode string `json:"postal_code" validate:"omitempty,len=5"`
	CityID     *int64 `json:"city_id" validate:"omitempty,min=1"`
	IsActive   *bool  `json:"is_active,omitempty"`
}

// CreateIndustryRequest for creating industry
type CreateIndustryRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Slug        string `json:"slug" validate:"required,min=2,max=255"`
	Description string `json:"description"`
	IconURL     string `json:"icon_url" validate:"omitempty,url"`
	IsActive    bool   `json:"is_active"`
}

// UpdateIndustryRequest for updating industry
type UpdateIndustryRequest struct {
	Name        string `json:"name" validate:"omitempty,min=2,max=255"`
	Slug        string `json:"slug" validate:"omitempty,min=2,max=255"`
	Description string `json:"description"`
	IconURL     string `json:"icon_url" validate:"omitempty,url"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

// CreateCompanySizeRequest for creating company size
type CreateCompanySizeRequest struct {
	Label        string `json:"label" validate:"required,min=2,max=255"`
	MinEmployees int    `json:"min_employees" validate:"required,min=0"`
	MaxEmployees *int   `json:"max_employees" validate:"omitempty,gtefield=MinEmployees"`
	IsActive     bool   `json:"is_active"`
}

// UpdateCompanySizeRequest for updating company size
type UpdateCompanySizeRequest struct {
	Label        string `json:"label" validate:"omitempty,min=2,max=255"`
	MinEmployees *int   `json:"min_employees" validate:"omitempty,min=0"`
	MaxEmployees *int   `json:"max_employees" validate:"omitempty,gtefield=MinEmployees"`
	IsActive     *bool  `json:"is_active,omitempty"`
}

// CreateJobTypeRequest for creating job type
type CreateJobTypeRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Code  string `json:"code" validate:"required,min=2,max=30"`
	Order int    `json:"order" validate:"omitempty"`
}

// UpdateJobTypeRequest for updating job type
type UpdateJobTypeRequest struct {
	Name  string `json:"name" validate:"omitempty,min=2,max=100"`
	Code  string `json:"code" validate:"omitempty,min=2,max=30"`
	Order *int   `json:"order,omitempty"`
}

// JobTypeResponse represents a job type response
type JobTypeResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Code  string `json:"code"`
	Order int    `json:"order"`
}

// ListResponse generic paginated list response
type ListResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}
