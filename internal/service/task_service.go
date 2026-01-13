package service

import (
	"context"
	"errors"

	"github.com/alexroel/gin-tasks-api/internal/domain"
	"github.com/alexroel/gin-tasks-api/internal/repository"
)

var (
	ErrTaskNotFound    = errors.New("tarea no encontrada")
	ErrTaskUnauthorized = errors.New("no tienes permiso para acceder a esta tarea")
)

// TaskService define las operaciones de negocio para tareas
type TaskService interface {
	Create(ctx context.Context, userID uint, req *domain.CreateTask) (*domain.Task, error)
	GetAll(ctx context.Context) ([]domain.Task, error)
	GetByID(ctx context.Context, id uint) (*domain.Task, error)
	GetByUserID(ctx context.Context, userID uint) ([]domain.Task, error)
	Update(ctx context.Context, id, userID uint, req *domain.UpdateTask) (*domain.Task, error)
	Delete(ctx context.Context, id, userID uint) error
	UpdateStatus(ctx context.Context, id, userID uint, completed bool) (*domain.Task, error)
}

type taskService struct {
	repo repository.TaskRepository
}

// NewTaskService crea una nueva instancia de TaskService
func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

// Create crea una nueva tarea para un usuario
func (s *taskService) Create(ctx context.Context, userID uint, req *domain.CreateTask) (*domain.Task, error) {
	task := &domain.Task{
		Title:  req.Title,
		UserID: userID,
	}

	if err := s.repo.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

// GetAll obtiene todas las tareas
func (s *taskService) GetAll(ctx context.Context) ([]domain.Task, error) {
	return s.repo.GetAll(ctx)
}

// GetByID obtiene una tarea por su ID
func (s *taskService) GetByID(ctx context.Context, id uint) (*domain.Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

// GetByUserID obtiene todas las tareas de un usuario
func (s *taskService) GetByUserID(ctx context.Context, userID uint) ([]domain.Task, error) {
	return s.repo.GetByUserID(ctx, userID)
}

// Update actualiza una tarea existente
func (s *taskService) Update(ctx context.Context, id, userID uint, req *domain.UpdateTask) (*domain.Task, error) {
	// Obtener tarea existente
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}

	// Verificar que la tarea pertenece al usuario
	if task.UserID != userID {
		return nil, ErrTaskUnauthorized
	}

	// Actualizar campos si se proporcionan
	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Completed != nil {
		task.Completed = *req.Completed
	}

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

// Delete elimina una tarea por su ID
func (s *taskService) Delete(ctx context.Context, id, userID uint) error {
	// Obtener tarea existente
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if task == nil {
		return ErrTaskNotFound
	}

	// Verificar que la tarea pertenece al usuario
	if task.UserID != userID {
		return ErrTaskUnauthorized
	}

	return s.repo.Delete(ctx, id)
}

// UpdateStatus actualiza el estado de completado de una tarea
func (s *taskService) UpdateStatus(ctx context.Context, id, userID uint, completed bool) (*domain.Task, error) {
	// Obtener tarea existente
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}

	// Verificar que la tarea pertenece al usuario
	if task.UserID != userID {
		return nil, ErrTaskUnauthorized
	}

	// Actualizar estado
	if err := s.repo.UpdateStatus(ctx, id, completed); err != nil {
		return nil, err
	}

	task.Completed = completed
	return task, nil
}
