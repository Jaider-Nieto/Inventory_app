package auth

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jaider-nieto/ecommerce-go/models"
)

func CreateJWT(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"first_name": user.FirstName,
		"email":      user.Email,
	})

	return token.SignedString([]byte(os.Getenv("SECRET_JWT")))
}
