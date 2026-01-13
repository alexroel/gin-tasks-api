package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexroel/gin-tasks-api/internal/domain"
	"github.com/alexroel/gin-tasks-api/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupAuthRouter(handler *AuthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestAuthHandler_SignUp_Success(t *testing.T) {
	// Arrange
	mockService := new(mocks.AuthServiceMock)
	handler := NewAuthHandler(mockService)
	router := setupAuthRouter(handler)
	router.POST("/auth/signup", handler.SignUpHandler)

	user := &domain.User{ID: 1, FullName: "Test User", Email: "test@example.com"}
	mockService.On("Register", mock.AnythingOfType("*domain.UserCreate")).Return(user, nil)

	reqBody := `{"full_name": "Test User", "email": "test@example.com", "password": "password123"}`
	req, _ := http.NewRequest("POST", "/auth/signup", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusCreated, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	assert.Equal(t, "Usuario registrado exitosamente", response["message"])
	mockService.AssertExpectations(t)
}

func TestAuthHandler_SignUp_InvalidBody(t *testing.T) {
	// Arrange
	mockService := new(mocks.AuthServiceMock)
	handler := NewAuthHandler(mockService)
	router := setupAuthRouter(handler)
	router.POST("/auth/signup", handler.SignUpHandler)

	reqBody := `{"invalid": "body"}`
	req, _ := http.NewRequest("POST", "/auth/signup", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestAuthHandler_SignUp_EmailExists(t *testing.T) {
	// Arrange
	mockService := new(mocks.AuthServiceMock)
	handler := NewAuthHandler(mockService)
	router := setupAuthRouter(handler)
	router.POST("/auth/signup", handler.SignUpHandler)

	mockService.On("Register", mock.AnythingOfType("*domain.UserCreate")).Return(nil, errors.New("Ya existe un usuario con ese email"))

	reqBody := `{"full_name": "Test User", "email": "test@example.com", "password": "password123"}`
	req, _ := http.NewRequest("POST", "/auth/signup", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	mockService.AssertExpectations(t)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	// Arrange
	mockService := new(mocks.AuthServiceMock)
	handler := NewAuthHandler(mockService)
	router := setupAuthRouter(handler)
	router.POST("/auth/login", handler.Login)

	user := &domain.User{ID: 1, FullName: "Test User", Email: "test@example.com"}
	mockService.On("Login", mock.AnythingOfType("*domain.UserLogin")).Return("fake-token", user, nil)

	reqBody := `{"email": "test@example.com", "password": "password123"}`
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	assert.Equal(t, "Inicio de sesión exitoso", response["message"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "fake-token", data["token"])
	mockService.AssertExpectations(t)
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	// Arrange
	mockService := new(mocks.AuthServiceMock)
	handler := NewAuthHandler(mockService)
	router := setupAuthRouter(handler)
	router.POST("/auth/login", handler.Login)

	mockService.On("Login", mock.AnythingOfType("*domain.UserLogin")).Return("", nil, errors.New("Credenciales inválidas"))

	reqBody := `{"email": "test@example.com", "password": "wrongpassword"}`
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	mockService.AssertExpectations(t)
}

func TestAuthHandler_Login_InvalidBody(t *testing.T) {
	// Arrange
	mockService := new(mocks.AuthServiceMock)
	handler := NewAuthHandler(mockService)
	router := setupAuthRouter(handler)
	router.POST("/auth/login", handler.Login)

	reqBody := `{"invalid": "body"}`
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestAuthHandler_Profile_Success(t *testing.T) {
	// Arrange
	mockService := new(mocks.AuthServiceMock)
	handler := NewAuthHandler(mockService)
	router := setupAuthRouter(handler)

	// Middleware para simular usuario autenticado
	router.GET("/auth/profile", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.Profile(c)
	})

	user := &domain.User{ID: 1, FullName: "Test User", Email: "test@example.com"}
	mockService.On("GetUserByID", uint(1)).Return(user, nil)

	req, _ := http.NewRequest("GET", "/auth/profile", nil)
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestAuthHandler_Profile_NotAuthenticated(t *testing.T) {
	// Arrange
	mockService := new(mocks.AuthServiceMock)
	handler := NewAuthHandler(mockService)
	router := setupAuthRouter(handler)
	router.GET("/auth/profile", handler.Profile)

	req, _ := http.NewRequest("GET", "/auth/profile", nil)
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestAuthHandler_Profile_UserNotFound(t *testing.T) {
	// Arrange
	mockService := new(mocks.AuthServiceMock)
	handler := NewAuthHandler(mockService)
	router := setupAuthRouter(handler)

	router.GET("/auth/profile", func(c *gin.Context) {
		c.Set("userID", uint(999))
		handler.Profile(c)
	})

	mockService.On("GetUserByID", uint(999)).Return(nil, errors.New("not found"))

	req, _ := http.NewRequest("GET", "/auth/profile", nil)
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockService.AssertExpectations(t)
}

func TestAuthHandler_UpdateProfile_Success(t *testing.T) {
	// Arrange
	mockService := new(mocks.AuthServiceMock)
	handler := NewAuthHandler(mockService)
	router := setupAuthRouter(handler)

	router.PUT("/auth/profile", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.UpdateProfile(c)
	})

	user := &domain.User{ID: 1, FullName: "Updated Name", Email: "test@example.com"}
	mockService.On("UpdateProfile", uint(1), mock.AnythingOfType("*domain.UserUpdate")).Return(user, nil)

	reqBody := `{"full_name": "Updated Name"}`
	req, _ := http.NewRequest("PUT", "/auth/profile", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestAuthHandler_DeleteAccount_Success(t *testing.T) {
	// Arrange
	mockService := new(mocks.AuthServiceMock)
	handler := NewAuthHandler(mockService)
	router := setupAuthRouter(handler)

	router.DELETE("/auth/profile", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.DeleteAccount(c)
	})

	mockService.On("DeleteAccount", uint(1)).Return(nil)

	req, _ := http.NewRequest("DELETE", "/auth/profile", nil)
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestAuthHandler_DeleteAccount_Error(t *testing.T) {
	// Arrange
	mockService := new(mocks.AuthServiceMock)
	handler := NewAuthHandler(mockService)
	router := setupAuthRouter(handler)

	router.DELETE("/auth/profile", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.DeleteAccount(c)
	})

	mockService.On("DeleteAccount", uint(1)).Return(errors.New("error deleting"))

	req, _ := http.NewRequest("DELETE", "/auth/profile", nil)
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	mockService.AssertExpectations(t)
}
