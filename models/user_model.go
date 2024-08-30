package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName string `gorm:"not null" json:"first_name"`
	LastName  string `gorm:"not null" json:"last_name"`
	Email     string `gorm:"not null;unique" json:"email"`
	Password  string `gorm:"not null" json:"password"`
	Tasks     []Task `json:"tasks"`
}
