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
}
