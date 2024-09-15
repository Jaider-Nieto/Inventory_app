package test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jaider-nieto/ecommerce-go/user-service/handlers"
	"github.com/jaider-nieto/ecommerce-go/user-service/models"
	"github.com/jaider-nieto/ecommerce-go/user-service/repository"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func setup() {
	if err := godotenv.Load("./../.env"); err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}

	dsn := os.Getenv("DSN_TEST")
	log.Printf("%v", dsn)
	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to conected database")
	}

	db.AutoMigrate(&models.User{})
}

func cleanUp() {
	db.Exec("TRUNCATE TABLE users RESTART IDENTITY")
}

func TestGetUsers(t *testing.T) {
	setup()
	defer cleanUp()

	userRepository := repository.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepository)

	users := []models.User{
		{
			FirstName: "Jaider",
			LastName:  "Nieto",
			Email:     "jaider123456@gmail.com",
			Password:  "jaider123",
		},
		{
			FirstName: "Augusto",
			LastName:  "Criollo",
			Email:     "augusto123456@gmail.com",
			Password:  "Augusto123",
		},
	}

	for _, user := range users {
		log.Printf("%v", users)
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("failed to create user: %v", err)
		}
	}

	req, _ := http.NewRequest(http.MethodGet, "/users", nil)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(userHandler.GetUsersHandler)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}
	var responseUsers []models.User
	if err := json.Unmarshal(rr.Body.Bytes(), &responseUsers); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	if len(responseUsers) != len(users) {
		t.Fatalf("expected %d users, got %d", len(users), len(responseUsers))
	}

	for i, user := range users {
		if responseUsers[i].FirstName != user.FirstName || responseUsers[i].LastName != user.LastName || responseUsers[i].Email != user.Email || responseUsers[i].Password != user.Password {
			t.Errorf("expected user %v, got %v", user, responseUsers[i])
		}
	}
}

func TestGetUser(t *testing.T) {
	setup()
	defer cleanUp()

	userRepository := repository.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepository)

	user := models.User{
		FirstName: "Jaider",
		LastName:  "Nieto",
		Email:     "jaiderlolxd@gmail.com",
		Password:  "jaider123",
	}

	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	req, _ := http.NewRequest(http.MethodGet, "/users/1", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "1",
	})

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(userHandler.GetUserHandler)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusCreated, rr.Code)
	}
	var responseUser models.User
	if err := json.Unmarshal(rr.Body.Bytes(), &responseUser); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}
	if responseUser.FirstName != user.FirstName || responseUser.LastName != user.LastName || responseUser.Email != user.Email || responseUser.Password != user.Password {
		t.Errorf("expected user %v, got %v", user, responseUser)
	}
}

func TestRegisterUser(t *testing.T) {
	setup()
	defer cleanUp()

	userRepository := repository.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepository)

	user := models.User{
		FirstName: "Jaider",
		LastName:  "Nieto",
		Email:     "jaiderlolxd@gmail.com",
		Password:  "jaider123",
	}

	body, err := json.Marshal(&user)
	if err != nil {
		t.Fatalf("failed to marshal response body: %v", err)

	}

	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(userHandler.RegisterUserHandlder)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		log.Println(rr.Body.String())
		t.Fatalf("expected status code %d, got %d", http.StatusCreated, rr.Code)
	}
	var responseUser models.User
	if err := json.Unmarshal(rr.Body.Bytes(), &responseUser); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}
	if responseUser.FirstName != user.FirstName || responseUser.LastName != user.LastName || responseUser.Email != user.Email {
		t.Errorf("expected user %v, got %v", user, responseUser)
	}
}

func TestLoginUser(t *testing.T) {
	setup()
	defer cleanUp()

	userRepository := repository.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepository)

	user := models.User{
		FirstName: "Jaider",
		LastName:  "Nieto",
		Email:     "jaiderlolxd@gmail.com",
		Password:  "$2a$10$pPGhl2x0uUR4QkKKMnQWz.JzSTkzI7.SNyGn7iW8cCYNByFUeGdq2",
	}

	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	userLogin := models.UserLogin{
		Email:    "jaiderlolxd@gmail.com",
		Password: "hashpassword",
	}
	body, err := json.Marshal(&userLogin)
	if err != nil {
		t.Fatalf("failed to marshal response body: %v", err)
	}

	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(userHandler.LoginUserHanlder)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		log.Println(rr.Body.String())
		t.Fatalf("expected status code %d, got %d", http.StatusCreated, rr.Code)
	}
	if rr.Body.String() != "user login" {
		t.Fatalf("failed to login user: %v", rr.Body.String())
	}
}

func TestDeleteUserHandler(t *testing.T) {
	setup()
	defer cleanUp()

	userRepository := repository.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepository)

	user := models.User{
		Model:     gorm.Model{ID: 1},
		FirstName: "Jaider",
		LastName:  "Nieto",
		Email:     "jaiderlolxd@gmail.com",
		Password:  "jaider123",
	}

	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	req, _ := http.NewRequest(http.MethodDelete, "/users/1", nil)

	req = mux.SetURLVars(req, map[string]string{
		"id": "1",
	})

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(userHandler.DeleteUserHandler)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		log.Println(rr.Body.String())
		t.Fatalf("expected status code %d, got %d", http.StatusCreated, rr.Code)
	}
	if rr.Body.String() != "user deleted" {
		t.Fatalf("failed to delete user: %v", rr.Body.String())
	}
}

func TestUpdateUserHandler(t *testing.T) {
	setup()
	defer cleanUp()

	userRepository := repository.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepository)

	user := models.User{
		Model:     gorm.Model{ID: 1},
		FirstName: "Jaider",
		LastName:  "Nieto",
		Email:     "jaiderlolxd@gmail.com",
		Password:  "jaider123",
	}

	updatedUser := models.UserUpdate{
		FirstName: "Augusto",
		LastName:  "Criollo",
		Email:     "jaider2002@gmail.com",
	}

	body, err := json.Marshal(&updatedUser)
	if err != nil {
		t.Fatalf("failed to marshal response body: %v", err)
	}

	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	req, _ := http.NewRequest(http.MethodPatch, "/users/1", bytes.NewBuffer(body))

	req = mux.SetURLVars(req, map[string]string{
		"id": "1",
	})

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(userHandler.PatchUserHandler)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		log.Println(rr.Body.String())
		t.Fatalf("expected status code %d, got %d", http.StatusCreated, rr.Code)
	}
	var responseUser models.User
	if err := json.Unmarshal(rr.Body.Bytes(), &responseUser); err != nil {
		log.Println(rr.Body.String())
		t.Fatalf("failed to unmarshal response body: %v", err)
	}
	if responseUser.FirstName != updatedUser.FirstName ||
		responseUser.LastName != updatedUser.LastName ||
		responseUser.Email != updatedUser.Email {
		t.Errorf("expected user %v, got %v", updatedUser, responseUser)
	}
}
