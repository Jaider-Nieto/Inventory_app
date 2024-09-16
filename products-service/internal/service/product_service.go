package service

import (
	"github.com/jaider-nieto/ecommerce-go/products-service/internal/models"
	"github.com/jaider-nieto/ecommerce-go/products-service/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductService struct {
	repository *repository.ProductRepository
}

func NewProductService(repository *repository.ProductRepository) *ProductService {
	return &ProductService{repository: repository}
}

func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	return s.repository.FindAll()
}
func (s *ProductService) GetOneProduct(id string) (models.Product, error) {
	return s.repository.FindOne(id)
}
func (s *ProductService) CreateProduct(product models.Product) (*mongo.InsertOneResult, error) {
	return s.repository.Create(product)
}
