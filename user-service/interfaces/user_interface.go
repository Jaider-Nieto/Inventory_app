package interfaces

import (
	"net/http"

	"github.com/jaider-nieto/ecommerce-go/user-service/models"
)

type UserHandelrInterface interface {
	RegisterUserHandler(w http.ResponseWriter, r *http.Request)
	LoginUserHanlder(w http.ResponseWriter, r *http.Request)
	GetUsersHandler(w http.ResponseWriter, r *http.Request)
	GetUserHandler(w http.ResponseWriter, r *http.Request)
	DeleteUserHandler(w http.ResponseWriter, r *http.Request)
	PatchUserHandler(w http.ResponseWriter, r *http.Request)
}

type UserRepositoryInterface interface {
	FindAllUsers() ([]models.User, error)
	FindUserByID(id string) (models.User, error)
	FindUserByEmail(email string) (models.User, error)
	CreateUser(user models.User) (models.User, error)
	DeleteUser(id string) error
	UpdateUser(user models.User) error
}
