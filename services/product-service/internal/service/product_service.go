// internal/service/product_service.go
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ntdt/product-service/internal/domain"
	"github.com/ntdt/product-service/internal/repository"
	"github.com/ntdt/product-service/pkg/cache"
	"github.com/ntdt/product-service/pkg/logger"
	"github.com/ntdt/product-service/pkg/messaging"
)

type ProductService interface {
	GetProducts(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, error)
	GetProductByID(ctx context.Context, id string) (*domain.Product, error)
	CreateProduct(ctx context.Context, product domain.Product) (*domain.Product, error)
	UpdateProduct(ctx context.Context, id string, product domain.Product) (*domain.Product, error)
	DeleteProduct(ctx context.Context, id string) error
}

type productService struct {
	repo       repository.ProductRepository
	cache      cache.RedisClient
	messageBus messaging.RabbitMQClient
	logger     logger.Logger
}

func NewProductService(repo repository.ProductRepository, cache cache.RedisClient, messageBus messaging.RabbitMQClient, logger logger.Logger) ProductService {
	return &productService{
		repo:       repo,
		cache:      cache,
		messageBus: messageBus,
		logger:     logger,
	}
}

func (s *productService) GetProducts(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, error) {
	return s.repo.FindAll(ctx, filter)
}

func (s *productService) GetProductByID(ctx context.Context, id string) (*domain.Product, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("product:%s", id)

	// Check if product exists in cache
	cachedProduct, err := s.cache.Get(ctx, cacheKey)
	if err == nil && cachedProduct != "" {
		var product domain.Product
		if err := json.Unmarshal([]byte(cachedProduct), &product); err == nil {
			return &product, nil
		}
	}

	// If not in cache, get from repository
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if product != nil {
		// Cache product
		productJSON, _ := json.Marshal(product)
		s.cache.Set(ctx, cacheKey, string(productJSON), 30*time.Minute)
	}

	return product, nil
}

func (s *productService) CreateProduct(ctx context.Context, product domain.Product) (*domain.Product, error) {
	newProduct, err := s.repo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	// Publish event to message bus
	err = s.publishProductEvent("product.created", newProduct)
	if err != nil {
		// Log error but don't fail the operation
		s.logger.Error("Failed to publish product created event", err)
	}

	return newProduct, nil
}

func (s *productService) UpdateProduct(ctx context.Context, id string, product domain.Product) (*domain.Product, error) {
	updatedProduct, err := s.repo.Update(ctx, id, product)
	if err != nil {
		return nil, err
	}

	if updatedProduct != nil {
		// Invalidate cache
		cacheKey := fmt.Sprintf("product:%s", id)
		s.cache.Delete(ctx, cacheKey)

		// Publish event to message bus
		err = s.publishProductEvent("product.updated", updatedProduct)
		if err != nil {
			s.logger.Error("Failed to publish product updated event", err)
		}
	}

	return updatedProduct, nil
}

func (s *productService) DeleteProduct(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("product:%s", id)
	s.cache.Delete(ctx, cacheKey)

	// Publish event to message bus
	deleteEvent := map[string]interface{}{
		"id":        id,
		"timestamp": time.Now(),
	}

	err = s.publishEvent("product.deleted", deleteEvent)
	if err != nil {
		s.logger.Error("Failed to publish product deleted event", err)
	}

	return nil
}

func (s *productService) publishProductEvent(eventType string, product *domain.Product) error {
	event := map[string]interface{}{
		"id":        product.ID.Hex(),
		"product":   product,
		"timestamp": time.Now(),
	}

	return s.publishEvent(eventType, event)
}

func (s *productService) publishEvent(eventType string, payload interface{}) error {
	eventJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return s.messageBus.Publish("product_exchange", eventType, eventJSON)
}
