package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jaider-nieto/ecommerce-go/products-service/docs"
	"github.com/jaider-nieto/ecommerce-go/products-service/internal/config"
	"github.com/jaider-nieto/ecommerce-go/products-service/internal/routes"
	"github.com/joho/godotenv"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// @title Tag Service API
// @version 1.0
// @description Tag service API in Go using Gin framework
// @host localhost:8082
// @basePath /
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	c := config.NewContainer()
	router := gin.Default()
	router.Use(gin.Logger())

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.ProductRoutes(router, c)
	router.Run(":" + os.Getenv("PORT"))
}
