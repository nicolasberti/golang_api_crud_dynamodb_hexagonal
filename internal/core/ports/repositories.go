package ports

import (
	"context"
	"github.com/tu-usuario/product-crud-hexagonal/internal/core/domain"
)

type ProductRepository interface {
	Save(ctx context.Context, product domain.Product) error
	GetByID(ctx context.Context, id string) (domain.Product, error)
	Update(ctx context.Context, product domain.Product) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.Product, error)
	ListWithFilters(ctx context.Context, filters ProductFilters) (*ProductListResult, error)
}

// ProductFilters represents filtering options for product queries
type ProductFilters struct {
	Name      string
	MinPrice  float64
	MaxPrice  float64
	SortBy    string
	SortOrder string
	Offset    int
	Limit     int
}

// ProductListResult contains the result of a filtered product query
type ProductListResult struct {
	Products   []domain.Product
	TotalItems int
}
