package routes

import (
	"github.com/gorilla/mux"
	"github.com/jaider-nieto/ecommerce-go/handlers"
)

func Routes() *mux.Router {
	r := mux.NewRouter()

	//Rutas User.
	r.HandleFunc("/users", handlers.GetUsersHandler).Methods("GET")
	r.HandleFunc("/users/{id}", handlers.GetUserHandler).Methods("GET")
	r.HandleFunc("/users", handlers.PostUserHanlder).Methods("POST")
	r.HandleFunc("/users/{id}", handlers.DeleteUserHandler).Methods("DELETE")

	//Rutas Task.
	r.HandleFunc("/tasks", handlers.GetTasksHandler).Methods("GET")
	r.HandleFunc("/tasks/{id}", handlers.GetTaskHandler).Methods("GET")
	r.HandleFunc("/tasks", handlers.PostTasksHandler).Methods("POST")
	r.HandleFunc("/tasks/{id}", handlers.DeleteTasksHandler).Methods("DELETE")

	return r
}
