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

func (s *service) ListWithFilters(ctx context.Context, filters ports.ProductFilters) (*ports.ProductListResult, error) {
	s.logger.Info("listing products with filters",
		"name", filters.Name,
		"min_price", filters.MinPrice,
		"max_price", filters.MaxPrice,
		"sort_by", filters.SortBy,
		"sort_order", filters.SortOrder,
		"offset", filters.Offset,
		"limit", filters.Limit,
	)

	result, err := s.repo.ListWithFilters(ctx, filters)
	if err != nil {
		s.logger.Error("failed to list products with filters", "error", err)
		return nil, err
	}

	s.logger.Info("successfully listed products", "count", len(result.Products), "total", result.TotalItems)
	return result, nil
}
