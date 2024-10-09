package service

import (
	"context"
	"log"

	"github.com/jaider-nieto/ecommerce-go/products-service/internal/models"
	"github.com/jaider-nieto/ecommerce-go/products-service/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductService struct {
	repository *repository.ProductRepository
	cache      *repository.ProductRedisRepository
}

func NewProductService(repository *repository.ProductRepository, cache *repository.ProductRedisRepository) *ProductService {
	return &ProductService{repository: repository, cache: cache}
}

func (s *ProductService) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	// Intenta obtener los productos del caché de Redis.
	cacheProducts, err := s.cache.GetAll(ctx, "products_all")
	if err != nil {
		return []models.Product{}, err
	}
	
	// Si hay productos en caché, se retornan.
	if len(cacheProducts) > 0 {
		log.Println("cache")
		return cacheProducts, nil
	}

	// Si no hay productos en caché, búscalos en la base de datos.
	products, err := s.repository.FindAll()
	if err != nil {
		return []models.Product{}, err
	}

	// Guarda los productos en Redis y maneja el error si lo hay.
	if err := s.cache.Set(ctx, "products_all", products); err != nil {
		return []models.Product{}, err
	}

	return products, nil
}
func (s *ProductService) GetOneProduct(ctx context.Context, id string) (*models.Product, error) {
	// Intenta obtener el producto del caché de Redis.
	cacheProduct, err := s.cache.GetOne(ctx, "product_"+id)
	if err != nil {
		return nil, err
	}
	// Si exite el producto en caché, se retorna.
	if cacheProduct != nil {
		log.Println("cache")
		return cacheProduct, nil
	}

	// Si no esta el producto en caché, se busca en la base de datos.
	product, err := s.repository.FindOne(id)
	if err != nil {
		return nil, err
	}

	// Guarda el producto en Redis y maneja el error si lo hay.
	if err := s.cache.Set(ctx, "product_"+id, product); err != nil {
		return nil, err
	}

	log.Println("db")
	return &product, nil
}
func (s *ProductService) CreateProduct(ctx context.Context, product models.Product) (*mongo.InsertOneResult, error) {
	return s.repository.Create(product)
}
