package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(email string) (string, error) {
	claims := jwt.MapClaims{
		"sub": email,
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte("SECRET_JWT"))
}
