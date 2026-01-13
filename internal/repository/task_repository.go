package repository

import (
	"errors"

	"github.com/alexroel/gin-tasks-api/internal/config"
	"github.com/alexroel/gin-tasks-api/internal/domain"
	"gorm.io/gorm"
)

// TaskRepository define las operaciones de base de datos para tareas
type TaskRepository interface {
	Create(task *domain.Task) error
	GetAll() ([]domain.Task, error)
	GetByID(id uint) (*domain.Task, error)
	GetByUserID(userID uint) ([]domain.Task, error)
	Update(task *domain.Task) error
	Delete(id uint) error
	UpdateStatus(id uint, completed bool) error
}

// taskRepository implementa TaskRepository
type taskRepository struct {
	db *gorm.DB
}

// NewTaskRepository crea una nueva instancia de TaskRepository
func NewTaskRepository() TaskRepository {
	return &taskRepository{db: config.DB}
}

// Create crea una nueva tarea en la base de datos
func (r *taskRepository) Create(task *domain.Task) error {
	return r.db.Create(task).Error
}

// GetAll obtiene todas las tareas
func (r *taskRepository) GetAll() ([]domain.Task, error) {
	var tasks []domain.Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}

// GetByID obtiene una tarea por su ID
func (r *taskRepository) GetByID(id uint) (*domain.Task, error) {
	var task domain.Task
	err := r.db.First(&task, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &task, err
}

// GetByUserID obtiene todas las tareas de un usuario
func (r *taskRepository) GetByUserID(userID uint) ([]domain.Task, error) {
	var tasks []domain.Task
	err := r.db.Where("user_id = ?", userID).Find(&tasks).Error
	return tasks, err
}

// Update actualiza una tarea existente
func (r *taskRepository) Update(task *domain.Task) error {
	return r.db.Save(task).Error
}

// Delete elimina una tarea por su ID (soft delete)
func (r *taskRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Task{}, id).Error
}

// UpdateStatus actualiza el estado de completado de una tarea
func (r *taskRepository) UpdateStatus(id uint, completed bool) error {
	return r.db.Model(&domain.Task{}).Where("id = ?", id).Update("completed", completed).Error
}
