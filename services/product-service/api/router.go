// api/router.go
package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ntdt/product-service/api/handlers"
	"github.com/ntdt/product-service/api/middleware"
	"github.com/ntdt/product-service/config"
	"github.com/ntdt/product-service/internal/service"
	"github.com/ntdt/product-service/pkg/logger"

	_ "github.com/ntdt/product-service/api/swagger" // swagger docs
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(productService service.ProductService, logger logger.Logger, cfg *config.Config) *gin.Engine {
	r := gin.New()

	// Middleware
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(logger))
	r.Use(middleware.Cors())

	// Security middleware
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.RateLimiter(cfg.RateLimit))

	// Health check
	r.GET("/health", handlers.HealthCheck())

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	v1 := r.Group("/api/v1")
	{
		v1.Use(middleware.Auth(cfg.Auth))

		products := v1.Group("/products")
		{
			h := handlers.NewProductHandler(productService, logger)
			products.GET("", h.ListProducts)
			products.GET("/:id", h.GetProduct)
			products.POST("", h.CreateProduct)
			products.PUT("/:id", h.UpdateProduct)
			products.DELETE("/:id", h.DeleteProduct)
		}
	}

	return r
}
