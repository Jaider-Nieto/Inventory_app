package interfaces

import "github.com/jaider-nieto/ecommerce-go/models"

type UserRepository interface {
	FindAllUsers() ([]models.User, error)
	FindUserByID(id string) (models.User, error)
	FindUserByEmail(email string) (models.User, error)
	CreateUser(user *models.User) error
	DeleteUser(id string) error
}
