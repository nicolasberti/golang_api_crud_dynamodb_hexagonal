package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tu-usuario/product-crud-hexagonal/internal/adapters/http/dto"
	"github.com/tu-usuario/product-crud-hexagonal/internal/core/domain"
	"github.com/tu-usuario/product-crud-hexagonal/internal/core/ports"
	"log/slog"
)

// MockProductService for testing
type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) Create(ctx context.Context, name, description string, price float64) (domain.Product, error) {
	args := m.Called(ctx, name, description, price)
	return args.Get(0).(domain.Product), args.Error(1)
}

func (m *MockProductService) Get(ctx context.Context, id string) (domain.Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Product), args.Error(1)
}

func (m *MockProductService) Update(ctx context.Context, id, name, description string, price float64) (domain.Product, error) {
	args := m.Called(ctx, id, name, description, price)
	return args.Get(0).(domain.Product), args.Error(1)
}

func (m *MockProductService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductService) List(ctx context.Context) ([]domain.Product, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Product), args.Error(1)
}

func (m *MockProductService) ListWithFilters(ctx context.Context, filters ports.ProductFilters) (*ports.ProductListResult, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(*ports.ProductListResult), args.Error(1)
}

func setupTestRouter() (*gin.Engine, *MockProductService) {
	gin.SetMode(gin.TestMode)

	mockService := &MockProductService{}
	logger := slog.Default()
	handler := NewProductHandler(mockService, logger)

	router := gin.New()
	v1 := router.Group("/api/v1")
	products := v1.Group("/products")
	{
		products.GET("", handler.List)
		products.POST("", handler.Create)
		products.GET("/:id", handler.Get)
		products.PUT("/:id", handler.Update)
		products.DELETE("/:id", handler.Delete)
	}

	return router, mockService
}

func TestProductHandler_List_WithDefaults(t *testing.T) {
	router, mockService := setupTestRouter()

	// Mock data
	now := time.Now().UTC()
	products := []domain.Product{
		{ID: "1", Name: "Test Product 1", Description: "Description 1", Price: 10.99, CreatedAt: now, UpdatedAt: now},
		{ID: "2", Name: "Test Product 2", Description: "Description 2", Price: 20.99, CreatedAt: now, UpdatedAt: now},
	}

	expectedResult := &ports.ProductListResult{
		Products:   products,
		TotalItems: 2,
	}

	mockService.On("ListWithFilters", mock.Anything, mock.MatchedBy(func(filters ports.ProductFilters) bool {
		return filters.Limit == 20 && filters.Page == 1 && filters.Offset == 0
	})).Return(expectedResult, nil)

	// Make request
	req, _ := http.NewRequest("GET", "/api/v1/products", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.ListProductsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response.Products, 2)
	assert.Equal(t, 1, response.Pagination.CurrentPage)
	assert.Equal(t, 20, response.Pagination.PerPage)
	assert.Equal(t, 2, response.Pagination.TotalItems)
	assert.Equal(t, 1, response.Pagination.TotalPages)
	assert.False(t, response.Pagination.HasNext)
	assert.False(t, response.Pagination.HasPrev)

	mockService.AssertExpectations(t)
}

