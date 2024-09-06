package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"gorm.io/gorm"

	"github.com/gorilla/mux"
	"github.com/jaider-nieto/ecommerce-go/middlewares"
	"github.com/jaider-nieto/ecommerce-go/models"
	"github.com/jaider-nieto/ecommerce-go/repository"
)

func initHandlerUsers(t *testing.T) *userHandler {
	t.Helper()

	userRepositoryMock := &repository.UserRepositoryMocked{}
	return NewUserHandler(userRepositoryMock)
}

func TestGetUsersHandler(t *testing.T) {
	testCases := []struct {
		Name           string
		ExpectedStatus int
		ExpectedUsers  []models.User
	}{
		{
			Name:           "Get all users",
			ExpectedStatus: http.StatusOK,
			ExpectedUsers: []models.User{
				{
					FirstName: "Jaider",
					LastName:  "Nieto",
					Email:     "email@example.com",
					Password:  "hashPassword",
				},
				{
					FirstName: "Augusto",
					LastName:  "Criollo",
					Email:     "email2@example.com",
					Password:  "hashPassword",
				},
			},
		},
	}

	h := initHandlerUsers(t)

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.Name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodGet, "/users", nil)
			rr := httptest.NewRecorder()

			h.GetUsersHandler(rr, req)

			if rr.Code != tc.ExpectedStatus {
				t.Errorf("expected status %v, got %v", tc.ExpectedStatus, rr.Code)
			}
			var gotUser []models.User
			if err := json.Unmarshal(rr.Body.Bytes(), &gotUser); err != nil {
				t.Fatalf("failed to unmarshal response body: %v", err)
			}
			if !reflect.DeepEqual(gotUser, tc.ExpectedUsers) {
				t.Errorf("unexpected response body: got %v, want %v", gotUser, tc.ExpectedUsers)
			}
		})
	}
}
func TestGetUserHandler(t *testing.T) {
	tc := []struct {
		Name           string
		UserID         string
		ExpectedStatus int
		ExpectedUser   models.User
		ExpectedError  string
	}{
		{
			Name:           "User exists",
			UserID:         "1",
			ExpectedStatus: http.StatusOK,
			ExpectedUser: models.User{
				Model:     gorm.Model{ID: 1},
				FirstName: "Jaider",
				LastName:  "Nieto",
				Email:     "email@example.com",
				Password:  "hashPassword",
				Tasks:     nil,
			},
		},
		{
			Name:           "Invalid user",
			UserID:         "-1",
			ExpectedStatus: http.StatusNotFound,
			ExpectedError:  "user not found",
		},
	}

	h := initHandlerUsers(t)

	for i := range tc {
		tc := tc[i]

		t.Run(tc.Name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/users/"+tc.UserID, nil)
			rr := httptest.NewRecorder()

			req = mux.SetURLVars(req, map[string]string{
				"id": tc.UserID,
			})

			h.GetUserHandler(rr, req)

			if rr.Code != tc.ExpectedStatus {
				t.Errorf("unexpected status: got %v, want %v", tc.ExpectedStatus, rr.Code)
			}

			if tc.ExpectedStatus == http.StatusOK {
				var gotUser models.User
				if err := json.Unmarshal(rr.Body.Bytes(), &gotUser); err != nil {
					t.Fatalf("failed to unmarshal response body: %v", err)
				}

				if !reflect.DeepEqual(gotUser, tc.ExpectedUser) {
					t.Errorf("unexpected response body: got %v, want %v", gotUser, tc.ExpectedUser)
				}
			} else if tc.ExpectedStatus == http.StatusNotFound {
				gotError := rr.Body.String()
				if gotError != tc.ExpectedError {
					t.Errorf("unexpected error message: got %v, want %v", gotError, tc.ExpectedError)
				}
			}

		})
	}
}
func TestRegisterUserHandlder(t *testing.T) {
	tc := []struct {
		Name           string
		ExpectedError  string
		ExpectedStatus int
		ExpectedUser   models.User
	}{
		{
			Name:           "Register valid user",
			ExpectedStatus: http.StatusCreated,
			ExpectedUser: models.User{
				Model:     gorm.Model{ID: 1},
				FirstName: "Jaider",
				LastName:  "Nieto",
				Email:     "email@example.com",
				Password:  "hashPassword",
			},
		},
		{
			Name:           "invalid email",
			ExpectedStatus: http.StatusBadRequest,
			ExpectedUser: models.User{
				Model:     gorm.Model{ID: 1},
				FirstName: "Jaider",
				LastName:  "Nieto",
				Email:     "email.com",
				Password:  "hashPassword",
			},
			ExpectedError: "Field validation error on Email: email",
		},
		{
			Name:           "invalid user",
			ExpectedStatus: http.StatusBadRequest,
			ExpectedUser: models.User{
				Model:     gorm.Model{ID: 1},
				FirstName: "Jaider",
				Email:     "email@example.com",
				Password:  "hashPassword",
			},
			ExpectedError: "Field validation error on LastName: required",
		},
		{
			Name:           "invalid password",
			ExpectedStatus: http.StatusBadRequest,
			ExpectedUser: models.User{
				FirstName: "Jaider",
				LastName:  "Nieto",
				Email:     "email@example.com",
				Password:  "1234567",
			},
			ExpectedError: "Field validation error on Password: min",
		},
	}

	h := initHandlerUsers(t)

	for i := range tc {
		tc := tc[i]

		t.Run(tc.Name, func(t *testing.T) {

			body, err := json.Marshal(tc.ExpectedUser)
			if err != nil {
				t.Fatalf("could not marshal json: %v", err)
			}
			handler := middlewares.ValidationMiddleware(
				http.HandlerFunc(h.RegisterUserHandlder),
				&models.User{},
			)

			req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()

			req.Header.Set("Content-Type", "application/json")

			handler.ServeHTTP(rr, req)

			if rr.Code != tc.ExpectedStatus {
				t.Errorf("unexpected status: got %v, want %v", rr.Code, tc.ExpectedStatus)
			}

			if tc.ExpectedStatus == http.StatusCreated {
				var gotUser models.User
				if err := json.Unmarshal(rr.Body.Bytes(), &gotUser); err != nil {
					t.Fatalf("failed to unmarshal response body: %v", err)
				}

				if !reflect.DeepEqual(gotUser, tc.ExpectedUser) {
					t.Errorf("unexpected response body: got %v, want %v", gotUser, tc.ExpectedUser)
				}
			} else if tc.ExpectedStatus == http.StatusBadRequest {
				if rr.Body.String() != tc.ExpectedError {
					t.Errorf("unexpected response error: got %v, want %v", tc.ExpectedError, rr.Body.String())
				}
			}
		})
	}
}
