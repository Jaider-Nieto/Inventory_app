package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jaider-nieto/ecommerce-go/user-service/handlers"
	"github.com/jaider-nieto/ecommerce-go/user-service/middlewares"
	"github.com/jaider-nieto/ecommerce-go/user-service/models"
	"github.com/jaider-nieto/ecommerce-go/user-service/repository"
	"gorm.io/gorm"
)

func Routes(db *gorm.DB) *mux.Router {
	r := mux.NewRouter()

	//Inicializa los repositorios.
	userReposiroy := repository.NewUserRepository(db)

	//Inicializa los handlers.
	handlerUsers := handlers.NewUserHandler(userReposiroy)

	//Rutas User.
	r.Handle("/register", middlewares.ValidationMiddleware(http.HandlerFunc(handlerUsers.RegisterUserHandlder), &models.User{})).Methods("POST")

	r.Handle("/login", middlewares.ValidationMiddleware(http.HandlerFunc(handlerUsers.LoginUserHanlder), &models.UserLogin{})).Methods("POST")

	r.HandleFunc("/users", handlerUsers.GetUsersHandler).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", handlerUsers.GetUserHandler).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", handlerUsers.DeleteUserHandler).Methods("DELETE")
	r.HandleFunc("/users/{id:[0-9]+}", handlerUsers.PatchUserHandler).Methods("PATCH")

	return r
}
