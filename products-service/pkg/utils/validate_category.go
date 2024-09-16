package utils

import (
	"strings"

	"github.com/jaider-nieto/ecommerce-go/products-service/internal/constants"
)

func IsValidCategory(category string) bool {
	category = strings.ToLower(category)
	for _, c := range constants.AllowCategories {
		if c == category {
			return true
		}
	}
	return false
}
