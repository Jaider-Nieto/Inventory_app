package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"gorm.io/gorm"

	"github.com/gorilla/mux"
	"github.com/jaider-nieto/ecommerce-go/user-service/middlewares"
	"github.com/jaider-nieto/ecommerce-go/user-service/models"
	"github.com/jaider-nieto/ecommerce-go/user-service/repository"
)

func initHandlerUsers(t *testing.T, shouldReturnError bool) *userHandler {
	t.Helper()

	// Configura el repositorio mockeado para simular errores si se solicita.
	userRepositoryMock := &repository.UserRepositoryMocked{ShouldReturnError: shouldReturnError}

	// Inicializa el handler con el repositorio mockeado.
	return NewUserHandler(userRepositoryMock)
}

func initRequest(method string, url string, body io.Reader) (*httptest.ResponseRecorder, *http.Request) {
	req := httptest.NewRequest(method, url, body)
	rr := httptest.NewRecorder()

	return rr, req
}

func TestGetUsersHandler(t *testing.T) {
	testCases := []struct {
		Name              string
		ExpectedStatus    int
		ExpectedError     string
		ExpectedUsers     []models.User
		ShouldReturnError bool
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
		{
			Name:              "Server error",
			ExpectedStatus:    http.StatusInternalServerError,
			ExpectedError:     "internal server error",
			ShouldReturnError: true,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.Name, func(t *testing.T) {
			h := initHandlerUsers(t, tc.ShouldReturnError)
			rr, req := initRequest(http.MethodGet, "/users", nil)

			h.GetUsersHandler(rr, req)

			if rr.Code != tc.ExpectedStatus {
				t.Errorf("expected status %v, got %v", tc.ExpectedStatus, rr.Code)
			}
			if rr.Code == http.StatusInternalServerError {
				if rr.Body.String() != tc.ExpectedError {
					t.Errorf("unexpected error: got %v, want %v", rr.Body.String(), tc.ExpectedUsers)
				}
			} else if rr.Code == http.StatusOK {
				var gotUser []models.User
				if err := json.Unmarshal(rr.Body.Bytes(), &gotUser); err != nil {
					t.Fatalf("failed to unmarshal response body: %v", err)
				}
				if !reflect.DeepEqual(gotUser, tc.ExpectedUsers) {
					t.Errorf("unexpected response body: got %v, want %v", gotUser, tc.ExpectedUsers)
				}
			}

		})
	}
}

