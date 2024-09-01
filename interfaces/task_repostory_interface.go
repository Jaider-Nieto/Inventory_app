package interfaces

import "github.com/jaider-nieto/ecommerce-go/models"

type TaskRepositoryInterface interface {
	FindAllTasks() ([]models.Task, error)
	FindTaskById(id string) (models.Task, error)
	CreateTask(task models.Task) (models.Task, error)
	DeleteTask(id string) error
	UpdateTask(task models.Task) error
}
