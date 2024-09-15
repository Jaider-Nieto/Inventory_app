package repository

import (
	"errors"

	"github.com/jaider-nieto/ecommerce-go/user-service/models"
	"gorm.io/gorm"
)

type UserRepositoryMocked struct {
	ShouldReturnError bool
}

func (rm *UserRepositoryMocked) FindAllUsers() ([]models.User, error) {
	if rm.ShouldReturnError {
		return nil, errors.New("internal server error")
	}
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
	if rm.ShouldReturnError && id == "1" {
		return models.User{}, errors.New("internal server error")
	}
	if id == "1" || id == "2" {
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

	if rm.ShouldReturnError && user.ID == 2 {
		return models.User{}, errors.New("internal server error")
	}

	userCreated := models.User{
		Model:     gorm.Model{ID: 1},
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  "hashPassword",
	}
	return userCreated, nil
}
func (rm *UserRepositoryMocked) FindUserByEmail(email string) (models.User, error) {
	if rm.ShouldReturnError && email != "email@valid.com" {
		return models.User{}, errors.New("internal server error")
	}
	if email == "email@valid.com" {
		return models.User{
			Model:     gorm.Model{ID: 1},
			FirstName: "Jaider",
			LastName:  "Nieto",
			Email:     "email@example.com",
			Password:  "$2a$10$pPGhl2x0uUR4QkKKMnQWz.JzSTkzI7.SNyGn7iW8cCYNByFUeGdq2",
		}, nil
	}

	return models.User{}, errors.New("record not found")
}
func (rm *UserRepositoryMocked) DeleteUser(id string) error {
	if rm.ShouldReturnError || id == "1" {
		return errors.New("internal server error")
	}
	if id == "2" {
		return nil
	}

	return errors.New("user not found")
}
func (rm *UserRepositoryMocked) UpdateUser(user models.User) error {
	if rm.ShouldReturnError {
		return errors.New("internal server error")
	}
	return nil
}
