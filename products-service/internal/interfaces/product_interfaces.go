package interfaces

import (
	"context"

	"github.com/jaider-nieto/ecommerce-go/products-service/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// ProductMongoRepositoryInterface define los métodos para interactuar con productos en MongoDB.
type ProductMongoRepositoryInterface interface {
	// FindAll retorna una lista de todos los productos.
	FindAll() ([]models.Product, error)
	// FindOne busca un producto por su ID y lo retorna.
	FindOne(id string) (*models.Product, error)
	// Create inserta un nuevo producto y retorna el resultado de la operación.
	Create(product models.Product) (*mongo.InsertOneResult, error)
}

// ProductRedisRepositoryInterface define métodos para interactuar con un repositorio de productos en cache (Redis).
type ProductRedisRepositoryInterface interface {
	// GetAll recupera todos los productos del cache usando la clave especificada.
	GetAll(ctx context.Context, key string) ([]models.Product, error)

	// GetOne recupera un producto específico del cache usando la clave proporcionada.
	GetOne(ctx context.Context, key string) (*models.Product, error)

	// Set almacena un producto en el cache (Redis) bajo la clave especificada.
	Set(ctx context.Context, key string, payload interface{}) error
}
