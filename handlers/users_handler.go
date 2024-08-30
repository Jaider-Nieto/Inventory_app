package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jaider-nieto/ecommerce-go/db"
	"github.com/jaider-nieto/ecommerce-go/models"
	"golang.org/x/crypto/bcrypt"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []models.User

	db.DB.Find(&users)

	json.NewEncoder(w).Encode(users)
}
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	params := mux.Vars(r)

	db.DB.First(&user, params["id"])

	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}
	db.DB.Model(&user).Association("Tasks").Find(&user.Tasks)

	json.NewEncoder(w).Encode(&user)
}
func PostUserHanlder(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error decoding request: " + err.Error()))
		return
	}

	var userExist models.User
	db.DB.Where("email = ?", user.Email).First(&userExist)

	if userExist.Email == user.Email {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error email duplicated"))
		return
	}

	hash, hashErr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error generating password hash: " + hashErr.Error()))
		return
	}

	user.Password = string(hash)

	createdUser := db.DB.Create(&user)
	dbErr := createdUser.Error
	if dbErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error saving user to database: " + dbErr.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	encodeErr := json.NewEncoder(w).Encode(&user)

	if encodeErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error encoding response: " + encodeErr.Error()))
		return
	}
}
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var params = mux.Vars(r)

	db.DB.First(&user, params["id"])
	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}

	db.DB.Delete(&user)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted"))
}
