package services

import (
	"context"
	"time"

	"github.com/tu-usuario/product-crud-hexagonal/internal/core/domain"
	"github.com/tu-usuario/product-crud-hexagonal/internal/core/ports"
	"log/slog"
)

type service struct {
	repo   ports.ProductRepository
	logger *slog.Logger
}

func NewProductService(repo ports.ProductRepository, logger *slog.Logger) ports.ProductService {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s *service) Create(ctx context.Context, name, description string, price float64) (domain.Product, error) {
	product, err := domain.NewProduct(name, description, price)
	if err != nil {
		s.logger.Warn("invalid product creation attempt", "error", err)
		return domain.Product{}, domain.ErrInvalidProduct
	}

	if err := s.repo.Save(ctx, *product); err != nil {
		s.logger.Error("failed to save product", "error", err)
		return domain.Product{}, err
	}

	return *product, nil
}

func (s *service) Get(ctx context.Context, id string) (domain.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) Update(ctx context.Context, id, name, description string, price float64) (domain.Product, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Product{}, err
	}

	existing.Name = name
	existing.Description = description
	existing.Price = price
	existing.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, existing); err != nil {
		s.logger.Error("failed to update product", "id", id, "error", err)
		return domain.Product{}, err
	}

	return existing, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) List(ctx context.Context) ([]domain.Product, error) {
	return s.repo.List(ctx)
}
