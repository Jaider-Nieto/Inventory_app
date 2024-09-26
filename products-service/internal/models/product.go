package models

import (
	"github.com/jaider-nieto/ecommerce-go/products-service/pkg/utils"
)

// swagger:model
type Product struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	Category    string `json:"category" bson:"category"`
	Price       uint   `json:"price" bson:"price"`
	Stock       uint   `json:"stock" bson:"stock"`
	Rating      []uint `json:"rating" bson:"rating"`
}

func (p *Product) IsValidCategory() bool {
	return utils.IsValidCategory(p.Category)
}
