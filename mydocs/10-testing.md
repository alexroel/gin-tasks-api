# Testing en Go

En esta sección aprenderemos todo lo necesario para implementar tests profesionales en aplicaciones Go con Gin Framework. Cubriremos desde tests unitarios básicos hasta tests de integración E2E.

---
## Introducción

El testing es una parte fundamental del desarrollo de software profesional. En Go, el testing está integrado en el lenguaje a través del paquete `testing` y la herramienta `go test`.

### ¿Por qué es importante el testing?

- **Confianza en el código**: Los tests garantizan que el código funciona como se espera
- **Refactorización segura**: Puedes modificar código con la seguridad de que los tests detectarán regresiones
- **Documentación viva**: Los tests sirven como documentación de cómo usar el código
- **Desarrollo más rápido**: A largo plazo, los tests reducen el tiempo de debugging

### Herramientas que usaremos

```go
// Paquetes principales
import (
    "testing"                           // Paquete estándar de Go
    "net/http/httptest"                 // Para testing HTTP
    "github.com/stretchr/testify/assert" // Assertions más expresivas
    "github.com/stretchr/testify/mock"   // Mocking
    "github.com/stretchr/testify/require" // Assertions que detienen el test
    "github.com/stretchr/testify/suite"  // Test suites
    "github.com/gin-gonic/gin"          // Framework web
)
```

---
## Testing en Go

### Convenciones básicas

1. **Nombre de archivos**: Los archivos de test terminan en `_test.go`
2. **Nombre de funciones**: Las funciones de test comienzan con `Test`
3. **Firma de funciones**: `func TestNombreDescriptivo(t *testing.T)`

### Estructura de un test básico

```go
package miPaquete

import "testing"

func TestSuma(t *testing.T) {
    // Arrange (Preparar)
    a := 2
    b := 3
    expected := 5

    // Act (Actuar)
    result := Suma(a, b)

    // Assert (Verificar)
    if result != expected {
        t.Errorf("Suma(%d, %d) = %d; esperado %d", a, b, result, expected)
    }
}
```

### Usando testify/assert

```go
package miPaquete

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestSuma(t *testing.T) {
    result := Suma(2, 3)
    
    assert.Equal(t, 5, result)
    assert.NotEqual(t, 0, result)
    assert.True(t, result > 0)
}
```

### assert vs require

```go
func TestConRequire(t *testing.T) {
    user, err := GetUser(1)
    
    // require detiene el test si falla (útil para precondiciones)
    require.NoError(t, err)
    require.NotNil(t, user)
    
    // assert continúa aunque falle (útil para verificaciones)
    assert.Equal(t, "John", user.Name)
    assert.Equal(t, "john@example.com", user.Email)
}
```

### Ejecutar tests

```bash
# Ejecutar todos los tests
go test ./...

# Con verbose
go test -v ./...

# Un paquete específico
go test ./internal/service

# Un test específico
go test -run TestSuma ./...

# Sin cache
go test -count=1 ./...
```

---
## Tests unitarios

Los tests unitarios prueban una unidad de código de forma aislada, sin dependencias externas.

### Ejemplo: Testing de funciones utilitarias

**Archivo**: `pkg/utils/password.go`

```go
package utils

import (
    "errors"
    "golang.org/x/crypto/bcrypt"
)

var (
    ErrPasswordTooShort = errors.New("la contraseña debe tener al menos 6 caracteres")
    ErrPasswordTooLong  = errors.New("la contraseña no puede exceder 72 caracteres")
)

func HashPassword(password string) (string, error) {
    if len(password) < 6 {
        return "", ErrPasswordTooShort
    }
    if len(password) > 72 {
        return "", ErrPasswordTooLong
    }
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func CheckPassword(hashedPassword, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}
```

**Archivo**: `pkg/utils/password_test.go`

```go
package utils

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestHashPassword_Success(t *testing.T) {
    password := "password123"

    hash, err := HashPassword(password)

    require.NoError(t, err)
    assert.NotEmpty(t, hash)
    assert.NotEqual(t, password, hash)
}

func TestHashPassword_TooShort(t *testing.T) {
    password := "12345" // 5 caracteres

    hash, err := HashPassword(password)

    assert.Error(t, err)
    assert.Equal(t, ErrPasswordTooShort, err)
    assert.Empty(t, hash)
}

func TestCheckPassword_Valid(t *testing.T) {
    password := "password123"
    hash, _ := HashPassword(password)

    result := CheckPassword(hash, password)

    assert.True(t, result)
}

func TestCheckPassword_Invalid(t *testing.T) {
    password := "password123"
    hash, _ := HashPassword(password)

    result := CheckPassword(hash, "wrongpassword")

    assert.False(t, result)
}
```

