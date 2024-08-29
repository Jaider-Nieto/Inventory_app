package main

import (
	"log"
	"net/http"

	"github.com/jaider-nieto/ecommerce-go/db"
	"github.com/jaider-nieto/ecommerce-go/models"
	"github.com/jaider-nieto/ecommerce-go/routes"
)

func main() {
	db.DBConnection()

	db.DB.AutoMigrate(models.User{}, models.Task{})

	err := http.ListenAndServe(":3001", routes.Routes())
	if err != nil {
		log.Fatal("Failed to start the server.")
	}
}
