package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jaider-nieto/ecommerce-go/db"
	"github.com/jaider-nieto/ecommerce-go/models"
	"github.com/jaider-nieto/ecommerce-go/routes"
	"github.com/joho/godotenv"
)


func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err.Error())
	}

	db.DBConnection(os.Getenv("DSN"))

	db.DB.AutoMigrate(models.User{})
	db.DB.AutoMigrate(models.Task{})

	http.ListenAndServe(os.Getenv("PORT"), routes.Routes(db.DB))
}
