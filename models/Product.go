package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string
	Price       uint `gorm:"not null"`
	Stock       uint `gorm:"not null; default:0"`
}
