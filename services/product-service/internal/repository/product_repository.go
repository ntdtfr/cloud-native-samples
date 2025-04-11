// internal/repository/product_repository.go
package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ntdt/product-service/internal/domain"
)

type ProductRepository interface {
	FindAll(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, error)
	FindByID(ctx context.Context, id string) (*domain.Product, error)
	Create(ctx context.Context, product domain.Product) (*domain.Product, error)
	Update(ctx context.Context, id string, product domain.Product) (*domain.Product, error)
	Delete(ctx context.Context, id string) error
}

type mongoProductRepository struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewProductRepository(client *mongo.Client, database string) ProductRepository {
	return &mongoProductRepository{
		client:     client,
		database:   database,
		collection: "products",
	}
}

func (r *mongoProductRepository) FindAll(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, error) {
	coll := r.client.Database(r.database).Collection(r.collection)

	findOptions := options.Find()
	findOptions.SetLimit(int64(filter.Limit))
	findOptions.SetSkip(int64(filter.Offset))

	if filter.SortBy != "" {
		order := 1
		if filter.SortOrder == "desc" {
			order = -1
		}
		findOptions.SetSort(bson.D{{Key: filter.SortBy, Value: order}})
	}

	filterBson := bson.M{}
	if filter.Name != "" {
		filterBson["name"] = bson.M{"$regex": primitive.Regex{Pattern: filter.Name, Options: "i"}}
	}
	if len(filter.Categories) > 0 {
		filterBson["categories"] = bson.M{"$in": filter.Categories}
	}
	if filter.MinPrice > 0 {
		filterBson["price"] = bson.M{"$gte": filter.MinPrice}
	}
	if filter.MaxPrice > 0 {
		if _, ok := filterBson["price"]; ok {
			filterBson["price"].(bson.M)["$lte"] = filter.MaxPrice
		} else {
			filterBson["price"] = bson.M{"$lte": filter.MaxPrice}
		}
	}

	cursor, err := coll.Find(ctx, filterBson, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []domain.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *mongoProductRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {
	coll := r.client.Database(r.database).Collection(r.collection)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var product domain.Product
	err = coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}

func (r *mongoProductRepository) Create(ctx context.Context, product domain.Product) (*domain.Product, error) {
	coll := r.client.Database(r.database).Collection(r.collection)

	product.ID = primitive.NewObjectID()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	_, err := coll.InsertOne(ctx, product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *mongoProductRepository) Update(ctx context.Context, id string, product domain.Product) (*domain.Product, error) {
	coll := r.client.Database(r.database).Collection(r.collection)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	product.UpdatedAt = time.Now()

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": product}

	result := coll.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))

	var updatedProduct domain.Product
	if err := result.Decode(&updatedProduct); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &updatedProduct, nil
}

func (r *mongoProductRepository) Delete(ctx context.Context, id string) error {
	coll := r.client.Database(r.database).Collection(r.collection)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := coll.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
