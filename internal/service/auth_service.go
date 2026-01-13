package service

import (
	"context"
	"errors"

	"github.com/alexroel/gin-tasks-api/internal/config"
	"github.com/alexroel/gin-tasks-api/internal/domain"
	"github.com/alexroel/gin-tasks-api/internal/repository"
	"github.com/alexroel/gin-tasks-api/pkg/jwt"
	"github.com/alexroel/gin-tasks-api/pkg/utils"
)

// AuthServiceInterface define las operaciones del servicio de autenticación
type AuthServiceInterface interface {
	Register(ctx context.Context, req *domain.UserCreate) (*domain.User, error)
	Login(ctx context.Context, req *domain.UserLogin) (string, *domain.User, error)
	GetUserByID(ctx context.Context, userID uint) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID uint, req *domain.UserUpdate) (*domain.User, error)
	DeleteAccount(ctx context.Context, userID uint) error
}

type AuthService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

// Register registra un nuevo usuario
func (s *AuthService) Register(ctx context.Context, req *domain.UserCreate) (*domain.User, error) {
	// Validar si el usuario ya existe
	ok, err := s.repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if ok {
		return nil, errors.New("Ya existe un usuario con ese email")
	}

	// Hashear la contraseña
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	// Asignar la contraseña hasheada al usuario
	user := &domain.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: hashedPassword,
	}

	// Crear el usuario
	err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login autentica a un usuario
func (s *AuthService) Login(ctx context.Context, req *domain.UserLogin) (string, *domain.User, error) {
	// Buscar el usuario por email
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return "", nil, errors.New("Credenciales inválidas")
	}
	// Verificar la contraseña
	if !utils.CheckPassword(user.Password, req.Password) {
		return "", nil, errors.New("Credenciales inválidas")
	}
	// Generar token
	token, err := jwt.GenerateToken(user.ID, user.Email, config.AppConfig.JWTSecret, config.AppConfig.JWTExpireIn)
	if err != nil {
		return "", nil, errors.New("Error al generar Token")
	}

	return token, user, nil
}

// GetUserByID obtiene un usuario por su ID
func (s *AuthService) GetUserByID(ctx context.Context, userID uint) (*domain.User, error) {
	return s.repo.GetByID(ctx, userID)
}

// UpdateProfile actualiza el perfil del usuario autenticado
func (s *AuthService) UpdateProfile(ctx context.Context, userID uint, req *domain.UserUpdate) (*domain.User, error) {
	// Obtener usuario existente
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("usuario no encontrado")
	}

	// Actualizar campos si se proporcionan
	if req.FullName != nil {
		user.FullName = *req.FullName
	}

	if req.Email != nil && *req.Email != user.Email {
		// Verificar si el nuevo email ya existe
		exists, err := s.repo.ExistsByEmail(ctx, *req.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("el email ya está registrado")
		}
		user.Email = *req.Email
	}

	if req.Password != nil {
		hashedPassword, err := utils.HashPassword(*req.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteAccount elimina la cuenta del usuario autenticado
func (s *AuthService) DeleteAccount(ctx context.Context, userID uint) error {
	// Verificar que el usuario existe
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("usuario no encontrado")
	}

	return s.repo.Delete(ctx, userID)
}
