package dto

import "time"

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// ErrorInfo represents error details in response
type ErrorInfo struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Meta represents metadata for paginated responses
type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page  int `json:"page" query:"page" validate:"omitempty,min=1"`
	Limit int `json:"limit" query:"limit" validate:"omitempty,min=1,max=100"`
}

// GetPage returns the page number, defaulting to 1
func (p *PaginationRequest) GetPage() int {
	if p.Page < 1 {
		return 1
	}
	return p.Page
}

// GetLimit returns the limit, defaulting to 10
func (p *PaginationRequest) GetLimit() int {
	if p.Limit < 1 {
		return 10
	}
	if p.Limit > 100 {
		return 100
	}
	return p.Limit
}

// GetOffset calculates the offset for database queries
func (p *PaginationRequest) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

// CalculateTotalPages calculates total pages from total items and limit
func CalculateTotalPages(totalItems int64, limit int) int {
	if limit <= 0 {
		limit = 10
	}
	totalPages := int(totalItems) / limit
	if int(totalItems)%limit > 0 {
		totalPages++
	}
	return totalPages
}

// SuccessResponse creates a success response
func SuccessResponse(message string, data interface{}) *Response {
	return &Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse creates an error response
func ErrorResponse(code, message string) *Response {
	return &Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	}
}

// ErrorResponseWithDetails creates an error response with details
func ErrorResponseWithDetails(code, message string, details map[string]interface{}) *Response {
	return &Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

// PaginatedResponse creates a paginated response
func PaginatedResponse(message string, data interface{}, page, limit int, totalItems int64) *Response {
	return &Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta: &Meta{
			Page:       page,
			Limit:      limit,
			TotalItems: totalItems,
			TotalPages: CalculateTotalPages(totalItems, limit),
		},
	}
}

// IDRequest represents a request with ID parameter
type IDRequest struct {
	ID int64 `json:"id" uri:"id" validate:"required,min=1"`
}

// SlugRequest represents a request with slug parameter
type SlugRequest struct {
	Slug string `json:"slug" uri:"slug" validate:"required,min=1"`
}

// UUIDRequest represents a request with UUID parameter
type UUIDRequest struct {
	UUID string `json:"uuid" uri:"uuid" validate:"required,uuid"`
}

// SortRequest represents sorting parameters
type SortRequest struct {
	SortBy    string `json:"sort_by" query:"sort_by" validate:"omitempty"`
	SortOrder string `json:"sort_order" query:"sort_order" validate:"omitempty,oneof=asc desc"`
}

// GetSortBy returns the sort field, defaulting to "created_at"
func (s *SortRequest) GetSortBy() string {
	if s.SortBy == "" {
		return "created_at"
	}
	return s.SortBy
}

// GetSortOrder returns the sort order, defaulting to "desc"
func (s *SortRequest) GetSortOrder() string {
	if s.SortOrder == "" {
		return "desc"
	}
	return s.SortOrder
}

// DateRangeRequest represents date range filter
type DateRangeRequest struct {
	StartDate *time.Time `json:"start_date" query:"start_date" validate:"omitempty"`
	EndDate   *time.Time `json:"end_date" query:"end_date" validate:"omitempty,gtefield=StartDate"`
}

// SearchRequest represents a search query
type SearchRequest struct {
	Query string `json:"query" query:"q" validate:"omitempty,min=1,max=255"`
}

// StatusFilter represents status filtering
type StatusFilter struct {
	Status string `json:"status" query:"status" validate:"omitempty"`
}

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrorResponse creates a validation error response
func ValidationErrorResponse(errors []ValidationError) *Response {
	details := make(map[string]interface{})
	for _, err := range errors {
		details[err.Field] = err.Message
	}

	return &Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    "VALIDATION_ERROR",
			Message: "Validation failed",
			Details: details,
		},
	}
}
