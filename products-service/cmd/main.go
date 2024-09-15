package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jaider-nieto/ecommerce-go/products-service/internal/config"
	"github.com/jaider-nieto/ecommerce-go/products-service/internal/controller"
	"github.com/jaider-nieto/ecommerce-go/products-service/internal/repository"
	"github.com/jaider-nieto/ecommerce-go/products-service/internal/routes"
	"github.com/jaider-nieto/ecommerce-go/products-service/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
        log.Fatalf("Error loading .env file")
    }

    // Validar las variables de entorno
    // config.ValidateEnvVars()

	mongoURI := os.Getenv("MONGO_URI")
	client := config.InitMongoDB(mongoURI)

	productRepository := repository.NewProductRepository(config.GetMongoCollection(client, "products_db", "products"))
	productService := service.NewProductService(productRepository)
	productController := controller.NewProductController(productService)

	router := gin.Default()

	routes.ProductRoutes(router, productController)
	router.Run(":8082")
}
