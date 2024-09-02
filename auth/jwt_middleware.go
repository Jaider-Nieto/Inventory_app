package auth

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jaider-nieto/ecommerce-go/utils"
)

func AuthMiddelware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenString, err := utils.TokenSeparator(authHeader)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if method := t.Method; method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("invalid signing method")
			}
			secret := []byte(os.Getenv("SECRET_JWT"))
			return secret, nil
		})

		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid token"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
