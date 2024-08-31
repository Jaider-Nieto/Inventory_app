package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jaider-nieto/ecommerce-go/interfaces"
	"github.com/jaider-nieto/ecommerce-go/models"
)

type taskHandler struct {
	*Handler
}

func NewTaskHandler(TaskRepository interfaces.TaskRepositoryInterface) *taskHandler {
	return &taskHandler{
		Handler: &Handler{
			taskRepository: TaskRepository,
		},
	}
}

func (h *taskHandler) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.taskRepository.FindAllTasks()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Tasks not found"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}
func (h *taskHandler) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	var params = mux.Vars(r)

	task, err := h.taskRepository.FindTaskById(params["id"])
	if task.ID == 0 || err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Task not found"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&task)
}
func (h *taskHandler) PostTasksHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	json.NewDecoder(r.Body).Decode(&task)

	task, err := h.taskRepository.CreateTask(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Task not created: " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}
func (h *taskHandler) DeleteTasksHandler(w http.ResponseWriter, r *http.Request) {
	var params = mux.Vars(r)

	err := h.taskRepository.DeleteTask(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Task not deleted: " + err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}
