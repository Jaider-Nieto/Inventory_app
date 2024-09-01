package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jaider-nieto/ecommerce-go/interfaces"
	"github.com/jaider-nieto/ecommerce-go/models"
	"github.com/jaider-nieto/ecommerce-go/utils"
	"github.com/jaider-nieto/ecommerce-go/validations"
	"gorm.io/gorm"
)

type userHandler struct {
	*Handler
}

func NewUserHandler(UserRepository interfaces.UserRepositoryInterface) *userHandler {
	return &userHandler{
		Handler: &Handler{
			userRepository: UserRepository,
		},
	}
}

func (h *userHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.userRepository.FindAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(users)
}
func (h *userHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	user, err := h.userRepository.FindUserByID(params["id"])
	if user.ID == 0 || err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}

	json.NewEncoder(w).Encode(&user)
}
func (h *userHandler) PostUserHanlder(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error decoding request: " + err.Error()))
		return
	}

	if !validations.EmailValidation(user.Email) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid email format"))
		return
	}

	userExist, err := h.userRepository.FindUserByEmail(user.Email)
	if userExist.Email == user.Email || err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error email duplicated"))
		return
	}

	hashPassword, hashErr := utils.HashPassword(user.Password)
	if hashErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error hash password: " + hashErr.Error()))
		return
	}

	user.Password = hashPassword

	user, dbErr := h.userRepository.CreateUser(user)
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
func (h *userHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var params = mux.Vars(r)

	user, err := h.userRepository.FindUserByID(params["id"])
	if user.ID == 0 || err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}

	deleteErr := h.userRepository.DeleteUser(params["id"])
	if deleteErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User not deleted: " + deleteErr.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted"))
}
func (h *userHandler) PatchUserHandler(w http.ResponseWriter, r *http.Request) {
	var params = mux.Vars(r)

	user, err := h.userRepository.FindUserByID(params["id"])
	if user.ID == 0 || err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}

	var input map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if firstName, ok := input["first_name"].(string); ok {
		user.FirstName = firstName
	}
	if lastName, ok := input["last_name"].(string); ok {
		user.LastName = lastName
	}
	if email, ok := input["email"].(string); ok {
		user.Email = email
	}

	if err := h.userRepository.UpdateUser(user); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&user)
}
