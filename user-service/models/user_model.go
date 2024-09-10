package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName string `gorm:"not null" json:"first_name" validate:"required"`
	LastName  string `gorm:"not null" json:"last_name" validate:"required"`
	Email     string `gorm:"not null;unique" json:"email" validate:"required,email"`
	Password  string `gorm:"not null" json:"password" validate:"required,min=8"`
}

type UserUpdate struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
}

type UserLogin struct {
	Email    string `gorm:"not null;unique" json:"email" validate:"required,email"`
	Password string `gorm:"not null" json:"password" validate:"required,min=8"`
}