### Ejemplo: Testing de JWT

**Archivo**: `pkg/jwt/jwt_test.go`

```go
package jwt

import (
    "testing"
    "time"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

const testSecret = "test-secret-key-12345"

func TestGenerateToken_Success(t *testing.T) {
    userID := uint(1)
    email := "test@example.com"

    token, err := GenerateToken(userID, email, testSecret, time.Hour)

    require.NoError(t, err)
    assert.NotEmpty(t, token)
}

func TestValidateToken_Success(t *testing.T) {
    token, _ := GenerateToken(1, "test@example.com", testSecret, time.Hour)

    claims, err := ValidateToken(token, testSecret)

    require.NoError(t, err)
    assert.Equal(t, uint(1), claims.UserID)
    assert.Equal(t, "test@example.com", claims.Email)
}

func TestValidateToken_Expired(t *testing.T) {
    // Token que ya expiró
    token, _ := GenerateToken(1, "test@example.com", testSecret, -time.Hour)

    claims, err := ValidateToken(token, testSecret)

    assert.Error(t, err)
    assert.Nil(t, claims)
}

func TestValidateToken_WrongSecret(t *testing.T) {
    token, _ := GenerateToken(1, "test@example.com", testSecret, time.Hour)

    claims, err := ValidateToken(token, "wrong-secret")

    assert.Error(t, err)
    assert.Nil(t, claims)
}
```

### Table-Driven Tests

Los table-driven tests permiten probar múltiples casos de forma organizada:

```go
func TestHashPassword_TableDriven(t *testing.T) {
    tests := []struct {
        name        string
        password    string
        wantErr     bool
        expectedErr error
    }{
        {
            name:     "contraseña válida mínima",
            password: "123456",
            wantErr:  false,
        },
        {
            name:        "contraseña muy corta",
            password:    "12345",
            wantErr:     true,
            expectedErr: ErrPasswordTooShort,
        },
        {
            name:        "contraseña muy larga - 73 chars",
            password:    strings.Repeat("a", 73),
            wantErr:     true,
            expectedErr: ErrPasswordTooLong,
        },
        {
            name:     "contraseña con caracteres especiales",
            password: "P@ssw0rd!#$%",
            wantErr:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            hash, err := HashPassword(tt.password)

            if tt.wantErr {
                assert.Error(t, err)
                assert.Equal(t, tt.expectedErr, err)
                assert.Empty(t, hash)
            } else {
                assert.NoError(t, err)
                assert.NotEmpty(t, hash)
                // Verificar que el hash funciona
                assert.True(t, CheckPassword(hash, tt.password))
            }
        })
    }
}
```

---
## Testing de Handlers con Gin

Para testear handlers de Gin usamos `httptest` para simular peticiones HTTP.

### Configuración del router de test

```go
package handler

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

// setupRouter configura Gin en modo test
func setupRouter(handler *TaskHandler) *gin.Engine {
    gin.SetMode(gin.TestMode)
    router := gin.New()
    return router
}
```

### Ejemplo: Testing de TaskHandler

```go
func TestTaskHandler_Create_Success(t *testing.T) {
    // Arrange
    mockService := new(mocks.TaskServiceMock)
    handler := NewTaskHandler(mockService)
    router := setupRouter(handler)

    // Configurar ruta con contexto de usuario
    router.POST("/tasks", func(c *gin.Context) {
        c.Set("userID", uint(1))
        handler.Create(c)
    })

    // Configurar mock
    task := &domain.Task{ID: 1, Title: "Nueva tarea", UserID: 1}
    mockService.On("Create", uint(1), mock.AnythingOfType("*domain.CreateTask")).Return(task, nil)

    // Crear request
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
    router := setupRouter(handler)
    router.POST("/tasks", handler.Create) // Sin userID en contexto

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
    router := setupRouter(handler)

    router.POST("/tasks", func(c *gin.Context) {
        c.Set("userID", uint(1))
        handler.Create(c)
    })

    // Body inválido - falta el campo title
    reqBody := `{"invalid": "body"}`
    req, _ := http.NewRequest("POST", "/tasks", bytes.NewBufferString(reqBody))
    req.Header.Set("Content-Type", "application/json")
    resp := httptest.NewRecorder()

    // Act
    router.ServeHTTP(resp, req)

    // Assert
    assert.Equal(t, http.StatusBadRequest, resp.Code)
}
```

