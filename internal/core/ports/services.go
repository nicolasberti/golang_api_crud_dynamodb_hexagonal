package ports

import (
	"context"
	"github.com/tu-usuario/product-crud-hexagonal/internal/core/domain"
)

type ProductService interface {
	Create(ctx context.Context, name, description string, price float64) (domain.Product, error)
	Get(ctx context.Context, id string) (domain.Product, error)
	Update(ctx context.Context, id, name, description string, price float64) (domain.Product, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.Product, error)
}
