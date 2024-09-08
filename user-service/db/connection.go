package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConnection(DNS string) {
	var err error

	DB, err = gorm.Open(postgres.Open(DNS), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB connected")
}
