package interfaces

import (
	"net/http"

	"github.com/jaider-nieto/ecommerce-go/models"
)

type TaskHandlerInterface interface {
	GetTasksHandler(w http.ResponseWriter, r *http.Request)
	GetTaskHandler(w http.ResponseWriter, r *http.Request)
	PostTaskHandler(w http.ResponseWriter, r *http.Request)
	DeleteTaskHandler(w http.ResponseWriter, r *http.Request)
	PatchTaskHandler(w http.ResponseWriter, r *http.Request)
}

type TaskRepositoryInterface interface {
	FindAllTasks() ([]models.Task, error)
	FindTaskById(id string) (models.Task, error)
	CreateTask(task models.Task) (models.Task, error)
	DeleteTask(id string) error
	UpdateTask(task models.Task) error
}