### Testing con parámetros de ruta

```go
func TestTaskHandler_GetByID_Success(t *testing.T) {
    mockService := new(mocks.TaskServiceMock)
    handler := NewTaskHandler(mockService)
    router := setupRouter(handler)

    router.GET("/tasks/:id", func(c *gin.Context) {
        c.Set("userID", uint(1))
        handler.GetByID(c)
    })

    task := &domain.Task{ID: 1, Title: "Tarea 1", UserID: 1}
    mockService.On("GetByID", uint(1)).Return(task, nil)

    req, _ := http.NewRequest("GET", "/tasks/1", nil)
    resp := httptest.NewRecorder()

    router.ServeHTTP(resp, req)

    assert.Equal(t, http.StatusOK, resp.Code)
    mockService.AssertExpectations(t)
}

func TestTaskHandler_GetByID_NotFound(t *testing.T) {
    mockService := new(mocks.TaskServiceMock)
    handler := NewTaskHandler(mockService)
    router := setupRouter(handler)

    router.GET("/tasks/:id", func(c *gin.Context) {
        c.Set("userID", uint(1))
        handler.GetByID(c)
    })

    mockService.On("GetByID", uint(999)).Return(nil, service.ErrTaskNotFound)

    req, _ := http.NewRequest("GET", "/tasks/999", nil)
    resp := httptest.NewRecorder()

    router.ServeHTTP(resp, req)

    assert.Equal(t, http.StatusNotFound, resp.Code)
}
```

---
## Mocking de repositorios

Los mocks nos permiten aislar la unidad bajo prueba simulando sus dependencias.

### Creando un Mock con testify/mock

**Archivo**: `internal/repository/mocks/user_repository_mock.go`

```go
package mocks

import (
    "github.com/alexroel/gin-tasks-api/internal/domain"
    "github.com/stretchr/testify/mock"
)

// MockUserRepository es un mock del repositorio de usuarios
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(user *domain.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uint) (*domain.User, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*domain.User, error) {
    args := m.Called(email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) ExistsByEmail(email string) (bool, error) {
    args := m.Called(email)
    return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) Update(user *domain.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
    args := m.Called(id)
    return args.Error(0)
}
```

### Mock de TaskRepository

**Archivo**: `internal/repository/mocks/task_repository_mock.go`

```go
package mocks

import (
    "github.com/alexroel/gin-tasks-api/internal/domain"
    "github.com/stretchr/testify/mock"
)

type TaskRepositoryMock struct {
    mock.Mock
}

func (m *TaskRepositoryMock) Create(task *domain.Task) error {
    args := m.Called(task)
    return args.Error(0)
}

func (m *TaskRepositoryMock) GetByID(id uint) (*domain.Task, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *TaskRepositoryMock) GetByUserID(userID uint) ([]domain.Task, error) {
    args := m.Called(userID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *TaskRepositoryMock) Update(task *domain.Task) error {
    args := m.Called(task)
    return args.Error(0)
}

func (m *TaskRepositoryMock) Delete(id uint) error {
    args := m.Called(id)
    return args.Error(0)
}
```

### Usando los Mocks en tests

```go
func TestAuthService_Register_Success(t *testing.T) {
    // Crear mock
    mockRepo := new(mocks.MockUserRepository)
    service := NewAuthService(mockRepo)

    req := &domain.UserCreate{
        FullName: "Test User",
        Email:    "test@example.com",
        Password: "password123",
    }

    // Configurar expectativas del mock
    mockRepo.On("ExistsByEmail", req.Email).Return(false, nil)
    mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

    // Ejecutar
    user, err := service.Register(req)

    // Verificar
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, req.Email, user.Email)
    
    // Verificar que se llamaron los métodos esperados
    mockRepo.AssertExpectations(t)
}

func TestAuthService_Register_EmailExists(t *testing.T) {
    mockRepo := new(mocks.MockUserRepository)
    service := NewAuthService(mockRepo)

    req := &domain.UserCreate{
        FullName: "Test User",
        Email:    "existing@example.com",
        Password: "password123",
    }

    // El email ya existe
    mockRepo.On("ExistsByEmail", req.Email).Return(true, nil)

    user, err := service.Register(req)

    assert.Error(t, err)
    assert.Nil(t, user)
    assert.Contains(t, err.Error(), "Ya existe un usuario")
    mockRepo.AssertExpectations(t)
}
```

