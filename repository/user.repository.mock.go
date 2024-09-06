package repository

import (
	"errors"

	"github.com/jaider-nieto/ecommerce-go/models"
	"gorm.io/gorm"
)

type UserRepositoryMocked struct{}

func (rm *UserRepositoryMocked) FindAllUsers() ([]models.User, error) {
	return []models.User{
		{
			FirstName: "Jaider",
			LastName:  "Nieto",
			Email:     "email@example.com",
			Password:  "hashPassword",
		},
		{
			FirstName: "Augusto",
			LastName:  "Criollo",
			Email:     "email2@example.com",
			Password:  "hashPassword",
		},
	}, nil
}
func (rm *UserRepositoryMocked) FindUserByID(id string) (models.User, error) {
	if id == "1" {
		user := models.User{
			Model:     gorm.Model{ID: 1},
			FirstName: "Jaider",
			LastName:  "Nieto",
			Email:     "email@example.com",
			Password:  "hashPassword",
		}
		return user, nil
	}

	return models.User{}, errors.New("user not found")
}
func (rm *UserRepositoryMocked) CreateUser(user models.User) (models.User, error) {
	userCreated := models.User{
		Model:     gorm.Model{ID: 1},
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  "hashPassword",
		Tasks:     nil,
	}
	return userCreated, nil
}
func (rm *UserRepositoryMocked) FindUserByEmail(email string) (models.User, error) {
	return models.User{}, nil
}
func (rm *UserRepositoryMocked) DeleteUser(id string) error {
	return nil
}
func (rm *UserRepositoryMocked) UpdateUser(user models.User) error {
	return nil
}
