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
	r.HandleFunc("/users/{id:[0-9]+}", handlerUsers.GetUserHandler).Methods("GET")
	r.HandleFunc("/users", handlerUsers.PostUserHanlder).Methods("POST")
	r.HandleFunc("/users/{id:[0-9]+}", handlerUsers.DeleteUserHandler).Methods("DELETE")
	r.HandleFunc("/users/{id:[0-9]+}", handlerUsers.PatchUserHandler).Methods("PATCH")

	//Rutas Task.
	r.HandleFunc("/tasks", handlerTask.GetTasksHandler).Methods("GET")
	r.HandleFunc("/tasks/{id:[0-9]+}", handlerTask.GetTaskHandler).Methods("GET")
	r.HandleFunc("/tasks", handlerTask.PostTasksHandler).Methods("POST")
	r.HandleFunc("/tasks/{id:[0-9]+}", handlerTask.DeleteTasksHandler).Methods("DELETE")

	return r
}
