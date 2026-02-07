package dto

import (
	"time"
)

// ListProductsRequest represents query parameters for listing products
type ListProductsRequest struct {
	// Pagination
	Page  int `form:"page" binding:"min=1"`
	Limit int `form:"limit" binding:"min=1,max=100"`

	// Filters
	Name     string  `form:"name"`
	MinPrice float64 `form:"min_price" binding:"min=0"`
	MaxPrice float64 `form:"max_price" binding:"min=0"`

	// Sorting
	SortBy    string `form:"sort_by" binding:"omitempty,oneof=name price created_at updated_at"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc"`

	// Field selection
	Fields string `form:"fields"`
}

// ListProductsResponse represents the response structure for listing products
type ListProductsResponse struct {
	Products       []ProductResponse `json:"products"`
	Pagination     PaginationInfo    `json:"pagination"`
	FiltersApplied FilterInfo        `json:"filters_applied,omitempty"`
}

// ProductResponse represents a product in API responses
type ProductResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PaginationInfo contains pagination metadata
type PaginationInfo struct {
	CurrentPage int  `json:"current_page"`
	PerPage     int  `json:"per_page"`
	TotalPages  int  `json:"total_pages"`
	TotalItems  int  `json:"total_items"`
	HasNext     bool `json:"has_next"`
	HasPrev     bool `json:"has_prev"`
}

// FilterInfo contains information about applied filters
type FilterInfo struct {
	Name     string  `json:"name,omitempty"`
	MinPrice float64 `json:"min_price,omitempty"`
	MaxPrice float64 `json:"max_price,omitempty"`
}

// SetDefaults sets default values for the request
func (r *ListProductsRequest) SetDefaults() {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.Limit <= 0 {
		r.Limit = 20
	}
	if r.SortBy == "" {
		r.SortBy = "created_at"
	}
	if r.SortOrder == "" {
		r.SortOrder = "desc"
	}
}

// GetOffset calculates the offset for database queries
func (r *ListProductsRequest) GetOffset() int {
	return (r.Page - 1) * r.Limit
}

// HasFilters returns true if any filter is applied
func (r *ListProductsRequest) HasFilters() bool {
	return r.Name != "" || r.MinPrice > 0 || r.MaxPrice > 0
}

// NewProductResponse creates a new product response from domain product
func NewProductResponse(id, name, description string, price float64, createdAt, updatedAt time.Time) ProductResponse {
	return ProductResponse{
		ID:          id,
		Name:        name,
		Description: description,
		Price:       price,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
