package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EngenMe/go-clean-arch-auth/internal/core/http/handler"
	"github.com/EngenMe/go-clean-arch-auth/internal/data/entity"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) Register(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserUseCase) Login(email, password string) (*entity.User, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUseCase) GetUserByID(id uint) (*entity.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUseCase) UpdateUser(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserUseCase) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserUseCase) ListUsers(page, limit int) (
	[]entity.UserResponse,
	error,
) {
	args := m.Called(page, limit)
	return args.Get(0).([]entity.UserResponse), args.Error(1)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	return router
}

func TestUserHandler_Register(t *testing.T) {
	mockUseCase := new(MockUserUseCase)
	userHandler := handler.NewUserHandler(mockUseCase)
	router := setupRouter()

	router.POST("/api/register", userHandler.Register)

	mockUseCase.On(
		"Register",
		mock.AnythingOfType("*entity.User"),
	).Return(nil).Once()

	registerBody := map[string]interface{}{
		"name":     "Test User",
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(registerBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(
		"POST",
		"/api/register",
		bytes.NewBuffer(jsonBody),
	)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "message")
	assert.Contains(t, response, "user")

	mockUseCase.AssertExpectations(t)

	w = httptest.NewRecorder()
	invalidBody := map[string]interface{}{
		"name":     "",
		"email":    "invalid-email",
		"password": "",
	}
	jsonBody, _ = json.Marshal(invalidBody)
	req, _ = http.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserHandler_Login(t *testing.T) {
	mockUseCase := new(MockUserUseCase)
	userHandler := handler.NewUserHandler(mockUseCase)
	router := setupRouter()

	router.POST("/api/login", userHandler.Login)

	testUser := &entity.User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "hashed_password",
	}
	mockUseCase.On("Login", "test@example.com", "password123").Return(
		testUser,
		nil,
	).Once()

	loginBody := map[string]interface{}{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(loginBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "message")
	assert.Contains(t, response, "token")
	assert.Contains(t, response, "user")

	mockUseCase.AssertExpectations(t)

	mockUseCase.On("Login", "test@example.com", "wrong-password").Return(
		nil,
		assert.AnError,
	).Once()

	loginBody = map[string]interface{}{
		"email":    "test@example.com",
		"password": "wrong-password",
	}
	jsonBody, _ = json.Marshal(loginBody)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	mockUseCase.AssertExpectations(t)
}
