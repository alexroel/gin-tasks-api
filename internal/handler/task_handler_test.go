package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexroel/gin-tasks-api/internal/domain"
	"github.com/alexroel/gin-tasks-api/internal/service"
	"github.com/alexroel/gin-tasks-api/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTaskRouter(handler *TaskHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestTaskHandler_Create_Success(t *testing.T) {
	// Arrange
	mockService := new(mocks.TaskServiceMock)
	handler := NewTaskHandler(mockService)
	router := setupTaskRouter(handler)

	router.POST("/tasks", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.Create(c)
	})

	task := &domain.Task{ID: 1, Title: "Nueva tarea", UserID: 1}
	mockService.On("Create", uint(1), mock.AnythingOfType("*domain.CreateTask")).Return(task, nil)

	reqBody := `{"title": "Nueva tarea"}`
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusCreated, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	assert.Equal(t, "Tarea creada exitosamente", response["message"])
	mockService.AssertExpectations(t)
}

func TestTaskHandler_Create_NotAuthenticated(t *testing.T) {
	// Arrange
	mockService := new(mocks.TaskServiceMock)
	handler := NewTaskHandler(mockService)
	router := setupTaskRouter(handler)
	router.POST("/tasks", handler.Create)

	reqBody := `{"title": "Nueva tarea"}`
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestTaskHandler_Create_InvalidBody(t *testing.T) {
	// Arrange
	mockService := new(mocks.TaskServiceMock)
	handler := NewTaskHandler(mockService)
	router := setupTaskRouter(handler)

	router.POST("/tasks", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.Create(c)
	})

	reqBody := `{"invalid": "body"}`
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestTaskHandler_GetAll_Success(t *testing.T) {
	// Arrange
	mockService := new(mocks.TaskServiceMock)
	handler := NewTaskHandler(mockService)
	router := setupTaskRouter(handler)

	router.GET("/tasks", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.GetAll(c)
	})

	tasks := []domain.Task{
		{ID: 1, Title: "Tarea 1", UserID: 1},
		{ID: 2, Title: "Tarea 2", UserID: 1},
	}
	mockService.On("GetByUserID", uint(1)).Return(tasks, nil)

	req, _ := http.NewRequest("GET", "/tasks", nil)
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.True(t, response["success"].(bool))
	mockService.AssertExpectations(t)
}

func TestTaskHandler_GetAll_NotAuthenticated(t *testing.T) {
	// Arrange
	mockService := new(mocks.TaskServiceMock)
	handler := NewTaskHandler(mockService)
	router := setupTaskRouter(handler)
	router.GET("/tasks", handler.GetAll)

	req, _ := http.NewRequest("GET", "/tasks", nil)
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestTaskHandler_GetByID_Success(t *testing.T) {
	// Arrange
	mockService := new(mocks.TaskServiceMock)
	handler := NewTaskHandler(mockService)
	router := setupTaskRouter(handler)

	router.GET("/tasks/:id", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.GetByID(c)
	})

	task := &domain.Task{ID: 1, Title: "Tarea 1", UserID: 1}
	mockService.On("GetByID", uint(1)).Return(task, nil)

	req, _ := http.NewRequest("GET", "/tasks/1", nil)
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestTaskHandler_GetByID_NotFound(t *testing.T) {
	// Arrange
	mockService := new(mocks.TaskServiceMock)
	handler := NewTaskHandler(mockService)
	router := setupTaskRouter(handler)

	router.GET("/tasks/:id", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.GetByID(c)
	})

	mockService.On("GetByID", uint(999)).Return(nil, service.ErrTaskNotFound)

	req, _ := http.NewRequest("GET", "/tasks/999", nil)
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockService.AssertExpectations(t)
}

func TestTaskHandler_GetByID_Forbidden(t *testing.T) {
	// Arrange
	mockService := new(mocks.TaskServiceMock)
	handler := NewTaskHandler(mockService)
	router := setupTaskRouter(handler)

	router.GET("/tasks/:id", func(c *gin.Context) {
		c.Set("userID", uint(2)) // Usuario diferente al dueño
		handler.GetByID(c)
	})

	task := &domain.Task{ID: 1, Title: "Tarea 1", UserID: 1} // Pertenece al usuario 1
	mockService.On("GetByID", uint(1)).Return(task, nil)

	req, _ := http.NewRequest("GET", "/tasks/1", nil)
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusForbidden, resp.Code)
	mockService.AssertExpectations(t)
}

