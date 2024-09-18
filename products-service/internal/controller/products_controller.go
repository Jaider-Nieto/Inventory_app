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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}
func (ctr *ProductController) GetProduct(c *gin.Context) {
	product, err := ctr.service.GetOneProduct(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}
func (ctrl *ProductController) PostProduct(c *gin.Context) {
	var product models.Product
	if err := c.BindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !product.IsValidCategory() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product category"})
		return
	}

	createdProduct, err := ctrl.service.CreateProduct(product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, createdProduct)
}
func (ctrl *ProductController) DeleteProduct(c *gin.Context){}