func TestGetUserHandler(t *testing.T) {
	tc := []struct {
		Name              string
		UserID            string
		ShouldReturnError bool
		ExpectedStatus    int
		ExpectedUser      models.User
		ExpectedError     string
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
			},
		},
		{
			Name:           "Invalid user",
			UserID:         "-1",
			ExpectedStatus: http.StatusNotFound,
			ExpectedError:  "user not found",
		},
		{
			Name:              "Server error",
			ShouldReturnError: true,
			UserID:            "1",
			ExpectedStatus:    http.StatusInternalServerError,
			ExpectedError:     "internal server error",
		},
	}

	for i := range tc {
		tc := tc[i]

		t.Run(tc.Name, func(t *testing.T) {
			h := initHandlerUsers(t, tc.ShouldReturnError)
			rr, req := initRequest(http.MethodGet, "/users/"+tc.UserID, nil)

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
			} else {
				if rr.Body.String() != tc.ExpectedError {
					t.Errorf("unexpected error message: got %v, want %v", rr.Body.String(), tc.ExpectedError)
				}
			}

		})
	}
}
func TestRegisterUserHandlder(t *testing.T) {
	tc := []struct {
		Name              string
		ExpectedError     string
		ExpectedStatus    int
		ExpectedUser      models.User
		ShouldReturnError bool
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
			Name: "email not found",
				ExpectedUser: models.User{
				Model:     gorm.Model{ID: 1},
				FirstName: "Jaider",
				LastName:  "Nieto",
				Email:     "email@example.com",
				Password:  "hashPassword",
			},
			ExpectedStatus:    http.StatusInternalServerError,
			ExpectedError:     "internal server error",
			ShouldReturnError: true,
		},
		{
			Name:           "invalid user",
			ExpectedStatus: http.StatusBadRequest,
			ExpectedUser: models.User{
				Model:     gorm.Model{ID: 1},
				FirstName: "Jaider",
				Email:     "email@valid.com",
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
				Email:     "email@valid.com",
				Password:  "1234567",
			},
			ExpectedError: "Field validation error on Password: min",
		},
		{
			Name:           "Server error",
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedUser: models.User{
				Model:     gorm.Model{ID: 2},
				FirstName: "Jaider",
				LastName:  "Nieto",
				Email:     "email@valid.com",
				Password:  "hashPassword",
			},
			ExpectedError:     "internal server error",
			ShouldReturnError: true,
		},
	}

	for i := range tc {
		tc := tc[i]

		t.Run(tc.Name, func(t *testing.T) {
			h := initHandlerUsers(t, tc.ShouldReturnError)

			body, err := json.Marshal(tc.ExpectedUser)
			if err != nil {
				t.Fatalf("could not marshal json: %v", err)
			}
			handler := middlewares.ValidationMiddleware(
				http.HandlerFunc(h.RegisterUserHandlder),
				&models.User{},
			)

			rr, req := initRequest(http.MethodPost, "/register", bytes.NewBuffer(body))

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
			} else {
				if rr.Body.String() != tc.ExpectedError {
					t.Errorf("unexpected response error: got %v, want %v", tc.ExpectedError, rr.Body.String())
				}
			}
		})
	}
}
func TestLoginUserHanlder(t *testing.T) {
	tc := []struct {
		Name              string
		ShouldReturnError bool
		ExpectedError     string
		ExpectedStatus    int
		ExpectedMessage   string
		UserLogin         models.UserLogin
		UserBad           any
	}{
		{
			Name:           "valid login",
			ExpectedStatus: http.StatusOK,
			UserLogin: models.UserLogin{
				Email:    "email@valid.com",
				Password: "hashpassword",
			},
			ExpectedMessage: "user login",
		},
		{
			Name:           "invalid email",
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  "Field validation error on Email: email",
			UserLogin: models.UserLogin{
				Email:    "invalid@email",
				Password: "hashpassword",
			},
		},
		{
			Name:           "user not found",
			ExpectedStatus: http.StatusNotFound,
			ExpectedError:  "email not found",
			UserLogin: models.UserLogin{
				Email:    "user@notfound.com",
				Password: "hashpassword",
			},
		},
		{
			Name:           "incorrect password",
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  "incorret password",
			UserLogin: models.UserLogin{
				Email:    "email@valid.com",
				Password: "aelkfnwlfnowfa",
			},
		},
		{
			Name: "Server error",
			UserLogin: models.UserLogin{
				Email:    "email@example.com",
				Password: "hashpassword",
			},
			ExpectedStatus:    http.StatusInternalServerError,
			ExpectedError:     "internal server error",
			ShouldReturnError: true,
		},
	}

	for i := range tc {
		tc := tc[i]

		t.Run(tc.Name, func(t *testing.T) {
			h := initHandlerUsers(t, tc.ShouldReturnError)

			body, err := json.Marshal(tc.UserLogin)
			if err != nil {
				t.Fatalf("could not marshal json: %v", err)
			}

			handler := middlewares.ValidationMiddleware(
				http.HandlerFunc(h.LoginUserHanlder),
				&models.UserLogin{},
			)

			rr, req := initRequest(http.MethodPost, "/login", bytes.NewBuffer(body))

			req.Header.Set("Content-Type", "application/json")

			handler.ServeHTTP(rr, req)

			if rr.Code != tc.ExpectedStatus {
				t.Fatalf("unexpected status: got %v want %v", rr.Code, tc.ExpectedStatus)
			}

			if rr.Code == http.StatusOK {
				if rr.Body.String() != tc.ExpectedMessage {
					t.Fatalf("unexpected message: got %v want %v",
						rr.Body.String(), tc.ExpectedMessage)
				}
			} else {
				if rr.Body.String() != tc.ExpectedError {
					t.Fatalf("unexpected error: got %v want %v", rr.Body.String(), tc.ExpectedError)
				}
			}
		})
	}
}
func TestDeleteUserHandler(t *testing.T) {
	tc := []struct {
		Name              string
		ExpectedError     string
		ExpectedMessage   string
		ExpectedStatus    int
		UserID            string
		ShouldReturnError bool
	}{
		{
			Name:            "Delete user",
			UserID:          "2",
			ExpectedStatus:  http.StatusOK,
			ExpectedMessage: "user deleted",
		},
		{
			Name:           "User not found",
			UserID:         "99",
			ExpectedStatus: http.StatusNotFound,
			ExpectedError:  "user not found",
		},
		{
			Name:              "Server error FindId",
			ExpectedStatus:    http.StatusInternalServerError,
			UserID:            "1",
			ExpectedError:     "internal server error",
			ShouldReturnError: true,
		},
		{
			Name:              "Server error delete",
			ExpectedStatus:    http.StatusInternalServerError,
			UserID:            "2",
			ExpectedError:     "internal server error",
			ShouldReturnError: true,
		},
	}

	for i := range tc {
		tc := tc[i]

		t.Run(tc.Name, func(t *testing.T) {
			h := initHandlerUsers(t, tc.ShouldReturnError)

			rr, req := initRequest(http.MethodDelete, "/users/"+tc.UserID, nil)

			req = mux.SetURLVars(req, map[string]string{
				"id": tc.UserID,
			})

			h.DeleteUserHandler(rr, req)

			if rr.Code != tc.ExpectedStatus {
				t.Fatalf("unexpected status: got %v want %v", rr.Code, tc.ExpectedStatus)
			}

			if rr.Code == http.StatusOK {
				if rr.Body.String() != tc.ExpectedMessage {
					t.Fatalf("unexpected message: got %v want %v", rr.Body.String(), tc.ExpectedMessage)
				}
			} else {
				if rr.Body.String() != tc.ExpectedError {
					t.Fatalf("unexpected error: got %v want %v", rr.Body.String(), tc.ExpectedError)
				}
			}
		})
	}
}
func TestPatchUserHandler(t *testing.T) {
	tc := []struct {
		Name              string
		ShouldReturnError bool
		ExpectedError     string
		ExpectedStatus    int
		UserID            string
		UserBody          models.UserUpdate
		ExpectedUser      models.User
	}{
		{
			Name:           "Patch user",
			ExpectedStatus: http.StatusOK,
			UserID:         "1",
			UserBody: models.UserUpdate{
				FirstName: "Jajaider",
				LastName:  "criollo",
				Email:     "jaiderlol@gmail.com",
			},
			ExpectedUser: models.User{
				Model:     gorm.Model{ID: 1},
				FirstName: "Jajaider",
				LastName:  "criollo",
				Email:     "jaiderlol@gmail.com",
				Password:  "hashPassword",
			},
		},
		{
			Name:           "User not found",
			ExpectedStatus: http.StatusNotFound,
			ExpectedError:  "user not found",
			UserID:         "99",
			UserBody: models.UserUpdate{
				FirstName: "Jajaider",
				Email:     "jaiderlol@gmail.com",
			},
		},
		{
			Name:              "Server error FindId",
			ExpectedStatus:    http.StatusInternalServerError,
			UserID:            "1",
			ExpectedError:     "internal server error",
			ShouldReturnError: true,
		},
		{
			Name:              "Server error delete",
			ExpectedStatus:    http.StatusInternalServerError,
			UserID:            "2",
			ExpectedError:     "internal server error",
			ShouldReturnError: true,
		},
	}

	for i := range tc {
		tc := tc[i]

		t.Run(tc.Name, func(t *testing.T) {
			h := initHandlerUsers(t, tc.ShouldReturnError)
			body, err := json.Marshal(tc.UserBody)
			if err != nil {
				t.Fatalf("could not marshal json: %v", err)
			}
			rr, req := initRequest(http.MethodPatch, "/users/"+tc.UserID, bytes.NewBuffer(body))

			req = mux.SetURLVars(req, map[string]string{
				"id": tc.UserID,
			})

			h.PatchUserHandler(rr, req)

			if rr.Code != tc.ExpectedStatus {
				t.Fatalf("unexpected status: got %v want %v", rr.Code, tc.ExpectedStatus)
			}

			if rr.Code == http.StatusOK {
				var gotUser models.User
				if err := json.Unmarshal(rr.Body.Bytes(), &gotUser); err != nil {
					t.Fatalf("failed to unmarshal response body: %v", err)
				}

				if !reflect.DeepEqual(gotUser, tc.ExpectedUser) {
					t.Errorf("unexpected response body: got %v, want %v", gotUser, tc.ExpectedUser)
				}
			} else {
				if rr.Body.String() != tc.ExpectedError {
					t.Errorf("unexpected response error: got %v, want %v", rr.Body.String(), tc.ExpectedError)
				}
			}
		})
	}
}
