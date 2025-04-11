// api/handlers/product_handler.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ntdt/product-service/internal/domain"
	"github.com/ntdt/product-service/internal/service"
	"github.com/ntdt/product-service/pkg/logger"
)

type ProductHandler struct {
	productService service.ProductService
	logger         logger.Logger
}

func NewProductHandler(productService service.ProductService, logger logger.Logger) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		logger:         logger,
	}
}

// ListProducts godoc
// @Summary List products
// @Description Get all products with filtering options
// @Tags products
// @Accept json
// @Produce json
// @Param name query string false "Product name (supports partial matching)"
// @Param categories query []string false "Product categories"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param sort_by query string false "Field to sort by"
// @Param sort_order query string false "Sort order (asc or desc)"
// @Param limit query int false "Number of records to return" default(10)
// @Param offset query int false "Number of records to skip" default(0)
// @Success 200 {array} domain.Product
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products [get]
// @Security BearerAuth
func (h *ProductHandler) ListProducts(c *gin.Context) {
	var filter domain.ProductFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		h.logger.Error("Failed to bind query parameters", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Invalid query parameters",
		})
		return
	}

	products, err := h.productService.GetProducts(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to get products", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Failed to retrieve products",
		})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProduct godoc
// @Summary Get product
// @Description Get a product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} domain.Product
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id} [get]
// @Security BearerAuth
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")

	product, err := h.productService.GetProductByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get product", err, logger.Fields{"productId": id})
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Failed to retrieve product",
		})
		return
	}

	if product == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Status: http.StatusNotFound,
			Error:  "Product not found",
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

// CreateProduct godoc
// @Summary Create product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param product body domain.Product true "Product data"
// @Success 201 {object} domain.Product
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products [post]
// @Security BearerAuth
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		h.logger.Error("Failed to bind request body", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Invalid request payload",
		})
		return
	}

	createdProduct, err := h.productService.CreateProduct(c.Request.Context(), product)
	if err != nil {
		h.logger.Error("Failed to create product", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Failed to create product",
		})
		return
	}

	c.JSON(http.StatusCreated, createdProduct)
}

// UpdateProduct godoc
// @Summary Update product
// @Description Update an existing product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body domain.Product true "Product data"
// @Success 200 {object} domain.Product
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id} [put]
// @Security BearerAuth
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		h.logger.Error("Failed to bind request body", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Status: http.StatusBadRequest,
			Error:  "Invalid request payload",
		})
		return
	}

	updatedProduct, err := h.productService.UpdateProduct(c.Request.Context(), id, product)
	if err != nil {
		h.logger.Error("Failed to update product", err, logger.Fields{"productId": id})
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Failed to update product",
		})
		return
	}

	if updatedProduct == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Status: http.StatusNotFound,
			Error:  "Product not found",
		})
		return
	}

	c.JSON(http.StatusOK, updatedProduct)
}

// DeleteProduct godoc
// @Summary Delete product
// @Description Delete a product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 204 "No Content"
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id} [delete]
// @Security BearerAuth
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	err := h.productService.DeleteProduct(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete product", err, logger.Fields{"productId": id})

		if err.Error() == "mongo: no documents in result" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Status: http.StatusNotFound,
				Error:  "Product not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: http.StatusInternalServerError,
			Error:  "Failed to delete product",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Check if service is healthy
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}
}

type ErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}