func TestProductHandler_List_WithPagination(t *testing.T) {
	router, mockService := setupTestRouter()

	products := []domain.Product{
		{ID: "1", Name: "Test Product 1", Description: "Description 1", Price: 10.99, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	expectedResult := &ports.ProductListResult{
		Products:   products,
		TotalItems: 50,
	}

	mockService.On("ListWithFilters", mock.Anything, mock.MatchedBy(func(filters ports.ProductFilters) bool {
		return filters.Limit == 10 && filters.Page == 2 && filters.Offset == 10
	})).Return(expectedResult, nil)

	req, _ := http.NewRequest("GET", "/api/v1/products?page=2&limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.ListProductsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 2, response.Pagination.CurrentPage)
	assert.Equal(t, 10, response.Pagination.PerPage)
	assert.Equal(t, 50, response.Pagination.TotalItems)
	assert.Equal(t, 5, response.Pagination.TotalPages)
	assert.True(t, response.Pagination.HasNext)
	assert.True(t, response.Pagination.HasPrev)

	mockService.AssertExpectations(t)
}

func TestProductHandler_List_WithFilters(t *testing.T) {
	router, mockService := setupTestRouter()

	products := []domain.Product{
		{ID: "1", Name: "Laptop", Description: "Gaming laptop", Price: 999.99, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	expectedResult := &ports.ProductListResult{
		Products:   products,
		TotalItems: 1,
	}

	mockService.On("ListWithFilters", mock.Anything, mock.MatchedBy(func(filters ports.ProductFilters) bool {
		return filters.Name == "Laptop" && filters.MinPrice == 500 && filters.MaxPrice == 1500
	})).Return(expectedResult, nil)

	req, _ := http.NewRequest("GET", "/api/v1/products?name=Laptop&min_price=500&max_price=1500", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.ListProductsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response.Products, 1)
	assert.NotNil(t, response.FiltersApplied)
	assert.Equal(t, "Laptop", response.FiltersApplied.Name)
	assert.Equal(t, 500.0, response.FiltersApplied.MinPrice)
	assert.Equal(t, 1500.0, response.FiltersApplied.MaxPrice)

	mockService.AssertExpectations(t)
}

func TestProductHandler_List_WithSorting(t *testing.T) {
	router, mockService := setupTestRouter()

	products := []domain.Product{
		{ID: "1", Name: "A Product", Description: "Description", Price: 10.99, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	expectedResult := &ports.ProductListResult{
		Products:   products,
		TotalItems: 1,
	}

	mockService.On("ListWithFilters", mock.Anything, mock.MatchedBy(func(filters ports.ProductFilters) bool {
		return filters.SortBy == "name" && filters.SortOrder == "asc"
	})).Return(expectedResult, nil)

	req, _ := http.NewRequest("GET", "/api/v1/products?sort_by=name&sort_order=asc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestProductHandler_List_InvalidPage(t *testing.T) {
	router, _ := setupTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/products?page=1001", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "page cannot exceed 1000", response["error"])
}

func TestProductHandler_List_InvalidPriceRange(t *testing.T) {
	router, _ := setupTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/products?min_price=100&max_price=50", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "min_price cannot be greater than max_price", response["error"])
}

func TestProductHandler_List_InvalidSortField(t *testing.T) {
	router, _ := setupTestRouter()

	req, _ := http.NewRequest("GET", "/api/v1/products?sort_by=invalid_field", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_List_ServiceError(t *testing.T) {
	router, mockService := setupTestRouter()

	mockService.On("ListWithFilters", mock.Anything, mock.Anything).Return(nil, assert.AnError)

	req, _ := http.NewRequest("GET", "/api/v1/products", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "internal server error", response["error"])

	mockService.AssertExpectations(t)
}

func TestListProductsRequest_SetDefaults(t *testing.T) {
	tests := []struct {
		name     string
		input    dto.ListProductsRequest
		expected dto.ListProductsRequest
	}{
		{
			name:  "all empty",
			input: dto.ListProductsRequest{},
			expected: dto.ListProductsRequest{
				Page:      1,
				Limit:     20,
				SortBy:    "created_at",
				SortOrder: "desc",
			},
		},
		{
			name: "partial values",
			input: dto.ListProductsRequest{
				Page: 5,
			},
			expected: dto.ListProductsRequest{
				Page:      5,
				Limit:     20,
				SortBy:    "created_at",
				SortOrder: "desc",
			},
		},
		{
			name: "all values set",
			input: dto.ListProductsRequest{
				Page:      3,
				Limit:     50,
				SortBy:    "name",
				SortOrder: "asc",
			},
			expected: dto.ListProductsRequest{
				Page:      3,
				Limit:     50,
				SortBy:    "name",
				SortOrder: "asc",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.input
			req.SetDefaults()
			assert.Equal(t, tt.expected, req)
		})
	}
}

func TestListProductsRequest_GetOffset(t *testing.T) {
	tests := []struct {
		name     string
		input    dto.ListProductsRequest
		expected int
	}{
		{"page 1 limit 20", dto.ListProductsRequest{Page: 1, Limit: 20}, 0},
		{"page 2 limit 10", dto.ListProductsRequest{Page: 2, Limit: 10}, 10},
		{"page 3 limit 5", dto.ListProductsRequest{Page: 3, Limit: 5}, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.GetOffset())
		})
	}
}

func TestListProductsRequest_HasFilters(t *testing.T) {
	tests := []struct {
		name     string
		input    dto.ListProductsRequest
		expected bool
	}{
		{"no filters", dto.ListProductsRequest{}, false},
		{"name filter", dto.ListProductsRequest{Name: "test"}, true},
		{"min_price filter", dto.ListProductsRequest{MinPrice: 10}, true},
		{"max_price filter", dto.ListProductsRequest{MaxPrice: 100}, true},
		{"multiple filters", dto.ListProductsRequest{Name: "test", MinPrice: 10}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.HasFilters())
		})
	}
}
