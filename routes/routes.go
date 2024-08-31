package routes

import (
	"github.com/gorilla/mux"
	"github.com/jaider-nieto/ecommerce-go/handlers"
	"github.com/jaider-nieto/ecommerce-go/repository"
	"gorm.io/gorm"
)

func Routes(db *gorm.DB) *mux.Router {
	r := mux.NewRouter()

	//Inicializa el repositorio base.
	repo := repository.NewRepository(db)

	//Inicializa los repositorios.
	userReposiroy := repository.NewUserRepository(repo)
	taskReposirory := repository.NewTaskRepository(repo)

	//Inicializa los handlers.
	handlerUsers := handlers.NewUserHandler(userReposiroy)
	handlerTask := handlers.NewTaskHandler(taskReposirory)

	//Rutas User.
	r.HandleFunc("/users", handlerUsers.GetUsersHandler).Methods("GET")
	r.HandleFunc("/users/{id}", handlerUsers.GetUserHandler).Methods("GET")
	r.HandleFunc("/users", handlerUsers.PostUserHanlder).Methods("POST")
	r.HandleFunc("/users/{id}", handlerUsers.DeleteUserHandler).Methods("DELETE")

	//Rutas Task.
	r.HandleFunc("/tasks", handlerTask.GetTasksHandler).Methods("GET")
	r.HandleFunc("/tasks/{id}", handlerTask.GetTaskHandler).Methods("GET")
	r.HandleFunc("/tasks", handlerTask.PostTasksHandler).Methods("POST")
	r.HandleFunc("/tasks/{id}", handlerTask.DeleteTasksHandler).Methods("DELETE")

	return r
}
