package repository

import (
	"github.com/jaider-nieto/ecommerce-go/models"
)

type TaskRepository struct {
	*Repository
}

func NewTaskRepository(repository *Repository) *TaskRepository {
	return &TaskRepository{Repository: repository}
}

func (r *TaskRepository) FindAllTasks() ([]models.Task, error) {
	var tasks []models.Task

	err := r.DB.Find(&tasks).Error

	return tasks, err
}
func (r TaskRepository) FindTaskById(id string) (models.Task, error) {
	var task models.Task

	err := r.DB.First(&task, id).Error

	return task, err
}
func (r TaskRepository) CreateTask(task models.Task) (models.Task, error) {
	err := r.DB.Create(&task).Error

	return task, err
}
func (r TaskRepository) DeleteTask(id string) error {
	var task models.Task

	err := r.DB.Delete(&task, id).Error

	return err
}
func (r TaskRepository) UpdateTask(task models.Task) error {
	err := r.DB.Save(&task).Error

	return err
}
