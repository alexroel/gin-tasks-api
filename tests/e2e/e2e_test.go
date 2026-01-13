package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexroel/gin-tasks-api/internal/config"
	"github.com/alexroel/gin-tasks-api/internal/domain"
	"github.com/alexroel/gin-tasks-api/internal/handler"
	"github.com/alexroel/gin-tasks-api/internal/middleware"
	repoMocks "github.com/alexroel/gin-tasks-api/internal/repository/mocks"
	"github.com/alexroel/gin-tasks-api/internal/service"
	"github.com/alexroel/gin-tasks-api/internal/service/mocks"
	"github.com/alexroel/gin-tasks-api/pkg/jwt"
	"github.com/alexroel/gin-tasks-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// setupE2EConfig configura el entorno para tests E2E
func setupE2EConfig() {
	config.AppConfig = &config.Config{
		JWTSecret:   "e2e-test-secret-key-12345",
		JWTExpireIn: 24 * time.Hour,
	}
}

// ========== Tests E2E - Flujos Completos ==========

func TestE2E_RegisterAndLogin_Flow(t *testing.T) {
	setupE2EConfig()

	// Mocks
	authServiceMock := new(mocks.AuthServiceMock)
	authHandler := handler.NewAuthHandler(authServiceMock)

	router := gin.New()
	gin.SetMode(gin.TestMode)
	router.POST("/api/v1/auth/signup", authHandler.SignUpHandler)
	router.POST("/api/v1/auth/login", authHandler.Login)

	// 1. Registro de usuario
	t.Run("1. Registro exitoso", func(t *testing.T) {
		newUser := &domain.User{
			ID:       1,
			FullName: "E2E User",
			Email:    "e2e@example.com",
		}

		authServiceMock.On("Register", mock.AnythingOfType("*domain.UserCreate")).Return(newUser, nil).Once()

		signupBody := `{"full_name": "E2E User", "email": "e2e@example.com", "password": "password123"}`
		req, _ := http.NewRequest("POST", "/api/v1/auth/signup", bytes.NewBufferString(signupBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)

		var response map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))
	})

	// 2. Login
	t.Run("2. Login exitoso", func(t *testing.T) {
		user := &domain.User{
			ID:       1,
			FullName: "E2E User",
			Email:    "e2e@example.com",
		}

		token, _ := jwt.GenerateToken(1, "e2e@example.com", config.AppConfig.JWTSecret, config.AppConfig.JWTExpireIn)
		authServiceMock.On("Login", mock.AnythingOfType("*domain.UserLogin")).Return(token, user, nil).Once()

		loginBody := `{"email": "e2e@example.com", "password": "password123"}`
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBufferString(loginBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var response map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		assert.NotEmpty(t, data["token"])
	})

	authServiceMock.AssertExpectations(t)
}

