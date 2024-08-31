package interfaces

import "net/http"

type UserHandelrInterface interface {
	GetUsersHandler(w http.ResponseWriter, r *http.Request)
	GetUserHandler(w http.ResponseWriter, r *http.Request)
	PostUserHandler(w http.ResponseWriter, r *http.Request)
	DeleteUserHandler(w http.ResponseWriter, r *http.Request)
}
