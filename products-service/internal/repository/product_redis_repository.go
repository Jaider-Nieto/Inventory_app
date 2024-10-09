package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jaider-nieto/ecommerce-go/products-service/internal/models"
	"github.com/redis/go-redis/v9"
)

type ProductRedisRepository struct {
	client *redis.Client
}

func NewProductRedisRepository(client *redis.Client) *ProductRedisRepository {
	return &ProductRedisRepository{client: client}
}

func (r *ProductRedisRepository) GetAll(ctx context.Context, key string) ([]models.Product, error) {
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var products []models.Product
	if err := json.Unmarshal([]byte(data), &products); err != nil {
		return nil, errors.New("failed to unmarshal products from cache")
	}

	return products, nil
}
func (r *ProductRedisRepository) GetOne(ctx context.Context, key string) (*models.Product, error) {
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var product models.Product
	if err := json.Unmarshal([]byte(data), &product); err != nil {
		return nil, err
	}

	return &product, nil
}
func (r *ProductRedisRepository) Set(ctx context.Context, key string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return errors.New("failed to marshal payload to JSON")
	}
	if err := r.client.Set(ctx, key, data, time.Second*60).Err(); err != nil {
		return err
	}
	return nil
}