### Métodos útiles de mock

```go
// Configurar retorno simple
mockRepo.On("GetByID", uint(1)).Return(user, nil)

// Usar matchers para argumentos
mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)
mockRepo.On("GetByEmail", mock.Anything).Return(nil, nil)

// Limitar número de llamadas
mockRepo.On("GetByID", uint(1)).Return(user, nil).Once()
mockRepo.On("GetByID", uint(1)).Return(user, nil).Times(3)

// Verificar expectativas
mockRepo.AssertExpectations(t)

// Verificar que se llamó un método específico
mockRepo.AssertCalled(t, "GetByID", uint(1))

// Verificar que NO se llamó un método
mockRepo.AssertNotCalled(t, "Delete", mock.Anything)
```

---
## Testing de casos de uso

Los casos de uso (services) contienen la lógica de negocio y son críticos de testear.

### Ejemplo: Testing de TaskService

```go
package service

import (
    "errors"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
)

func TestTaskService_Create_Success(t *testing.T) {
    mockRepo := new(mocks.TaskRepositoryMock)
    service := NewTaskService(mockRepo)

    userID := uint(1)
    req := &domain.CreateTask{Title: "Nueva tarea"}

    mockRepo.On("Create", mock.AnythingOfType("*domain.Task")).Return(nil)

    task, err := service.Create(userID, req)

    require.NoError(t, err)
    assert.Equal(t, "Nueva tarea", task.Title)
    assert.Equal(t, userID, task.UserID)
    assert.False(t, task.Completed)
    mockRepo.AssertExpectations(t)
}

func TestTaskService_Update_Success(t *testing.T) {
    mockRepo := new(mocks.TaskRepositoryMock)
    service := NewTaskService(mockRepo)

    existingTask := &domain.Task{ID: 1, Title: "Tarea original", UserID: 1}
    newTitle := "Tarea actualizada"
    req := &domain.UpdateTask{Title: &newTitle}

    mockRepo.On("GetByID", uint(1)).Return(existingTask, nil)
    mockRepo.On("Update", mock.AnythingOfType("*domain.Task")).Return(nil)

    task, err := service.Update(1, 1, req)

    require.NoError(t, err)
    assert.Equal(t, "Tarea actualizada", task.Title)
    mockRepo.AssertExpectations(t)
}

func TestTaskService_Update_NotFound(t *testing.T) {
    mockRepo := new(mocks.TaskRepositoryMock)
    service := NewTaskService(mockRepo)

    newTitle := "Tarea actualizada"
    req := &domain.UpdateTask{Title: &newTitle}

    mockRepo.On("GetByID", uint(999)).Return(nil, nil)

    task, err := service.Update(999, 1, req)

    assert.Error(t, err)
    assert.Equal(t, ErrTaskNotFound, err)
    assert.Nil(t, task)
}

func TestTaskService_Update_Unauthorized(t *testing.T) {
    mockRepo := new(mocks.TaskRepositoryMock)
    service := NewTaskService(mockRepo)

    // Tarea pertenece al usuario 1
    existingTask := &domain.Task{ID: 1, Title: "Tarea", UserID: 1}
    newTitle := "Tarea actualizada"
    req := &domain.UpdateTask{Title: &newTitle}

    mockRepo.On("GetByID", uint(1)).Return(existingTask, nil)

    // Usuario 2 intenta actualizar
    task, err := service.Update(1, 2, req)

    assert.Error(t, err)
    assert.Equal(t, ErrTaskUnauthorized, err)
    assert.Nil(t, task)
}

func TestTaskService_Delete_Success(t *testing.T) {
    mockRepo := new(mocks.TaskRepositoryMock)
    service := NewTaskService(mockRepo)

    existingTask := &domain.Task{ID: 1, Title: "Tarea", UserID: 1}

    mockRepo.On("GetByID", uint(1)).Return(existingTask, nil)
    mockRepo.On("Delete", uint(1)).Return(nil)

    err := service.Delete(1, 1)

    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

### Testing de AuthService con hash real

```go
func TestAuthService_Login_Success(t *testing.T) {
    // Configurar AppConfig para generar tokens
    config.AppConfig = &config.Config{
        JWTSecret:   "test-secret-key-12345",
        JWTExpireIn: 24 * time.Hour,
    }

    mockRepo := new(mocks.MockUserRepository)
    service := NewAuthService(mockRepo)

    // Generar hash real
    realPassword := "password123"
    hashedPassword, _ := utils.HashPassword(realPassword)

    existingUser := &domain.User{
        ID:       1,
        FullName: "Test User",
        Email:    "test@example.com",
        Password: hashedPassword,
    }

    req := &domain.UserLogin{
        Email:    "test@example.com",
        Password: realPassword,
    }

    mockRepo.On("GetByEmail", req.Email).Return(existingUser, nil)

    token, user, err := service.Login(req)

    assert.NoError(t, err)
    assert.NotEmpty(t, token)
    assert.Equal(t, existingUser.Email, user.Email)
    mockRepo.AssertExpectations(t)
}
```

---
## Coverage de tests

El coverage indica qué porcentaje del código está cubierto por tests.

### Ejecutar tests con coverage

```bash
# Coverage básico
go test ./... -cover

