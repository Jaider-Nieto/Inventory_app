package config

import (
	"os"

	"github.com/jaider-nieto/ecommerce-go/products-service/internal/controller"
	"github.com/jaider-nieto/ecommerce-go/products-service/internal/repository"
	"github.com/jaider-nieto/ecommerce-go/products-service/internal/service"
)

func NewContainer() *controller.ProductController {

	mongoURI := os.Getenv("MONGO_URI")
	clientMongo := InitMongoDB(mongoURI)
	clientRedis := InitRedisClient(os.Getenv("REDIS_ADR"), os.Getenv("REDIS_PASSWORD"))

	productCacheRepository := repository.NewProductRedisRepository(clientRedis)
	productRepository := repository.NewProductRepository(GetMongoCollection(clientMongo, "products_db", "products"))
	productService := service.NewProductService(productRepository, productCacheRepository)
	productController := controller.NewProductController(productService)

	return productController
}
