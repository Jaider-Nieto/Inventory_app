package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jaider-nieto/ecommerce-go/products-service/internal/models"
	"github.com/redis/go-redis/v9"
)

// ProductRedisRepository interactúa con Redis para el almacenamiento en caché de productos
type ProductRedisRepository struct {
	client *redis.Client
}

// NewProductRedisRepository inicializa un nuevo repositorio Redis para productos
func NewProductRedisRepository(client *redis.Client) *ProductRedisRepository {
	return &ProductRedisRepository{client: client}
}

// GetAll obtiene todos los productos del caché usando una clave
func (r *ProductRedisRepository) GetAll(ctx context.Context, key string) ([]models.Product, error) {
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		// No se encontró en el caché
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var products []models.Product
	// Deserializa los datos almacenados en caché a una lista de productos
	if err := json.Unmarshal([]byte(data), &products); err != nil {
		return nil, errors.New("failed to unmarshal products from cache")
	}

	return products, nil
}

// GetOne obtiene un solo producto del caché usando una clave
func (r *ProductRedisRepository) GetOne(ctx context.Context, key string) (*models.Product, error) {
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		// No se encontró en el caché
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var product models.Product
	// Deserializa los datos almacenados en caché a un producto
	if err := json.Unmarshal([]byte(data), &product); err != nil {
		return nil, errors.New("failed to unmarshal product from cache")
	}

	return &product, nil
}

// Set almacena en caché los datos proporcionados (producto) en Redis
func (r *ProductRedisRepository) Set(ctx context.Context, key string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return errors.New("failed to marshal payload to JSON")
	}
	// Almacena los datos en caché con una expiración de 60 segundos
	if err := r.client.Set(ctx, key, data, time.Second*60).Err(); err != nil {
		return err
	}
	return nil
}

// Clean elimina todos los datos almacenados en caché de Redis
func (r *ProductRedisRepository) Clean(ctx context.Context) error {
	if err := r.client.FlushAll(ctx).Err(); err != nil {
		return errors.New("failed to flush all Redis data")
	}
	return nil
}
