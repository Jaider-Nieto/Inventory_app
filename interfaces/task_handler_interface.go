package interfaces

import "net/http"

type TaskHandlerInterface interface {
	GetTasksHandler(w http.ResponseWriter, r *http.Request)
	GetTaskHandler(w http.ResponseWriter, r *http.Request)
	PostTaskHandler(w http.ResponseWriter, r *http.Request)
	DeleteTaskHandler(w http.ResponseWriter, r *http.Request)
	PatchTaskHandler(w http.ResponseWriter, r *http.Request)
}