package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jaider-nieto/ecommerce-go/products-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ProductRepository gestiona las operaciones de base de datos relacionadas con los productos
type ProductRepository struct {
	collection *mongo.Collection
}

// NewProductRepository crea una nueva instancia de ProductRepository con la colección especificada
func NewProductRepository(collection *mongo.Collection) *ProductRepository {
	return &ProductRepository{collection: collection}
}

// FindAll obtiene todos los productos de la colección
func (r *ProductRepository) FindAll() ([]models.Product, error) {
	var products []models.Product
	// Establece un contexto con timeout de 5 segundos para la operación de búsqueda
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Realiza la búsqueda de todos los productos
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	// Decodifica todos los productos encontrados en la variable 'products'
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

// FindOne busca un producto por su ID en la colección
func (r *ProductRepository) FindOne(id string) (*models.Product, error) {
	var product models.Product

	// Contexto con un límite de 5 segundos para la operación de búsqueda
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Convierte el ID de string a ObjectID de MongoDB
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// Error si el ID no es válido
		return nil, fmt.Errorf("invalid product ID: %v", err)
	}

	// Realiza la búsqueda del producto en la colección usando el ObjectID
	result := r.collection.FindOne(ctx, bson.M{"_id": objID})
	if result.Err() != nil {
		// Si no se encuentra ningún producto, devuelve un error descriptivo
		if result.Err() == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("product not found with ID: %s", id)
		}
		// Otros errores se manejan con más contexto
		return nil, fmt.Errorf("error finding product: %v", result.Err())
	}

	// Decodifica el documento encontrado en la estructura 'Product'
	if err := result.Decode(&product); err != nil {
		return nil, fmt.Errorf("error decoding product: %v", err)
	}

	// Retorna el producto encontrado
	return &product, nil
}

// Create inserta un nuevo producto en la colección y devuelve el resultado de la operación
func (r *ProductRepository) Create(product models.Product) (*mongo.InsertOneResult, error) {
	// Contexto con timeout de 5 segundos para la operación de inserción
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Inserta el producto en la colección
	result, err := r.collection.InsertOne(ctx, product)
	if err != nil {
		return nil, err
	}

	// Retorna el resultado de la inserción
	return result, nil
}