func TestE2E_TaskCRUD_Flow(t *testing.T) {
	setupE2EConfig()

	// Mocks
	taskServiceMock := new(mocks.TaskServiceMock)
	taskHandler := handler.NewTaskHandler(taskServiceMock)

	router := gin.New()
	gin.SetMode(gin.TestMode)

	// Rutas protegidas
	tasks := router.Group("/api/v1/tasks")
	tasks.Use(middleware.AuthMiddleware(config.AppConfig.JWTSecret))
	{
		tasks.POST("", taskHandler.Create)
		tasks.GET("", taskHandler.GetAll)
		tasks.GET("/:id", taskHandler.GetByID)
		tasks.PUT("/:id", taskHandler.Update)
		tasks.DELETE("/:id", taskHandler.Delete)
	}

	// Generar token válido
	token, _ := jwt.GenerateToken(1, "e2e@example.com", config.AppConfig.JWTSecret, config.AppConfig.JWTExpireIn)

	// 1. Crear tarea
	var createdTaskID uint = 1
	t.Run("1. Crear tarea", func(t *testing.T) {
		task := &domain.Task{
			ID:     createdTaskID,
			Title:  "Mi primera tarea E2E",
			UserID: 1,
		}

		taskServiceMock.On("Create", uint(1), mock.AnythingOfType("*domain.CreateTask")).Return(task, nil).Once()

		body := `{"title": "Mi primera tarea E2E"}`
		req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
	})

	// 2. Obtener todas las tareas
	t.Run("2. Obtener tareas del usuario", func(t *testing.T) {
		tasks := []domain.Task{
			{ID: 1, Title: "Tarea 1", UserID: 1},
			{ID: 2, Title: "Tarea 2", UserID: 1},
		}

		taskServiceMock.On("GetByUserID", uint(1)).Return(tasks, nil).Once()

		req, _ := http.NewRequest("GET", "/api/v1/tasks", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var response map[string]interface{}
		json.Unmarshal(resp.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))
	})

	// 3. Obtener tarea por ID
	t.Run("3. Obtener tarea por ID", func(t *testing.T) {
		task := &domain.Task{ID: 1, Title: "Tarea 1", UserID: 1}

		taskServiceMock.On("GetByID", uint(1)).Return(task, nil).Once()

		req, _ := http.NewRequest("GET", "/api/v1/tasks/1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	// 4. Actualizar tarea
	t.Run("4. Actualizar tarea", func(t *testing.T) {
		updatedTask := &domain.Task{ID: 1, Title: "Tarea actualizada", UserID: 1}

		taskServiceMock.On("Update", uint(1), uint(1), mock.AnythingOfType("*domain.UpdateTask")).Return(updatedTask, nil).Once()

		body := `{"title": "Tarea actualizada"}`
		req, _ := http.NewRequest("PUT", "/api/v1/tasks/1", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	// 5. Eliminar tarea
	t.Run("5. Eliminar tarea", func(t *testing.T) {
		taskServiceMock.On("Delete", uint(1), uint(1)).Return(nil).Once()

		req, _ := http.NewRequest("DELETE", "/api/v1/tasks/1", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	taskServiceMock.AssertExpectations(t)
}

func TestE2E_UnauthorizedAccess(t *testing.T) {
	setupE2EConfig()

	taskServiceMock := new(mocks.TaskServiceMock)
	taskHandler := handler.NewTaskHandler(taskServiceMock)

	router := gin.New()
	gin.SetMode(gin.TestMode)

	tasks := router.Group("/api/v1/tasks")
	tasks.Use(middleware.AuthMiddleware(config.AppConfig.JWTSecret))
	{
		tasks.GET("", taskHandler.GetAll)
	}

	tests := []struct {
		name           string
		authorization  string
		expectedStatus int
	}{
		{
			name:           "sin token",
			authorization:  "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "token inválido",
			authorization:  "Bearer invalid-token",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "formato incorrecto",
			authorization:  "Basic some-token",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/v1/tasks", nil)
			if tt.authorization != "" {
				req.Header.Set("Authorization", tt.authorization)
			}
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
		})
	}
}

func TestE2E_TaskNotFound(t *testing.T) {
	setupE2EConfig()

	taskServiceMock := new(mocks.TaskServiceMock)
	taskHandler := handler.NewTaskHandler(taskServiceMock)

	router := gin.New()
	gin.SetMode(gin.TestMode)

	tasks := router.Group("/api/v1/tasks")
	tasks.Use(middleware.AuthMiddleware(config.AppConfig.JWTSecret))
	{
		tasks.GET("/:id", taskHandler.GetByID)
	}

	token, _ := jwt.GenerateToken(1, "test@example.com", config.AppConfig.JWTSecret, config.AppConfig.JWTExpireIn)

	taskServiceMock.On("GetByID", uint(999)).Return(nil, service.ErrTaskNotFound).Once()

	req, _ := http.NewRequest("GET", "/api/v1/tasks/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	taskServiceMock.AssertExpectations(t)
}

func TestE2E_ValidationErrors(t *testing.T) {
	setupE2EConfig()

	authServiceMock := new(mocks.AuthServiceMock)
	authHandler := handler.NewAuthHandler(authServiceMock)

	router := gin.New()
	gin.SetMode(gin.TestMode)
	router.POST("/api/v1/auth/signup", authHandler.SignUpHandler)

	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "email inválido",
			body:           `{"full_name": "Test", "email": "invalid-email", "password": "password123"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "password muy corta",
			body:           `{"full_name": "Test", "email": "test@example.com", "password": "123"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "campos faltantes",
			body:           `{"full_name": "Test"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "body vacío",
			body:           `{}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/v1/auth/signup", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
		})
	}
}

// ========== Tests de Integración de Services ==========

func TestIntegration_AuthService_RegisterAndLogin(t *testing.T) {
	setupE2EConfig()

	mockUserRepo := new(repoMocks.MockUserRepository)
	authService := service.NewAuthService(mockUserRepo)

	// 1. Registro
	t.Run("registro y login integrado", func(t *testing.T) {
		// Setup para registro
		mockUserRepo.On("ExistsByEmail", "integration@example.com").Return(false, nil).Once()
		mockUserRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil).Once()

		registerReq := &domain.UserCreate{
			FullName: "Integration User",
			Email:    "integration@example.com",
			Password: "password123",
		}

		user, err := authService.Register(registerReq)
		require.NoError(t, err)
		require.NotNil(t, user)

		// Guardar el hash generado para el login
		hashedPassword := user.Password
		user.ID = 1

		// Setup para login - devolver el usuario con el hash correcto
		mockUserRepo.On("GetByEmail", "integration@example.com").Return(&domain.User{
			ID:       1,
			FullName: "Integration User",
			Email:    "integration@example.com",
			Password: hashedPassword,
		}, nil).Once()

		loginReq := &domain.UserLogin{
			Email:    "integration@example.com",
			Password: "password123",
		}

		token, loggedUser, err := authService.Login(loginReq)
		require.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.Equal(t, "Integration User", loggedUser.FullName)

		mockUserRepo.AssertExpectations(t)
	})
}

func TestIntegration_TaskService_CRUD(t *testing.T) {
	mockTaskRepo := new(repoMocks.TaskRepositoryMock)
	taskService := service.NewTaskService(mockTaskRepo)

	userID := uint(1)

	// 1. Crear
	t.Run("crear tarea", func(t *testing.T) {
		mockTaskRepo.On("Create", mock.AnythingOfType("*domain.Task")).Return(nil).Once()

		task, err := taskService.Create(userID, &domain.CreateTask{Title: "Nueva tarea"})
		require.NoError(t, err)
		assert.Equal(t, "Nueva tarea", task.Title)
		assert.Equal(t, userID, task.UserID)
	})

	// 2. Obtener
	t.Run("obtener tarea", func(t *testing.T) {
		expectedTask := &domain.Task{ID: 1, Title: "Nueva tarea", UserID: userID}
		mockTaskRepo.On("GetByID", uint(1)).Return(expectedTask, nil).Once()

		task, err := taskService.GetByID(1)
		require.NoError(t, err)
		assert.Equal(t, "Nueva tarea", task.Title)
	})

	// 3. Actualizar
	t.Run("actualizar tarea", func(t *testing.T) {
		existingTask := &domain.Task{ID: 1, Title: "Nueva tarea", UserID: userID}
		mockTaskRepo.On("GetByID", uint(1)).Return(existingTask, nil).Once()
		mockTaskRepo.On("Update", mock.AnythingOfType("*domain.Task")).Return(nil).Once()

		newTitle := "Tarea actualizada"
		task, err := taskService.Update(1, userID, &domain.UpdateTask{Title: &newTitle})
		require.NoError(t, err)
		assert.Equal(t, "Tarea actualizada", task.Title)
	})

	// 4. Eliminar
	t.Run("eliminar tarea", func(t *testing.T) {
		existingTask := &domain.Task{ID: 1, Title: "Tarea actualizada", UserID: userID}
		mockTaskRepo.On("GetByID", uint(1)).Return(existingTask, nil).Once()
		mockTaskRepo.On("Delete", uint(1)).Return(nil).Once()

		err := taskService.Delete(1, userID)
		require.NoError(t, err)
	})

	mockTaskRepo.AssertExpectations(t)
}

// ========== Benchmark Tests ==========

func BenchmarkPasswordHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		utils.HashPassword("benchmark-password-123")
	}
}

func BenchmarkTokenGeneration(b *testing.B) {
	setupE2EConfig()
	for i := 0; i < b.N; i++ {
		jwt.GenerateToken(1, "benchmark@example.com", config.AppConfig.JWTSecret, config.AppConfig.JWTExpireIn)
	}
}

func BenchmarkTokenValidation(b *testing.B) {
	setupE2EConfig()
	token, _ := jwt.GenerateToken(1, "benchmark@example.com", config.AppConfig.JWTSecret, config.AppConfig.JWTExpireIn)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jwt.ValidateToken(token, config.AppConfig.JWTSecret)
	}
}
