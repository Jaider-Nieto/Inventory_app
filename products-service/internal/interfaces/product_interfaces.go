package interfaces

import (
	"context"

	"github.com/jaider-nieto/ecommerce-go/products-service/internal/models"
)

// ProductMongoRepositoryInterface define los métodos para interactuar con productos en MongoDB.
type ProductMongoRepositoryInterface interface {
	// FindAll retorna una lista de todos los productos.
	FindAll(ctx context.Context, page, size int) ([]models.Product, error)
	// FindOne busca un producto por su ID y lo retorna.
	FindOne(ctx context.Context, id string) (*models.Product, error)
	// Create inserta un nuevo producto y retorna el resultado de la operación.
	Create(ctx context.Context, product models.Product) error
	// Delete Elimina un producto por su ID.
	Delete(ctx context.Context, id string) error
	// Exist Determina si un producto existe en la base de datos
	Update(ctx context.Context, id string, product map[string]interface{}) error
}

// ProductRedisRepositoryInterface define métodos para interactuar con un repositorio de productos en cache (Redis).
type ProductRedisRepositoryInterface interface {
	// GetAll recupera todos los productos del cache usando la clave especificada.
	GetAll(ctx context.Context, key string) ([]models.Product, error)

	// GetOne recupera un producto específico del cache usando la clave proporcionada.
	GetOne(ctx context.Context, key string) (*models.Product, error)

	// Set almacena un producto en el cache (Redis) bajo la clave especificada.
	Set(ctx context.Context, key string, payload interface{}) error

	// Clean elimina el cache existente
	Clean(ctx context.Context) error
}
