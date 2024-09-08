package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jaider-nieto/ecommerce-go/user-service/db"
	"github.com/jaider-nieto/ecommerce-go/user-service/models"
	"github.com/jaider-nieto/ecommerce-go/user-service/routes"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err.Error())
	}

	db.DBConnection(os.Getenv("DSN"))

	db.DB.AutoMigrate(models.User{})

	http.ListenAndServe(os.Getenv("PORT"), routes.Routes(db.DB))
}
