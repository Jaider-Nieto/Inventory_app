package auth

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthLogin(c *gin.Context) {
	var creds Creds
	if err := c.Bind(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
	}
	// Validar credenciales con el Servicio de Usuarios
	if !validCreds(creds.Email, creds.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := CreateJWT(creds.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func validCreds(email, password string) bool {
	res, err := http.Post("http://localhost:8080/login", "application/json", bytes.NewBuffer([]byte(`{"email":"`+email+`", "password":"`+password+`"}`)))

	if err != nil || res.StatusCode != http.StatusOK {
		return false
	}
	return true
}
