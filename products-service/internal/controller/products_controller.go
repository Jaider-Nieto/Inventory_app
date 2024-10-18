package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jaider-nieto/ecommerce-go/products-service/internal/models"
	"github.com/jaider-nieto/ecommerce-go/products-service/internal/service"
)

// ProductController maneja las solicitudes relacionadas con productos.
type ProductController struct {
	service service.ProductService // Servicio para manejar la lógica de negocio de productos
}

// NewProductController crea una nueva instancia de ProductController.
func NewProductController(service *service.ProductService) *ProductController {
	return &ProductController{service: *service}
}

// GetProduct maneja la solicitud para obtener un producto específico por ID.
// @Summary Get a product
// @Description Retrieve a product by user_id from the database
// @Tags products
// @Accept json
// @Produce json
// @Param user_id path string true "User ID" example("12345")
// @Success 200 {object} models.Product "Product data"
// @Failure 400 {object} map[string]string "Error description"
// @Failure 404 {object} map[string]string "Product not found"
// @Router /products/{user_id} [get]
func (ctrl *ProductController) GetProducts(c *gin.Context) {
	page := c.DefaultQuery("page", "1")      // Si no se pasa el parámetro, usa 1
	pageSize := c.DefaultQuery("size", "10") // Si no se pasa el parámetro, usa 10

	// Convierte los parámetros a enteros
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeInt < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid size parameter"})
		return
	}

	// Llama al servicio para obtener todos los productos
	products, err := ctrl.service.GetAllProducts(c.Request.Context(), pageInt, pageSizeInt)
	if err != nil {
		// Retorna un error si ocurre al obtener productos
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products) // Retorna los productos en formato JSON
}

// GetProduct maneja la solicitud para obtener un producto específico por ID.
// @Summary Get a product
// @Description Retrieve a product by user_id from the database
// @Tags products
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]string "error"
// @Router /products/{user_id} [get]
func (ctr *ProductController) GetProduct(c *gin.Context) {
	// Llama al servicio para obtener un producto por ID
	product, err := ctr.service.GetOneProduct(c.Request.Context(), c.Param("user_id"))
	if err != nil {
		// Retorna un error si ocurre al obtener el producto
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product) // Retorna el producto en formato JSON
}

// PostProduct maneja la solicitud para crear un nuevo producto.
// @Summary Create product
// @Description Create a new product in MongoDB
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.Product true "Product Data"
// @Success 200 {object} string
// @Failure 400 {object} error "error"
// @Router /products [post]
func (ctrl *ProductController) PostProduct(c *gin.Context) {
	var product models.Product

	// Vincula el cuerpo de la solicitud a la estructura Product
	if err := c.BindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Retorna error si falla la vinculación
		return
	}

	// Verifica si la categoría del producto es válida
	if !product.IsValidCategory() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product category"})
		return
	}

	// Llama al servicio para crear el producto
	err := ctrl.service.CreateProduct(c.Request.Context(), product)
	if err != nil {
		// Retorna un error si ocurre al crear el producto
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "created product")
}

// DeleteProduct maneja la solicitud para borrar un producto.
// @Summary Delete a product
// @Description Delete a product by its user_id
// @Tags products
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} string "deleted product"
// @Failure 400 {object} map[string]string "error"
// @Router /products/{user_id} [delete]
func (ctrl *ProductController) DeleteProduct(c *gin.Context) {
	// Llama al servicio para borrar el producto utilizando el user_id del parámetro de la URL
	if err := ctrl.service.DeleteProduct(c.Request.Context(), c.Param("user_id")); err != nil {
		// Si ocurre un error al borrar, retorna un mensaje de error con el código 400
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// Si el borrado es exitoso, retorna un mensaje de confirmación con el código 200
	c.JSON(http.StatusOK, "deleted product")
}

// UpdateProduct maneja la solicitud para actualizar un producto.
// @Summary Update a product
// @Description Update a product's fields using its user_id
// @Tags products
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param updates body map[string]interface{} true "Product fields to update"
// @Success 200 {object} string "updated product"
// @Failure 400 {object} map[string]string "error"
// @Router /products/{user_id} [patch]
func (ctrl *ProductController) UpdateProduct(c *gin.Context) {
	// Declara un mapa para almacenar los campos a actualizar enviados en el cuerpo de la solicitud
	var updates map[string]interface{}

	// Intenta enlazar los datos JSON del cuerpo de la solicitud al mapa
	if err := c.ShouldBindJSON(&updates); err != nil {
		// Si hay un error en el formato del cuerpo, retorna un error con el código 400
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Llama al servicio para actualizar el producto utilizando el user_id y los campos proporcionados
	if err := ctrl.service.UpdateProduct(c.Request.Context(), c.Param("user_id"), updates); err != nil {
		// Si ocurre un error al actualizar, retorna un mensaje de error con el código 400
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// Si la actualización es exitosa, retorna un mensaje de confirmación con el código 200
	c.JSON(http.StatusOK, "updated product")
}
