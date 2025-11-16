package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yizhinailong/api-demo/internal/model"
	"github.com/yizhinailong/api-demo/internal/service"
)

// MockUserService is a mock implementation of UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUser(ctx context.Context, input *service.GetUserInput) (*model.User, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) CreateUser(ctx context.Context, input *service.CreateUserInput) (*model.User, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestUserHandler_CreateUser(t *testing.T) {
	mockService := new(MockUserService)
	handler := &UserHandler{
		userService: mockService,
	}

	router := setupTestRouter()
	handler.RegisterRoutes(router)

	t.Run("create user successfully", func(t *testing.T) {
		input := service.CreateUserInput{
			Username: "testuser",
			Email:    "test@example.com",
			Database: "mysql",
		}

		expectedUser := &model.User{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
		}

		mockService.On("CreateUser", mock.Anything, &input).Return(expectedUser, nil)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/users/create", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "user created", response["status"])
		assert.Equal(t, float64(1), response["userID"])

		mockService.AssertExpectations(t)
	})

	t.Run("create user with invalid input", func(t *testing.T) {
		input := map[string]string{
			"invalid": "data",
		}

		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/users/create", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("create user with service error", func(t *testing.T) {
		input := service.CreateUserInput{
			Username: "testuser",
			Email:    "test@example.com",
			Database: "mysql",
		}

		mockService.On("CreateUser", mock.Anything, &input).Return(nil, assert.AnError)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/users/create", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_GetUser(t *testing.T) {
	mockService := new(MockUserService)
	handler := &UserHandler{
		userService: mockService,
	}

	router := setupTestRouter()
	handler.RegisterRoutes(router)

	t.Run("get user successfully", func(t *testing.T) {
		input := service.GetUserInput{
			ID:       1,
			Database: "mysql",
		}

		expectedUser := &model.User{
			ID:        1,
			Username:  "testuser",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockService.On("GetUser", mock.Anything, &input).Return(expectedUser, nil)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest("GET", "/users/get", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response model.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.ID, response.ID)
		assert.Equal(t, expectedUser.Username, response.Username)
		assert.Equal(t, expectedUser.Email, response.Email)

		mockService.AssertExpectations(t)
	})

	t.Run("get user with invalid input", func(t *testing.T) {
		input := map[string]string{
			"invalid": "data",
		}

		body, _ := json.Marshal(input)
		req := httptest.NewRequest("GET", "/users/get", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("get non-existent user", func(t *testing.T) {
		input := service.GetUserInput{
			ID:       999,
			Database: "mysql",
		}

		mockService.On("GetUser", mock.Anything, &input).Return(nil, assert.AnError)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest("GET", "/users/get", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_RegisterRoutes(t *testing.T) {
	handler := &UserHandler{
		userService: &MockUserService{},
	}

	router := setupTestRouter()
	handler.RegisterRoutes(router)

	// Test that routes are registered
	routes := router.Routes()
	assert.NotEmpty(t, routes)

	// Check for our specific routes
	var foundCreate, foundGet bool
	for _, route := range routes {
		if route.Path == "/users/create" && route.Method == "POST" {
			foundCreate = true
		}
		if route.Path == "/users/get" && route.Method == "GET" {
			foundGet = true
		}
	}

	assert.True(t, foundCreate, "Create user route should be registered")
	assert.True(t, foundGet, "Get user route should be registered")
}
