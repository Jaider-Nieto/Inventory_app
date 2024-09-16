package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/jaider-nieto/ecommerce-go/products-service/internal/models"
)

type ProductRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(collection *mongo.Collection) *ProductRepository {
	return &ProductRepository{collection: collection}
}

func (r *ProductRepository) FindAll() ([]models.Product, error) {
	var products []models.Product
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) FindOne(id string) (models.Product, error) {
	var product models.Product

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Product{}, err
	}

	result := r.collection.FindOne(ctx, bson.M{"_id": objID})
	if result.Err() != nil {
		return models.Product{}, result.Err()
	}

	result.Decode(&product)

	return product, nil
}

func (r *ProductRepository) Create(product models.Product) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	result, err := r.collection.InsertOne(ctx, product)
	if err != nil {
		return nil, err
	}

	return result, nil
}
