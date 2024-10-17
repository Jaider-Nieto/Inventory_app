package repository

import (
	"context"
	"fmt"

	"github.com/jaider-nieto/ecommerce-go/products-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ProductRepository gestiona las operaciones de base de datos relacionadas con los productos.
type ProductRepository struct {
	collection *mongo.Collection
}

// NewProductRepository crea una nueva instancia de ProductRepository con la colección especificada.
func NewProductRepository(collection *mongo.Collection) *ProductRepository {
	return &ProductRepository{collection: collection}
}

// FindAll obtiene todos los productos de la colección con paginación.
func (r *ProductRepository) FindAll(ctx context.Context, page, size int) ([]models.Product, error) {
	var products []models.Product

	// Calcula el número de documentos a omitir según la página y el tamaño.
	skip := (page - 1) * size

	// Realiza la búsqueda de productos en la colección, aplicando la paginación.
	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find().SetSkip(int64(skip)).SetLimit(int64(size)))
	if err != nil {
		return nil, fmt.Errorf("error finding products: %v", err)
	}
	defer cursor.Close(ctx) // Asegura que el cursor se cierre después de su uso.

	// Decodifica todos los productos encontrados en la variable 'products'.
	if err = cursor.All(ctx, &products); err != nil {
		return nil, fmt.Errorf("error decoding products: %v", err)
	}

	if len(products) == 0 {
		return nil, fmt.Errorf("no products found") // Mensaje informativo si no hay productos.
	}

	return products, nil
}

// FindOne busca un producto por su ID en la colección.
func (r *ProductRepository) FindOne(ctx context.Context, id string) (*models.Product, error) {
	var product models.Product

	// Convierte el ID de string a ObjectID de MongoDB.
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %v", err)
	}

	// Realiza la búsqueda del producto en la colección usando el ObjectID.
	result := r.collection.FindOne(ctx, bson.M{"_id": objID})
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("product not found with ID: %s", id)
		}
		return nil, fmt.Errorf("error finding product: %v", err)
	}

	// Decodifica el documento encontrado en la estructura 'Product'.
	if err := result.Decode(&product); err != nil {
		return nil, fmt.Errorf("error decoding product: %v", err)
	}

	return &product, nil
}

// Create inserta un nuevo producto en la colección y devuelve el resultado de la operación.
func (r *ProductRepository) Create(ctx context.Context, product models.Product) error {
	_, err := r.collection.InsertOne(ctx, product)
	if err != nil {
		return fmt.Errorf("error inserting product: %v", err) // Mensaje de error informativo.
	}

	return nil
}

// Delete elimina un producto por ID en la colección y devuelve un error si lo hay.
func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid product ID: %v", err)
	}

	// Elimina el producto en la colección
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return fmt.Errorf("error deleting product: %v", err)
	}
	//Si es 0 retorna error ya que no fue encontrado el producto.
	if result.DeletedCount == 0 {
		return fmt.Errorf("no product found with ID: %s", id)
	}

	return nil
}

func (r *ProductRepository) Update(ctx context.Context, id string, product map[string]interface{}) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid product ID: %v", err)
	}

	result, err := r.collection.UpdateByID(ctx, objID, bson.M{"$set": product})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("no product found with ID: %s", id)
	}

	return nil
}
