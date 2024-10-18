package service

import (
	"context"
	"strconv"

	"github.com/jaider-nieto/ecommerce-go/products-service/internal/interfaces"
	"github.com/jaider-nieto/ecommerce-go/products-service/internal/models"
)

type ProductService struct {
	repository interfaces.ProductMongoRepositoryInterface
	cache      interfaces.ProductRedisRepositoryInterface
}

func NewProductService(repository interfaces.ProductMongoRepositoryInterface, cache interfaces.ProductRedisRepositoryInterface) *ProductService {
	return &ProductService{repository: repository, cache: cache}
}

func (s *ProductService) GetAllProducts(ctx context.Context, page, size int) ([]models.Product, error) {
	pageStr := strconv.Itoa(page)

	// Intenta obtener los productos del caché de Redis.
	cacheProducts, err := s.cache.GetAll(ctx, "products_all_page_"+pageStr)
	if err != nil {
		return []models.Product{}, err
	}

	// Si hay productos en caché, se retornan.
	if len(cacheProducts) > 0 {
		return cacheProducts, nil
	}

	// Si no hay productos en caché, búscalos en la base de datos.
	products, err := s.repository.FindAll(ctx, page, size)
	if err != nil {
		return []models.Product{}, err
	}

	// Guarda los productos en Redis y maneja el error si lo hay.
	if err := s.cache.Set(ctx, "products_all_page_"+pageStr, products); err != nil {
		return []models.Product{}, err
	}

	return products, nil
}

func String(page int) {
	panic("unimplemented")
}
func (s *ProductService) GetOneProduct(ctx context.Context, id string) (*models.Product, error) {
	// Intenta obtener el producto del caché de Redis.
	cacheProduct, err := s.cache.GetOne(ctx, "product_"+id)
	if err != nil {
		return nil, err
	}
	// Si exite el producto en caché, se retorna.
	if cacheProduct != nil {
		return cacheProduct, nil
	}

	// Si no esta el producto en caché, se busca en la base de datos.
	product, err := s.repository.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}

	// Guarda el producto en Redis y maneja el error si lo hay.
	if err := s.cache.Set(ctx, "product_"+id, product); err != nil {
		return nil, err
	}

	return product, nil
}
func (s *ProductService) CreateProduct(ctx context.Context, product models.Product) error {
	// Limpia el cache existente y maneja el error si lo hay.
	if err := s.cache.Clean(ctx); err != nil {
		return err
	}

	return s.repository.Create(ctx, product)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	// Elimina el producto en la base de datos y maneja el error si lo hay.
	if err := s.repository.Delete(ctx, id); err != nil {
		return err
	}

	// Limpia el cache existente y maneja el error si lo hay.
	if err := s.cache.Clean(ctx); err != nil {
		return err
	}

	return nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, id string, product map[string]interface{}) error {
	if err := s.repository.Update(ctx, id, product); err != nil {
		return err
	}

	// Limpia el cache existente y maneja el error si lo hay.
	if err := s.cache.Clean(ctx); err != nil {
		return err
	}

	return nil
}
