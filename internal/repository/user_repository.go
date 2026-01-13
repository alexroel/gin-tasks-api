package repository

import (
	"errors"

	"github.com/alexroel/gin-tasks-api/internal/config"
	"github.com/alexroel/gin-tasks-api/internal/domain"
	"gorm.io/gorm"
)

// UserRepository define las operaciones de base de datos para usuarios
type UserRepository interface {
	Create(user *domain.User) error
	GetAll() ([]domain.User, error)
	GetByID(id uint) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id uint) error
	ExistsByEmail(email string) (bool, error)
}

// userRepository implementa UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository crea una nueva instancia de UserRepository
func NewUserRepository() UserRepository {
	return &userRepository{db: config.DB}
}

// Create crea un nuevo usuario en la base de datos
func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

// GetAll obtiene todos los usuarios
func (r *userRepository) GetAll() ([]domain.User, error) {
	var users []domain.User
	err := r.db.Find(&users).Error
	return users, err
}

// GetByID obtiene un usuario por su ID
func (r *userRepository) GetByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// GetByEmail obtiene un usuario por su email
func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// Update actualiza un usuario existente
func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

// Delete elimina un usuario por su ID (soft delete)
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&domain.User{}, id).Error
}

// ExistsByEmail verifica si existe un usuario con el email dado
func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}
