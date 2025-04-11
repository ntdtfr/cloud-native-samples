// internal/domain/product.go
package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name" binding:"required"`
	Description string             `json:"description" bson:"description"`
	Price       float64            `json:"price" bson:"price" binding:"required,gt=0"`
	SKU         string             `json:"sku" bson:"sku" binding:"required"`
	Inventory   int                `json:"inventory" bson:"inventory" binding:"min=0"`
	Categories  []string           `json:"categories" bson:"categories"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type ProductFilter struct {
	Name       string   `form:"name"`
	Categories []string `form:"categories"`
	MinPrice   float64  `form:"min_price"`
	MaxPrice   float64  `form:"max_price"`
	SortBy     string   `form:"sort_by"`
	SortOrder  string   `form:"sort_order"`
	Limit      int      `form:"limit,default=10"`
	Offset     int      `form:"offset,default=0"`
}
