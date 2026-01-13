package repository

import (
	"context"
	"errors"

	"github.com/alexroel/gin-tasks-api/internal/config"
	"github.com/alexroel/gin-tasks-api/internal/domain"
	"gorm.io/gorm"
)

// TaskRepository define las operaciones de base de datos para tareas
type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	GetAll(ctx context.Context) ([]domain.Task, error)
	GetByID(ctx context.Context, id uint) (*domain.Task, error)
	GetByUserID(ctx context.Context, userID uint) ([]domain.Task, error)
	Update(ctx context.Context, task *domain.Task) error
	Delete(ctx context.Context, id uint) error
	UpdateStatus(ctx context.Context, id uint, completed bool) error
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
func (r *taskRepository) Create(ctx context.Context, task *domain.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

// GetAll obtiene todas las tareas
func (r *taskRepository) GetAll(ctx context.Context) ([]domain.Task, error) {
	var tasks []domain.Task
	err := r.db.WithContext(ctx).Find(&tasks).Error
	return tasks, err
}

// GetByID obtiene una tarea por su ID
func (r *taskRepository) GetByID(ctx context.Context, id uint) (*domain.Task, error) {
	var task domain.Task
	err := r.db.WithContext(ctx).First(&task, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &task, err
}

// GetByUserID obtiene todas las tareas de un usuario
func (r *taskRepository) GetByUserID(ctx context.Context, userID uint) ([]domain.Task, error) {
	var tasks []domain.Task
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&tasks).Error
	return tasks, err
}

// Update actualiza una tarea existente
func (r *taskRepository) Update(ctx context.Context, task *domain.Task) error {
	return r.db.WithContext(ctx).Save(task).Error
}

// Delete elimina una tarea por su ID (soft delete)
func (r *taskRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Task{}, id).Error
}

// UpdateStatus actualiza el estado de completado de una tarea
func (r *taskRepository) UpdateStatus(ctx context.Context, id uint, completed bool) error {
	return r.db.WithContext(ctx).Model(&domain.Task{}).Where("id = ?", id).Update("completed", completed).Error
}