# Generar archivo de coverage
go test ./... -coverprofile=coverage.out

# Ver coverage por función
go tool cover -func=coverage.out

# Ver coverage en HTML (muy útil)
go tool cover -html=coverage.out -o coverage.html
```

### Ejemplo de salida de coverage

```
github.com/alexroel/gin-tasks-api/internal/config/config.go:34:      LoadConfig      100.0%
github.com/alexroel/gin-tasks-api/internal/config/config.go:73:      validateConfig  100.0%
github.com/alexroel/gin-tasks-api/internal/config/config.go:87:      getEnv          100.0%
github.com/alexroel/gin-tasks-api/internal/handler/task_handler.go:  Create          85.7%
github.com/alexroel/gin-tasks-api/internal/service/task_service.go:  Update          90.0%
total:                                                               (statements)    75.3%
```

### Objetivos de coverage recomendados

| Capa | Coverage Mínimo | Ideal |
|------|-----------------|-------|
| Utils/Pkg | 80% | 90%+ |
| Services | 80% | 90%+ |
| Handlers | 70% | 80%+ |
| Middleware | 90% | 100% |
| Repository | N/A* | N/A* |

*Los repositorios generalmente se testean con integración, no unitarios.

### Coverage por paquete

```bash
# Ver coverage de cada paquete
go test ./internal/... ./pkg/... -cover

