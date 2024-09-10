package test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

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

	req, _ := http.NewRequest(http.MethodGet, "/user", nil)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(userHandler.GetUsersHandler)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusCreated, rr.Code)
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
