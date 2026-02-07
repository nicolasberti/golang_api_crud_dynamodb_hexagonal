package http

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tu-usuario/product-crud-hexagonal/internal/adapters/http/dto"
	"github.com/tu-usuario/product-crud-hexagonal/internal/core/domain"
	"github.com/tu-usuario/product-crud-hexagonal/internal/core/ports"
	"log/slog"
)

type ProductHandler struct {
	service ports.ProductService
	logger  *slog.Logger
}

func NewProductHandler(service ports.ProductService, logger *slog.Logger) *ProductHandler {
	return &ProductHandler{
		service: service,
		logger:  logger,
	}
}

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.service.Create(c.Request.Context(), req.Name, req.Description, req.Price)
	if err != nil {
		if err == domain.ErrInvalidProduct {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("failed to create product", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) Get(c *gin.Context) {
	id := c.Param("id")
	product, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("failed to get product", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) List(c *gin.Context) {
	var req dto.ListProductsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Warn("invalid query parameters", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	// Set defaults
	req.SetDefaults()

	// Additional validations
	if req.Page > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page cannot exceed 1000"})
		return
	}

	if req.MinPrice > 0 && req.MaxPrice > 0 && req.MinPrice > req.MaxPrice {
		c.JSON(http.StatusBadRequest, gin.H{"error": "min_price cannot be greater than max_price"})
		return
	}

	// Build filters for service
	filters := ports.ProductFilters{
		Name:      req.Name,
		MinPrice:  req.MinPrice,
		MaxPrice:  req.MaxPrice,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
		Offset:    req.GetOffset(),
		Limit:     req.Limit,
	}

	result, err := h.service.ListWithFilters(c.Request.Context(), filters)
	if err != nil {
		h.logger.Error("failed to list products with filters", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Build response
	response := dto.ListProductsResponse{
		Products: make([]dto.ProductResponse, len(result.Products)),
		Pagination: dto.PaginationInfo{
			CurrentPage: req.Page,
			PerPage:     req.Limit,
			TotalItems:  result.TotalItems,
			TotalPages:  int(math.Ceil(float64(result.TotalItems) / float64(req.Limit))),
			HasNext:     req.Page*req.Limit < result.TotalItems,
			HasPrev:     req.Page > 1,
		},
	}

	// Convert domain products to DTOs
	for i, product := range result.Products {
		response.Products[i] = dto.NewProductResponse(
			product.ID,
			product.Name,
			product.Description,
			product.Price,
			product.CreatedAt,
			product.UpdatedAt,
		)
	}

	// Add filter info if filters were applied
	if req.HasFilters() {
		response.FiltersApplied = dto.FilterInfo{
			Name:     req.Name,
			MinPrice: req.MinPrice,
			MaxPrice: req.MaxPrice,
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.service.Update(c.Request.Context(), id, req.Name, req.Description, req.Price)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("failed to update product", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("failed to delete product", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}
