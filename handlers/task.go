package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jaider-nieto/ecommerce-go/db"
	"github.com/jaider-nieto/ecommerce-go/models"
)

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task
	db.DB.Find(&tasks)

	json.NewEncoder(w).Encode(tasks)

}
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	var params = mux.Vars(r)

	db.DB.First(&task, params["id"])

	if task.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Task not found"))
	}

	json.NewEncoder(w).Encode(&task)

}
func PostTasksHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	json.NewDecoder(r.Body).Decode(&task)

	createdTask := db.DB.Create(&task)

	err := createdTask.Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	json.NewEncoder(w).Encode(task)
}
func DeleteTasksHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	var params = mux.Vars(r)

	db.DB.Unscoped().Delete(&task, params["id"])

}
