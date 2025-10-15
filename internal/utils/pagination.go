package utils

import (
	"math"
)

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalRows  int64 `json:"total_rows"`
	TotalPages int   `json:"total_pages"`
}

type PaginationInput struct {
	Page  int
	Limit int
}

func NewPagination(page, limit int, totalRows int64) *Pagination {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))

	return &Pagination{
		Page:       page,
		Limit:      limit,
		TotalRows:  totalRows,
		TotalPages: totalPages,
	}
}

func CalculateOffset(page, limit int) int {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return (page - 1) * limit
}

func ValidatePagination(page, limit, maxLimit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	return page, limit
}

func GetPaginationMeta(page, limit int, totalRows int64) map[string]interface{} {
	pagination := NewPagination(page, limit, totalRows)
	return map[string]interface{}{
		"page":        pagination.Page,
		"limit":       pagination.Limit,
		"total_rows":  pagination.TotalRows,
		"total_pages": pagination.TotalPages,
		"has_next":    pagination.Page < pagination.TotalPages,
		"has_prev":    pagination.Page > 1,
	}
}
