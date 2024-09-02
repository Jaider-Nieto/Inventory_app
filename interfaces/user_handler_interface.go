package interfaces

import "net/http"

type UserHandelrInterface interface {
	RegisterUserHandler(w http.ResponseWriter, r *http.Request)
	LoginUserHanlder(w http.ResponseWriter, r *http.Request)
	GetUsersHandler(w http.ResponseWriter, r *http.Request)
	GetUserHandler(w http.ResponseWriter, r *http.Request)
	DeleteUserHandler(w http.ResponseWriter, r *http.Request)
	PatchUserHandler(w http.ResponseWriter, r *http.Request)
}
