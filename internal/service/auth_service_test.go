package service

import (
	"errors"
	"testing"

	"github.com/alexroel/gin-tasks-api/internal/config"
	"github.com/alexroel/gin-tasks-api/internal/domain"
	"github.com/alexroel/gin-tasks-api/internal/repository/mocks"
	"github.com/alexroel/gin-tasks-api/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ========== Tests para Register ==========

func TestAuthService_Register_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewAuthService(mockRepo)

	req := &domain.UserCreate{
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	mockRepo.On("ExistsByEmail", req.Email).Return(false, nil)
	mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

	user, err := service.Register(req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.FullName, user.FullName)
	assert.Equal(t, req.Email, user.Email)
	assert.NotEqual(t, req.Password, user.Password) // Password debe estar hasheado
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

	mockRepo.On("ExistsByEmail", req.Email).Return(true, nil)

	user, err := service.Register(req)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "Ya existe un usuario")
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Register_ExistsByEmailError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewAuthService(mockRepo)

	req := &domain.UserCreate{
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	mockRepo.On("ExistsByEmail", req.Email).Return(false, errors.New("db error"))

	user, err := service.Register(req)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Register_PasswordTooShort(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewAuthService(mockRepo)

	req := &domain.UserCreate{
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "123", // Muy corta
	}

	mockRepo.On("ExistsByEmail", req.Email).Return(false, nil)

	user, err := service.Register(req)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Register_CreateError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewAuthService(mockRepo)

	req := &domain.UserCreate{
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	mockRepo.On("ExistsByEmail", req.Email).Return(false, nil)
	mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(errors.New("db error"))

	user, err := service.Register(req)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

// ========== Tests para Login ==========

func setupConfigForLoginTest() {
	// Configurar AppConfig necesario para generar tokens
	config.AppConfig = &config.Config{
		JWTSecret:   "test-secret-key-12345",
		JWTExpireIn: 24 * 60 * 60 * 1000000000, // 24 horas en nanosegundos
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	setupConfigForLoginTest()
	
	mockRepo := new(mocks.MockUserRepository)
	service := NewAuthService(mockRepo)

	// Generar un hash real para la contraseña "password123"
	realPassword := "password123"
	hashedPassword, err := utils.HashPassword(realPassword)
	assert.NoError(t, err)

	existingUser := &domain.User{
		FullName: "Test User",
		Email:    "test@example.com",
		Password: hashedPassword,
	}
	existingUser.ID = 1

	req := &domain.UserLogin{
		Email:    "test@example.com",
		Password: realPassword,
	}

	mockRepo.On("GetByEmail", req.Email).Return(existingUser, nil)

	token, user, err := service.Login(req)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotNil(t, user)
	assert.Equal(t, existingUser.Email, user.Email)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewAuthService(mockRepo)

	req := &domain.UserLogin{
		Email:    "notfound@example.com",
		Password: "password123",
	}

	mockRepo.On("GetByEmail", req.Email).Return(nil, nil)

	token, user, err := service.Login(req)

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "Credenciales inválidas")
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewAuthService(mockRepo)

	existingUser := &domain.User{
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "$2a$10$validhashhere", // Hash que no coincidirá
	}
	existingUser.ID = 1

	req := &domain.UserLogin{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	mockRepo.On("GetByEmail", req.Email).Return(existingUser, nil)

	token, user, err := service.Login(req)

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "Credenciales inválidas")
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_DBError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewAuthService(mockRepo)

	req := &domain.UserLogin{
		Email:    "test@example.com",
		Password: "password123",
	}

	mockRepo.On("GetByEmail", req.Email).Return(nil, errors.New("db connection error"))

	token, user, err := service.Login(req)

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

// ========== Tests para GetUserByID ==========

func TestAuthService_GetUserByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewAuthService(mockRepo)

	expectedUser := &domain.User{
		FullName: "Test User",
		Email:    "test@example.com",
	}
	expectedUser.ID = 1

	mockRepo.On("GetByID", uint(1)).Return(expectedUser, nil)

	user, err := service.GetUserByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.FullName, user.FullName)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_GetUserByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewAuthService(mockRepo)

	mockRepo.On("GetByID", uint(999)).Return(nil, nil)

	user, err := service.GetUserByID(999)

	assert.NoError(t, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_GetUserByID_Error(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewAuthService(mockRepo)

	mockRepo.On("GetByID", uint(1)).Return(nil, errors.New("db error"))

	user, err := service.GetUserByID(1)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

// ========== Tests para UpdateProfile ==========

func TestAuthService_UpdateProfile_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewAuthService(mockRepo)

	existingUser := &domain.User{
		FullName: "Old Name",
		Email:    "old@example.com",
		Password: "hashedpassword",
	}
	existingUser.ID = 1

	newName := "New Name"
	req := &domain.UserUpdate{
		FullName: &newName,
	}

	mockRepo.On("GetByID", uint(1)).Return(existingUser, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.User")).Return(nil)

	user, err := service.UpdateProfile(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, newName, user.FullName)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_UpdateProfile_ChangeEmail(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewAuthService(mockRepo)

	existingUser := &domain.User{
		FullName: "Test User",
		Email:    "old@example.com",
		Password: "hashedpassword",
	}
	existingUser.ID = 1

	newEmail := "new@example.com"
	req := &domain.UserUpdate{
		Email: &newEmail,
	}

	mockRepo.On("GetByID", uint(1)).Return(existingUser, nil)
	mockRepo.On("ExistsByEmail", newEmail).Return(false, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.User")).Return(nil)

	user, err := service.UpdateProfile(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, newEmail, user.Email)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_UpdateProfile_EmailExists(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewAuthService(mockRepo)

	existingUser := &domain.User{
		FullName: "Test User",
		Email:    "old@example.com",
		Password: "hashedpassword",
	}
	existingUser.ID = 1

	newEmail := "taken@example.com"
	req := &domain.UserUpdate{
		Email: &newEmail,
	}

	mockRepo.On("GetByID", uint(1)).Return(existingUser, nil)
	mockRepo.On("ExistsByEmail", newEmail).Return(true, nil)

	user, err := service.UpdateProfile(1, req)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "email ya está registrado")
	mockRepo.AssertExpectations(t)
}

func TestAuthService_UpdateProfile_ChangePassword(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewAuthService(mockRepo)

	existingUser := &domain.User{
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "oldhash",
	}
	existingUser.ID = 1

	newPassword := "newpassword123"
	req := &domain.UserUpdate{
		Password: &newPassword,
	}

	mockRepo.On("GetByID", uint(1)).Return(existingUser, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.User")).Return(nil)

	user, err := service.UpdateProfile(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEqual(t, "oldhash", user.Password)    // Password cambió
	assert.NotEqual(t, newPassword, user.Password)  // No es texto plano
	mockRepo.AssertExpectations(t)
}

func TestAuthService_UpdateProfile_UserNotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewAuthService(mockRepo)

	newName := "New Name"
	req := &domain.UserUpdate{
		FullName: &newName,
	}

	mockRepo.On("GetByID", uint(999)).Return(nil, nil)

	user, err := service.UpdateProfile(999, req)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "usuario no encontrado")
	mockRepo.AssertExpectations(t)
}

func TestAuthService_UpdateProfile_SameEmail(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewAuthService(mockRepo)

	existingUser := &domain.User{
		FullName: "Test User",
		Email:    "same@example.com",
		Password: "hashedpassword",
	}
	existingUser.ID = 1

	sameEmail := "same@example.com"
	req := &domain.UserUpdate{
		Email: &sameEmail,
	}

	mockRepo.On("GetByID", uint(1)).Return(existingUser, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.User")).Return(nil)

	user, err := service.UpdateProfile(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	// No debe verificar ExistsByEmail si es el mismo email
	mockRepo.AssertExpectations(t)
}

func TestAuthService_UpdateProfile_UpdateError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewAuthService(mockRepo)

	existingUser := &domain.User{
		FullName: "Test User",
		Email:    "test@example.com",
	}
	existingUser.ID = 1

	newName := "New Name"
	req := &domain.UserUpdate{
		FullName: &newName,
	}

	mockRepo.On("GetByID", uint(1)).Return(existingUser, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.User")).Return(errors.New("db error"))

	user, err := service.UpdateProfile(1, req)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

// ========== Tests para DeleteAccount ==========

func TestAuthService_DeleteAccount_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewAuthService(mockRepo)

	existingUser := &domain.User{
		FullName: "Test User",
		Email:    "test@example.com",
	}
	existingUser.ID = 1

	mockRepo.On("GetByID", uint(1)).Return(existingUser, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	err := service.DeleteAccount(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_DeleteAccount_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewAuthService(mockRepo)

	mockRepo.On("GetByID", uint(999)).Return(nil, nil)

	err := service.DeleteAccount(999)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "usuario no encontrado")
	mockRepo.AssertExpectations(t)
}

func TestAuthService_DeleteAccount_GetByIDError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewAuthService(mockRepo)

	mockRepo.On("GetByID", uint(1)).Return(nil, errors.New("db error"))

	err := service.DeleteAccount(1)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_DeleteAccount_DeleteError(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	service := NewAuthService(mockRepo)

	existingUser := &domain.User{
		FullName: "Test User",
		Email:    "test@example.com",
	}
	existingUser.ID = 1

	mockRepo.On("GetByID", uint(1)).Return(existingUser, nil)
	mockRepo.On("Delete", uint(1)).Return(errors.New("db error"))

	err := service.DeleteAccount(1)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