func TestTaskHandler_GetByID_InvalidID(t *testing.T) {
	// Arrange
	mockService := new(mocks.TaskServiceMock)
	handler := NewTaskHandler(mockService)
	router := setupTaskRouter(handler)

	router.GET("/tasks/:id", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.GetByID(c)
	})

	req, _ := http.NewRequest("GET", "/tasks/invalid", nil)
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestTaskHandler_Update_Success(t *testing.T) {
	// Arrange
	mockService := new(mocks.TaskServiceMock)
	handler := NewTaskHandler(mockService)
	router := setupTaskRouter(handler)

	router.PUT("/tasks/:id", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.Update(c)
	})

	task := &domain.Task{ID: 1, Title: "Tarea actualizada", UserID: 1}
	mockService.On("Update", uint(1), uint(1), mock.AnythingOfType("*domain.UpdateTask")).Return(task, nil)

	reqBody := `{"title": "Tarea actualizada"}`
	req, _ := http.NewRequest("PUT", "/tasks/1", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestTaskHandler_Update_NotFound(t *testing.T) {
	// Arrange
	mockService := new(mocks.TaskServiceMock)
	handler := NewTaskHandler(mockService)
	router := setupTaskRouter(handler)

	router.PUT("/tasks/:id", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.Update(c)
	})

	mockService.On("Update", uint(999), uint(1), mock.AnythingOfType("*domain.UpdateTask")).Return(nil, service.ErrTaskNotFound)

	reqBody := `{"title": "Tarea actualizada"}`
	req, _ := http.NewRequest("PUT", "/tasks/999", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockService.AssertExpectations(t)
}

func TestTaskHandler_Update_Unauthorized(t *testing.T) {
	// Arrange
	mockService := new(mocks.TaskServiceMock)
	handler := NewTaskHandler(mockService)
	router := setupTaskRouter(handler)

	router.PUT("/tasks/:id", func(c *gin.Context) {
		c.Set("userID", uint(2))
		handler.Update(c)
	})

	mockService.On("Update", uint(1), uint(2), mock.AnythingOfType("*domain.UpdateTask")).Return(nil, service.ErrTaskUnauthorized)

	reqBody := `{"title": "Tarea actualizada"}`
	req, _ := http.NewRequest("PUT", "/tasks/1", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusForbidden, resp.Code)
	mockService.AssertExpectations(t)
}

func TestTaskHandler_Delete_Success(t *testing.T) {
	// Arrange
	mockService := new(mocks.TaskServiceMock)
	handler := NewTaskHandler(mockService)
	router := setupTaskRouter(handler)

	router.DELETE("/tasks/:id", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.Delete(c)
	})

	mockService.On("Delete", uint(1), uint(1)).Return(nil)

	req, _ := http.NewRequest("DELETE", "/tasks/1", nil)
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestTaskHandler_Delete_NotFound(t *testing.T) {
	// Arrange
	mockService := new(mocks.TaskServiceMock)
	handler := NewTaskHandler(mockService)
	router := setupTaskRouter(handler)

	router.DELETE("/tasks/:id", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.Delete(c)
	})

	mockService.On("Delete", uint(999), uint(1)).Return(service.ErrTaskNotFound)

	req, _ := http.NewRequest("DELETE", "/tasks/999", nil)
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockService.AssertExpectations(t)
}

func TestTaskHandler_ToggleStatus_Success(t *testing.T) {
	// Arrange
	mockService := new(mocks.TaskServiceMock)
	handler := NewTaskHandler(mockService)
	router := setupTaskRouter(handler)

	router.PATCH("/tasks/:id/status", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.ToggleStatus(c)
	})

	task := &domain.Task{ID: 1, Title: "Tarea", UserID: 1, Completed: true}
	mockService.On("UpdateStatus", uint(1), uint(1), true).Return(task, nil)

	reqBody := `{"completed": true}`
	req, _ := http.NewRequest("PATCH", "/tasks/1/status", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestTaskHandler_ToggleStatus_NotFound(t *testing.T) {
	// Arrange
	mockService := new(mocks.TaskServiceMock)
	handler := NewTaskHandler(mockService)
	router := setupTaskRouter(handler)

	router.PATCH("/tasks/:id/status", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.ToggleStatus(c)
	})

	mockService.On("UpdateStatus", uint(999), uint(1), true).Return(nil, service.ErrTaskNotFound)

	reqBody := `{"completed": true}`
	req, _ := http.NewRequest("PATCH", "/tasks/999/status", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockService.AssertExpectations(t)
}
