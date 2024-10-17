package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jaider-nieto/ecommerce-go/products-service/internal/controller"
)

func ProductRoutes(router *gin.Engine, productsController *controller.ProductController) {
	productGroup := router.Group("/products")
	{
		productGroup.GET("/", productsController.GetProducts)
		productGroup.GET("/:user_id", productsController.GetProduct)
		productGroup.POST("/", productsController.PostProduct)
		productGroup.PATCH("/:user_id", productsController.UpdateProduct)
		productGroup.DELETE("/:user_id", productsController.DeleteProduct)
	}

}
