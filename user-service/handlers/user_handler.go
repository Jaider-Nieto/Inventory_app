package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jaider-nieto/ecommerce-go/user-service/interfaces"
	"github.com/jaider-nieto/ecommerce-go/user-service/models"
	"github.com/jaider-nieto/ecommerce-go/user-service/utils"
	"golang.org/x/crypto/bcrypt"
)

type userHandler struct {
	userRepository interfaces.UserRepositoryInterface
}

func NewUserHandler(UserRepository interfaces.UserRepositoryInterface) *userHandler {
	return &userHandler{
		userRepository: UserRepository,
	}
}

func (h *userHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.userRepository.FindAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
func (h *userHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	user, err := h.userRepository.FindUserByID(params["id"])
	if err != nil {
		if err.Error() == "user not found" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(&user)
}
func (h *userHandler) RegisterUserHandlder(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error decoding request: " + err.Error()))
		return
	}

	userExist, err := h.userRepository.FindUserByEmail(user.Email)
	if userExist.Email == user.Email || err != nil && err.Error() != "record not found" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}

	hashPassword, hashErr := utils.HashPassword(user.Password)
	if hashErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error hash password: " + hashErr.Error()))
		return
	}

	user.Password = hashPassword

	user, dbErr := h.userRepository.CreateUser(user)
	if dbErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&user)
}
func (h *userHandler) LoginUserHanlder(w http.ResponseWriter, r *http.Request) {
	var userLogin models.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&userLogin); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error decoding request: " + err.Error()))
		return
	}

	user, err := h.userRepository.FindUserByEmail(userLogin.Email)
	if err != nil {
		if err.Error() == "record not found" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("incorret password"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("user login"))
}
func (h *userHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var params = mux.Vars(r)

	_, err := h.userRepository.FindUserByID(params["id"])
	if err != nil {
		if err.Error() == "user not found" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	deleteErr := h.userRepository.DeleteUser(params["id"])
	if deleteErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(deleteErr.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("user deleted"))
}
func (h *userHandler) PatchUserHandler(w http.ResponseWriter, r *http.Request) {
	var params = mux.Vars(r)

	user, err := h.userRepository.FindUserByID(params["id"])
	if err != nil {
		if err.Error() == "user not found" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&user)
}
