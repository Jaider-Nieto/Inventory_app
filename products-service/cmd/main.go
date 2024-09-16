package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jaider-nieto/ecommerce-go/products-service/internal/config"
	"github.com/jaider-nieto/ecommerce-go/products-service/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	c := config.NewContainer()
	router := gin.Default()

	routes.ProductRoutes(router, c)
	router.Run(":" + os.Getenv("PORT"))
}
