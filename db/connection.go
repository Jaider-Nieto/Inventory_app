package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DSN = "host=localhost user=jaider password=1005716614 dbname=ecommerce port=5432"
var DB *gorm.DB

func DBConnection() {
	var err error

	DB, err = gorm.Open(postgres.Open(DSN), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}
	
	log.Println("DB connected")
}
