package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jaider-nieto/ecommerce-go/products-service/internal/models"
	"github.com/jaider-nieto/ecommerce-go/products-service/internal/service"
)

type ProductController struct {
	service service.ProductService
}

func NewProductController(service *service.ProductService) *ProductController {
	return &ProductController{service: *service}
}

func (ctrl *ProductController) GetProducts(c *gin.Context) {
	products, err := ctrl.service.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, products)
}
func (ctrl *ProductController) PostProduct(c *gin.Context) {
	var product models.Products
	if err := c.BindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	createdProduct, err := ctrl.service.CreateProduct(product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, createdProduct)
}
