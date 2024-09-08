package handlers

import (
	"bytes"
	"encoding/json"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jaider-nieto/ecommerce-go/user-service/middlewares"
	"github.com/jaider-nieto/ecommerce-go/user-service/models"
	"github.com/jaider-nieto/ecommerce-go/user-service/repository"
)

func initHandlerUsers(t *testing.T) *userHandler {
	t.Helper()

	userRepositoryMock := &repository.UserRepositoryMocked{}
	return NewUserHandler(userRepositoryMock)
}
func initRequest(method string, url string, body io.Reader) (*httptest.ResponseRecorder, *http.Request) {
	req := httptest.NewRequest(method, url, body)
	rr := httptest.NewRecorder()

	return rr, req
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
			rr, req := initRequest(http.MethodGet, "/users", nil)

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

			rr, req := initRequest(http.MethodPost, "/register", bytes.NewBuffer(body))

			req.Header.Set("Content-Type", "application/json")

			handler.ServeHTTP(rr, req)

			if rr.Code != tc.ExpectedStatus {
				t.Errorf("unexpected status: got %v, want %v", rr.Code, tc.ExpectedStatus)
			}

			if tc.ExpectedStatus == http.StatusCreated {
				var gotUser models.User
				if err := json.Unmarshal(rr.Body.Bytes(), &gotUser); err != nil {
					log.Panicf("%v", rr.Body.String())
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
func TestLoginUserHanlder(t *testing.T) {
	tc := []struct {
		Name           string
		ExpectedError  string
		ExpectedStatus int
		ExpectedToken  string
		UserLogin      models.UserLogin
	}{
		{
			Name:           "valid login",
			ExpectedStatus: http.StatusOK,
			ExpectedToken:  "BearereyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImVtYWlsQHZhbGlkLmNvbSIsImZpcnN0X25hbWUiOiJKYWlkZXIifQ.O8hk_iNG08quNDiqtBAX2WLLIAEu5phS8DIG2wPWgP8",
			UserLogin: models.UserLogin{
				Email:    "email@valid.com",
				Password: "hashpassword",
			},
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
			ExpectedError:  "user not found",
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
	}

	h := initHandlerUsers(t)

	for i := range tc {
		tc := tc[i]

		t.Run(tc.Name, func(t *testing.T) {

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

			// log.Printf("error %v", rr.Body.String())
			if rr.Code != tc.ExpectedStatus {
				// log.Printf("%v", rr.Header().Get("Authorization"))
				t.Fatalf("unexpected status: got %v want %v", rr.Code, tc.ExpectedStatus)
			}

			if rr.Code == http.StatusOK {
				if rr.Header().Get("Authorization") != tc.ExpectedToken {
					t.Fatalf("unexpected token: got %v want %v",
						rr.Header().Get("Authorization"), tc.ExpectedToken)
				}
			} else if rr.Code == http.StatusBadRequest || rr.Code == http.StatusNotFound {
				if rr.Body.String() != tc.ExpectedError {
					t.Fatalf("unexpected error: got %v want %v", rr.Body.String(), tc.ExpectedError)
				}
			}
		})
	}
}
func TestDeleteUserHandler(t *testing.T) {
	tc := []struct {
		Name            string
		ExpectedError   string
		ExpectedMessage string
		ExpectedStatus  int
		UserID          string
	}{
		{
			Name:            "Delete user",
			UserID:          "1",
			ExpectedStatus:  http.StatusOK,
			ExpectedMessage: "user deleted",
		},
		{
			Name:           "User not found",
			UserID:         "99",
			ExpectedStatus: http.StatusNotFound,
			ExpectedError:  "user not found",
		},
	}

	h := initHandlerUsers(t)

	for i := range tc {
		tc := tc[i]

		t.Run(tc.Name, func(t *testing.T) {

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
			} else if rr.Code == http.StatusBadRequest || rr.Code == http.StatusNotFound {
				if rr.Body.String() != tc.ExpectedError {
					t.Fatalf("unexpected error: got %v want %v", rr.Body.String(), tc.ExpectedError)
				}
			}
		})
	}
}
func TestPatchUserHandler(t *testing.T) {
	tc := []struct {
		Name           string
		ExpectedError  string
		ExpectedStatus int
		UserID         string
		UserBody       models.UserUpdate
		ExpectedUser   models.User
	}{
		{
			Name:           "Patch user",
			ExpectedStatus: http.StatusOK,
			UserID:         "1",
			UserBody: models.UserUpdate{
				FirstName: "Jajaider",
				Email:     "jaiderlol@gmail.com",
			},
			ExpectedUser: models.User{
				Model:     gorm.Model{ID: 1},
				FirstName: "Jajaider",
				LastName:  "Nieto",
				Email:     "jaiderlol@gmail.com",
				Password:  "hashPassword",
			},
		},
		{
			Name:           "Not found user",
			ExpectedStatus: http.StatusNotFound,
			UserID:         "99",
			UserBody: models.UserUpdate{
				FirstName: "Jajaider",
				Email:     "jaiderlol@gmail.com",
			},
		},
	}

	h := initHandlerUsers(t)

	for i := range tc {
		tc := tc[i]

		t.Run(tc.Name, func(t *testing.T) {
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
				log.Printf("%v", rr.Body.String())
				t.Fatalf("unexpected status: got %v want %v", rr.Code, tc.ExpectedStatus)
			}

			if rr.Code == http.StatusOK {
				var gotUser models.User
				if err := json.Unmarshal(rr.Body.Bytes(), &gotUser); err != nil {
					log.Panicf("%v", rr.Body.String())
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