# Ejemplo de salida:
# ok   internal/config       coverage: 42.9% of statements
# ok   internal/handler      coverage: 75.3% of statements
# ok   internal/middleware   coverage: 100.0% of statements
# ok   internal/service      coverage: 86.7% of statements
# ok   pkg/jwt               coverage: 83.3% of statements
# ok   pkg/utils             coverage: 75.0% of statements
```

---
## Testing de API REST

Los tests E2E (End-to-End) prueban flujos completos de la API.

### Estructura de tests E2E

```go
package e2e

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func setupE2EConfig() {
    config.AppConfig = &config.Config{
        JWTSecret:   "e2e-test-secret-key-12345",
        JWTExpireIn: 24 * time.Hour,
    }
}
```

### Test E2E: Flujo de registro y login

```go
func TestE2E_RegisterAndLogin_Flow(t *testing.T) {
    setupE2EConfig()

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
    })

    // 2. Login
    t.Run("2. Login exitoso", func(t *testing.T) {
        user := &domain.User{ID: 1, FullName: "E2E User", Email: "e2e@example.com"}
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
        data := response["data"].(map[string]interface{})
        assert.NotEmpty(t, data["token"])
    })

    authServiceMock.AssertExpectations(t)
}
```

### Test E2E: CRUD de tareas

```go
func TestE2E_TaskCRUD_Flow(t *testing.T) {
    setupE2EConfig()

    taskServiceMock := new(mocks.TaskServiceMock)
    taskHandler := handler.NewTaskHandler(taskServiceMock)

    router := gin.New()
    gin.SetMode(gin.TestMode)

    tasks := router.Group("/api/v1/tasks")
    tasks.Use(middleware.AuthMiddleware(config.AppConfig.JWTSecret))
    {
        tasks.POST("", taskHandler.Create)
        tasks.GET("", taskHandler.GetAll)
        tasks.GET("/:id", taskHandler.GetByID)
        tasks.PUT("/:id", taskHandler.Update)
        tasks.DELETE("/:id", taskHandler.Delete)
    }

    // Token válido para las peticiones
    token, _ := jwt.GenerateToken(1, "e2e@example.com", config.AppConfig.JWTSecret, config.AppConfig.JWTExpireIn)

    // 1. Crear tarea
    t.Run("1. Crear tarea", func(t *testing.T) {
        task := &domain.Task{ID: 1, Title: "Mi tarea E2E", UserID: 1}
        taskServiceMock.On("Create", uint(1), mock.AnythingOfType("*domain.CreateTask")).Return(task, nil).Once()

        body := `{"title": "Mi tarea E2E"}`
        req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBufferString(body))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Authorization", "Bearer "+token)
        resp := httptest.NewRecorder()

        router.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusCreated, resp.Code)
    })

    // 2. Obtener tareas
    t.Run("2. Obtener tareas", func(t *testing.T) {
        tasks := []domain.Task{{ID: 1, Title: "Tarea 1", UserID: 1}}
        taskServiceMock.On("GetByUserID", uint(1)).Return(tasks, nil).Once()

        req, _ := http.NewRequest("GET", "/api/v1/tasks", nil)
        req.Header.Set("Authorization", "Bearer "+token)
        resp := httptest.NewRecorder()

        router.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusOK, resp.Code)
    })

    // 3. Actualizar tarea
    t.Run("3. Actualizar tarea", func(t *testing.T) {
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

    // 4. Eliminar tarea
    t.Run("4. Eliminar tarea", func(t *testing.T) {
        taskServiceMock.On("Delete", uint(1), uint(1)).Return(nil).Once()

        req, _ := http.NewRequest("DELETE", "/api/v1/tasks/1", nil)
        req.Header.Set("Authorization", "Bearer "+token)
        resp := httptest.NewRecorder()

        router.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusOK, resp.Code)
    })

    taskServiceMock.AssertExpectations(t)
}
```

### Test de acceso no autorizado

```go
func TestE2E_UnauthorizedAccess(t *testing.T) {
    setupE2EConfig()

    taskServiceMock := new(mocks.TaskServiceMock)
    taskHandler := handler.NewTaskHandler(taskServiceMock)

    router := gin.New()
    gin.SetMode(gin.TestMode)
    
    tasks := router.Group("/api/v1/tasks")
    tasks.Use(middleware.AuthMiddleware(config.AppConfig.JWTSecret))
    tasks.GET("", taskHandler.GetAll)

    tests := []struct {
        name           string
        authorization  string
        expectedStatus int
    }{
        {"sin token", "", http.StatusUnauthorized},
        {"token inválido", "Bearer invalid-token", http.StatusUnauthorized},
        {"formato incorrecto", "Basic some-token", http.StatusUnauthorized},
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
```

### Benchmarks

```go
func BenchmarkPasswordHash(b *testing.B) {
    for i := 0; i < b.N; i++ {
        utils.HashPassword("benchmark-password-123")
    }
}

func BenchmarkTokenGeneration(b *testing.B) {
    setupE2EConfig()
    for i := 0; i < b.N; i++ {
        jwt.GenerateToken(1, "bench@example.com", config.AppConfig.JWTSecret, config.AppConfig.JWTExpireIn)
    }
}

func BenchmarkTokenValidation(b *testing.B) {
    setupE2EConfig()
    token, _ := jwt.GenerateToken(1, "bench@example.com", config.AppConfig.JWTSecret, config.AppConfig.JWTExpireIn)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        jwt.ValidateToken(token, config.AppConfig.JWTSecret)
    }
}
```

Ejecutar benchmarks:

```bash
go test ./tests/e2e -bench=. -benchmem
```

---
## Buenas prácticas de testing

### 1. Patrón AAA (Arrange-Act-Assert)

```go
func TestExample(t *testing.T) {
    // Arrange - Preparar datos y mocks
    mockRepo := new(mocks.MockRepository)
    service := NewService(mockRepo)
    mockRepo.On("GetByID", uint(1)).Return(expectedData, nil)

    // Act - Ejecutar la acción
    result, err := service.GetByID(1)

    // Assert - Verificar resultados
    assert.NoError(t, err)
    assert.Equal(t, expectedData, result)
    mockRepo.AssertExpectations(t)
}
```

### 2. Nombres descriptivos

```go
// ✅ Bueno - Describe qué se prueba y el escenario
func TestAuthService_Register_Success(t *testing.T) {}
func TestAuthService_Register_EmailAlreadyExists(t *testing.T) {}
func TestAuthService_Register_PasswordTooShort(t *testing.T) {}

// ❌ Malo - No describe el escenario
func TestRegister(t *testing.T) {}
func TestRegister2(t *testing.T) {}
```

### 3. Tests independientes

```go
// ✅ Cada test es independiente
func TestTaskService_Create(t *testing.T) {
    mockRepo := new(mocks.TaskRepositoryMock) // Mock nuevo
    service := NewTaskService(mockRepo)
    // ...
}

func TestTaskService_Update(t *testing.T) {
    mockRepo := new(mocks.TaskRepositoryMock) // Mock nuevo
    service := NewTaskService(mockRepo)
    // ...
}
```

### 4. Usar require para precondiciones

```go
func TestWithPreconditions(t *testing.T) {
    user, err := service.GetUser(1)
    
    // Si falla aquí, no tiene sentido continuar
    require.NoError(t, err)
    require.NotNil(t, user)
    
    // Estas son las verificaciones principales
    assert.Equal(t, "John", user.Name)
    assert.Equal(t, "john@example.com", user.Email)
}
```

### 5. Table-driven tests para múltiples casos

```go
func TestValidation_TableDriven(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"válido", "valid@email.com", false},
        {"sin @", "invalid.com", true},
        {"vacío", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### 6. Cleanup de recursos

```go
func TestWithCleanup(t *testing.T) {
    // Setup
    tempFile, _ := os.CreateTemp("", "test")
    
    // Cleanup al final del test
    t.Cleanup(func() {
        os.Remove(tempFile.Name())
    })
    
    // Test...
}
```

### 7. No testear implementaciones internas

```go
// ✅ Testear comportamiento público
func TestService_GetUser(t *testing.T) {
    user, err := service.GetUser(1)
    assert.NoError(t, err)
    assert.Equal(t, "John", user.Name)
}

// ❌ No testear métodos privados directamente
// func TestService_validateUser(t *testing.T) {} // Privado
```

### 8. Estructura de carpetas recomendada

```
internal/
├── handler/
│   ├── task_handler.go
│   └── task_handler_test.go    # Tests junto al código
├── service/
│   ├── task_service.go
│   ├── task_service_test.go
│   └── mocks/
│       └── task_service_mock.go
├── repository/
│   ├── task_repository.go
│   └── mocks/
│       └── task_repository_mock.go
pkg/
├── jwt/
│   ├── jwt.go
│   └── jwt_test.go
├── utils/
│   ├── password.go
│   └── password_test.go
tests/
└── e2e/
    └── e2e_test.go             # Tests de integración
```

### 9. Comandos útiles

```bash
# Ejecutar todos los tests
go test ./...

# Con verbose y coverage
go test -v -cover ./...

# Sin cache (forzar re-ejecución)
go test -count=1 ./...

# Solo un test específico
go test -run TestTaskService_Create ./internal/service

# Generar reporte de coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Ejecutar benchmarks
go test -bench=. -benchmem ./...

# Detectar race conditions
go test -race ./...
```

---

## Resumen

| Tema | Herramientas/Conceptos |
|------|------------------------|
| Testing básico | `testing`, `go test`, AAA pattern |
| Assertions | `testify/assert`, `testify/require` |
| Mocking | `testify/mock`, interfaces |
| HTTP Testing | `httptest`, `gin.TestMode` |
| Coverage | `-cover`, `-coverprofile` |
| Table-driven | Subtests con `t.Run()` |
| E2E | Flujos completos, benchmarks |

### Métricas del proyecto

```
✅ 120+ tests pasando
✅ 0 tests skipped
✅ Coverage promedio: ~75%
✅ config.go: 100% coverage
✅ middleware: 100% coverage
✅ services: 86.7% coverage
```