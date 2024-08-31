package routes

import (
	"github.com/gorilla/mux"
	"github.com/jaider-nieto/ecommerce-go/handlers"
	"github.com/jaider-nieto/ecommerce-go/repository"
	"gorm.io/gorm"
)

func Routes(db *gorm.DB) *mux.Router {
	r := mux.NewRouter()

	repo := repository.NewRepository(db)

	userReposiroy := repository.NewUserRepository(repo)

	handlerUsers := handlers.NewUserHandler(userReposiroy)

	//Rutas User.
	r.HandleFunc("/users", handlerUsers.GetUsersHandler).Methods("GET")
	r.HandleFunc("/users/{id}", handlerUsers.GetUserHandler).Methods("GET")
	r.HandleFunc("/users", handlerUsers.PostUserHanlder).Methods("POST")
	r.HandleFunc("/users/{id}", handlerUsers.DeleteUserHandler).Methods("DELETE")

	//Rutas Task.
	r.HandleFunc("/tasks", handlers.GetTasksHandler).Methods("GET")
	r.HandleFunc("/tasks/{id}", handlers.GetTaskHandler).Methods("GET")
	r.HandleFunc("/tasks", handlers.PostTasksHandler).Methods("POST")
	r.HandleFunc("/tasks/{id}", handlers.DeleteTasksHandler).Methods("DELETE")

	return r
}
