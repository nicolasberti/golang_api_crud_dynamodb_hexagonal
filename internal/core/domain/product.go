package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidProduct = errors.New("invalid product data")
	ErrNotFound       = errors.New("product not found")
)

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewProduct Factory para crear un producto válido
func NewProduct(name, description string, price float64) (*Product, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if price < 0 {
		return nil, errors.New("price cannot be negative")
	}

	now := time.Now().UTC()
	return &Product{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Price:       price,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}
