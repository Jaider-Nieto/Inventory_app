package config

import (
	"os"

	"github.com/jaider-nieto/ecommerce-go/products-service/internal/controller"
	"github.com/jaider-nieto/ecommerce-go/products-service/internal/repository"
	"github.com/jaider-nieto/ecommerce-go/products-service/internal/service"
)

func NewContainer() *controller.ProductController {

	mongoURI := os.Getenv("MONGO_URI")
	client := InitMongoDB(mongoURI)

	productRepository := repository.NewProductRepository(GetMongoCollection(client, "products_db", "products"))
	productService := service.NewProductService(productRepository)
	productController := controller.NewProductController(productService)

	return productController
}
